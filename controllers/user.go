/*
Copyright 2014 Jiang Le

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/naokij/go-sendcloud"
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
	"time"
)

type UserEditForm struct {
	Username    string `form:"Username,text"valid:"Required"`
	Email       string `form:"Email,text"valid:"Required;Email;"`
	PublicEmail bool   `form:"PublicEmail,checkbox""`
	Nickname    string `form:"Nickname,text"`
	Info        string `form:"Info,textarea"`
	Company     string `form:"Company,text"`
	Location    string `form:"Location,text"`
	Url         string `form:"Url,text"`
	Qq          string `form:"Qq,text"`
	WeChat      string `form:"WeChat,file"`
	Weibo       string `form:"Weibo,text"`
}

type UserPasswordForm struct {
	CurrentPassword string `Valid:"Required"`
	Password        string `Valid:"MinSize(6);"`
	PasswordConfirm string `Valid:"Required;"`
}

//用户控制器
type UserController struct {
	BaseController
}

func (this *UserController) Profile() {
	this.Layout = "layout.html"
	this.TplNames = "user_profile.html"
	user, err := this.getUserFromRequest()
	if err != nil {
		this.Abort("404")
	}
	this.Data["TheUser"] = &user
	return
}

func (this *UserController) Edit() {
	this.Layout = "layout.html"
	this.TplNames = "user_edit.html"
	this.Data["PageTitle"] = fmt.Sprintf("用户设置 | %s", setting.AppName)
	//获取用户，并且判断是否有权限执行此操作
	user, err := this.getUserFromRequest()
	if err != nil {
		this.Abort("404")
	}
	if !this.canEdit(user) {
		this.Abort("403")
	}
	if this.Ctx.Input.IsPost() {
		action := this.GetString("action")
		switch action {
		case "UpdateUser":
			this.processUserEditForm(&user)
		case "UpdatePassword":
			this.processUserPasswordForm(&user)
		case "UploadAvatar":
			this.processUploadAvatar(&user)
		}
	}
	this.Data["TheUser"] = &user
	if this.Data["UserEditForm"] == nil {
		this.Data["UserEditForm"] = &user
	}

}

func (this *UserController) getUserFromRequest() (user models.User, err error) {
	user = models.User{Username: this.Ctx.Input.Param(":username")}
	if errRead := user.Read("Username"); errRead != nil {
		err = errors.New("404")
	}
	return user, err
}

func (this *UserController) canEdit(user models.User) bool {
	if !this.IsLogin {
		return false
	} else if user.Id != this.User.Id && !this.User.IsAdmin {
		return false
	}
	return true
}

func (this *UserController) processUserEditForm(user *models.User) {
	valid := validation.Validation{}
	var usernameChanged, emailChanged bool
	userEditForm := UserEditForm{}
	if err := this.ParseForm(&userEditForm); err != nil {
		beego.Error(err)
	}
	_, err := valid.Valid(userEditForm)
	if err != nil {
		beego.Error(err)
		this.Abort("400")
	}
	if user.Username != userEditForm.Username {
		usernameChanged = true
		if time.Since(user.Created).Hours() <= 720 {
			tmpUser := models.User{Username: userEditForm.Username}
			if err := tmpUser.ValidUsername(); err != nil {
				valid.SetError("Username", err.Error())
			}
			if tmpUser.Read("Username") == nil {
				valid.SetError("Username", "用户名已经被使用")
			}
		} else {
			valid.SetError("Username", "注册超过30天后无法修改用户名")
		}
	}
	if user.Email != userEditForm.Email {
		emailChanged = true
		tmpUser := models.User{Email: userEditForm.Email}
		if err := tmpUser.Read("Email"); err == nil {
			valid.SetError("Email", "电子邮件地址已经被使用")
		}
	}
	user.Url = userEditForm.Url
	if err := user.ValidateUrl(); user.Url != "" && err != nil {
		valid.SetError("Url", err.Error())
	}
	this.Data["UserEditForm"] = &userEditForm
	if len(valid.Errors) > 0 {
		this.Data["UserEditFormValidErrors"] = valid.Errors
		beego.Trace(fmt.Sprint(valid.Errors))
	} else {
		if usernameChanged {
			user.Username = userEditForm.Username
		}
		if emailChanged {
			user.Email = userEditForm.Email
			user.IsActive = false
		}
		user.PublicEmail = userEditForm.PublicEmail
		user.Nickname = userEditForm.Nickname
		user.Info = userEditForm.Info
		user.Company = userEditForm.Company
		user.Location = userEditForm.Location
		user.Url = userEditForm.Url
		user.Qq = userEditForm.Qq
		user.Weibo = userEditForm.Weibo
		user.WeChat = userEditForm.WeChat
		if err := user.Update(); err != nil {
			this.Abort("500")
		}
		if usernameChanged && this.User.Id == user.Id {
			this.LogUserIn(user, false)
		}
		if emailChanged {
			//发验证邮件
			this.resendValidation(user)
			this.FlashWrite("notice", fmt.Sprintf("资料已经更新。由于修改了Email地址，我们向%s发送了一封验证邮件，请重新验证。", user.Email))
		} else {
			this.FlashWrite("notice", "资料已更新！")
		}
		redirectUrl := beego.UrlFor("UserController.Edit", ":username", user.Username)
		this.Redirect(redirectUrl, 302)
	}
}

func (this *UserController) processUserPasswordForm(user *models.User) {
	valid := validation.Validation{}
	userPasswordForm := UserPasswordForm{}
	if err := this.ParseForm(&userPasswordForm); err != nil {
		beego.Error(err)
	}
	_, err := valid.Valid(userPasswordForm)
	if err != nil {
		beego.Error(err)
		this.Abort("400")
	}
	if !user.VerifyPassword(userPasswordForm.CurrentPassword) {
		valid.SetError("CurrentPassword", "当前密码错误")
	}
	if len(valid.Errors) > 0 {
		this.Data["UserPasswordFormValidErrors"] = valid.Errors
		beego.Trace(fmt.Sprint(valid.Errors))
	} else {
		user.SetPassword(userPasswordForm.Password)
		if err := user.Update(); err != nil {
			this.Abort("500")
		}
		this.FlashWrite("notice", "密码已更新！")
		this.Redirect(this.Ctx.Request.RequestURI, 302)
	}
}

func (this *UserController) processUploadAvatar(user *models.User) {
	valid := validation.Validation{}
	avatarFile, header, err := this.GetFile("Avatar")
	if err != nil {
		this.Abort("400")
	}
	err = user.ValidateAndSetAvatar(avatarFile, header.Filename)
	if err != nil {
		valid.SetError("Avatar", err.Error())
		this.Data["UserAvatarFormValidErrors"] = valid.Errors
	} else {
		if err := user.Update("Avatar"); err != nil {
			this.Abort("500")
		}
		this.FlashWrite("notice", "头像已更新！")
		this.Redirect(this.Ctx.Request.RequestURI, 302)
	}
}

func (this *UserController) resendValidation(user *models.User) {
	//发验证邮件
	sub := sendcloud.NewSubstitution()
	sub.AddTo(user.Email)
	sub.AddSub("%name%", user.Username)
	sub.AddSub("%appname%", setting.AppName)
	code, err := user.GenerateActivateCode()
	if err != nil {
		beego.Trace(err)
		this.Abort("500")
	}
	sub.AddSub("%url%", setting.AppUrl+beego.UrlFor("AuthController.Activate", ":code", code))
	if err := setting.Sendcloud.SendTemplate("gotalk_revalidate", setting.AppName+"邮件验证", setting.From, setting.FromName, sub); err != nil {
		beego.Error(err)
	}
}

func (this *UserController) ResendValidation() {
	//获取用户，并且判断是否有权限执行此操作
	user, err := this.getUserFromRequest()
	if err != nil {
		this.Abort("404")
	}
	if !this.canEdit(user) {
		this.Abort("403")
	}
	this.resendValidation(&user)
	this.FlashWrite("notice", fmt.Sprintf("验证邮件已经发送，请登录%s进行验证。", user.Email))
	redirectUrl := beego.UrlFor("UserController.Edit", ":username", user.Username)
	this.Redirect(redirectUrl, 302)
}

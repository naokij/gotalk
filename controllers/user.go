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
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
)

type UserEditForm struct {
	Email       string `form:"Email,text"valid:"Required;Email;"`
	PublicEmail bool   `form:"PublicEmail,checkbox""`
	Nickname    string `form:"Nickname,text"`
	Info        string `form:"Info,textarea"`
	Company     string `form:"Company,text"`
	Location    string `form:"Location,text"`
	Url         string `form:"Url,text"`
	Qq          int    `form:"Qq,text"`
	WeChat      string `form:"WeChat,file"`
	Weibo       string `form:"Weibo,text"`
}

//用户控制器
type UserController struct {
	BaseController
}

func (this *UserController) Profile() {
	this.Layout = "layout.html"
	this.TplNames = "user_profile.html"
	user := models.User{Username: this.Ctx.Input.Param(":username")}
	if err := user.Read("Username"); err != nil {
		this.Abort("404")
	}
	this.Data["User"] = &user
	return
}

func (this *UserController) Edit() {
	this.Layout = "layout.html"
	this.TplNames = "user_edit.html"
	this.Data["PageTitle"] = fmt.Sprintf("用户设置 | %s", setting.AppName)
	//获取用户，并且判断是否有权限执行此操作
	user, err := this.getUserFromRequest()
	if err != nil {
		this.Abort(err.Error())
	}
	this.Data["TheUser"] = &user
	if this.Ctx.Input.IsPost() {
		action := this.GetString("action")
		switch action {
		case "UpdateUser":
			this.processUserEditForm(&user)
		case "UpdatePassword":
			this.processUserPasswordForm(&user)
		}
	}
	if this.Data["UserEditForm"] == nil {
		this.Data["UserEditForm"] = &user
	}

}

func (this *UserController) getUserFromRequest() (user models.User, err error) {
	user = models.User{Username: this.Ctx.Input.Param(":username")}
	if errRead := user.Read("Username"); errRead != nil {
		err = errors.New("404")
	}
	if !this.IsLogin {
		err = errors.New("403")
		beego.Trace("Not logged in, can not edit user info.")
	} else if user.Id != this.User.Id && !this.User.IsAdmin {
		beego.Trace("Can't edit this user, not owner or admin!")
		err = errors.New("403")
	}
	return user, err
}

func (this *UserController) processUserEditForm(user *models.User) {
	valid := validation.Validation{}
	userEditForm := UserEditForm{}
	if err := this.ParseForm(&userEditForm); err != nil {
		beego.Error(err)
	}
	_, err := valid.Valid(userEditForm)
	if err != nil {
		beego.Error(err)
		this.Abort("400")
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
		user.Email = userEditForm.Email
		user.PublicEmail = userEditForm.PublicEmail
		user.Nickname = userEditForm.Nickname
		user.Info = userEditForm.Info
		user.Company = userEditForm.Company
		user.Location = userEditForm.Location
		user.Url = userEditForm.Url
		user.Qq = userEditForm.Qq
		user.Weibo = userEditForm.Weibo
		user.WeChat = userEditForm.WeChat
		user.Update()
		beego.Trace("User info updated!")
	}
}

func (this *UserController) processUserPasswordForm(user *models.User) {

}

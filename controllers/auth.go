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
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/naokij/go-sendcloud"
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
	"net/url"
)

//登录表单
type LoginForm struct {
	Username string `form:"username,text"valid:"Required;"`
	Password string `form:"password,password"valid:"Required;"`
	Remember string `form:"remember,text"`
}

//注册表单
type RegisterForm struct {
	Username        string `form:"username,text"valid:"Required;"`
	Password        string `form:"password,password"valid:"Required;MinSize(6);"`
	PasswordConfirm string `form:"password_confirm,password"valid:"Required;"`
	Email           string `form:"email,text"valid:"Required;Email;"`
	CaptchaId       string `form:"captcha-id,hidden"valid:"Required;"`
	Captcha         string `form:"captcha,text"valid:"Required"`
}

//登录控制器
type AuthController struct {
	BaseController
}

func (this *AuthController) Login() {
	returnUrl := this.Ctx.Input.Param(":returnurl")
	if returnUrl != "" {
		u, err := url.Parse(returnUrl)
		if err == nil {
			if u.Host == setting.AppHost {
				this.SetSession("ReturnUrl", returnUrl)
			}
		}
	}
	this.Data["PageTitle"] = fmt.Sprintf("登录 | %s", setting.AppName)
	this.Layout = "layout.html"
	this.TplNames = "login.html"
}

//显示登录页面
func (this *AuthController) loginPageWithErrors(form LoginForm, errors []*validation.ValidationError) {
	this.Data["PageTitle"] = fmt.Sprintf("登录 | %s", setting.AppName)
	this.Layout = "layout.html"
	this.TplNames = "login.html"
	this.Data["form"] = form
	this.Data["errors"] = errors
	this.Data["HasError"] = true
	beego.Trace(errors[0])
}
func (this *AuthController) DoLogin() {
	valid := validation.Validation{}
	form := LoginForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	b, err := valid.Valid(form)
	if err != nil {
		beego.Error(err)
	}
	if !b {
		this.loginPageWithErrors(form, valid.Errors)
		return
	}
	//用户不存在？
	user := models.User{Username: form.Username, Email: form.Username}
	if err := user.Read("Username"); err != nil {
		if err2 := user.Read("Email"); err2 != nil {
			errMsg := fmt.Sprintf("用户 %s 不存在!", form.Username)
			beego.Trace(errMsg)
			valid.SetError("username", errMsg)
			this.loginPageWithErrors(form, valid.Errors)
			return
		}
	}
	//用户被禁止?
	if user.IsBanned {
		beego.Trace(fmt.Sprintf("用户%s被禁用，不能登录！", user.Username))
		valid.SetError("username", "抱歉，您被禁止登录！")
		this.loginPageWithErrors(form, valid.Errors)
		return
	}
	//检查密码
	if !user.VerifyPassword(form.Password) {
		beego.Trace(fmt.Sprintf("%s 登录失败！", form.Username))
		valid.SetError("password", "密码错误")
		this.loginPageWithErrors(form, valid.Errors)
		return
	}
	//验证全部通过
	var remember bool
	if form.Remember != "" {
		remember = true
	}
	this.LogUserIn(&user, remember)
	if returnUrl, ok := this.GetSession("ReturnUrl").(string); returnUrl != "" && ok {
		this.DelSession("ReturnUrl")
		this.Redirect(returnUrl, 302)
	} else {
		this.Redirect("/", 302)
	}
	return
}

func (this *AuthController) Logout() {
	this.LogUserOut()
	this.Redirect("/", 302)
	return
}

func (this *AuthController) Register() {
	this.Data["PageTitle"] = fmt.Sprintf("注册 | %s", setting.AppName)
	this.Layout = "layout.html"
	this.TplNames = "register.html"
}

//注册页面，现实错误信息
func (this *AuthController) registerPageWithErrors(form RegisterForm, errors []*validation.ValidationError) {
	this.Data["PageTitle"] = fmt.Sprintf("注册 | %s", setting.AppName)
	this.Layout = "layout.html"
	this.TplNames = "register.html"
	this.Data["form"] = form
	this.Data["errors"] = errors
	this.Data["HasError"] = true
	beego.Trace(errors[0])
}

func (this *AuthController) DoRegister() {
	this.Layout = "layout.html"
	this.TplNames = "register.html"
	valid := validation.Validation{}
	form := RegisterForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	b, err := valid.Valid(form)
	if err != nil {
		beego.Error(err)
	}
	if !b {
		this.registerPageWithErrors(form, valid.Errors)
		return
	}
	//验证用户名
	user := models.User{Username: form.Username}
	if err := user.ValidUsername(); err != nil {
		valid.SetError("username", err.Error())
		this.registerPageWithErrors(form, valid.Errors)
		return
	} else {
		if user.Read("Username") == nil {
			valid.SetError("username", fmt.Sprintf("%s已被使用，请使用其他用户名！", form.Username))
			this.registerPageWithErrors(form, valid.Errors)
			return
		}
	}
	//验证email未被注册
	user.Email = form.Email
	if user.Read("Email") == nil {
		valid.SetError("email", "已被使用，请直接使用此电邮登录")
		this.registerPageWithErrors(form, valid.Errors)
		return
	}
	//通过所有验证
	actCode, _ := user.GenerateActivateCode()
	user.SetPassword(form.Password)
	if err := user.Insert(); err != nil {
		beego.Error(err)
		this.Abort("500")
		return
	}
	sub := sendcloud.NewSubstitution()
	sub.AddTo(user.Email)
	sub.AddSub("%appname%", setting.AppName)
	sub.AddSub("%name%", user.Username)
	sub.AddSub("%url%", setting.AppUrl+beego.UrlFor("AuthController.Activate", ":code", actCode))
	if err := setting.Sendcloud.SendTemplate("gotalk_register", setting.AppName+"欢迎你", setting.From, setting.FromName, sub); err != nil {
		beego.Error(err)
	}
	this.FlashWrite("notice", fmt.Sprintf("注册成功！欢迎你, %s。建议你再花点时间上传头像、验证电子邮件！", user.Username))
	this.LogUserIn(&user, false)
	userEditUrl := beego.UrlFor("UserController.Edit", ":username", user.Username)
	this.Redirect(userEditUrl, 302)
	return
}

func (this *AuthController) ValidateUsername() {
	username := this.GetString("username")
	user := models.User{Username: username}
	if err := user.ValidUsername(); err != nil {
		this.Data["json"] = err.Error()
	} else {
		if user.Read("Username") == nil {
			//这个用户名已经存在
			this.Data["json"] = fmt.Sprintf("%s已被使用，请使用其他用户名！", username)
		} else {
			this.Data["json"] = true
		}
	}
	this.ServeJson()
}

func (this *AuthController) ValidateEmail() {
	email := this.GetString("email")
	user := models.User{Email: email}
	if user.Read("Email") == nil {
		this.Data["json"] = "已被使用，请直接使用此电邮登录"
	} else {
		this.Data["json"] = true
	}
	this.ServeJson()
}

func (this *AuthController) ValidateCaptcha() {
	captcha := this.GetString("captcha")
	captchaId := this.GetString("captchaid")
	this.Data["json"] = setting.Captcha.Verify(captchaId, captcha)
	this.ServeJson()
}

func (this *AuthController) Activate() {
	this.Data["PageTitle"] = fmt.Sprintf("用户激活 | %s", setting.AppName)
	code := this.Ctx.Input.Param(":code")
	user := models.User{}
	if user.VerifyActivateCode(code) {
		user.IsActive = true
		user.Update()
		this.FlashWrite("notice", "谢谢，你的电子邮件已经验证！")
	} else {
		this.FlashWrite("error", "糟糕，无法验证你的电子邮件！")
	}
	this.Redirect("/", 302)
}

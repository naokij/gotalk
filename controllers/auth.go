package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/naokij/gotalk/models"
)

//登录表单
type LoginForm struct {
	Username string `form:"username,text"valid:"Required;"`
	Password string `form:"password,password"valid:"Required;"`
	Remember string `form:"remember,text"`
}

//登录控制器
type AuthController struct {
	BaseController
}

func (this *AuthController) Login() {
	this.Layout = "layout.html"
	this.TplNames = "login.html"
}

//显示登录页面
func (this *AuthController) loginPageWithErrors(form LoginForm, errors []*validation.ValidationError) {
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
		// validation does not pass
		// blabla...
		this.Data["errors"] = valid.Errors
		for _, err := range valid.Errors {
			beego.Info(err.Key, err.Message)
		}
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
	this.Redirect("/", 302)
	return
}

func (this *AuthController) Logout() {
	this.LogUserOut()
	this.Redirect("/", 302)
	return
}

func (this *AuthController) Register() {

}

func (this *AuthController) DoRegister() {

}

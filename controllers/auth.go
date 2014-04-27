package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/naokij/gotalk/models"
)

//登录表单
type loginForm struct {
	Username string `form:"username,text"valid:"Required;"`
	Password string `form:"password,password"valid:"Required;"`
	Remember string `form:"remember,text"`
}

//登录控制器
type AuthController struct {
	BaseController
}

func (this *AuthController) Get() {
	this.Layout = "layout.html"
	this.TplNames = "login.html"
	if _, ok := this.FlashRead("error"); ok {
		this.Data["HasError"] = true
	}
	this.Data["loginForm"] = this.GetSession("loginForm")
	beego.Trace(fmt.Sprint(this.GetSession("loginForm")))
}

func (this *AuthController) Login() {
	flash := beego.NewFlash()
	valid := validation.Validation{}
	form := loginForm{}
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
		for _, err := range valid.Errors {
			beego.Info(err.Key, err.Message)
		}
		this.Redirect("/login", 302)
		return
	}
	//用户不存在？
	user := models.User{Username: form.Username, Email: form.Username}
	if err := user.Read("Username"); err != nil {
		if err2 := user.Read("Email"); err2 != nil {
			beego.Trace(fmt.Sprintf("用户 %s 不存在!", form.Username))
			flash.Error(fmt.Sprintf("用户 %s 不存在!", form.Username))
			flash.Store(&this.Controller)
			this.SetSession("loginForm", form)
			this.Redirect("/login", 302)
			return
		}
	}
	//用户被禁止?
	if user.IsBanned {
		beego.Trace(fmt.Sprintf("用户%s被禁用，不能登录！", user.Username))
		flash.Error(fmt.Sprintf("抱歉，您被禁止登录！", user.Username))
		flash.Store(&this.Controller)
		this.SetSession("loginForm", form)
		this.Redirect("/login", 302)
		return
	}
	//检查密码
	if !user.VerifyPassword(form.Password) {
		beego.Trace(fmt.Sprintf("%s 登录失败！", form.Username))
		flash.Error("密码错误！")
		flash.Store(&this.Controller)
		this.SetSession("loginForm", "good")
		this.Redirect("/login", 302)
		return
	}
	//验证全部通过
	this.Ctx.WriteString(fmt.Sprint(user))

}

func (this *AuthController) Logout() {

}

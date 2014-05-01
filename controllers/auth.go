package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
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
	this.Redirect("/", 302)
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
	beego.Trace(user)
	user.SetPassword(form.Password)
	if err := user.Insert(); err != nil {
		beego.Error(err)
		this.Abort("500")
		return
	}
	this.Redirect("/welcome", 302)
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

func (this *AuthController) Welcome() {

}

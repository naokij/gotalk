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
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
	"github.com/naokij/go-sendcloud"
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
	"github.com/naokij/gotalk/utils"
	"github.com/naokij/social-auth"
	"github.com/naokij/social-auth/apps"
	//"net/http"
)

var (
	GithubAuth *apps.Github
	QQAuth     *apps.QQ
	WeiboAuth  *apps.Weibo
	SocialAuth *social.SocialAuth
)

type SocialAuthLoginForm struct {
	Username string `form:"Username,text"valid:"Required;"`
	Password string `form:"Password,password"valid:"Required;"`
}

type SocialAuthRegisterForm struct {
	Username        string `form:"Username,text"valid:"Required;"`
	Password        string `form:"Password,password"valid:"Required;MinSize(6);"`
	PasswordConfirm string `form:"PasswordConfirm,password"valid:"Required;"`
	Email           string `form:"Email,text"valid:"Required;Email;"`
}

type socialAuther struct {
}

func (p *socialAuther) IsUserLogin(ctx *context.Context) (int, bool) {
	if username, ok := ctx.Input.CruSession.Get("AuthUsername").(string); username != "" && ok {
		user := models.User{Username: username}
		if user.Read("Username") == nil {
			return user.Id, true
		}
	}
	return 0, false
}

func (p *socialAuther) LoginUser(ctx *context.Context, uid int) (string, error) {
	user := models.User{Id: uid}
	if user.Read() == nil {
		ctx.Input.CruSession.Set("AuthUsername", user.Username)
	}
	return GetLoginRedirectUrl(ctx), nil
}

var SocialAuther social.SocialAuther = new(socialAuther)

//社交帐号登录控制器
type SocialAuthController struct {
	BaseController
}

func OAuthRedirect(ctx *context.Context) {
	redirect, err := SocialAuth.OAuthRedirect(ctx)
	if err != nil {
		beego.Error("OAuthRedirect", err)
	}

	if len(redirect) > 0 {
		ctx.Redirect(302, redirect)
	}
}

func OAuthAccess(ctx *context.Context) {
	redirect, _, err := SocialAuth.OAuthAccess(ctx)
	if err != nil {
		beego.Error("OAuthAccess", err)
	}

	if len(redirect) > 0 {
		ctx.Redirect(302, redirect)
	}
}

func (this *SocialAuthController) Connect() {
	this.Data["PageTitle"] = fmt.Sprintf("社交帐号登录 | %s", setting.AppName)
	this.Layout = "layout.html"
	this.TplNames = "social-login.html"
	if this.IsLogin {
		this.Redirect("/", 302)
	}
	//检查社交帐号登录是否正常
	var socialType social.SocialType
	if !this.canConnect(&socialType) {
		beego.Error(this.GetString("error_description"))
		this.Abort("500")
		this.Redirect(SocialAuth.LoginURL, 302)
		return
	}
	p, _ := social.GetProviderByType(socialType)
	if p == nil {
		beego.Error("unknown provider")
	}
	var socialUserLogin, socialUserEmail, socialUserAvatarUrl string
	var ok bool
	if socialUserLogin, ok = this.GetSession("social_user_login").(string); !ok {
		beego.Error("error while reading session ")
		this.Abort("500")
	}
	if socialUserEmail, ok = this.GetSession("social_user_email").(string); !ok {
		beego.Error("error while reading session ")
		this.Abort("500")
	}
	if socialUserAvatarUrl, ok = this.GetSession("social_user_avatar_url").(string); !ok {
		beego.Error("error while reading session ")
		this.Abort("500")
	}
	this.Data["SocialType"] = p.GetName()
	this.Data["SocialUserLogin"] = socialUserLogin
	this.Data["SocialUserEmail"] = socialUserEmail
	this.Data["SocialUserAvatarUrl"] = socialUserAvatarUrl
	//准备注册表格初始数据
	registerForm := SocialAuthRegisterForm{}
	var user models.User
	if this.Ctx.Input.IsGet() {
		user = models.User{Username: socialUserLogin}
		if user.Read("Username") == nil {
			registerForm.Username = socialUserLogin + utils.GetRandomString(3)
		} else {
			registerForm.Username = socialUserLogin
		}
		if socialUserEmail != "" {
			user = models.User{Email: socialUserEmail}
			if user.Read("Email") == nil {
				registerForm.Email = ""
			} else {
				registerForm.Email = socialUserEmail
			}
		}
		this.Data["RegisterForm"] = registerForm
	}

	if this.Ctx.Input.IsPost() {
		action := this.GetString("action")
		switch action {
		case "Register":
			this.processRegisterForm(socialType, registerForm, socialUserAvatarUrl)
		case "Login":
			this.processLoginForm(socialType)
		}
	}

}

func (this *SocialAuthController) processRegisterForm(socialType social.SocialType, form SocialAuthRegisterForm, socialUserAvatarUrl string) {
	valid := validation.Validation{}
	var user models.User
	var actCode string
	var sub *sendcloud.Substitution
	//var resp *http.Response
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	b, err := valid.Valid(form)
	if err != nil {
		beego.Error(err)
	}
	if !b {
		goto showRegisterErrors
	}
	//验证用户名
	user = models.User{Username: form.Username}
	if err := user.ValidUsername(); err != nil {
		valid.SetError("Username", err.Error())
		goto showRegisterErrors
	} else {
		if user.Read("Username") == nil {
			valid.SetError("Username", fmt.Sprintf("%s已被使用，请使用其他用户名！", form.Username))
			goto showRegisterErrors
		}
	}
	//验证email未被注册
	user.Email = form.Email
	if user.Read("Email") == nil {
		valid.SetError("Email", "已被使用，请直接使用此电邮登录")
		goto showRegisterErrors
	}
	//通过所有验证
	actCode, _ = user.GenerateActivateCode()
	user.SetPassword(form.Password)
	if err := user.Insert(); err != nil {
		beego.Error(err)
		this.Abort("500")
		return
	}
	sub = sendcloud.NewSubstitution()
	sub.AddTo(user.Email)
	sub.AddSub("%appname%", setting.AppName)
	sub.AddSub("%name%", user.Username)
	sub.AddSub("%url%", setting.AppUrl+beego.UrlFor("AuthController.Activate", ":code", actCode))
	if err := setting.Sendcloud.SendTemplate("gotalk_register", setting.AppName+"欢迎你", setting.From, setting.FromName, sub); err != nil {
		beego.Error(err)
	}
	//复制头像
	// if resp, err = http.Get(socialUserAvatarUrl); err != nil {
	// 	beego.Error(fmt.Sprintf("Error opening url:%s", socialUserAvatarUrl))
	// 	this.Abort("500")
	// 	return
	// }
	// defer resp.Body.Close()
	// user.ValidateAndSetAvatar(resp.Body, filename)

	this.FlashWrite("notice", fmt.Sprintf("注册成功！欢迎你, %s。建议你再花点时间验证电子邮件！", user.Username))
	if loginRedirect, _, err := SocialAuth.ConnectAndLogin(this.Ctx, socialType, user.Id); err != nil {
		beego.Error("ConnectAndLogin:", err)
		goto showRegisterErrors
	} else {
		beego.Trace("Let's redirect ", loginRedirect)
		this.Redirect(loginRedirect, 302)
		return
	}
showRegisterErrors:
	this.Data["RegisterForm"] = form
	this.Data["RegisterormErrors"] = valid.Errors
	return
}
func (this *SocialAuthController) processLoginForm(socialType social.SocialType) {
	valid := validation.Validation{}
	form := SocialAuthLoginForm{}
	var user models.User
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	b, err := valid.Valid(form)
	if err != nil {
		beego.Error(err)
	}
	if !b {
		goto showLoginErrors
	}
	//用户不存在？
	user = models.User{Username: form.Username, Email: form.Username}
	if err := user.Read("Username"); err != nil {
		if err2 := user.Read("Email"); err2 != nil {
			errMsg := fmt.Sprintf("用户 %s 不存在!", form.Username)
			valid.SetError("Username", errMsg)
			goto showLoginErrors
		}
	}
	//用户被禁止?
	if user.IsBanned {
		beego.Trace(fmt.Sprintf("用户%s被禁用，不能登录！", user.Username))
		valid.SetError("Username", "抱歉，您被禁止登录！")
		goto showLoginErrors
	}
	//检查密码
	if !user.VerifyPassword(form.Password) {
		beego.Trace(fmt.Sprintf("%s 登录失败！", form.Username))
		valid.SetError("Password", "密码错误")
		goto showLoginErrors
	}
	//验证全部通过
	if loginRedirect, _, err := SocialAuth.ConnectAndLogin(this.Ctx, socialType, user.Id); err != nil {
		beego.Error("ConnectAndLogin:", err)
		goto showLoginErrors
	} else {
		beego.Trace("Let's redirect ", loginRedirect)
		this.Redirect(loginRedirect, 302)
		return
	}
showLoginErrors:
	this.Data["LoginForm"] = form
	this.Data["LoginFormErrors"] = valid.Errors
	return
}

func (this *SocialAuthController) canConnect(socialType *social.SocialType) bool {
	if st, ok := SocialAuth.ReadyConnect(this.Ctx); !ok {
		return false
	} else {
		*socialType = st
	}
	return true
}

func SocialInit() {
	social.DefaultAppUrl = setting.AppUrl + "/"
	SocialAuth = social.NewSocial("/login/", SocialAuther)
	SocialAuth.ConnectSuccessURL = "/"
	SocialAuth.ConnectFailedURL = "/"
	SocialAuth.ConnectRegisterURL = "/register/connect"
	SocialAuth.LoginURL = "/login"
	WeiboAuth = apps.NewWeibo(setting.WeiboClientId, setting.WeiboClientSecret)
	GithubAuth = apps.NewGithub(setting.GithubClientId, setting.GithubClientSecret)
	QQAuth = apps.NewQQ(setting.QQClientId, setting.QQClientSecret)

	if err := social.RegisterProvider(GithubAuth); err != nil {
		beego.Error(err)
	}
	if err := social.RegisterProvider(QQAuth); err != nil {
		beego.Error(err)
	}
	if err := social.RegisterProvider(WeiboAuth); err != nil {
		beego.Error(err)
	}
}

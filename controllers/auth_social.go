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
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
	"github.com/naokij/social-auth"
	"github.com/naokij/social-auth/apps"
)

var (
	GithubAuth *apps.Github
	QQAuth     *apps.QQ
	WeiboAuth  *apps.Weibo
	SocialAuth *social.SocialAuth
)

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
	var socialUserInfo social.UserInfo
	var ok bool
	if socialUserInfo, ok = this.GetSession("social_user_info").(social.UserInfo); !ok {
		beego.Error("error getting social_user_info session")
		this.Abort("500")
	}
	this.Data["SocialType"] = p.GetName()
	this.Data["SocialUser"] = socialUserInfo
}

func (this *SocialAuthController) DoConnect() {

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

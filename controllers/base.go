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
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
	"github.com/naokij/gotalk/utils"
	"time"
)

type BaseController struct {
	beego.Controller
	User    *models.User
	IsLogin bool
}

//通过session获取登录信息，并且登录
func (this *BaseController) loginViaSession() bool {
	if username, ok := this.GetSession("AuthUsername").(string); username != "" && ok {
		//beego.Trace("loginViaSession pass 1 Session[AuthUsername]" + username)
		user := models.User{Username: username}
		if user.Read("Username") == nil {
			this.User = &user
			//beego.Trace("loginViaSession pass 2 ")
			return true
		}
		beego.Trace("loginViaSession pass 2 failed ")
	}
	//beego.Trace("loginViaSession failed ")
	return false
}

//通过remember cookie获取登录信息，并且登录
func (this *BaseController) loginViaRememberCookie() (success bool) {
	username := this.Ctx.GetCookie(setting.CookieUserName)
	if len(username) == 0 {
		return false
	}

	defer func() {
		if !success {
			this.DeleteRememberCookie()
		}
	}()

	user := models.User{Username: username}
	if err := user.Read("Username"); err != nil {
		return false
	}

	secret := utils.EncodeMd5(user.Salt + user.Password)
	value, _ := this.Ctx.GetSecureCookie(secret, setting.CookieRememberName)
	if value != username {
		return false
	}
	this.User = &user
	this.LogUserIn(&user, true)

	return true
}

//删除记忆登录cookie
func (this *BaseController) DeleteRememberCookie() {
	this.Ctx.SetCookie(setting.CookieUserName, "", -1)
	this.Ctx.SetCookie(setting.CookieRememberName, "", -1)
}

//登录用户
func (this *BaseController) LogUserIn(user *models.User, remember bool) {
	this.SessionRegenerateID()
	this.SetSession("AuthUsername", user.Username)
	if remember {
		secret := utils.EncodeMd5(user.Salt + user.Password)
		days := 86400 * 30
		this.Ctx.SetCookie(setting.CookieUserName, user.Username, days)
		this.SetSecureCookie(secret, setting.CookieRememberName, user.Username, days)
	}
}

//登出用户
func (this *BaseController) LogUserOut() {
	this.DeleteRememberCookie()
	this.DelSession("AuthUsername")
	this.DestroySession()
}

func (this *BaseController) Prepare() {
	if setting.ConfigBroken {
		this.Abort("500")
	}
	this.checkDeny()

	// page start time
	this.Data["PageStartTime"] = time.Now()
	this.Data["AppName"] = setting.AppName
	this.Data["AppVer"] = setting.AppVer
	this.Data["PageTitle"] = setting.AppName

	// start session
	this.StartSession()
	//从session中读取登录信息
	switch {
	// save logined user if exist in session
	case this.loginViaSession():
		this.IsLogin = true
	// save logined user if exist in remember cookie
	case this.loginViaRememberCookie():
		this.IsLogin = true
	}

	if this.IsLogin {
		this.Data["User"] = &this.User
		this.Data["IsLogin"] = this.IsLogin

		// if user forbided then do logout
		if this.User.IsBanned {
			this.LogUserOut()
			this.FlashWrite("error", "您的帐号被禁用，无法为您登录！")
			this.Redirect("/login", 302)
			return
		}
	}

	// read flash message
	beego.ReadFromRequest(&this.Controller)

	//检查once token，防止表单多次提交
	if this.Ctx.Request.Method == "POST" || this.Ctx.Request.Method == "DELETE" || this.Ctx.Request.Method == "PUT" ||
		(this.Ctx.Request.Method == "POST" && (this.Ctx.Request.Form.Get("_method") == "delete" || this.Ctx.Request.Form.Get("_method") == "put")) {
		this.CheckOnceToken()
	}

}

// read beego flash message
func (this *BaseController) FlashRead(key string) (string, bool) {
	if data, ok := this.Data["flash"].(map[string]string); ok {
		value, ok := data[key]
		return value, ok
	}
	return "", false
}

// write beego flash message
func (this *BaseController) FlashWrite(key string, value string) {
	flash := beego.NewFlash()
	flash.Data[key] = value
	flash.Store(&this.Controller)
}

//验证防重复提交token
func (this *BaseController) CheckOnceToken() {
	token := this.GetString("_once")
	if token == "" {
		this.Abort("Once")
	}
	if check := cache.GetInt(setting.Cache.Get("Once_" + token)); check == 1 {
		setting.Cache.Delete("Once_" + token)
	} else {
		this.Abort("Once")
	}
}

func (this *BaseController) CheckRequestFrequency(minutes, warningLevel, banLevel int64) {
	var key, denyKey string
	key = "CRF-" + utils.EncodeMd5(this.Ctx.Input.Uri()+this.Ctx.Input.IP())
	denyKey = "deny" + this.Ctx.Input.IP()
	if !setting.Cache.IsExist(key) {
		setting.Cache.Put(key, 1, minutes*60)
	} else {
		setting.Cache.Incr(key)
	}
	counterValue := cache.GetInt64(setting.Cache.Get(key))
	//beego.Trace("CRF ", key, ": ", counterValue)
	if counterValue >= banLevel {
		setting.Cache.Put(denyKey, 1, 86400)
		this.Abort("IPBan")
		return
	}
	if counterValue >= warningLevel {
		this.FlashWrite("warning", "你的提交频率不太正常，持续这样频率的反复提交可能导致你的ip被封锁！")
	}
}

func (this *BaseController) checkDeny() {
	denyKey := "deny" + this.Ctx.Input.IP()
	if setting.Cache.IsExist(denyKey) {
		this.Abort("IPBan")
		return
	}
}

func (this *BaseController) SetPaginator(per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(this.Ctx.Request, per, nums)
	this.Data["paginator"] = p
	return p
}

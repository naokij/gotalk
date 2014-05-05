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
	//"github.com/astaxie/beego"
	"github.com/naokij/gotalk/models"
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
	Wechat      string `form:"Wechat,file"`
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
	user := models.User{Username: this.Ctx.Input.Param(":username")}
	if err := user.Read("Username"); err != nil {
		this.Abort("404")
	}
	if !this.IsLogin {
		this.Abort("403")
	} else if user.Id != this.User.Id && !this.User.IsAdmin {
		this.Abort("403")
	}
	this.Data["UserEditForm"] = &user
}

func (this *UserController) DoEdit() {
	this.Layout = "layout.html"
	this.TplNames = "user_edit.html"
}

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
	"time"
)

type BaseController struct {
	beego.Controller
	User    models.User
	IsLogin bool
}

func (this *BaseController) Prepare() {
	// page start time
	this.Data["PageStartTime"] = time.Now()
	this.Data["AppName"] = setting.AppName
	this.Data["AppVer"] = setting.AppVer
	this.Data["PageTitle"] = setting.AppName
}

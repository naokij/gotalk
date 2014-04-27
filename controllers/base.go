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

	// start session
	this.StartSession()

	// read flash message
	beego.ReadFromRequest(&this.Controller)
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

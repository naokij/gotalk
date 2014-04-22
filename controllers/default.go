package controllers

import ()

type MainController struct {
	BaseController
}

func (this *MainController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.Layout = "layout.html"
	this.TplNames = "index.html"
}

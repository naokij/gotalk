package controllers

import (
//"github.com/astaxie/beego"
)

type AuthController struct {
	BaseController
}

func (this *AuthController) Get() {
	this.Layout = "layout.html"
	this.TplNames = "login.html"
}

func (this *AuthController) Login() {

}

func (this *AuthController) Logout() {

}

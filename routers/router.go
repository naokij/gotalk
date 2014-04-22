package routers

import (
	"github.com/astaxie/beego"
	"github.com/naokij/gotalk/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	authRouter := new(controllers.AuthController)
	beego.Router("/login", authRouter, "get:Get;post:Login")
	beego.Router("/logout", authRouter, "get:Logout")
}

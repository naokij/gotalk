package routers

import (
	"github.com/astaxie/beego"
	"github.com/naokij/gotalk/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	authRouter := new(controllers.AuthController)
	beego.Router("/login", authRouter, "get:Login;post:DoLogin")
	beego.Router("/logout", authRouter, "get:Logout")
	beego.Router("/register", authRouter, "get:Register;post:DoRegister")
	beego.Router("/register/validate-username", authRouter, "get:ValidateUsername")
	beego.Router("/register/validate-email", authRouter, "get:ValidateEmail")
}

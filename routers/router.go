package routers

import (
	"github.com/astaxie/beego"
	"github.com/naokij/gotalk/controllers"
)

func init() {

	beego.Router("/", &controllers.MainController{})
	authController := new(controllers.AuthController)
	beego.Router("/login", authController, "get:Login;post:DoLogin")
	beego.Router("/logout", authController, "get:Logout")
	beego.Router("/register", authController, "get:Register;post:DoRegister")
	beego.Router("/register/validate-username", authController, "get:ValidateUsername")
	beego.Router("/register/validate-email", authController, "get:ValidateEmail")
	beego.Router("/register/validate-captcha", authController, "get:ValidateCaptcha")
	beego.Router("/welcome", authController, "get:Welcome")
}

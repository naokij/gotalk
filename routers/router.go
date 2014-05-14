package routers

import (
	"github.com/astaxie/beego"
	"github.com/naokij/gotalk/controllers"
)

func init() {
	beego.Errorhandler("404", controllers.Error404)
	beego.Errorhandler("403", controllers.Error403)
	beego.Errorhandler("500", controllers.Error500)
	beego.Errorhandler("Once", controllers.ErrorOnce)

	beego.Router("/", &controllers.MainController{})
	authController := new(controllers.AuthController)
	beego.Router("/login", authController, "get:Login;post:DoLogin")
	beego.Router("/login/:returnurl(.+)", authController, "get:Login")
	beego.Router("/logout", authController, "get:Logout")
	beego.Router("/register", authController, "get:Register;post:DoRegister")
	beego.Router("/register/validate-username", authController, "get:ValidateUsername")
	beego.Router("/register/validate-email", authController, "get:ValidateEmail")
	beego.Router("/register/validate-captcha", authController, "get:ValidateCaptcha")
	userController := new(controllers.UserController)
	beego.Router("/user/:username(.+)/edit", userController, "get:Edit;post:Edit")
	beego.Router("/user/:username(.+)", userController, "get:Profile")
	beego.Router("/user/validate-password", userController, "post:ValidatePassword")
}

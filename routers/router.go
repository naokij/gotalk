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

	//登录
	authController := new(controllers.AuthController)
	beego.Router("/login", authController, "get:Login;post:DoLogin")
	beego.Router("/login/:returnurl(.+)", authController, "get:Login")
	beego.Router("/forget-password", authController, "get:ForgetPassword;post:ForgetPassword")
	beego.Router("/reset-password/:code([0-9a-zA-Z]+)", authController, "get:ResetPassword;post:ResetPassword")
	beego.Router("/logout", authController, "get:Logout")
	beego.Router("/register", authController, "get:Register;post:DoRegister")
	beego.Router("/register/validate-username", authController, "get:ValidateUsername")
	beego.Router("/register/validate-email", authController, "get:ValidateEmail")
	beego.Router("/register/validate-captcha", authController, "get:ValidateCaptcha")
	beego.Router("/activate/:code([0-9a-zA-Z]+)", authController, "get:Activate")
	//社交帐号登录
	beego.AddFilter("/login/:/access", "BeforeRouter", controllers.OAuthAccess)
	beego.AddFilter("/login/:", "BeforeRouter", controllers.OAuthRedirect)
	socialAuthController := new(controllers.SocialAuthController)
	beego.Router("/register/connect", socialAuthController, "get:Connect;post:DoConnect")

	userController := new(controllers.UserController)
	beego.Router("/user/:username(.+)/edit", userController, "get:Edit;post:Edit")
	beego.Router("/user/:username(.+)/resend-validation", userController, "get:ResendValidation")
	beego.Router("/user/:username(.+)/change-username", userController, "get:ChangeUsername;post:ChangeUsername")
	beego.Router("/user/:username(.+)", userController, "get:Profile")
}

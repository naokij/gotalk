package routers

import (
	"github.com/naokij/gotalk/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}

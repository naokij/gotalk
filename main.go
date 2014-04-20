package main

import (
	//"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/naokij/gotalk/setting"
	//"github.com/naokij/gotalk/models"
	_ "github.com/naokij/gotalk/routers"
)

func init() {

}

func main() {

	beego.Run()
}

package main

import (
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/naokij/gotalk/models"
	_ "github.com/naokij/gotalk/routers"
)

func init() {
	orm.RegisterDriver("mysql", orm.DR_MySQL)

	orm.RegisterDataBase("default", "mysql", "root:@/gotalk?charset=utf8")
}

func main() {

	beego.Run()
}

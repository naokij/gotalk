package main

import (
	"fmt"
	_ "github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/naokij/gotalk/models"
	_ "github.com/naokij/gotalk/routers"
)

func init() {
	orm.RegisterDriver("mysql", orm.DR_MySQL)

	orm.RegisterDataBase("default", "mysql", "root:@/gotalk?charset=utf8")
}

func main() {
	orm.Debug = true
	o := orm.NewOrm()
	fmt.Println(o)

	o.Using("default") // 默认使用 default，你可以指定为其他数据库
	user := new(models.User)
	user.Username = "naoki"
	user.Nickname = "江乐"
	user.SetPassword("jiangle")
	fmt.Println(o.Insert(user))
	fmt.Println("Pasword verify:", user.VerifyPassword("jiangle"))

	//beego.Run()
}

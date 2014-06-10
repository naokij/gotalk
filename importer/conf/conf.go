package conf

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	// "github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
	"runtime"
	// "time"
)

var DiscuzDb string
var Orm orm.Ormer
var OrmGotalk orm.Ormer
var Workers = runtime.NumCPU()
var WorkerLoad = 2500
var AvatarPath string

func init() {
	Workers = 8
	runtime.GOMAXPROCS(Workers)
	beego.AppConfigPath = "../conf/app.conf"
	beego.ParseConfig()
	setting.ReadConfig()
	DiscuzDb = beego.AppConfig.String("importer::discuzdb")
	AvatarPath = beego.AppConfig.String("importer::avatarpath")
	if err := orm.RegisterDataBase("discuz", "mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", setting.MySQLUser, setting.MySQLPassword, setting.MySQLHost, DiscuzDb)+"&loc=Asia%2FShanghai", 30); err != nil {
		fmt.Println("MySQL error:", err.Error())
		return
	}
	OrmGotalk = orm.NewOrm()
	if db, err := orm.GetDB("discuz"); err == nil {
		Orm, err = orm.NewOrmWithDB("mysql", "discuz", db)
		if err != nil {
			fmt.Println("MySQL error:", err.Error())
			return
		}
	} else {
		fmt.Println("MySQL error:", err.Error())
		return
	}
	orm.RunSyncdb("default", true, false)
}

package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/naokij/gotalk/models"
	"github.com/naokij/gotalk/setting"
	"time"
)

var DiscuzDb string
var Orm orm.Ormer
var OrmGotalk orm.Ormer

func Users() {
	type discuzType struct {
		Uid      int
		Username string
		Password string
		Email    string
		Regdate  time.Time
		Salt     string
		Qq       string
		Site     string
		Bio      string
	}
	var pos int64
	var discuzRes []discuzType
	var end bool
	fmt.Println("Converting Users")
	for !end {
		sql := "SELECT u.uid, u.username, u.password, u.email, u.regdate, u.salt,p.qq, p.site,p.bio from uc_members as u left join pre_common_member_profile as p on u.uid=p.uid limit ?,10000"
		num, err := Orm.Raw(sql, pos).QueryRows(&discuzRes)
		if err != nil {
			fmt.Println("MySQL error:", err.Error())
			return
		}
		if num == 0 {
			end = true
		}
		for _, discuzRow := range discuzRes {
			//fmt.Printf("%#v\n", discuzRow)
			gotalkRes := new(models.User)
			gotalkRes.Id = discuzRow.Uid
			gotalkRes.Username = discuzRow.Username
			gotalkRes.Password = discuzRow.Password
			gotalkRes.Email = discuzRow.Email
			gotalkRes.Salt = discuzRow.Salt
			gotalkRes.Created = discuzRow.Regdate
			gotalkRes.Url = discuzRow.Site
			gotalkRes.Info = discuzRow.Bio
			gotalkRes.Qq = discuzRow.Qq
			gotalkRes.Insert()
			//fmt.Println(discuzRow.Uid, gotalkRes.Id)
			if _, err := OrmGotalk.Raw("Update user set id = ? where id = ?", discuzRow.Uid, gotalkRes.Id).Exec(); err != nil {
				fmt.Println(err)
			}
		}
		fmt.Print(".")
		pos += 10000
	}
	fmt.Println("")
}

func main() {
	beego.AppConfigPath = "../conf/app.conf"
	beego.ParseConfig()
	setting.ReadConfig()
	DiscuzDb = beego.AppConfig.String("mysql::discuzdb")
	if err := orm.RegisterDataBase("discuz", "mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", setting.MySQLUser, setting.MySQLPassword, setting.MySQLHost, DiscuzDb)+"&loc=Asia%2FShanghai"); err != nil {
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
	Users()
}

/*
Copyright 2014 Jiang Le

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package setting

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/captcha"
	_ "github.com/go-sql-driver/mysql"
	"labix.org/v2/mgo"
)

var (
	AppName            string
	AppHost            string
	AppUrl             string
	AppLogo            string
	SecretKey          string
	CookieUserName     string
	CookieRememberName string
	MongodbSession     *mgo.Session
	MongodbHost        string
	MongodbName        string
	MySQLHost          string
	MySQLUser          string
	MySQLPassword      string
	MySQLDB            string
	XSRFKey            string
	Cache              cache.Cache
	Captcha            *captcha.Captcha
)

const (
	AppVer = "VERSION 0.0.1"
)

func ReadConfig() {
	var err error

	AppName = beego.AppConfig.String("appname")
	AppHost = beego.AppConfig.String("apphost")
	AppUrl = beego.AppConfig.String("appurl")
	AppLogo = beego.AppConfig.String("applogo")
	CookieUserName = beego.AppConfig.String("cookieusername")
	CookieRememberName = beego.AppConfig.String("CookieRememberName")
	MySQLHost = beego.AppConfig.String("mysql::host")
	MySQLUser = beego.AppConfig.String("mysql::user")
	MySQLPassword = beego.AppConfig.String("mysql::password")
	MySQLDB = beego.AppConfig.String("mysql::db")
	MongodbHost = beego.AppConfig.String("mongodb::host")
	MongodbName = beego.AppConfig.String("mongodb::name")
	SecretKey = beego.AppConfig.String("security::secret_key")
	XSRFKey = beego.AppConfig.String("security::xsrfkey")

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", MySQLUser, MySQLPassword, MySQLHost, MySQLDB)+"&loc=Asia%2FShanghai")

	MongodbSession, err = mgo.Dial(MongodbHost)
	if err != nil {
		beego.Error(err)
	}

	beego.EnableXSRF = true
	beego.XSRFKEY = XSRFKey
	beego.XSRFExpire = 360

	// cache system
	Cache, err = cache.NewCache("memory", `{"interval":360}`)

	Captcha = captcha.NewWithFilter("/captcha/", Cache)
	Captcha.FieldIdName = "captcha-id"
	Captcha.FieldCaptchaName = "captcha"

	beego.SessionOn = true
	beego.SessionProvider = "memory"
	beego.SessionCookieLifeTime = 0
	beego.SessionGCMaxLifetime = 86400
	//todo 更好的利用mongodb session
}

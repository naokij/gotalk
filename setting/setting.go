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
	_ "github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego/utils/captcha"
	_ "github.com/go-sql-driver/mysql"
	"github.com/naokij/GoStopForumSpam/stopforumspam"
	"github.com/naokij/go-sendcloud"
	_ "github.com/naokij/gotalk/cache/redis"
	"github.com/naokij/gotalk/filestore"
	"labix.org/v2/mgo"
)

var (
	AppName string
	AppHost string
	AppUrl  string
	AppLogo string
	TmpPath string
)

var (
	CookieUserName     string
	CookieRememberName string
)

var (
	MongodbHost   string
	MongodbName   string
	MySQLHost     string
	MySQLUser     string
	MySQLPassword string
	MySQLDB       string
	RedisHost     string
	RedisPort     string
)

var (
	SendcloudDomain string
	SendcloudUser   string
	SendcloudKey    string
	From            string
	FromName        string
)

var (
	// OAuth
	GithubClientId     string
	GithubClientSecret string
	WeiboClientId      string
	WeiboClientSecret  string
	QQClientId         string
	QQClientSecret     string
)

var (
	ConfigBroken bool
)

var (
	Cache          cache.Cache
	MongodbSession *mgo.Session
	Captcha        *captcha.Captcha
	AvatarFSM      *filestore.Manager
	AttachmentFSM  *filestore.Manager
	Sendcloud      *sendcloud.Client
	StopForumSpam  *stopforumspam.Client
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
	TmpPath = beego.AppConfig.String("tmppath")
	CookieUserName = beego.AppConfig.String("cookieusername")
	CookieRememberName = beego.AppConfig.String("CookieRememberName")
	MySQLHost = beego.AppConfig.String("mysql::host")
	MySQLUser = beego.AppConfig.String("mysql::user")
	MySQLPassword = beego.AppConfig.String("mysql::password")
	MySQLDB = beego.AppConfig.String("mysql::db")
	MongodbHost = beego.AppConfig.String("mongodb::host")
	MongodbName = beego.AppConfig.String("mongodb::name")
	RedisHost = beego.AppConfig.String("redis::host")
	RedisPort = beego.AppConfig.String("redis::port")
	SendcloudDomain = beego.AppConfig.String("sendcloud::domain")
	SendcloudUser = beego.AppConfig.String("sendcloud::user")
	SendcloudKey = beego.AppConfig.String("sendcloud::key")
	From = beego.AppConfig.String("sendcloud::from")
	FromName = beego.AppConfig.String("sendcloud::fromname")

	if err := orm.RegisterDriver("mysql", orm.DR_MySQL); err != nil {
		beego.Error("mysql", err)
		ConfigBroken = true
	}
	if err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", MySQLUser, MySQLPassword, MySQLHost, MySQLDB)+"&loc=Asia%2FShanghai", 30); err != nil {
		beego.Error("mysql", err)
		ConfigBroken = true
	}

	MongodbSession, err = mgo.Dial(MongodbHost)
	if err != nil {
		beego.Error("mongodb", err)
		ConfigBroken = true
	}

	// cache system
	Cache, err = cache.NewCache("redis", fmt.Sprintf(`{"conn":"%s:%s"}`, RedisHost, RedisPort))
	if err != nil {
		beego.Error("cache", err)
		ConfigBroken = true
	}

	Captcha = captcha.NewWithFilter("/captcha/", Cache)
	Captcha.FieldIdName = "captcha-id"
	Captcha.FieldCaptchaName = "captcha"

	beego.SessionOn = true
	beego.SessionProvider = "redis"
	beego.SessionSavePath = RedisHost + ":" + RedisPort
	beego.SessionCookieLifeTime = 0
	beego.SessionGCMaxLifetime = 86400
	//todo 更好的利用mongodb session

	avatarstore := beego.AppConfig.String("avatar::store")
	switch avatarstore {
	case "local":
		avatarConfig := filestore.Config{FSPath: beego.AppConfig.String("avatar::local_path"), UrlPrefix: beego.AppConfig.String("avatar::url")}
		AvatarFSM, err = filestore.NewManager(avatarstore, avatarConfig)
	}
	if err != nil {
		beego.Error(err)
	}

	//初始化sendcloud，用来发送邮件
	Sendcloud = sendcloud.New()
	Sendcloud.AddDomain(SendcloudDomain, SendcloudUser, SendcloudKey)
	//StopForumSpam
	StopForumSpam = stopforumspam.New(beego.AppConfig.String("stopforumspam::key"))

	//社交帐号登录
	GithubClientId = beego.AppConfig.String("oauth::github_client_id")
	GithubClientSecret = beego.AppConfig.String("oauth::github_client_secret")
	QQClientId = beego.AppConfig.String("oauth::qq_client_id")
	QQClientSecret = beego.AppConfig.String("oauth::qq_client_secret")
	WeiboClientId = beego.AppConfig.String("oauth::weibo_client_id")
	WeiboClientSecret = beego.AppConfig.String("oauth::weibo_client_secret")

}

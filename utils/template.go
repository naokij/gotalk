package utils

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/astaxie/beego"
	"github.com/naokij/gotalk/setting"
	"html/template"
	"strings"
	"time"
)

func loadtimes(t time.Time) int {
	return int(time.Since(t).Nanoseconds() / 1e6)
}

func nl2br(text string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
}

//生成防止重复提交token
func OnceToken() (token string) {
	token = strings.Replace(uuid.New(), "-", "", -1)
	setting.Cache.Put("Once_"+token, 1, 86400)
	return token
}

func OnceFormHtml() template.HTML {
	return template.HTML("<input type=\"hidden\" name=\"_once\" value=\"" +
		OnceToken() + "\"/>")
}

func LoginUrlFor(endpoint string, values ...string) string {
	return beego.UrlFor("AuthController.Login", ":returnurl", template.URLQueryEscaper(beego.UrlFor(endpoint, values...)))
}

func init() {
	// Register template functions.
	beego.AddFuncMap("loadtimes", loadtimes)
	beego.AddFuncMap("jsescape", template.JSEscapeString)
	beego.AddFuncMap("nl2br", nl2br)
	beego.AddFuncMap("oncetoken", OnceToken)
	beego.AddFuncMap("onceformhtml", OnceFormHtml)
	beego.AddFuncMap("loginurl", LoginUrlFor)
}

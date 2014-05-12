package utils

import (
	"github.com/astaxie/beego"
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

func init() {
	// Register template functions.
	beego.AddFuncMap("loadtimes", loadtimes)
	beego.AddFuncMap("jsescape", template.JSEscapeString)
	beego.AddFuncMap("nl2br", nl2br)
}

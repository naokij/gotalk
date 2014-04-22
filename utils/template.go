package utils

import (
	"github.com/astaxie/beego"
	"time"
)

func loadtimes(t time.Time) int {
	return int(time.Since(t).Nanoseconds() / 1e6)
}

func init() {
	// Register template functions.
	beego.AddFuncMap("loadtimes", loadtimes)
}

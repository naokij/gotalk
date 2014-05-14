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

package controllers

import (
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/naokij/gotalk/setting"
	"html/template"
	"net/http"
)

func Error404(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("404.html").ParseFiles(beego.ViewsPath + "/errors/404.html")
	data := make(map[string]interface{})
	data["Title"] = "页面未找到 | " + setting.AppName
	t.Execute(rw, data)
}

func Error403(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("403.html").ParseFiles(beego.ViewsPath + "/errors/403.html")
	data := make(map[string]interface{})
	data["Title"] = "禁止访问 | " + setting.AppName
	t.Execute(rw, data)
}

func Error500(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("50x.html").ParseFiles(beego.ViewsPath + "/errors/50x.html")
	data := make(map[string]interface{})
	data["Title"] = "服务器内部错误 | " + setting.AppName
	t.Execute(rw, data)
}

func ErrorOnce(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("once.html").ParseFiles(beego.ViewsPath + "/errors/once.html")
	data := make(map[string]interface{})
	data["Title"] = "重复提交 | " + setting.AppName
	t.Execute(rw, data)
}

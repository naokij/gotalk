package models

import (
	//"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
)

func init() {
	orm.RegisterDataBase("default", "mysql", "root@/gotalk?charset=utf8&loc=Asia%2FShanghai", 30)
	orm.RunSyncdb("default", true, false)
}

func TestPasswordVerify(t *testing.T) {

	var user User
	user.Username = "username"
	user.Email = "username@gotalk.local"
	user.SetPassword("password")
	if user.VerifyPassword("password") {
		t.Log("密码验证正确！")
	} else {
		t.Error("密码不对!")
	}
	var bannedUser User
	bannedUser.Username = "baduser"
	bannedUser.Email = "baduser@gotalk.local"
	bannedUser.SetPassword("passwordforbaduser")
	bannedUser.IsBanned = true
	bannedUser.Insert()

	if user.VerifyPassword("password1") {
		t.Error("伪造的密码居然能通过测试？")
	} else {
		t.Log("很好，假密码没有通过！")
	}
	user.Insert()
}

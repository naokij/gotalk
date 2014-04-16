package models

import (
	//"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
)

func init() {
	orm.RegisterDataBase("default", "mysql", "root@/gotalk?charset=utf8", 30)
	orm.RunSyncdb("default", true, false)
}

func TestTopicCreate(t *testing.T) {

}

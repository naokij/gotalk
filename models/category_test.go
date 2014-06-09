package models

import (
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/naokij/gotalk/setting"
	"testing"
)

func init() {
	//orm.RegisterDataBase("default", "mysql", "root@/gotalk?charset=utf8", 30)
	beego.AppConfigPath = "../conf/app.conf"
	beego.ParseConfig()
	setting.ReadConfig()
	orm.RunSyncdb("default", true, false)
}

func TestTreeCreate(t *testing.T) {

	node1 := Category{Id: 1, ParentCategory: nil, Name: "ELECTRONICS", UrlCode: "electronics"}
	node2 := Category{Id: 2, ParentCategory: &node1, Name: "TELEVISIONS", UrlCode: "televisions"}
	node3 := Category{Id: 3, ParentCategory: &node2, Name: "TUBE", UrlCode: "tube"}
	node4 := Category{Id: 4, ParentCategory: &node2, Name: "LCD", UrlCode: "lcd"}
	node5 := Category{Id: 5, ParentCategory: &node2, Name: "PLASMA", UrlCode: "plasma"}
	node6 := Category{Id: 6, ParentCategory: &node1, Name: "PORTABLE ELECTRONICS", UrlCode: "portable_electronics"}
	node7 := Category{Id: 7, ParentCategory: &node6, Name: "MP3 PLAYERS", UrlCode: "mp3_players"}
	node8 := Category{Id: 8, ParentCategory: &node7, Name: "FLASH", UrlCode: "flash"}
	node9 := Category{Id: 9, ParentCategory: &node6, Name: "CD PLAYERS", UrlCode: "cd_players"}
	node10 := Category{Id: 10, ParentCategory: &node6, Name: "2 WAY RADIOS", UrlCode: "2_way_radios"}
	err := node1.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node1.Name)
	}
	err = node2.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node2.Name)
	}
	err = node3.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node3.Name)
	}
	err = node4.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node4.Name)
	}
	err = node5.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node5.Name)
	}
	err = node6.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node6.Name)
	}
	err = node7.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node7.Name)
	}
	err = node8.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node8.Name)
	}
	err = node9.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node9.Name)
	}
	err = node10.Insert()
	if err != nil {
		t.Error("Unable to insert category node " + node10.Name)
	}
}

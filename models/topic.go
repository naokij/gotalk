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

package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/naokij/gotalk/setting"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

//话题/帖子
type Topic struct {
	Id                int
	Title             string    `orm:"size(255)"`
	ContentHex        string    `orm:"size(24)"` //mongodb对象编号
	Content           *Content  `orm:"-"`
	User              *User     `orm:"rel(fk);index"`
	Username          string    `orm:"size(30)`
	Category          *Category `orm:"rel(fk);index"`
	PvCount           int       `orm:"index"`
	CommentCount      int       `orm:"index"`
	BookmarkCount     int       `orm:"index"`
	IsExcellent       bool      `orm:"index"`
	IsClosed          bool      `orm:""`
	LastReplyUsername string    `orm:"size(30)`
	LastReplyAt       time.Time `orm:""`
	Created           time.Time `orm:"auto_now_add"`
	Updated           time.Time `orm:"auto_now"`
	Ip                string    `orm:"size(39)"`
}

func (m *Topic) Insert() error {
	var err error
	if m.Content != nil {
		m.ContentHex, err = m.Content.Insert()
		if err != nil {
			return err
		}
	}
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Topic) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	if m.ContentHex != "" && m.ContentHex != "0" {
		fmt.Println("ContentHex", m.ContentHex)
		content := Content{}
		m.Content = &content
		m.Content.Id = bson.ObjectIdHex(m.ContentHex)
		err := m.Content.Read()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Topic) Update(fields ...string) error {
	if m.Content != nil {
		err := m.Content.Update()
		if err != nil {
			return err
		}
		m.ContentHex = m.Content.Id.Hex()
	}
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Topic) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	if m.Content != nil {
		m.Content.Delete()
	}
	return nil
}

func (m *Topic) TableEngine() string {
	return "INNODB"
}

//留言/回帖
type Comment struct {
	Id         int
	Topic      *Topic    `orm:"rel(fk);index"`
	ContentHex string    `orm:"size(24)"`
	Content    *Content  `orm:"-"`
	User       *User     `orm:"rel(fk)"`
	Created    time.Time `orm:"auto_now_add"`
	Updated    time.Time `orm:"auto_now"`
}

func (m *Comment) Insert() error {
	var err error
	if m.Content != nil {
		m.ContentHex, err = m.Content.Insert()
		if err != nil {
			return err
		}
	}
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Comment) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	if m.ContentHex != "" && m.ContentHex != "0" {
		content := Content{}
		m.Content = &content
		m.Content.Id = bson.ObjectIdHex(m.ContentHex)
		err := m.Content.Read()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Comment) Update(fields ...string) error {
	if m.Content != nil {
		err := m.Content.Update()
		if err != nil {
			return err
		}
		m.ContentHex = m.Content.Id.Hex()
	}
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Comment) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	if m.Content != nil {
		m.Content.Delete()
	}
	return nil
}

func (m *Comment) TableEngine() string {
	return "INNODB"
}

type Content struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	TopicId   int
	CommentId int
	Message   string
}

func (m *Content) Session() *mgo.Session {
	return setting.MongodbSession.Clone()
}

func (m *Content) Insert() (string, error) {
	objectId := bson.NewObjectId()
	m.Id = objectId
	session := m.Session()
	defer func() {
		session.Close()
	}()
	c := session.DB(setting.MongodbName).C("Content")
	err := c.Insert(m)
	return m.Id.Hex(), err
}

func (m *Content) Read() error {
	session := m.Session()
	defer func() {
		session.Close()
	}()
	c := session.DB(setting.MongodbName).C("Content")
	err := c.FindId(m.Id).One(&m)
	return err
}

func (m *Content) Update() error {
	session := m.Session()
	defer func() {
		session.Close()
	}()
	c := session.DB(setting.MongodbName).C("Content")
	err := c.UpdateId(m.Id, m)
	return err
}

func (m *Content) Delete() error {
	session := m.Session()
	defer func() {
		session.Close()
	}()
	c := session.DB(setting.MongodbName).C("Content")
	err := c.RemoveId(m.Id)
	return err
}
func init() {
	orm.RegisterModel(new(Topic), new(Comment))
}

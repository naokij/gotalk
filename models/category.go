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
	"github.com/astaxie/beego/orm"
)

//话题分类
type Category struct {
	Id             int
	ParentCategory *Category `orm:"null;rel(fk);on_delete(set_null);index"`
	Depth          int       `orm:"index"`
	Sort           int       `orm:"index"`
	UrlCode        string    `orm:"size(30)";unique`
	Name           string    `orm:"size(50)"`
	Description    string    `orm:"size(255)"`
	TopicCount     int       `orm:"index"`
	CommentCount   int       `orm:"index"`
}

func (m *Category) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Category) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Category) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Category) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Category) TableEngine() string {
	return "INNODB"
}

func init() {
	orm.RegisterModel(new(Category))
}

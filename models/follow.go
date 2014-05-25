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
	"time"
)

type Follow struct {
	Id         int
	User       *User `orm:"rel(fk)"`
	FollowUser *User `orm:"rel(fk)"`
	Mutual     bool
	Created    time.Time `orm:"auto_now_add"`
}

func (*Follow) TableUnique() [][]string {
	return [][]string{
		[]string{"User", "FollowUser"},
	}
}

func (m *Follow) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Follow) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Follow) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Follow) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Follow) TableEngine() string {
	return "INNODB"
}

func Follows() orm.QuerySeter {
	return orm.NewOrm().QueryTable("follow")
}

func init() {
	orm.RegisterModel(new(Follow))
}

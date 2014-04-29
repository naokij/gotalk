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
	//"fmt"
	"encoding/hex"
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/naokij/gotalk/utils"
	"regexp"
	"time"
)

type User struct {
	Id          int
	Username    string    `orm:"size(30);unique"`
	Nickname    string    `orm:"size(30)"`
	Password    string    `orm:"size(128)"`
	Url         string    `orm:"size(100)"`
	Company     string    `orm:"size(30)"`
	Location    string    `orm:"size(30)"`
	Email       string    `orm:"size(80);unique"`
	Avatar      string    `orm:"size(32)"`
	Info        string    ``
	Weibo       string    `orm:"size(30)"`
	WeChat      string    `orm:"size(20)"`
	Qq          int       ``
	PublicEmail bool      ``
	Followers   int       ``
	Following   int       ``
	FavTopics   int       ``
	IsAdmin     bool      `orm:"index"`
	IsActive    bool      `orm:"index"`
	IsBanned    bool      `orm:"index"`
	Salt        string    `orm:"size(6)"`
	Created     time.Time `orm:"auto_now_add"`
	Updated     time.Time `orm:"auto_now"`
}

const (
	activeCodeLife = 180
	resetPasswordCodeLife
	UsernameRegex = `^[\p{Han}a-zA-Z0-9]+$`
)

func (m *User) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *User) ValidUsername() (err error) {
	reg := regexp.MustCompile(UsernameRegex)
	if !reg.MatchString(m.Username) {
		err = errors.New("只能包含英文、数字和汉字")
	} else {
		if !(utils.HZStringLength(m.Username) >= 3 && utils.HZStringLength(m.Username) <= 16) {
			err = errors.New("长度3-16（汉字长度按2计算）")
		}
	}
	return err
}

func (m *User) SetPassword(password string) error {
	m.Salt = utils.GetRandomString(6)
	m.Password = utils.EncodeMd5(utils.EncodeMd5(password) + m.Salt)
	return nil
}

func (m *User) VerifyPassword(password string) bool {
	if m.Password == utils.EncodeMd5(utils.EncodeMd5(password)+m.Salt) {
		return true
	}
	return false
}

func (m *User) GenerateActivateCode() string {
	data := utils.ToStr(m.Id) + m.Email + m.Username + m.Password + m.Salt
	code := utils.CreateTimeLimitCode(data, activeCodeLife, nil)

	// add tail hex username
	code += hex.EncodeToString([]byte(m.Username))
	return code
}

func (m *User) verifyCodePass1(code string) bool {
	if len(code) <= utils.TimeLimitCodeLength {
		return false
	}

	// use tail hex username query user
	hexStr := code[utils.TimeLimitCodeLength:]
	if b, err := hex.DecodeString(hexStr); err == nil {
		if m.Username == string(b) {
			return true
		}
	}

	return false
}

func (m *User) VerifyActivateCode(code string) bool {

	if m.verifyCodePass1(code) {
		// time limit code
		prefix := code[:utils.TimeLimitCodeLength]
		data := utils.ToStr(m.Id) + m.Email + m.Username + m.Password + m.Salt

		return utils.VerifyTimeLimitCode(data, activeCodeLife, prefix)
	}

	return false
}

func (u *User) TableEngine() string {
	return "INNODB"
}

func init() {
	orm.RegisterModel(new(User))
}

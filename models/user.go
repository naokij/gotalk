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
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/disintegration/imaging"
	"github.com/naokij/gotalk/setting"
	"github.com/naokij/gotalk/utils"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	Qq          string    `orm:"size(20)"`
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

func (m *User) AvatarUrl() (url string) {
	if m.Avatar == "" {
		url = m.gravatarUrl(48)
	} else {
		url = setting.AvatarFSM.GetConfig().UrlPrefix + string(m.Avatar[0]) + "/" + string(m.Avatar[1]) + "/" + m.Avatar + "-m.png"
	}
	return url
}

func (m *User) LargeAvatarUrl() (url string) {
	if m.Avatar == "" {
		url = m.gravatarUrl(220)
	} else {
		url = setting.AvatarFSM.GetConfig().UrlPrefix + string(m.Avatar[0]) + "/" + string(m.Avatar[1]) + "/" + m.Avatar + "-l.png"
	}
	return
}

func (m *User) gravatarUrl(size int) (url string) {
	hash := utils.EncodeMd5(strings.ToLower(m.Email))
	url = fmt.Sprintf("http://gravatar.duoshuo.com/avatar/%s?d=identicon&size=%d", hash, size)
	return url
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

func (m *User) ValidateUrl() (err error) {
	u, err := url.Parse(m.Url)
	if err != nil {
		err = errors.New("网址无效")
		return err
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		err = errors.New("只接受http和https协议的网址")
		return err
	}
	return nil
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

func (m *User) GenerateActivateCode() (code string, err error) {
	code = strings.Replace(uuid.New(), "-", "", -1)
	if err := setting.Cache.Put("activation:"+code, m.Username, 3600); err != nil {
		beego.Error("cache", err)
		return "", err
	}
	return code, nil
}

//验证激活码
//如果验证通过User对象会变成激活码对应的User
func (m *User) VerifyActivateCode(code string) bool {
	b := m.TestActivateCode(code)
	if b {
		if err := m.ConsumeActivateCode(code); err != nil {
			beego.Error(err)
		}
	}
	return b
}

//测试激活码
//测试完后不会删除
func (m *User) TestActivateCode(code string) bool {
	usernameFromCache := cache.GetString(setting.Cache.Get("activation:" + code))
	if usernameFromCache == "" {
		return false
	}
	m.Username = usernameFromCache
	if err := m.Read("Username"); err != nil {
		return false
	}
	return true
}

//使用激活码
func (m *User) ConsumeActivateCode(code string) error {
	return setting.Cache.Delete("activation:" + code)
}

func (m *User) ValidateAndSetAvatar(avatarFile io.Reader, filename string) error {
	var img image.Image
	var err error
	var ext string
	var tmpFile *os.File
	if tmpFile, err = ioutil.TempFile(setting.TmpPath, "uploaded-avatar-"); err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	if _, err = io.Copy(tmpFile, avatarFile); err != nil {
		return err
	}
	tmpFile.Seek(0, 0)
	if filename != "" {
		ext = strings.ToLower(filepath.Ext(filename))
	} else {
		ext = utils.GetImageFormat(tmpFile)
		tmpFile.Seek(0, 0)
	}
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return errors.New("只允许jpg, png, gif类型的图片")
	}
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(tmpFile)
		if err != nil {
			return errors.New("无法识别此jpg文件")
		}
	case ".png":
		img, err = png.Decode(tmpFile)
		if err != nil {
			return errors.New("无法识别此png文件")
		}
	case ".gif":
		img, err = gif.Decode(tmpFile)
		if err != nil {
			return errors.New("无法识别此gif文件")
		}
	}

	//crop正方形
	bound := img.Bounds()
	if w, h := bound.Dx(), bound.Dy(); w > h {
		img = imaging.CropCenter(img, h, h)
	} else if w < h {
		img = imaging.CropCenter(img, w, w)
	}
	//制作缩略图
	imgL := imaging.Resize(img, 220, 220, imaging.Lanczos)
	imgM := imaging.Resize(img, 48, 48, imaging.Lanczos)
	imgS := imaging.Resize(img, 24, 24, imaging.Lanczos)
	uuid := strings.Replace(uuid.New(), "-", "", -1)
	imgLName, imgMName, imgSName := setting.TmpPath+uuid+"-l.png", setting.TmpPath+uuid+"-m.png", setting.TmpPath+uuid+"-s.png"
	errL, errM, errS := imaging.Save(imgL, imgLName), imaging.Save(imgM, imgMName), imaging.Save(imgS, imgSName)
	if errL != nil || errM != nil || errS != nil {
		return errors.New("无法保存头像临时文件")
	}
	defer os.Remove(imgLName)
	defer os.Remove(imgMName)
	defer os.Remove(imgSName)

	_, errL = setting.AvatarFSM.PutFile(imgLName, string(uuid[0])+"/"+string(uuid[1])+"/"+uuid+"-l.png")
	_, errM = setting.AvatarFSM.PutFile(imgMName, string(uuid[0])+"/"+string(uuid[1])+"/"+uuid+"-m.png")
	_, errS = setting.AvatarFSM.PutFile(imgSName, string(uuid[0])+"/"+string(uuid[1])+"/"+uuid+"-s.png")

	if errL != nil || errM != nil || errS != nil {
		return errors.New("无法保存头像")
	}
	if m.Avatar != "" {
		errL = setting.AvatarFSM.Delete(m.Avatar + "-l.png")
		errM = setting.AvatarFSM.Delete(m.Avatar + "-m.png")
		errS = setting.AvatarFSM.Delete(m.Avatar + "-s.png")
	}
	m.Avatar = uuid
	//errL, errM, errS = setting.AvatarFSM.Delete
	return nil
}

func (m *User) FollowingUsers() orm.QuerySeter {
	return Follows().Filter("User", m.Id)
}

func (m *User) FollowerUsers() orm.QuerySeter {
	return Follows().Filter("FollowUser", m.Id)
}

func (u *User) Follow(who *User) (err error) {
	if err = who.Read(); err == nil {
		var mutual bool

		reverseFollow := Follow{User: who, FollowUser: u}
		if err := reverseFollow.Read("User", "FollowUser"); err == nil {
			mutual = true
		}

		follow := Follow{User: u, FollowUser: who, Mutual: mutual}
		if err := follow.Insert(); err != nil && mutual {
			reverseFollow.Mutual = mutual
			reverseFollow.Update("Mutual")
		}

		if nums, err := u.FollowingUsers().Count(); err == nil {
			u.Following = int(nums)
			u.Update("Following")
		}

		if nums, err := who.FollowerUsers().Count(); err == nil {
			who.Followers = int(nums)
			who.Update("Followers")
		}

	} else {
		err = fmt.Errorf("%s must be saved before he/she can be followed!")
	}
	return
}

func (u *User) UnFollow(who *User) (err error) {
	num, _ := u.FollowingUsers().Filter("FollowUser", who.Id).Delete()
	if num > 0 {
		who.FollowingUsers().Filter("FollowUser", u.Id).Update(orm.Params{
			"Mutual": false,
		})

		if nums, err := u.FollowingUsers().Count(); err == nil {
			u.Following = int(nums)
			u.Update("Following")
		}

		if nums, err := who.FollowerUsers().Count(); err == nil {
			who.Followers = int(nums)
			who.Update("Followers")
		}
	} else {
		err = fmt.Errorf("%s not following %s", u.Username, who.Username)
	}
	return
}

func (u *User) TableEngine() string {
	return "INNODB"
}

func init() {
	orm.RegisterModel(new(User))
}

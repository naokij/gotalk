package converters

import (
	"fmt"
	"github.com/naokij/gotalk/importer/conf"
	"github.com/naokij/gotalk/models"
	"os"
	"time"
)

type discuzUser struct {
	Uid      int
	Username string
	Password string
	Email    string
	Regdate  int64
	Salt     string
	Qq       string
	Site     string
	Bio      string
	Authstr  string
}

func cutString(str string, length int) string {
	chars := []rune(str)
	if len(chars) <= length {
		return str
	}
	return string(chars[0:length])
}

func Users() {

	var pos int64
	batchLoad := int64(conf.WorkerLoad * conf.Workers)
	discuzRes := make([]discuzUser, batchLoad)
	var end bool
	done := make(chan bool, conf.Workers)
	fmt.Println("Converting Users")
	os.RemoveAll("../avatars/")
	conf.OrmGotalk.Raw("ALTER TABLE `user` CHANGE `username` `username` VARCHAR( 30 ) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL").Exec()
	conf.OrmGotalk.Raw("ALTER TABLE `user` DROP INDEX `username`").Exec()
	conf.OrmGotalk.Raw("ALTER TABLE `user` DROP INDEX `email`").Exec()
	for !end {
		sql := "SELECT u.uid, u.username, u.password, u.email, u.regdate, u.salt,p.qq, p.site,p.bio, f.authstr from uc_members as u left join pre_common_member_profile as p on u.uid=p.uid left join pre_common_member_field_forum as f on u.uid=f.uid limit ?,?"
		num, err := conf.Orm.Raw(sql, pos, batchLoad).QueryRows(&discuzRes)
		if err != nil {
			fmt.Println("MySQL error:", err.Error())
			return
		}
		if num == 0 {
			end = true
			continue
		}
		rows := len(discuzRes)
		for i := 0; i < conf.Workers; i++ {
			if rows == 0 {
				break
			}
			startPos := i * conf.WorkerLoad
			var load int
			if conf.WorkerLoad > rows {
				load = rows
			} else {
				load = conf.WorkerLoad
			}
			endPos := startPos + load
			jobData := discuzRes[startPos:endPos]
			go UsersJob(jobData, done)
			rows -= load
		}
		for i := 0; i < conf.Workers; i++ {
			<-done
		}
		fmt.Print(".")
		pos += batchLoad
	}
	if _, err := conf.OrmGotalk.Raw("ALTER TABLE `user` ADD UNIQUE (`username`);").Exec(); err != nil {
		fmt.Println(err)
	}
	if _, err := conf.OrmGotalk.Raw("ALTER TABLE `user` ADD INDEX ( `email` ) ;").Exec(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("")
}

func UsersJob(discuzRes []discuzUser, done chan bool) {
	defer func() {
		done <- true
	}()
	for _, discuzRow := range discuzRes {
		gotalkRes := new(models.User)
		gotalkRes.Id = discuzRow.Uid
		gotalkRes.Username = discuzRow.Username
		gotalkRes.Password = discuzRow.Password
		gotalkRes.Email = discuzRow.Email
		gotalkRes.Salt = discuzRow.Salt
		gotalkRes.Created = time.Unix(discuzRow.Regdate, 0)
		gotalkRes.Url = cutString(discuzRow.Site, 100)
		gotalkRes.Info = cutString(discuzRow.Bio, 255)
		gotalkRes.Qq = discuzRow.Qq
		if discuzRow.Authstr == "" {
			gotalkRes.IsActive = true
		} else {
			gotalkRes.IsActive = false
		}
		p, err := conf.OrmGotalk.Raw(` INSERT INTO user (
			id ,
			username ,
			nickname ,
			password ,
			url ,
			company ,
			location ,
			email ,
			avatar ,
			info ,
			weibo ,
			we_chat ,
			qq ,
			public_email ,
			followers ,
			following ,
			fav_topics ,
			is_admin ,
			is_active ,
			is_banned ,
			salt ,
			created ,
			updated 
			)
			VALUES (
			? , ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`).Prepare()
		if err != nil {
			fmt.Println(err)
		}
		_, err = p.Exec(gotalkRes.Id,
			gotalkRes.Username,
			"",
			gotalkRes.Password,
			gotalkRes.Url,
			"",
			"",
			gotalkRes.Email,
			"",
			gotalkRes.Info,
			"",
			"",
			gotalkRes.Qq,
			0,
			0,
			0,
			0,
			0,
			gotalkRes.IsActive,
			0,
			gotalkRes.Salt,
			gotalkRes.Created,
		)
		if err != nil {
			fmt.Println(err)
		}
		idstr := fmt.Sprintf("%09d", gotalkRes.Id)

		avatarFileName := fmt.Sprintf("%s/%s/%s/%s/%s_avatar_big.jpg", conf.AvatarPath, idstr[:3], idstr[3:5], idstr[5:7], idstr[7:9])
		if avatarFile, err := os.Open(avatarFileName); err == nil {
			defer avatarFile.Close()
			gotalkRes.ValidateAndSetAvatar(avatarFile, "avatar_big.jpg")
			gotalkRes.Update("Avatar")
		}
		p.Close()
	}

}

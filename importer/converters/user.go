package converters

import (
	"fmt"
	"github.com/naokij/gotalk/importer/conf"
	"github.com/naokij/gotalk/models"
	"os"
	"time"
)

type DiscuzUser struct {
	Uid         int
	Username    string
	Password    string
	Email       string
	Regdate     int64
	Salt        string
	Qq          string
	Site        string
	Bio         string
	Authstr     string
	Groupid     int
	Extcredits1 int
	Extcredits8 int
	Posts       int
	Threads     int
	Digestposts int
}

func Users() {
	var working int
	var pos int64
	done := make(chan bool)
	rows, err := NumOfRows(conf.Orm, "SELECT count(u.uid) as rows from uc_members as u left join pre_common_member_profile as p on u.uid=p.uid left join pre_common_member_field_forum as f on u.uid=f.uid")
	if err != nil {
		fmt.Println("Users Error:", err)
	}
	fmt.Println("Converting Users")

	conf.OrmGotalk.Raw("ALTER TABLE `user` CHANGE `username` `username` VARCHAR( 30 ) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL").Exec()
	conf.OrmGotalk.Raw("ALTER TABLE `user` DROP INDEX `username`").Exec()
	conf.OrmGotalk.Raw("ALTER TABLE `user` DROP INDEX `email`").Exec()

	defer func() {
		if _, err := conf.OrmGotalk.Raw("ALTER TABLE `user` ADD UNIQUE (`username`);").Exec(); err != nil {
			fmt.Println(err)
		}
		if _, err := conf.OrmGotalk.Raw("ALTER TABLE `user` ADD INDEX ( `email` ) ;").Exec(); err != nil {
			fmt.Println(err)
		}
	}()

	for i := 0; i < conf.Workers; i++ {
		if rows == 0 {
			break
		}
		var load int
		if int64(conf.WorkerLoad) > rows {
			load = int(rows)
		} else {
			load = conf.WorkerLoad
		}
		go UsersWorker(pos, conf.WorkerLoad, done)
		pos += int64(load)
		rows -= int64(load)
		working++
	}
	for {
		<-done
		working--
		if rows == 0 {
			if working == 0 {
				break
			}
		} else {
			var load int
			if int64(conf.WorkerLoad) > rows {
				load = int(rows)
			} else {
				load = conf.WorkerLoad
			}
			go UsersWorker(pos, conf.WorkerLoad, done)
			working++
			pos += int64(load)
			rows -= int64(load)
		}
	}
}

func UsersWorker(pos int64, limit int, done chan bool) {
	defer func() {
		fmt.Println("User", pos, "to", pos+int64(limit), "done")
		done <- true
	}()
	var discuzUsers []DiscuzUser
	sql := "SELECT u.uid, u.username, u.password, u.email, u.regdate, u.salt,p.qq, p.site,p.bio, f.authstr, cm.groupid, mc.extcredits1, mc.extcredits8, mc.posts, mc.threads, mc.digestposts from uc_members as u left join pre_common_member_profile as p on u.uid=p.uid left join pre_common_member_field_forum as f on u.uid=f.uid left join pre_common_member as cm on u.uid=cm.uid left join pre_common_member_count as mc on u.uid=mc.uid limit ?,?"
	if _, err := conf.Orm.Raw(sql, pos, limit).QueryRows(&discuzUsers); err != nil {
		fmt.Println("Users Worker Error:", err)
		return
	}

	for _, discuzUser := range discuzUsers {
		var isActive int
		var isBanned int
		if discuzUser.Authstr == "" {
			isActive = 1
		} else {
			isActive = 0
		}
		if discuzUser.Groupid == 4 || discuzUser.Groupid == 5 || discuzUser.Groupid == 6 || discuzUser.Groupid == 17 {
			isBanned = 1
		} else {
			isBanned = 0
		}
		insertData := map[string]interface{}{
			"id":               discuzUser.Uid,
			"username":         discuzUser.Username,
			"nickname":         "",
			"password":         discuzUser.Password,
			"url":              cutString(discuzUser.Site, 100),
			"company":          "",
			"location":         "",
			"email":            discuzUser.Email,
			"avatar":           "",
			"info":             cutString(discuzUser.Bio, 255),
			"weibo":            "",
			"we_chat":          "",
			"qq":               discuzUser.Qq,
			"public_email":     0,
			"followers":        0,
			"following":        0,
			"fav_topics":       0,
			"topics":           discuzUser.Threads,
			"comments":         (discuzUser.Posts - discuzUser.Threads),
			"reputation":       discuzUser.Extcredits1,
			"credits":          discuzUser.Extcredits8,
			"excellent_topics": discuzUser.Digestposts,
			"is_admin":         0,
			"is_active":        isActive,
			"is_banned":        isBanned,
			"salt":             discuzUser.Salt,
			"created":          time.Unix(discuzUser.Regdate, 0),
			"updated":          "0000-00-00 00:00:00"}
		if err := Map2InsertSql(conf.OrmGotalk, "user", insertData); err != nil {
			fmt.Println(err)
		}
	}
	//处理头像
	for _, discuzUser := range discuzUsers {
		user := &models.User{Id: discuzUser.Uid}
		idstr := fmt.Sprintf("%09d", discuzUser.Uid)
		avatarFileName := fmt.Sprintf("%s/%s/%s/%s/%s_avatar_big.jpg", conf.AvatarPath, idstr[:3], idstr[3:5], idstr[5:7], idstr[7:9])
		if avatarFile, err := os.Open(avatarFileName); err == nil {
			defer avatarFile.Close()
			user.ValidateAndSetAvatar(avatarFile, "avatar_big.jpg")
			user.Update("Avatar")
		}
	}
}

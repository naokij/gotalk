package converters

import (
	"fmt"
	"github.com/naokij/gotalk/importer/conf"
	"github.com/naokij/gotalk/models"
)

type DiscuzFollow struct {
	Uid       int
	Fuid      int
	Fusername string
	Dateline  int
}

func Follows() {
	fmt.Println("Converting Follow")
	var discuzFollows []DiscuzFollow
	sql := "SELECT * from pre_home_friend;"
	if _, err := conf.Orm.Raw(sql).QueryRows(&discuzFollows); err != nil {
		fmt.Println("Follows error:", err)
		return
	}

	for _, discuzFollow := range discuzFollows {
		user := &models.User{Id: discuzFollow.Uid}
		user.Read()
		followUser := &models.User{Id: discuzFollow.Fuid, Username: discuzFollow.Fusername}
		user.Follow(followUser)
	}
	fmt.Println("Follow done")
}

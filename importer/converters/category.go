package converters

import (
	"fmt"
	"github.com/naokij/gotalk/importer/conf"
	"github.com/naokij/gotalk/models"
	"strconv"
	"strings"
	"time"
)

type DiscuzCategory struct {
	Fid          int
	Fup          int
	Type         string
	Name         string
	Status       int
	Displayorder int
	Threads      int
	Posts        int
	Lastpost     string
	Description  string
	Rules        string
}

func Categories() {
	fmt.Println("Converting Categories")
	sql := "SELECT f.fid, f.fup, f.type, f.name, f.status, f.displayorder, f.threads, f.posts, f.lastpost,ff.description, ff.rules from pre_forum_forum as f left join pre_forum_forumfield as ff on f.fid=ff.fid;"
	var discuzCategories []DiscuzCategory
	if _, err := conf.Orm.Raw(sql).QueryRows(&discuzCategories); err != nil {
		fmt.Println("Categories error:", err)
		return
	}
	for _, discuzCategory := range discuzCategories {
		//lastpost:513148	求助联机问题	1403109764	enigmadudu
		type lastReply struct {
			Username   string
			CommentId  int
			TopicTitle string
			At         time.Time
		}
		var lastReplyParsed lastReply
		if discuzCategory.Lastpost != "" {
			lastPostData := strings.Split(discuzCategory.Lastpost, "\t")
			ts, _ := strconv.ParseInt(lastPostData[2], 10, 0)
			commentId, _ := strconv.ParseInt(lastPostData[0], 10, 0)
			lastReplyParsed.Username = lastPostData[3]
			lastReplyParsed.TopicTitle = lastPostData[1]
			lastReplyParsed.CommentId = int(commentId)
			lastReplyParsed.At = time.Unix(ts, 0)
		}

		insertData := map[string]interface{}{
			"id":                     discuzCategory.Fid,
			"parent_category_id":     discuzCategory.Fup,
			"depth":                  0,
			"sort":                   discuzCategory.Displayorder,
			"url_code":               "",
			"name":                   discuzCategory.Name,
			"description":            discuzCategory.Description,
			"rules":                  discuzCategory.Rules,
			"topic_count":            discuzCategory.Threads,
			"comment_count":          discuzCategory.Posts,
			"is_read_only":           0,
			"is_mod_only":            0,
			"is_hidden":              0,
			"last_reply_username":    lastReplyParsed.Username,
			"last_reply_comment_id":  lastReplyParsed.CommentId,
			"last_reply_topic_title": lastReplyParsed.TopicTitle,
			"last_reply_at":          lastReplyParsed.At}
		if err := Map2InsertSql(conf.OrmGotalk, "category", insertData); err != nil {
			fmt.Println(err)
		}
	}
	//处理深度问题
	categoryDepth(0, 0)

	fmt.Println("Categories done!")
}
func categoryDepth(parentId int, depth int) {
	if depth == 4 {
		return
	}
	categories := []models.Category{}
	depth++
	qs := conf.OrmGotalk.QueryTable("category")
	qs.Filter("parent_category_id", parentId).All(&categories)
	for _, category := range categories {
		category.Depth = depth
		category.Update("depth")
		categoryDepth(category.Id, depth)
	}
	return
}

package converters

import (
	"fmt"
	"github.com/naokij/bbcode"
	"github.com/naokij/gotalk/importer/conf"
	"github.com/naokij/gotalk/importer/csvdatafile"
	"github.com/naokij/gotalk/models"
	"time"
)

type DiscuzTopic struct {
	Tid        int
	Subject    string
	Message    string
	Authorid   int
	Author     string
	Fid        int
	Views      int
	Replies    int
	Favtimes   int
	Digest     bool
	Closed     bool
	Lastposter string
	Lastpost   int64
	Dateline   int64
	Useip      string
}

func Topics() {
	var working int
	var pos int64
	done := make(chan bool)
	rows, err := NumOfRows(conf.Orm, "pre_forum_thread")
	if err != nil {
		fmt.Println("Topics Error:", err)
	}
	fmt.Println("Converting Topics")

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
		go TopicsWorker(pos, conf.WorkerLoad, done)
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
			go TopicsWorker(pos, conf.WorkerLoad, done)
			working++
			pos += int64(load)
			rows -= int64(load)
		}
	}
}

func TopicsWorker(pos int64, limit int, done chan bool) {
	defer func() {
		fmt.Println("Topic", pos, "to", pos+int64(limit), "done")
		done <- true
	}()
	bbcodeCompiler := bbcode.NewCompiler(true, true)
	var discuzTopics []DiscuzTopic
	csvfile := csvdatafile.New("topic", fmt.Sprintf("/tmp/topic-%d.csv", pos))
	if err := csvfile.Create(); err != nil {
		fmt.Println("Topic Worker Create File:", err)
	}
	sql := "select t.tid, t.subject, p.message, t.authorid, t.author, t.fid, t.views, t.replies, t.favtimes, t.digest, t.closed, t.lastposter, t.lastpost, t.dateline,p.useip from pre_forum_thread as t left join pre_forum_post as p on t.tid=p.tid where p.first = 1 limit ?,?"
	if _, err := conf.Orm.Raw(sql, pos, limit).QueryRows(&discuzTopics); err != nil {
		fmt.Println("TopicsWorker Error:", err)
		return
	}
	csvfile.Fields = []string{
		"id",
		"title",
		"content_hex",
		"user_id",
		"username",
		"category_id",
		"pv_count",
		"comment_count",
		"bookmark_count",
		"is_excellent",
		"is_closed",
		"last_reply_username",
		"last_reply_at",
		"created",
		"updated",
		"ip"}

	for _, discuzTopic := range discuzTopics {
		content := &models.Content{Message: bbcodeCompiler.Compile(discuzTopic.Message)}
		var contentHex string
		var err error
		if contentHex, err = content.Insert(); err != nil {
			fmt.Println("Topic Worker Content Insert:", err)
		}
		if err := csvfile.AppendRow(
			discuzTopic.Tid,
			discuzTopic.Subject,
			contentHex,
			discuzTopic.Authorid,
			discuzTopic.Author,
			discuzTopic.Fid,
			discuzTopic.Views,
			discuzTopic.Replies,
			discuzTopic.Favtimes,
			Btoi(discuzTopic.Digest),
			Btoi(discuzTopic.Closed),
			discuzTopic.Lastposter,
			time.Unix(discuzTopic.Lastpost, 0).Format("2006-01-02 15:04:05"),
			time.Unix(discuzTopic.Dateline, 0).Format("2006-01-02 15:04:05"),
			time.Unix(discuzTopic.Lastpost, 0).Format("2006-01-02 15:04:05"),
			discuzTopic.Useip); err != nil {
			fmt.Println(err)
		}
		if err := csvfile.Flush(); err != nil {
			fmt.Println(err)
		}
	}
	csvfile.Close()
	if err := csvfile.LoadToMySQL(conf.OrmGotalk); err != nil {
		fmt.Println(err)
	}
	csvfile.Remove()
}

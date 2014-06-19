package converters

import (
	"fmt"
	"github.com/naokij/bbcode"
	"github.com/naokij/gotalk/importer/conf"
	"github.com/naokij/gotalk/importer/csvdatafile"
	"github.com/naokij/gotalk/models"
	"time"
)

type DiscuzComment struct {
	Pid      int
	Tid      int
	Message  string
	Authorid int
	Author   string
	Useip    string
	Dateline int64
}

func Comments() {
	var working int
	var pos int64
	done := make(chan bool)
	rows, err := NumOfRows(conf.Orm, "select count(pid) as rows  from pre_forum_post where first = 0")
	if err != nil {
		fmt.Println("Comments Error:", err)
	}
	fmt.Println("Converting Comments")

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
		go CommentsWorker(pos, conf.WorkerLoad, done)
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
			go CommentsWorker(pos, conf.WorkerLoad, done)
			working++
			pos += int64(load)
			rows -= int64(load)
		}
	}
}

func CommentsWorker(pos int64, limit int, done chan bool) {
	defer func() {
		fmt.Println("Comment", pos, "to", pos+int64(limit), "done")
		done <- true
	}()
	bbcodeCompiler := bbcode.NewCompiler(true, true)
	var discuzComments []DiscuzComment
	csvfile := csvdatafile.New("comment", fmt.Sprintf("/tmp/comment-%d.csv", pos))
	if err := csvfile.Create(); err != nil {
		fmt.Println("Coment Worker Create File:", err)
	}
	sql := "select pid, tid, message, authorid, author, useip, dateline  from pre_forum_post where first = 0 limit ?,?"
	if _, err := conf.Orm.Raw(sql, pos, limit).QueryRows(&discuzComments); err != nil {
		fmt.Println("CommentsWorker Error:", err)
		return
	}
	csvfile.Fields = []string{
		"id",
		"topic_id",
		"content_hex",
		"user_id",
		"username",
		"ip",
		"created"}

	for _, discuzComment := range discuzComments {
		content := &models.Content{Message: bbcodeCompiler.Compile(discuzComment.Message), TopicId: discuzComment.Tid, CommentId: discuzComment.Pid}
		var contentHex string
		var err error
		if contentHex, err = content.Insert(); err != nil {
			fmt.Println("Comments Worker Content Insert:", err)
		}
		if err := csvfile.AppendRow(
			discuzComment.Pid,
			discuzComment.Tid,
			contentHex,
			discuzComment.Authorid,
			discuzComment.Author,
			discuzComment.Useip,
			time.Unix(discuzComment.Dateline, 0).Format("2006-01-02 15:04:05"),
		); err != nil {
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

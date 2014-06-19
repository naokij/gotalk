package main

import (
	"fmt"
	"github.com/naokij/gotalk/importer/converters"
	"time"
)

func main() {
	startTime := time.Now()
	//converters.Users()
	converters.Topics()
	converters.Comments()
	timeUsed := time.Since(startTime)
	fmt.Printf("任务耗时: %4.2f分钟\n", timeUsed.Minutes())
}

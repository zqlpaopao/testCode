package main

import (
	"fmt"
	"github.com/robfig/cron"
	"time"
)

func main(){
	//cronTask := cron.New()
	//cronTask.

	//秒 分 时  天 月 周
	t , err :=cron.Parse("*/2 35-37 10 * * *")
	//2021-07-09 09:01:00 +0800 CST
	fmt.Println(1111,t.Next(time.Now()))
	fmt.Println(err)

}
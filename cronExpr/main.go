/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/15 上午11:10
 */

package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"os"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)
	// 秒 分  时  天  月  年 星期 支持到秒 7位
	//第一种   支持7位，只写5位就是分钟了
	if expr, err = cronexpr.Parse("*/2 * * * * *"); nil != err {
		fmt.Println(err)
	}
	a := cronexpr.MustParse("*/2 * * * * *")
	cronexpr.Expression{}
	fmt.Println(expr.)
	fmt.Println(expr.Next(time.Now()))
	os.Exit(3)
	//第二种MustParse() 返回*Expression 没有错误，认为表达式正确

	//获取当前时间
	now = time.Now()
	//传入当前时间，返回下次调度时间
	// 2020-05-02 22:25:00 +0800 CST
	nextTime = expr.Next(now)

	//定时器超时触发
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了")
	})

	//防止主函数退出
	time.Sleep(10 * time.Second)
	fmt.Println(nextTime)

}

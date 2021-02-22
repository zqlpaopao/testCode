/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/2/22 下午4:20
 */

package main

import (
	"fmt"
	"runtime"
	"test/goroutinePool/src"
	"time"
)

func main(){
	//创建任务
	t:= src.InitTask(func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	})

	//创建协程池
	p := src.NewPool(4)

	//模拟任务添加
	tN := 0//测试数量
	go func() {
		for{
			p.AddTask(t)
			tN++
			if 0 == 1000%10 {
				fmt.Println(runtime.NumGoroutine())
				time.Sleep(time.Second)
			}
		}
	}()

	//启动协程池
	p.Run()
	fmt.Println(runtime.NumGoroutine())
}
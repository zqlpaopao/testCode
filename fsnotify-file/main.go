package main

import (
	"log"
	"github.com/fsnotify/fsnotify"
)
/**
https://mp.weixin.qq.com/s/0GcjkYEiO6A7Wah3fuhUgA

*/
func main() {
	// 创建文件/目录监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 打印监听事件
				log.Println("event:", event)
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()
	// 监听当前目录
	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
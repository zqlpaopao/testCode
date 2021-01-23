/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/23 上午10:45
 */


package main

import (
"fmt"
"time"
)

type token struct{}
/**
* 哪些操作会使channel panic
* 1、关闭一个nil的channel
* 2、关闭一个已经关闭的channel
* 3、像一个已经关闭的channel发送数据
* 4、像为初始化的channel 发送数据会deadlock!
 */
func main() {
	num := 4
	var chs []chan token
	// 4 个work
	for i := 0; i < num; i++ {
		chs = append(chs, make(chan token))
	}
	for j := 0; j < num; j++ {
		go worker(j, chs[j], chs[(j+1)%num])
	}
	// 先把令牌交给第一个
	chs[0] <- struct{}{}
	select {}
}

func worker(id int, ch chan token, next chan token) {
	for {
		// 对应work 取得令牌
		token := <-ch
		fmt.Println(id + 1)
		time.Sleep(1 * time.Second)
		// 传递给下一个
		next <- token
	}
}

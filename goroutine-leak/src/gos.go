/**
 * @Author: zhangsan
 * @Description:
 * @File:  gos
 * @Version: 1.0.0
 * @Date: 2021/3/25 下午2:11
 */

package src

import (
	"fmt"
	"time"
)

func Go1() {
	c := make(chan int)
	go func() {
		<-c
	}()
	time.Sleep(2*time.Second)
	fmt.Println("G1 exit")
}
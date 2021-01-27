/**
 * @Author: zhangsan
 * @Description:
 * @File:  ticker
 * @Version: 1.0.0
 * @Date: 2021/1/27 上午11:29
 */

package pk

import (
	"fmt"
	"time"
)

/////////////////////////////////////////
//--> @Description Ticker
//--> @Param
//--> @return
////////////////////////////////////////
func Ticker(){
	// fmt.Println("公众号：Golang来啦")
	fmt.Println(time.Now())
	ticker1 := time.NewTicker(2*time.Second)
	for {
		curTime := <-ticker1.C
		fmt.Println(curTime)
	}
}

/////////////////////////////////////////
//--> @Description Ticker
//--> @Param
//--> @return
////////////////////////////////////////
func SelectTicker(){
	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("done")
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	ticker.Stop()   // 取消定时器
	done <- true
	fmt.Println("Ticker stopped")
}

/////////////////////////////////////////
//--> @Description Tick
//--> @Param
//--> @return
////////////////////////////////////////
func Tick(){
	for {
		ss := time.Tick(3 * time.Second)
		time.Sleep(1 *time.Second)
		fmt.Println(ss)

	}

}
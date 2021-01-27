/**
 * @Author: zhangsan
 * @Description:
 * @File:  timer
 * @Version: 1.0.0
 * @Date: 2021/1/27 上午10:49
 */

package pk

import (
	"fmt"
	"time"
)

/////////////////////////////////////////
//--> @Description time.Timer
//--> @Param
//--> @return
//now time 2021-01-27 10:45:54
//end time 2021-01-27 10:45:57
////////////////////////////////////////
func timer(){
	fmt.Println("now time",time.Now().Format("2006-01-02 15:04:05"))
	t := time.NewTimer(3 *time.Second)
	<- t.C
	fmt.Println("end time",time.Now().Format("2006-01-02 15:04:05"))
}

/////////////////////////////////////////
//--> @Description 停止定时器
//--> @Param
//--> @return
////////////////////////////////////////
func StopTimer(){
	fmt.Println("start time: ", time.Now().Format("2006-01-02 15:04:05"))
	timer1 := time.NewTimer(3 * time.Second)
	go func() {
		//协程会阻塞,Timer.C会返回到期的时间
		curTime := <-timer1.C
		// 定时时间到了之后执行打印操作
		fmt.Println("Timer 1 here, current time: ", curTime.Format("2006-01-02 15:04:05"))
	}()


	timer2 := time.NewTimer(3 * time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 here")
	}()
	// 取消定时器,返回true OR false 代表定时器是否取消成功
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stop")
	}
	time.Sleep(8 * time.Second)
}

/////////////////////////////////////////
//--> @Description 重置定时器，会再次执行一次
//--> @Param
//--> @return
////////////////////////////////////////
func ResetTimer(){
	fmt.Println("start time: ", time.Now().Format("2006-01-02 15:04:05"))
	t1 := time.NewTimer(5 * time.Second)

	ok := t1.Reset(3 * time.Second)
	// 重新设置为 3s
	fmt.Println("ok: ", ok)
	curTime := <-t1.C

	fmt.Println("now time: ", curTime.Format("2006-01-02 15:04:05"))
}

/////////////////////////////////////////
//--> @Description 指定执行次数
//--> @Param
//--> @return
////////////////////////////////////////
func RunTimes(){
	timer1 := time.NewTimer(5 * time.Second)
	fmt.Println("start time: ", time.Now().Format("2006-01-02 15:04:05"))

	go func() {
		count := 0
		for {
			<-timer1.C
			fmt.Println("timer 1", time.Now().Format("2006-01-02 15:04:05"))

			count++
			fmt.Println("调用 Reset() 重新设置过期时间，将时间修改为 2s")
			ok := timer1.Reset(2*time.Second)
			fmt.Println("ok",ok)
			if count > 2 {
				fmt.Println("调用 Stop() 停止定时器")
				timer1.Stop()
			}
		}
	}()

	time.Sleep(15 * time.Second)
	fmt.Println("end time：", time.Now().Format("2006-01-02 15:04:05"))
}
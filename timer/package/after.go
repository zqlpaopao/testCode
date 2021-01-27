/**
 * @Author: zhangsan
 * @Description:
 * @File:  after
 * @Version: 1.0.0
 * @Date: 2021/1/27 下午1:36
 */

package pk

import (
	"fmt"
	"time"
)

/////////////////////////////////////////
//--> @Description
//--> @Param
//--> @return
////////////////////////////////////////
func After(){
	fmt.Println("start time",time.Now().Format("2006-01-02 15:04:05"))
	after1 := time.After(2*time.Second)
	fmt.Println("读取前, time: ",time.Now().Format("2006-01-02 15:04:05"))
	curTime1 := <-after1    // 读取
	fmt.Println("读取后 , time: ",time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("读取的内容",curTime1)

	fmt.Println()


	fmt.Println("start time",time.Now().Format("2006-01-02 15:04:05"))
	after2 := time.After(2*time.Second)
	time.Sleep(3*time.Second)         // 为了使定时时间过期
	fmt.Println("读取前, time: ",time.Now().Format("2006-01-02 15:04:05"))
	curTime2 := <-after2   // 读取
	fmt.Println("读取后 , time: ",time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("读取的内容",curTime2)
}

/////////////////////////////////////////
//--> @Description 超时控制
//--> @Param
//--> @return
////////////////////////////////////////
func TimeoutAfter(){
	ch := make(chan string, 1)
	go func() {
		// 假设我们在这里执行一个外部调用，2秒之后将结果写入 ch
		time.Sleep(time.Second * 2)
		ch <- "success"
	}()

	select {
	case res := <-ch:
		fmt.Println(res)
	case <-time.After(time.Second * 3):
		fmt.Println("timeout 1")
	}

	time.Sleep(5 *time.Second)
}

/////////////////////////////////////////
//--> @Description
//--> @Param
//--> @return
////////////////////////////////////////
func AfterFunc(){
	ch := make(chan int, 1)

	time.AfterFunc(10*time.Second, func() {
		fmt.Println("10 seconds over....")
		ch <- 8
	})

	for {
		select {
		case n := <-ch:
			fmt.Println(n, "is arriving")
			fmt.Println("Done!")
			return
		default:
			time.Sleep(3 * time.Second)
			fmt.Println("time to wait")
		}
	}
}

/////////////////////////////////////////
//--> @Description 重置定时器
//--> @Param
//--> @return
////////////////////////////////////////
func StopAfterFunc(){
	fmt.Println("start time",time.Now().Format("2006-01-02 15:04:05"))

	timer1 :=time.AfterFunc(3*time.Second, func() {
		fmt.Println("10 seconds over....")
		fmt.Println("run time",time.Now().Format("2006-01-02 15:04:05"))

	})
	time.Sleep(1*time.Second)
	fmt.Println("sleep end  time",time.Now().Format("2006-01-02 15:04:05"))

	timer1.Reset(6*time.Second)
	fmt.Println("reset time",time.Now().Format("2006-01-02 15:04:05"))


	time.Sleep(15*time.Second)

}
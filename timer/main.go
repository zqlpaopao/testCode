/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/27 上午10:30
 */

package main

import pk "test/timer/package"

/////////////////////////////////////////
//--> @Description time.Sleep 延迟时间
//--> @Param
//--> @return
////////////////////////////////////////
func timeSleep(){

}

func main(){

	//1.time.Sleep

	//2.time.Timer
	//now time 2021-01-27 10:45:54
	//end time 2021-01-27 10:45:57
	//timer()

	//停止定时器
	//pk.StopTimer()
	//重置定时器
	//pk.ResetTimer()
	//设置重试次数
	//pk.RunTimes()



	//pk.Ticker()
	//结合select use
	//pk.SelectTicker()


	//pk.Tick()

	//pk.After()
	//pk.TimeoutAfter()



	pk.AfterFunc()
	//pk.StopAfterFunc()

	//t := time.AfterFunc(8*time.Second, func() {
	//	fmt.Println("Golang来啦")
	//})
	//for {
	//	select {
	//	case <-t.C:
	//		fmt.Println("seekload")
	//		break
	//	}
	//}
}



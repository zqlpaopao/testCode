/**
 * @Author: zHangSan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/4/16 下午4:44
 */

package main

import "runtime"

func main(){
	//goGC()

	//gOTraceBack() 查看crash的时候相信信息

	goDebugGc()

}
/*****************************GODEBUG-异常发生打印详细信息****************************************/
func goDebugGc(){
	gc()
}





/*****************************GOTRACEBACK-异常发生打印详细信息****************************************/
func gOTraceBack(){
	panic("kerboom")

}






/*****************************goGC--修改GC频率****************************************/
func goGC(){
	gc()
	runtime.GC()
	gc()
	gc()
	gc()
}

func gc(){
	for i := 0; i< 10000;i++{
		g := make(map[string]string)
		g = g
	}
}

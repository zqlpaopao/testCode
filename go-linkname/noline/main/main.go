/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/3/1 上午10:03
 */

package main

//go:noinline
func appStr(w string)string{
	return w + " word"
}

func main(){
	appStr("hello")
}

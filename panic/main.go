/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/26 下午5:51
 */

package main

import (
	"fmt"
)

func main() {
	defer Recovers()()
	panic("ff")
	fmt.Println(1)
}

func Recovers() func(){
	return func() {
		if err := recover();err  != nil {
			fmt.Println("vv", err)
		}
	}
}
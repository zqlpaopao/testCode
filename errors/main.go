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
	//"errors"
)

func main() {
	defer Recovers()()

	panic("ff")
	//e := errors.New("原始错误e")
	//w := fmt.Errorf("Wrap了一个错误%w", e)
	//fmt.Println(w)
	//fmt.Println(errors.Unwrap(w))
	fmt.Println(1)
}

func Recovers() func(){
	return func() {
		if err := recover();err  != nil {
			fmt.Println("vv", err)
		}
	}
}
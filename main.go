/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/29 下午6:42
 */

package main

import "fmt"

const (
	aa = 1<< iota
)

func main() {
	var a = map[string]string{"aa":"bb"}

	for v,v1 := range a{
		fmt.Println(v)
		fmt.Println(v1)
	}

}
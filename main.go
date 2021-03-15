/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/29 下午6:42
 */

package main

import "fmt"

func main() {
	a1 :=[]int64{1,2,3,4,5,6,7}
	a2 := a1[3:6]
	fmt.Println(a2)

	var v int64
	for _,v = range a1{
		//fmt.Println(i)
		fmt.Println(v)
	}
}
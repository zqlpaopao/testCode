/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/3/17 下午5:57
 */

package main

import (
	"fmt"
	"unsafe"
)
type c struct {

}
type b struct {
	namea string
	//aa int64
	bb string
}

type a struct {
	c
	b
	name string
	age int64
}

func main(){
	var s a
	fmt.Println(unsafe.Offsetof(s.namea))
	fmt.Println(unsafe.Offsetof(s.name))
	fmt.Println(unsafe.Offsetof(s.age))


}

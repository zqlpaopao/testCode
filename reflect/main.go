/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/12 下午5:13
 */

package main

import (
	"fmt"
	"reflect"
)

type Struc struct {
	name string
	age int
}

type in interface {
	set()
}

func ( s *Struc) set(){

}

func main(){
	//获取具体表层类型
	var a Struc
	ty := reflect.TypeOf(a) //main.Struc

	//返回真实的数据结构类型
	fmt.Println(ty.Kind().String())  //struct

	//获取对应的值
	fmt.Println(reflect.ValueOf(a))//{ 0}


}
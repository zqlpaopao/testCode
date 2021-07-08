package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan int64,10)
	ch<-1

	go is(ch)

	time.Sleep(10*time.Second)
}

//-- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func is(ch chan int64){
	v ,ok := <- ch
	fmt.Println(v)
	fmt.Println(ok)
}
package main

import (
	"fmt"
	"test/unsafe/src"
	"unsafe"
)

type model struct {
	users unsafe.Pointer
}


type user struct {
	name string
	age string
}
type user1 struct {
	name string
	age int
}

func main(){
	u := src.Usera
	type m src.User


	un := unsafe.Pointer(&u)
	fmt.Println(uintptr(un))
	//结构体转换
	fmt.Println((*m)(un))

}

func(u *user1)Get()string{
	return u.name
}
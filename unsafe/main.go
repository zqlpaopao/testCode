package main

import (
	"fmt"
	"unsafe"
)
// 64位平台，对齐参数是8
type User1 struct {
	A int32 // 4     8
	B []int32 //24   24
	C string //16    16
	D bool //1       8

}

type User2 struct {
	B []int32 //24    24
	A int32  //4      4
	D bool  //1       4
	C string //16      16
}

type User3 struct {
	D bool//1       8
	B []int32//24   24
	A int32//4      4 + 4
	C string//16    16
}

func main()  {
	var u1 User1
	var u2 User2
	var u3 User3

	fmt.Println("u1 size is ",unsafe.Sizeof(u1))
	fmt.Println("u2 size is ",unsafe.Sizeof(u2))
	fmt.Println("u3 size is ",unsafe.Sizeof(u3))
}

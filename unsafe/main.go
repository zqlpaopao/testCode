package main

import (
	"fmt"
	"unsafe"
)
func main()  {
	fmt.Println(unsafe.Sizeof(test1{})) // 8
	fmt.Println(unsafe.Sizeof(test2{})) // 4
}
type test1 struct {
	a int32
	b struct{}
}

type test2 struct {
	a struct{}
	b int32
}


package main

/*
#include <stdio.h>
#include "awesome.h"
#cgo linux CFLAGS: -L./ -I./
#cgo linux LDFLAGS: -L./ -I./ -lcshared
*/
import "C"

import (
	"fmt"
)

func main() {

	//s := "hello"
	var i int
	i = C.Add(1, 2)
	//C.Test()
	fmt.Println(i)
}

func test(a, b int) int {
	fmt.Println("test", a, b)
	return a + b
}

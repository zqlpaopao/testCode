package main

import (
	"fmt"
	"reflect"
)

type st struct {
	aa string
	bb int64
	C  string
	d  []int64
}

func main() {
	var a = st{
		aa: "al",
		bb: 2,
		C:  "gv",
		d:  []int64{1, 3, 5, 6},
	}

	as := reflect.ValueOf(a)
	c := as.FieldByName("d").InterfaceData()
	//for _, v := range as.FieldByName("d").([]int64) {
	//
	//}

	fmt.Println(c)
}

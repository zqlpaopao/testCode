package main

import (
	"fmt"
	"strconv"
)

type IntModel int

const (
	One IntModel = iota
	Two
	Three
)


func main(){
	p := Two
	var is int
	is = 1
	if p == IntModel(is){
		fmt.Println(p.String())
	}
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[One-0]
	_ = x[Two-1]
	_ = x[Three-2]
}

const _IntModel_name = "OneTwoThree"

var _IntModel_index = [...]uint8{0, 3, 6, 11}

func (i IntModel) String() string {
	if i < 0 || i >= IntModel(len(_IntModel_index)-1) {
		return "IntModel(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _IntModel_name[_IntModel_index[i]:_IntModel_index[i+1]]
}

func Switch(i int )string{
	switch i {
	case 0:
		return "one"
	case 1:
		return "two"
	case 2:
		return "three"
	default:
		return ""
	}
}
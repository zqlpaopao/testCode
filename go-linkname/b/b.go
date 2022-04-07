package b

import (
	_ "unsafe"
)

//go:linkname sayHi  b.sayHi
//go:nosplit
func sayHi()string{
	return "k"
}


//go:linkname aaaStr test/go-linkname/a.Aa
var aaaStr = "111"

//go:linkname aBStr test/go-linkname/a.Ab
var aBStr = "11s1"

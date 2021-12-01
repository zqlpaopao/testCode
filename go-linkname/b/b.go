package b

import (
	_ "unsafe"
)

//go:linkname sayHi  mysql/go-linkname/a.hello
//go:nosplit
func sayHi()string{
	return "k"
}


//go:linkname aaaStr mysql/go-linkname/a.Aa
var aaaStr = "111"

//go:linkname aBStr mysql/go-linkname/a.Ab
var aBStr = "11s1"
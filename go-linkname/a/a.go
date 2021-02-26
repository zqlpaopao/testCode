package a

import (
	_ "unsafe"
	_"test/go-linkname/b"
)

//变量是可以的
var Aa = ""

//常来是不可以的
const Ab = ""

//只能通过其他正常函数调用才生效，直接调用没反应
func Hello()string

func Greet() string {
	return Hello()
}


package a

import (
	_ "test/go-linkname/b"
	_ "unsafe"
)

//变量是可以的
var Aa = ""

//常来是不可以的
const Ab = ""

//只能通过其他正常函数调用才生效，直接调用没反应
//即使是小写的也可以在外部在重新实现
//go:linkname hello b.sayHi
func hello()string

func Greet() string {
	return hello()
}
//go:generate  echo 1
func Adds(){

}

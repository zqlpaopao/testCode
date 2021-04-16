package main

import (
	//_ "embed"
)

////go:embed test.txt
var testString string  // 当前目录，解析为string类型

////go:embed test.txt
var testByte []byte  // 当前目录，解析为[]byte类型

////go:embed test/test.txt
var testAbsolutePath string  // 子目录，解析为string类型

////go:embed notExistsFile
var testErr0 string // 文件不存在，编译报错：pattern notExistsFile: no matching files found

////go:embed dir
var testErr1 string // dir是目录，编译报错：pattern dir: cannot embed directory dir: contains no embeddable files

////go:embed ../test.txt
//var testErr2 string // 相对路径，不是当前目录或子目录，编译报错：pattern ../test.txt: invalid pattern syntax

////go:embed D:\test.txt
//var testErr3 string // 绝对路径，编译报错：pattern D:\test.txt: no matching files found

func main() {
	println(testString)
	println(testByte)
	println(string(testByte))
	println(testAbsolutePath)
}

//https://studygolang.com/articles/33398?fr=sidebar

//go.mod 要改成是1.16,不然报错
//embed/main.go:10:3: go:embed requires go1.16 or later (-lang was set to go1.14; check go.mod)
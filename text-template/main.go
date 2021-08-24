package main

import (
	"fmt"
	"os"
	"text/template"
)

/**

template 包是数据驱动的文本输出模板，其实就是在写好的模板中填充数据。

*/

func main(){
	// 模板定义
	tepl := "My name is {{ . }}"

	// 解析模板
	tmpl, err := template.New("test").Parse(tepl)
	fmt.Println(err)

	// 数据驱动模板
	data := "jack"
	err = tmpl.Execute(os.Stdout, data)
}
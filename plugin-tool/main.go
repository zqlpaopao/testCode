package main

import (
	"fmt"
	"log"
	"plugin"
)

var In int = 65

func Func()string{
	return "zhangSan"
}


func main(){
	p ,err := plugin.Open("plugin.so")
	if err != nil{
		log.Panicln(err)
	}

	v ,err := p.Lookup("In")
	fmt.Println(*v.(*int))

	f , err := p.Lookup("Func")

	fn := f.(func()string)
	fmt.Println(fn())
}
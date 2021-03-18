package main

import "fmt"
type Person struct {
	Name string
	Age int64
}

func main()  {
	p := Person{
		Name: "zhangsan",
		Age: int64(8),
	}
	fmt.Printf("原始struct地址是：%p\n",&p)
	modifiedAge(p)
	fmt.Println(p)
}


func modifiedAge(p Person)  {
	fmt.Printf("函数里接收到struct的内存地址是：%p\n",&p)
	p.Age = 10
}


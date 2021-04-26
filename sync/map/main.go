package main

import (
	"fmt"
	"sync"
)

func main(){
	var m sync.Map

	//数据写入
	fmt.Println("写入操作")
	m.Store("name","zhangsan")
	m.Store("age",23)

	fmt.Println("读取操作")
	//数据读取 ok 存在就是true，不村子啊是false
	name,ok := m.Load("name1")
	fmt.Println(name)
	fmt.Println(ok)

	fmt.Println("遍历操作")
	//便利数据
	m.Range(func(k,v interface{})bool{
		if k == "name"{
			return false
		}
		fmt.Println(k)
		fmt.Println(v)
		return true
	})


	fmt.Println("删除操作")
	//删除
	//m.Delete("age")
	name,ok = m.Load("age")
	fmt.Println(name)
	fmt.Println(ok)

	fmt.Println("读写操作")
	//存在就读取，不存在就插入
	m.LoadOrStore("name","zhangSan")
	name ,ok = m.Load("name")
	fmt.Println(name)
	fmt.Println(ok)




	fmt.Println("获取mlen")

	var i int
	m.Range(func(k,v interface{})bool{
		i ++
		return true
	})
	fmt.Println(i)

	fmt.Println("对比")
	if name ,ok = m.Load("name");ok{
		//if name.(string) =="zhangsan"{
		if name =="zhangsan"{
			fmt.Println("==")
		}else {
			fmt.Println("!=")
		}
	}
	if name ,ok = m.Load("age");ok{
		if name ==23{
			//if name =="zhangsan"{
			fmt.Println(true)
		}else{
			fmt.Println("false")
		}
	}


}




















































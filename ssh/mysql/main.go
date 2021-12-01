package main

import (
	"fmt"
	"test/ssh/mysql/src"
)

func main(){
	//连接反代服务器
	client, err := src.DialWithPasswd(&src.Config{
		Addr:   "11.191.161.27:221",
		User:   "root",
		Passwd: "Noah0b2",
	})
	if err != nil {
		panic(err)
	}
	s ,err := client.GetClient().NewSession()
	fmt.Println(err)
	//测试输出
	out ,err := src.TestSShServer(s,client,"more ./test.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

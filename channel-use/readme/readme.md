# 1、作为控制开关

用无缓冲通道顺序打印1，2，3  ...100

```go
func print(){
	c:= make(chan struct{})
	go oddNum(c)
	go evenNum(c)
	time.Sleep(time.Second)
	fmt.Println("done")
}


func oddNum (c chan struct{}) {
	for i := 1; i <= 100; i++ {
		<-c
		if i % 2 == 1 {
			fmt.Println("奇数：", i)
		}
	}
}

func evenNum (c chan struct{}) {
	for i := 1; i <= 100; i++ {
		c <- struct{}{}
		if i % 2 == 0 {
			fmt.Println("偶数：", i)
		}
	}
}
```



# 2、协程间信号通信

当一个chan呗close的时候，读取会收到struct{}{},继续向下走，可用于协程间通信，类似于sync.Cond的通知

```go
package main

import (
	"fmt"
	"time"
)

func main(){
	//打印奇数偶数
	//print()//顺序打印奇数偶数

	single()
	time.Sleep(5 * time.Second)
}

//////////////////////////////////////////打印奇数偶数/////////////////////////////
func single(){
	c := make(chan struct{})
	go func(){
		for i := 0;i <10;i++{
			fmt.Println(i)
		}
		close(c)
	}()
	g := <-c
	fmt.Printf("%#v",g)
	fmt.Println("done")
	fmt.Println(g == struct{}{})
}
```

```go
0
1
2
3
4
5
6
7
8
9
struct {}{}done
true

```
































































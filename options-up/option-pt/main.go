package main

import "fmt"

//函数选项模式的普通版本
//可以发现只要是第一次定义了参数 ，后续没有可变性

type options struct {
	count int
	maxSize int
}


type Fn func(log *Log)

type Log struct {
	option options
}

func(l *Log)clone(f...Fn)*Log{
	for _ ,v := range f{
		v(l)
	}
	return l
}

func NewLog(f... Fn)*Log{
	l := &Log{}
	return l.clone(f...)
}

func SetMaxCount(c int)Fn{
	return func(log *Log) {
		log.option.count = c
	}
}

func SetMaxSize(c int)Fn{
	return func(log *Log) {
		log.option.maxSize = c
	}
}

func main(){
	l := NewLog(SetMaxCount(4),SetMaxSize(5))

	fmt.Println(l.option.maxSize)
	fmt.Println(l.option.count)

}

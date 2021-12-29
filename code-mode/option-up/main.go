package main

import "fmt"

//函数选项的升级版
//后续参数具有可变性，而且不影响之前的操作

// An Option configures a Logger.
type Option interface {
	apply(*Log)
}
// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Log)

func (f optionFunc) apply(log *Log) {
	f(log)
}


//Log 最终操作实体
type Log struct {
	count int
	maxSize int
}

func NewLog(f... Option)*Log{
	log := &Log{
		count:   1,
		maxSize: 2,
	}
	return log.WithOptions(f...)
}

func (l *Log) clone() *Log {
	copy := *l
	return &copy
}

//WithOptions 为什么没有取址，为这个样就可以复制多个l，并且之间没冲突
func (l Log)WithOptions(opts ...Option) *Log {
	c := l.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

//NewCount 实际的复制操作
func NewCount(c int)optionFunc{
	return func(log *Log) {
		log.count = c
	}
}
func NewMaxSize(c int)optionFunc{
	return func(log *Log) {
		log.maxSize = c
	}
}

func main(){

	l := NewLog(
		NewCount(2),
		NewMaxSize(3))

	l1 := l.WithOptions(NewCount(23),
		NewMaxSize(33))

	fmt.Println(l.maxSize)
	fmt.Println(l.count)

	fmt.Println(l1.maxSize)
	fmt.Println(l1.count)


}
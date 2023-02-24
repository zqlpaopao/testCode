package main

import "fmt"

//观察者模式
// https://mp.weixin.qq.com/s/4NqjkXVqFPamEc_QsyRipA

/*
	观察者模式 (Observer Pattern)，定义对象间的一种一对多依赖关系，
	使得每当一个对象状态发生改变时，其相关依赖对象皆得到通知，依赖对象在收到通知后，
	可自行调用自身的处理程序，实现想要干的事情，比如更新自己的状态。

*/

// Subject 接口，它相当于是发布者的定义
type Subject interface {
	Subscribe(observer Observer)
	Notify(msg string)
}

// Observer 观察者接口
type Observer interface {
	Update(msg string)
}

// Subject 实现

type SubjectImpl struct {
	observers []Observer
}

// Subscribe 添加观察者（订阅者）
func (sub *SubjectImpl) Subscribe(observer Observer) {
	sub.observers = append(sub.observers, observer)
}

// Notify 发布通知
func (sub *SubjectImpl) Notify(msg string) {
	for _, o := range sub.observers {
		o.Update(msg + "Notify")
	}
}

// Observer1 Observer1
type Observer1 struct{}

// Update 实现观察者接口
func (Observer1) Update(msg string) {
	fmt.Printf("Observer1: %s\n", msg)
}

// Observer2 Observer2
type Observer2 struct{}

// Update 实现观察者接口
func (Observer2) Update(msg string) {
	fmt.Printf("Observer2: %s\n", msg)
}

func main() {
	sub := &SubjectImpl{}
	sub.Subscribe(&Observer1{})
	sub.Subscribe(&Observer2{})
	sub.Notify("Hello")
}

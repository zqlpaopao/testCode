package main

import (
	"container/list"
	"fmt"
	"log"
)

func main(){


	var stack = list.New()
	stack.PushBack(1)
	stack.PushBack("string")

	//获取双向链表长度
	fmt.Println(stack.Len())


	//获取双向链表的第一个元素，如果没有返回nil
	fmt.Println(stack.Front().Value)
	fmt.Println(stack.Front().Next().Value)
	//fmt.Println(stack.Front().Prev().Value)

	//获取最后一个元素，如果没有为nil
	fmt.Println(stack.Back().Value)
	fmt.Println(stack.Back().Prev().Value)

	//插入元素在第一个元素之前
	//log.Println(stack.PushFront("first").Value)
	fmt.Println(stack.Front().Value)

	//PushFrontList创建链表other的拷贝，并将拷贝的最后一个位置连接到链表l的第一个位置。
	//stack.PushFrontList(stack)
	log.Println(stack.Len())

	//PushBack将一个值为v的新元素插入链表的最后一个位置，返回生成的新元素。
	stack.PushBack("end")
	log.Println(stack.Back().Value)

	//PushBack创建链表other的拷贝，并将链表l的最后一个位置连接到拷贝的第一个位置。
	//stack.PushBackList(stack)
	log.Println(stack.Len())

	//InsertBefore将一个值为v的新元素插入到mark前面，并返回生成的新元素。如果mark不是l的元素，l不会被修改。
	stack.InsertAfter("string",&list.Element{Value: "end"})
	log.Println(stack.InsertAfter("string",stack.Front()).Prev().Value)

	//InsertAfter将一个值为v的新元素插入到mark后面，并返回新生成的元素。如果mark不是l的元素，l不会被修改。

	//MoveToFront将元素e移动到链表的第一个位置，如果e不是l的元素，l不会被修改。
	stack.MoveToFront(stack.Back())
	fmt.Println(stack.Front().Value)

	//MoveToBack将元素e移动到链表的最后一个位置，如果e不是l的元素，l不会被修改。
	stack.MoveToBack(stack.Front())

	//MoveBefore将元素e移动到mark的前面。如果e或mark不是l的元素，或者e==mark，l不会被修改。
	//MoveAfter将元素e移动到mark的后面。如果e或mark不是l的元素，或者e==mark，l不会被修改。
	stack.MoveAfter(stack.Back(),stack.Front())
	stack.MoveBefore(stack.Back(),stack.Front())

	//移除元素
	stack.Remove(stack.Front())
	log.Println(stack.Front().Value)
	//清空链表
	stack.Init()
	fmt.Println(stack.Len())
}

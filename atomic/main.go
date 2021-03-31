/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/3/29 下午5:05
 */

package main

import (
	"fmt"
	"sync/atomic"
)

func main(){
	//原子性获取值
	var i int32 =5
	fmt.Println(atomic.LoadInt32(&i))//5

	var is int64 = 6
	fmt.Println(atomic.LoadInt64(&is))//6

	//atomic.LoadUint32()
	//atomic.LoadUint64()
	//atomic.LoadPointer()
	//atomic.LoadUintptr()

	//原子性的保存值
	var i1 int32
	atomic.StoreInt32(&i1,10)
	fmt.Println(i1)

	//atomic.StoreInt64(&i1,10)
	//atomic.StoreUint32()
	//atomic.StoreUint64()
	//atomic.StoreUintptr()
	//atomic.StorePointer()

	fmt.Println("原子性的增机值")
	var i2 int64 = 10
	atomic.AddInt64(&i2,5)

	atomic.AddInt64(&i2,5)
	fmt.Println(i2)//20

	//atomic.AddInt32()
	//atomic.AddUint32()
	//atomic.AddUint64()
	//atomic.AddUintptr()

	fmt.Println("原子性的将新值加到原地址，并返回旧值")
	var i3 int64 = 10
	old := atomic.SwapInt64(&i3,20)
	fmt.Println(i3)//20
	fmt.Println(old)//10

	//atomic.SwapInt32()
	//atomic.SwapUint32()
	//atomic.SwapUint64()
	//atomic.SwapUintptr()
	//atomic.SwapPointer()
	//atomic.SwapPointer()

	fmt.Println("原子性的比较")
	var i4 int64 = 9
	//此值是否等于9，等于9 置换为10
	fmt.Println(atomic.CompareAndSwapInt64(&i4,9,10))
	fmt.Println(i4)
	//atomic.CompareAndSwapInt32()
	//atomic.CompareAndSwapUint32()
	//atomic.CompareAndSwapUint64()
	//atomic.CompareAndSwapPointer()
	//atomic.CompareAndSwapUintptr()




}
# go 语言编译器的 "//go:" 详解

在看go的源码的时候常常会遇到

```go
//go:linkname newFunc oldFunc

例如
// Implemented in runtime.
func runtime_registerPoolCleanup(cleanup func())
func runtime_procPin() int
func runtime_procUnpin()
```

这样的写法，只有方法名和参数，连方法体都没有，一开始很懵逼，自己尝试照着写，为啥人家的不报错，我的会报错呢？

然后各种资料查找研究，于是有了下面这篇文章

## 1、 //go:oninline

**`noinline` 顾名思义，不要内联。**

Inline 内联

是在编译期间发生的，将函数调用调用处替换为被调用函数主体的一种编译器优化手段。

## 1、优缺点

> 优点：
>
> 1. 减少函数调用的开销，提高程序执行速度
> 2. 复制后的更大函数为其他变异优化带来的可能性
> 3. 消除分支，函数执行是从一个函数的入口指针跳到另一个指针处，减少函数调用消耗
>
> 缺点：
>
> 1. 会增加代码空间的增长
> 2. 如果有很多重复的代码，会增加CPU负担，降低缓存命中率

对于短小而且工作较少的函数，使用内联是有效益的。

## 2、案例






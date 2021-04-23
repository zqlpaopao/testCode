# ==sync.Cond==

## 1、数据结构

```go
//每个cond有一个相关连的Locker（通常实际一个*Mutex或者*RWMutex），当改变条件或者点用wait方法时
//第一次使用后不可以被复制
// Each Cond has an associated Locker L (often a *Mutex or *RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
//
// A Cond must not be copied after first use.
type Cond struct {
        noCopy noCopy

        // L is held while observing or changing the condition
        L Locker

        notify  notifyList
        checker copyChecker
}
```

## 2、方法汇总

### 2.1 创建一个实例

NewCond 创建 Cond 实例时，需要关联一个锁。
```go
func NewCond(l Locker) *Cond
```

### 2.2 Signal 唤醒一个协程

Signal 只唤醒任意1个等待条件变量c的goroutine，无需锁保护
```go
// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Signal()
```

### 2.3 唤醒所有协程

```go
// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast()
```
Broadcast 唤醒所有等待条件变量c的goroutine，无需锁保护

### 2.4 Wait等待

```go
// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//    c.L.Lock()
//    for !condition() {
//        c.Wait()
//    }
//    ... make use of condition ...
//    c.L.Unlock()
//
func (c *Cond) Wait()
```

注释的意思是，当调用wait的时候，会自动释放c.L,并挂起调用者所在的goroutine，因此当前协程会阻塞在Wait的地方
，当收到Signal或者Broadcast的信号后，那么wait方法在结束阻塞的时候，会调用c.L加锁，并继续执行

还有注释，演示了使用的方法，使用for condition，而不是使用if condition，因为可能第一次没能满足，需要再一次等待下一次唤醒

## 3、基本使用

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	var f bool
	go func() {
		time.Sleep(time.Second * 5)
		cond.L.Lock()
		f = true
		cond.L.Unlock()

		cond.Signal()

	}()

	fmt.Println("waiting")

	cond.L.Lock()
	for !f {
		cond.Wait()
	}
	cond.L.Unlock()

	fmt.Println("done")
}
```



## 4、源码流程和实现原理

Wait原理及流程

![image-20210423172627906](readme.assets/image-20210423172627906.png)



Signal原理及流程

![image-20210423173853265](readme.assets/image-20210423173853265.png)



Broadcast

![image-20210423174409247](readme.assets/image-20210423174409247.png)





## 5、问题回顾

### cond.Wait 是阻塞的吗？是如何阻塞的？ 

是阻塞的。不过不是 sleep 这样阻塞的。

调用 `goparkunlock` 解除当前 goroutine 的 m 的绑定关系，将当前 goroutine 状态机切换为等待状态。等待后续 goready 函数时候能够恢复现场。

### cond.Signal 是如何通知一个等待的 goroutine ?

1. 判断是否有没有被唤醒的 goroutine，如果都已经唤醒了，直接就返回了
2. 将已通知 goroutine 的数量加1
3. 从等待唤醒的 goroutine 队列中，获取 head 指针指向的 goroutine，将其重新加入调度
4. 被阻塞的 goroutine 可以继续执行

### cond.Broadcast 是如何通知等待的 goroutine 的？

1. 判断是否有没有被唤醒的 goroutine，如果都已经唤醒了，直接就返回了
2. 将等待通知的 goroutine 数量和已经通知过的 goroutine 数量设置成相等
3. 遍历等待唤醒的 goroutine 队列，将所有的等待的 goroutine 都重新加入调度
4. 所有被阻塞的 goroutine 可以继续执行

### cond.Wait本身就是阻塞状态，为什么 cond.Wait 需要在循环内 ？

我们能注意到，调用 cond.Wait 的位置，使用的是 for 的方式来调用 wait 函数，而不是使用 if 语句。

这是由于 wait 函数被唤醒时，存在虚假唤醒等情况，导致唤醒后发现，条件依旧不成立。因此需要使用 for 语句来循环地进行等待，直到条件成立为止。



### 为什么不能 sync.Cond 不能复制 ？

sync.Cond 不能被复制的原因，并不是因为 sync.Cond 内部嵌套了 Locker。因为 NewCond 时传入的 Mutex/RWMutex 指针，对于 Mutex 指针复制是没有问题的。

主要原因是 sync.Cond 内部是维护着一个 notifyList。如果这个队列被复制的话，那么就在并发场景下导致不同 goroutine 之间操作的 notifyList.wait、notifyList.notify 并不是同一个，这会导致出现有些 goroutine 会一直堵塞。
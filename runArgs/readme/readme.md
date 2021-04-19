# 1、GOGC
GOGC 用于控制GC的处发频率， 其值默认为100,

意为直到自上次垃圾回收后heap size已经增长了100%时GC才触发运行。即是GOGC=100意味着live heap size 每增长一倍，GC触发运行一次。

如设定GOGC=200, 则live heap size 自上次垃圾回收后，增长2倍时，GC触发运行， 总之，其值越大则GC触发运行频率越低， 反之则越高，

如果GOGC=off 则关闭GC.

```go
/**
 * @Author: zHangSan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/4/16 下午4:44
 */

package main

import "runtime"

func main(){
	gc()
	runtime.GC()
	gc()
	gc()
	gc()

}

func gc(){
	for i := 0; i< 10000;i++{
		g := make(map[string]string)
		g = g
	}
}
````

```go
GOGC=100 GODEBUG=gctrace=1 go run ./main.go
gc 1 @0.017s 0%: 0.008+0.48+0.004 ms clock, 0.032+0.33/0.22/0.33+0.019 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 2 @0.031s 0%: 0.003+0.33+0.002 ms clock, 0.014+0.27/0.17/0.57+0.010 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 3 @0.057s 0%: 0.012+0.35+0.003 ms clock, 0.050+0/0.27/0.57+0.012 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
# command-line-arguments
gc 1 @0.002s 9%: 0.023+3.1+0.036 ms clock, 0.093+0.20/2.0/4.1+0.14 ms cpu, 4->7->6 MB, 5 MB goal, 4 P
gc 2 @0.015s 8%: 0.004+4.0+0.003 ms clock, 0.017+1.4/2.8/2.0+0.013 ms cpu, 10->13->13 MB, 12 MB goal, 4 P
gc 1 @0.000s 3%: 0.008+0.13+0.002 ms clock, 0.032+0/0.045/0.16+0.008 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
```
可以看到第二次的GC的时候是上次的两倍

# 2、GOTRACEBACK 打印异常信息

GOTRACEBACK用于控制当异常发生时，系统提供信息的详细程度， 在go 1.5， GOTRACEBACK有4个值。


GOTRACEBACK=0 只输出panic异常信息。
GOTRACEBACK=1 此为go的默认设置值， 输出所有goroutine的stack traces, 除去与go runtime相关的stack frames.
GOTRACEBACK=2 在GOTRACEBACK=1的基础上， 还输出与go runtime相关的stack frames,从而了解哪些goroutines是由go runtime启动运行的。
GOTRACEBACK=crash, 在GOTRACEBACK=2的基础上，go runtime处发进程segfault错误，从而生成core dump, 当然要操作系统允许的情况下， 而不是调用os.Exit。

```go
package main

import "runtime"

func main(){
	//goGC()

	gOTraceBack()

}
/*****************************GOTRACEBACK-异常发生打印详细信息****************************************/
func gOTraceBack(){
	panic("kerboom")

}
```
```go
GOTRACEBACK=2 go run main.go
```
```go
-> % GOTRACEBACK=2 go run main.go               
panic: kerboom

goroutine 1 [running]:
panic(0x105ff60, 0x10839b0)
        /usr/local/go/src/runtime/panic.go:1064 +0x46d fp=0xc000030768 sp=0xc0000306b0 pc=0x1028a2d
main.gOTraceBack(...)
        /Users/zhangsan/Documents/GitHub/testCode/runArgs/main.go:21
main.main()
        /Users/zhangsan/Documents/GitHub/testCode/runArgs/main.go:16 +0x39 fp=0xc000030788 sp=0xc000030768 pc=0x1056e89
runtime.main()
        /usr/local/go/src/runtime/proc.go:203 +0x212 fp=0xc0000307e0 sp=0xc000030788 pc=0x102b452
runtime.goexit()
        /usr/local/go/src/runtime/asm_amd64.s:1373 +0x1 fp=0xc0000307e8 sp=0xc0000307e0 pc=0x10529d1

goroutine 2 [force gc (idle)]:
runtime.gopark(0x1077928, 0x10cd3a0, 0x1411, 0x1)
        /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc000030fb0 sp=0xc000030f90 pc=0x102b820
runtime.goparkunlock(...)
        /usr/local/go/src/runtime/proc.go:310
runtime.forcegchelper()
        /usr/local/go/src/runtime/proc.go:253 +0xb7 fp=0xc000030fe0 sp=0xc000030fb0 pc=0x102b6d7
runtime.goexit()
        /usr/local/go/src/runtime/asm_amd64.s:1373 +0x1 fp=0xc000030fe8 sp=0xc000030fe0 pc=0x10529d1
created by runtime.init.6
        /usr/local/go/src/runtime/proc.go:242 +0x35

goroutine 3 [GC sweep wait]:
runtime.gopark(0x1077928, 0x10cd460, 0x140c, 0x1)
        /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc0000317a8 sp=0xc000031788 pc=0x102b820
runtime.goparkunlock(...)
        /usr/local/go/src/runtime/proc.go:310
runtime.bgsweep(0xc00001c070)
        /usr/local/go/src/runtime/mgcsweep.go:70 +0x9c fp=0xc0000317d8 sp=0xc0000317a8 pc=0x101aa7c
runtime.goexit()
        /usr/local/go/src/runtime/asm_amd64.s:1373 +0x1 fp=0xc0000317e0 sp=0xc0000317d8 pc=0x10529d1
created by runtime.gcenable
        /usr/local/go/src/runtime/mgc.go:214 +0x5c

goroutine 4 [GC scavenge wait]:
runtime.gopark(0x1077928, 0x10cd420, 0x140d, 0x1)
        /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc000031f78 sp=0xc000031f58 pc=0x102b820
runtime.goparkunlock(...)
        /usr/local/go/src/runtime/proc.go:310
runtime.bgscavenge(0xc00001c070)
        /usr/local/go/src/runtime/mgcscavenge.go:237 +0xd0 fp=0xc000031fd8 sp=0xc000031f78 pc=0x1018fd0
runtime.goexit()
        /usr/local/go/src/runtime/asm_amd64.s:1373 +0x1 fp=0xc000031fe0 sp=0xc000031fd8 pc=0x10529d1
created by runtime.gcenable
        /usr/local/go/src/runtime/mgc.go:215 +0x7e
exit status 2

```

GOTRACEBACK 在go 1.6中的变化


GOTRACEBACK=none 只输出panic异常信息。
GOTRACEBACK=single 只输出被认为引发panic异常的那个goroutine的相关信息。
GOTRACEBACK=all 输出所有goroutines的相关信息，除去与go runtime相关的stack frames.
GOTRACEBACK=system 输出所有goroutines的相关信息，包括与go runtime相关的stack frames,从而得知哪些goroutine是go runtime启动运行的。
GOTRACEBACK=crash 与go 1.5相同， 未变化。


为了与go 1.5兼容，0 对应 none, 1 对应 all, 以及 2 对应 system.


注意： 在go 1.6中， 默认，只输出引发panci异常的goroutine的stack trace.

# 3、控制CPU数量GOMAXPROCS
GOMAXPROCS


GOMAXPROCS 大家比较熟悉， 用于控制操作系统的线程数量， 这些线程用于运行go程序中的goroutines.
到go 1.5的时候， GOMAXPROCS的默认值就是我们的go程序启动时可见的操作系统认为的CPU个数。


注意： 在我们的go程序中使用的操作系统线程数量，也包括：正服务于cgo calls的线程, 阻塞于操作系统calls的线程，
所以go 程序中使用的操作系统线程数量可能大于GOMAXPROCS的值。

# 4、GODEBUG
GODEBUG=gctrace=1,schedtrace=1000 godoc -http=:8080

## 1、gctrace
gctrace用途主要是用于跟踪GC的不同阶段的耗时与GC前后的内存量对比。
```go
package main

import "runtime"

func main(){
	//goGC()

	//gOTraceBack() 查看crash的时候相信信息

	goDebugGc()

}
/*****************************GODEBUG-异常发生打印详细信息****************************************/
func goDebugGc(){
	gc()
}
func gc(){
	for i := 0; i< 10000;i++{
		g := make(map[string]string)
		g = g
	}
}
```

```go
GODEBUG=gctrace=1 go run  main.go                            
gc 1 @0.109s 0%: 0.023+0.68+0.005 ms clock, 0.094+0.66/0.44/0+0.021 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 2 @0.136s 0%: 0.004+0.37+0.003 ms clock, 0.019+0.29/0.16/0.68+0.012 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 3 @0.283s 0%: 0.020+0.47+0.003 ms clock, 0.081+0/0.31/0.66+0.014 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
# command-line-arguments
gc 1 @0.019s 6%: 0.057+8.1+0.086 ms clock, 0.23+1.5/5.0/10+0.34 ms cpu, 4->4->3 MB, 5 MB goal, 4 P
gc 4 @0.407s 0%: 0.010+1.0+0.059 ms clock, 0.043+0.34/0.18/0.60+0.23 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
# command-line-arguments
gc 1 @0.002s 12%: 0.007+3.8+0.038 ms clock, 0.029+0.17/2.8/1.8+0.15 ms cpu, 4->6->6 MB, 5 MB goal, 4 P
gc 2 @0.021s 9%: 0.005+8.8+0.003 ms clock, 0.020+0/8.0/1.5+0.014 ms cpu, 9->13->13 MB, 12 MB goal, 4 P
```

1. gc 1 : 代表第一次gc
2. @0.109 ： 这是gc的markTermination阶段完成，距离runtime到现在 STW的最后工作，关闭内存屏障，停止后台标记及辅助标记，做一些清理工作
3. 0% ： 当前为止，gc的标记工作（包括两次mark阶段和STW和并发标记）多占用的CPU时间占PU的百分比
4. 0.023+0.68+0.005 ms clock ： 三部分 0。023表示整个过程在mar阶段STW停顿时间（单个P的）；
0.68表示并发标记时间（所有P的）；0.005表示markTermination的STW时间（单个P）的

5. 0.094+0.66/0.44/0+0.021 ms cpu ： 0.094三部分，0.094表示整个进程在mark阶段STW停顿时间（0.023 * 8）;
0.66/0.44/0 三块信息，0.66是mutator assists占用的时间，0。44是dedicated mark worker+ fractional mark worker占用的时间，0是idle mark workers占用的时间，这些时间接近
0.68 * P的个数，+0。021ms 表示整个进程在markTermination阶段Ste停顿时间0。005 * P的个数

6. 4->4->0 MB 4表示开始mark阶段前的heap_live大小；4表示markTermination阶段之前的heap_live大小，0表示被标记对象的大小

7. 5MB goal： 表示下一次触发GC的内存占用阀值是5MB，向上取整

8. 4P ：本次gc共有多少P

- heap_live要结合go的内存管理来理解。因为go按照不同的对象大小，会分配不同页数的span。span是对内存页进行管理的基本单元，每页8k大小。所以肯定会出现span中有内存是空闲着没被用上的。

不过怎么用go先不管，反正是把它划分给程序用了。而heap_live就表示所有span的大小。

而程序到底用了多少呢？就是在GC扫描对象时，扫描到的存活对象大小就是已用的大小。对应上面就是8MB。

- mark worker分为三种，dedicated、fractional和idle。分别表示标记工作干活时的专注程度。dedicated最专注，除非被抢占打断，否则一直干活。idle最偷懒，干一点活就退出，控制权让给出别的goroutine。它们都是并发标记工作里的worker。

## 2、schedtrace 调度跟踪
```go
GODEBUG=schedtrace=1000  go run  main.go 
SCHED 0ms: gomaxprocs=4 idleprocs=2 threads=4 spinningthreads=1 idlethreads=0 runqueue=0 [0 0 0 0]
# command-line-arguments
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=4 spinningthreads=1 idlethreads=0 runqueue=0 [1 0 0 0]
# command-line-arguments
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=5 spinningthreads=1 idlethreads=0 runqueue=0 [0 0 0 0]
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=4 spinningthreads=1 idlethreads=0 runqueue=0 [1 0 0 0]
```

可以使运行时在每 X 毫秒发出一次详细的多行信息，信息内容主要包括调度程序、处理器、OS 线程 和 Goroutine 的状态。

sched：每一行都代表调度器的调试信息，后面提示的毫秒数表示启动到现在的运行时间，输出的时间间隔受 schedtrace 的值影响。

gomaxprocs：当前的 CPU 核心数（GOMAXPROCS 的当前值）。

idleprocs：空闲的处理器数量，后面的数字表示当前的空闲数量。

threads：OS 线程数量，后面的数字表示当前正在运行的线程数量。

spinningthreads：自旋状态的 OS 线程数量。

idlethreads：空闲的线程数量。

runqueue：全局队列中中的 Goroutine 数量，而后面的 [0 0 1 1] 则分别代表这 4 个 P 的本地队列正在运行的 Goroutine 数量。




详细信息
```go
GODEBUG=scheddetail=1,schedtrace=1000  go run  main.go

SCHED 0ms: gomaxprocs=4 idleprocs=2 threads=4 spinningthreads=1 idlethreads=0 runqueue=0 gcwaiting=0 nmidlelocked=0 stopwait=0 sysmonwait=0
  P0: status=1 schedtick=0 syscalltick=0 m=0 runqsize=0 gfreecnt=0 timerslen=0
  P1: status=1 schedtick=0 syscalltick=0 m=3 runqsize=0 gfreecnt=0 timerslen=0
  P2: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0 timerslen=0
  P3: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0 timerslen=0
  M3: p=1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
  M2: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
  M1: p=-1 curg=17 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=false lockedg=17
  M0: p=0 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=1
  G1: status=1() m=-1 lockedm=0
  G17: status=6() m=1 lockedm=1
  G2: status=1() m=-1 lockedm=-1
# command-line-arguments
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=4 spinningthreads=1 idlethreads=0 runqueue=0 [1 0 0 0]
# command-line-arguments
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=5 spinningthreads=1 idlethreads=0 runqueue=0 gcwaiting=0 nmidlelocked=1 stopwait=0 sysmonwait=0
  P0: status=1 schedtick=0 syscalltick=0 m=3 runqsize=0 gfreecnt=0 timerslen=0
  P1: status=1 schedtick=1 syscalltick=0 m=2 runqsize=0 gfreecnt=0 timerslen=0
  P2: status=1 schedtick=0 syscalltick=0 m=4 runqsize=0 gfreecnt=0 timerslen=0
  P3: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0 timerslen=0
  M4: p=2 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
  M3: p=0 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
  M2: p=1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=true blocked=false lockedg=-1
  M1: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
  M0: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=1
  G1: status=1(chan receive) m=-1 lockedm=0
  G2: status=4(force gc (idle)) m=-1 lockedm=-1
  G3: status=1() m=-1 lockedm=-1
  G4: status=4(GC scavenge wait) m=-1 lockedm=-1
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=5 spinningthreads=1 idlethreads=0 runqueue=0 gcwaiting=0 nmidlelocked=1 stopwait=0 sysmonwait=0
  P0: status=1 schedtick=0 syscalltick=0 m=3 runqsize=0 gfreecnt=0 timerslen=0
  P1: status=1 schedtick=2 syscalltick=0 m=2 runqsize=0 gfreecnt=0 timerslen=0
  P2: status=1 schedtick=0 syscalltick=0 m=4 runqsize=0 gfreecnt=0 timerslen=0
  P3: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0 timerslen=0
  M4: p=2 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=true blocked=false lockedg=-1
  M3: p=0 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=true blocked=false lockedg=-1
  M2: p=1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
  M1: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
  M0: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=1
  G1: status=1(chan receive) m=-1 lockedm=0
  G2: status=4(force gc (idle)) m=-1 lockedm=-1
  G3: status=4(GC sweep wait) m=-1 lockedm=-1
  G4: status=4(GC scavenge wait) m=-1 lockedm=-1
```
status：G 的运行状态。

m：隶属哪一个 M。

lockedm：是否有锁定 M。

在第一点中我们有提到 G 的运行状态，这对于分析内部流转非常的有用，共涉及如下 9 种状态：

状态	值	含义
_Gidle	0	刚刚被分配，还没有进行初始化。
_Grunnable	1	已经在运行队列中，还没有执行用户代码。
_Grunning	2	不在运行队列里中，已经可以执行用户代码，此时已经分配了 M 和 P。
_Gsyscall	3	正在执行系统调用，此时分配了 M。
_Gwaiting	4	在运行时被阻止，没有执行用户代码，也不在运行队列中，此时它正在某处阻塞等待中。
_Gmoribund_unused	5	尚未使用，但是在 gdb 中进行了硬编码。
_Gdead	6	尚未使用，这个状态可能是刚退出或是刚被初始化，此时它并没有执行用户代码，有可能有也有可能没有分配堆栈。
_Genqueue_unused	7	尚未使用。
_Gcopystack	8	正在复制堆栈，并没有执行用户代码，也不在运行队列中。
在理解了各类的状态的意思后，我们结合上述案例看看，如下：

G1: status=4(semacquire) m=-1 lockedm=-1
G2: status=4(force gc (idle)) m=-1 lockedm=-1
G3: status=4(GC sweep wait) m=-1 lockedm=-1
G17: status=1() m=-1 lockedm=-1
G18: status=2() m=4 lockedm=-1
在这个片段中，G1 的运行状态为 _Gwaiting，并没有分配 M 和锁定。这时候你可能好奇在片段中括号里的是什么东西呢，其实是因为该 status=4 是表示 Goroutine 在运行时时被阻止，而阻止它的事件就是 semacquire 事件，是因为 semacquire 会检查信号量的情况，在合适的时机就调用 goparkunlock 函数，把当前 Goroutine 放进等待队列，并把它设为 _Gwaiting 状态。

那么在实际运行中还有什么原因会导致这种现象呢，我们一起看看，如下：

    waitReasonZero                                    // ""
    waitReasonGCAssistMarking                         // "GC assist marking"
    waitReasonIOWait                                  // "IO wait"
    waitReasonChanReceiveNilChan                      // "chan receive (nil chan)"
    waitReasonChanSendNilChan                         // "chan send (nil chan)"
    waitReasonDumpingHeap                             // "dumping heap"
    waitReasonGarbageCollection                       // "garbage collection"
    waitReasonGarbageCollectionScan                   // "garbage collection scan"
    waitReasonPanicWait                               // "panicwait"
    waitReasonSelect                                  // "select"
    waitReasonSelectNoCases                           // "select (no cases)"
    waitReasonGCAssistWait                            // "GC assist wait"
    waitReasonGCSweepWait                             // "GC sweep wait"
    waitReasonChanReceive                             // "chan receive"
    waitReasonChanSend                                // "chan send"
    waitReasonFinalizerWait                           // "finalizer wait"
    waitReasonForceGGIdle                             // "force gc (idle)"
    waitReasonSemacquire                              // "semacquire"
    waitReasonSleep                                   // "sleep"
    waitReasonSyncCondWait                            // "sync.Cond.Wait"
    waitReasonTimerGoroutineIdle                      // "timer goroutine (idle)"
    waitReasonTraceReaderBlocked                      // "trace reader (blocked)"
    waitReasonWaitForGCCycle                          // "wait for GC cycle"
    waitReasonGCWorkerIdle                            // "GC worker (idle)"
我们通过以上 waitReason 可以了解到 Goroutine 会被暂停运行的原因要素，也就是会出现在括号中的事件。

M
p：隶属哪一个 P。

curg：当前正在使用哪个 G。

runqsize：运行队列中的 G 数量。

gfreecnt：可用的G（状态为 Gdead）。

mallocing：是否正在分配内存。

throwing：是否抛出异常。

preemptoff：不等于空字符串的话，保持 curg 在这个 m 上运行。

P
status：P 的运行状态。

schedtick：P 的调度次数。

syscalltick：P 的系统调用次数。

m：隶属哪一个 M。

runqsize：运行队列中的 G 数量。

gfreecnt：可用的G（状态为 Gdead）。

状态	值	含义
_Pidle	0	刚刚被分配，还没有进行进行初始化。
_Prunning	1	当 M 与 P 绑定调用 acquirep 时，P 的状态会改变为 _Prunning。
_Psyscall	2	正在执行系统调用。
_Pgcstop	3	暂停运行，此时系统正在进行 GC，直至 GC 结束后才会转变到下一个状态阶段。
_Pdead	4	废弃，不再使用。


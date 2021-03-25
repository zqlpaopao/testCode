# 1、协程泄漏问题
如果你启动了一个goroutine，但并没有按照预期的一样退出，直到程序结束，此goroutine才结束，这种情况就是 goroutine 泄露。当 goroutine 泄露发生时，该 goroutine 的栈一直被占用而不能释放，goroutine 里的函数在堆上申请的空间也不能被垃圾回收器回收。这样，在程序运行期间，内存占用持续升高，可用内存越来也少，最终将导致系统崩溃。

# 2、channel阻塞引起泄漏
```go
/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/3/25 下午1:37
 */

package main

import (
	"fmt"
	"runtime"
	"time"
)

func Go1() {
	c := make(chan int)
	go func() {
		<-c
	}()
	time.Sleep(2*time.Second)
	fmt.Println("G1 exit")
}

func main() {
	go Go1()

	c := time.Tick(time.Second)
	for range c {
		fmt.Printf("goroutine [nums]: %d\n", runtime.NumGoroutine())
	}
}
```

结果
```go
go run main.go                         
goroutine [nums]: 3
goroutine [nums]: 3
G1 exit
goroutine [nums]: 2
goroutine [nums]: 2
goroutine [nums]: 2
goroutine [nums]: 2

```

# 3、检测工具goleak
```go
go get -u go.uber.org/goleak
```

编写test用例
```go
/**
 * @Author: zhangsan
 * @Description:
 * @File:  gos_test.go
 * @Version: 1.0.0
 * @Date: 2021/3/25 下午2:12
 */

package src

import (
	"testing"
	"go.uber.org/goleak"
)

func TestGo1(t *testing.T) {
	defer goleak.VerifyNone(t)
	Go1()
}
```

```go
time go test -v
```

```go
=== RUN   TestGo1
G1 exit
    TestGo1: leaks.go:78: found unexpected goroutines:
        [Goroutine 19 in state chan receive, with test/goroutine-leak/src.Go1.func1 on top of the stack:
        goroutine 19 [chan receive]:
        test/goroutine-leak/src.Go1.func1(0xc0001022a0)
                /Users/zhangsan/Documents/GitHub/testCode/goroutine-leak/src/gos.go:19 +0x34
        created by test/goroutine-leak/src.Go1
                /Users/zhangsan/Documents/GitHub/testCode/goroutine-leak/src/gos.go:18 +0x52
        ]
--- FAIL: TestGo1 (2.45s)
FAIL
exit status 1
FAIL    test/goroutine-leak/src 2.458s
go test -v  0.44s user 0.30s system 24% cpu 2.987 total
```
goroutine19 chan receive 失败


# 4、gops
```go
go get -u github.com/google/gops
```

程序
```go
package main

import (
	"log"
	"runtime"
	"time"

	"github.com/google/gops/agent"
)

func main() {
	if err := agent.Listen(agent.Options{
		Addr:            "0.0.0.0:8848",
		// ConfigDir:       "/home/centos/gopsconfig", // 最好使用默认
		ShutdownCleanup: true}); err != nil {
		log.Fatal(err)
	}

    // 测试代码
	_ = make([]int, 1000, 1000)
	runtime.GC()

	_ = make([]int, 1000, 2000)
	runtime.GC()

	time.Sleep(time.Hour)
}
```
agent Option选项
agent有3个配置：

Addr：agent要监听的ip和端口，默认ip为环回地址，端口随机分配。
ConfigDir：该目录存放的不是agent的配置，而是每一个使用了agent的go进程信息，文件以pid命名，内容是该pid进程所监听的端口号，所以其中文件的目的是形成pid到端口的映射。默认值为~/.config/gops
ShutdownCleanup：进程退出时，是否清理ConfigDir中的文件，默认值为false，不清理
通常可以把Addr设置为要监听的IP，把ShutdownCleanup设置为ture，进程退出后，残留在ConfigDir目录的文件不再有用，最好清除掉。

# gops原理
gops的原理是，代码中导入gops/agent，建立agent服务，gops命令连接agent读取进程信息。

gops

agent的实现原理可以查看agent/handle函数。

使用go标准库中原生接口实现相关功能，如同你要在自己的程序中开启pprof类似，只不过这部分功能由gops/agent实现了：

使用runtime.MemStats获取内存情况
使用runtime/pprof获取调用栈、cpu profile和memory profile
使用runtime/trace获取trace
使用runtime获取stats信息
使用runtime/debug、GC设置和启动GC
再谈ConfigDir。从源码上看，ConfigDir对agent并没有用途，对gops有用。当gops和ConfigDir在一台机器上时，即gops查看本机的go进程信息，gops可以通过其中的文件，快速找到agent服务的端口。能够实现：gops <sub-cmd> pid到gops <sub-cmd> 127.0.0.1:port的转换。

如果代码中通过ConfigDir指定了其他目录，使用gops时，需要添加环境变量GOPS_CONFIG_DIR指向ConfigDir使用的目录。

## 1、读取进程信息
```go
gops            
79    1     com.docker.vmnetd              go1.13.10 /Library/PrivilegedHelperTools/com.docker.vmnetd
936   924   vpnkit-bridge                  go1.13.10 /Applications/Docker.app/Contents/MacOS/vpnkit-bridge
79708 19650 gops                           go1.14.3  /Users/zhangsan/go/src/bin/gops
925   907   com.docker.backend             go1.13.10 /Applications/Docker.app/Contents/MacOS/com.docker.backend
940   924   docker-mutagen                 go1.13    /Applications/Docker.app/Contents/Resources/cli-plugins/docker-mutagen
79092 2663  go                             go1.14.3  /usr/local/go/bin/go
937   924   com.docker.driver.amd64-linux  go1.13.10 /Applications/Docker.app/Contents/MacOS/com.docker.driver.amd64-linux
79106 79092 main                         * go1.14.3  /private/var/folders/69/m2hx0d9532q3cnkt80_ndtxh0000gn/T/go-build288894134/b001/exe/main
924   907   com.docker.supervisor          go1.13.10 /Applications/Docker.app/Contents/MacOS/com.docker.supervisor
```

<font color=res size=5x>**带*的是程序中使用了gops/agent，不带*的是普通的go程序。**</font>

## 2、进程树
```go
gops tree       
...
├── 924
│   ├── 936 (vpnkit-bridge) {go1.13.10}
│   ├── 940 (docker-mutagen) {go1.13}
│   └── 937 (com.docker.driver.amd64-linux) {go1.13.10}
├── 2663
│   └── 79092 (go) {go1.14.3}
│       └── [*]  79106 (main) {go1.14.3}
├── 907
│   └── 925 (com.docker.backend) {go1.13.10}
├── 1
│   └── 79 (com.docker.vmnetd) {go1.13.10}
└── 19650
    └── 79862 (gops) {go1.14.3}

```

## 3、进程概要信息
```go
gops 79092   
parent PID:     2663
threads:        13
memory usage:   0.132%
cpu usage:      0.128%
username:       zhangsan
cmd+args:       go run main.go
elapsed time:   05:13
```

## 4、当前调用栈信息
```go
gops stack 79106    
goroutine 19 [running]:
runtime/pprof.writeGoroutineStacks(0x1190a00, 0xc000010008, 0x0, 0x0)
        /usr/local/go/src/runtime/pprof/pprof.go:665 +0x9d
runtime/pprof.writeGoroutine(0x1190a00, 0xc000010008, 0x2, 0x0, 0x0)
        /usr/local/go/src/runtime/pprof/pprof.go:654 +0x44
runtime/pprof.(*Profile).WriteTo(0x1273d20, 0x1190a00, 0xc000010008, 0x2, 0xc000010008, 0x0)
        /usr/local/go/src/runtime/pprof/pprof.go:329 +0x3da
github.com/google/gops/agent.handle(0x17b5840, 0xc000010008, 0xc000196000, 0x1, 0x1, 0x0, 0x0)
        /Users/zhangsan/Documents/GitHub/testCode/vendor/github.com/google/gops/agent/agent.go:201 +0x1af
github.com/google/gops/agent.listen()
        /Users/zhangsan/Documents/GitHub/testCode/vendor/github.com/google/gops/agent/agent.go:145 +0x2bf
created by github.com/google/gops/agent.Listen
        /Users/zhangsan/Documents/GitHub/testCode/vendor/github.com/google/gops/agent/agent.go:123 +0x36a

goroutine 1 [sleep, 6 minutes]:
time.Sleep(0x34630b8a000)
        /usr/local/go/src/runtime/time.go:188 +0xba
main.main()
        /Users/zhangsan/Documents/GitHub/testCode/goroutine-leak/main.go:26 +0xad

goroutine 6 [syscall, 6 minutes]:
os/signal.signal_recv(0x0)
        /usr/local/go/src/runtime/sigqueue.go:144 +0x96
os/signal.loop()
        /usr/local/go/src/os/signal/signal_unix.go:23 +0x22
created by os/signal.Notify.func1
        /usr/local/go/src/os/signal/signal.go:127 +0x44

goroutine 18 [chan receive, 6 minutes]:
github.com/google/gops/agent.gracefulShutdown.func1(0xc00007e120)
        /Users/zhangsan/Documents/GitHub/testCode/vendor/github.com/google/gops/agent/agent.go:158 +0x41
created by github.com/google/gops/agent.gracefulShutdown
        /Users/zhangsan/Documents/GitHub/testCode/vendor/github.com/google/gops/agent/agent.go:156 +0xd7
```

## 5、内存使用情况memstats
```go
 gops memstats 79106
alloc: 100.68KB (103096 bytes) // 当前分配出去未收回的内存总量
total-alloc: 2.11MB (2217136 bytes) // 已分配出去的内存总量
sys: 69.83MB (73220352 bytes) //当前进程从OS获取的内存总量
lookups: 0
mallocs: 416 // 分配的对象数量
frees: 83 82 // 释放的对象数量
heap-alloc: 100.68KB (103096 bytes) // 当前分配出去未收回的堆内存总量
heap-sys: 63.59MB (66682880 bytes) // 当前堆从OS获取的内存
heap-idle: 63.15MB (66215936 bytes)// 当前堆中空闲的内存量
heap-in-use: 456.00KB (466944 bytes)// 当前堆使用中的内存量
heap-released: 62.04MB (65052672 bytes)
heap-objects: 333 // 堆中对象数量
stack-in-use: 416.00KB (425984 bytes)// 栈使用中的内存量 
stack-sys: 416.00KB (425984 bytes)
stack-mspan-inuse: 36.12KB (36992 bytes)
stack-mspan-sys: 48.00KB (49152 bytes)
stack-mcache-inuse: 6.78KB (6944 bytes)
stack-mcache-sys: 16.00KB (16384 bytes)
other-sys: 1004.42KB (1028530 bytes)
gc-sys: 3.41MB (3574024 bytes)
next-gc: when heap-alloc >= 4.00MB (4194304 bytes)// 下次GC的条件
last-gc: 2021-03-25 15:09:34.705688 +0800 CST// 上次GC的时间
gc-pause-total: 144.926µs// GC总暂停时间
gc-pause: 42877// 上次GC暂停时间，单位纳秒
gc-pause-end: 1616656174705688000
num-gc: 5// 已进行的GC次数
num-forced-gc: 2
gc-cpu-fraction: 1.0046411080700944e-06
enable-gc: true // 是否开始GC
debug-gc: false
```
## 6、运行时信息 system信息
```go
gops stats 79106 
goroutines: 4
OS threads: 10
GOMAXPROCS: 4
num CPU: 4
```

## 6、运行当前5s的trace
```go
gops trace 79106
Tracing now, will take 5 secs...
Trace dump saved to: /var/folders/69/m2hx0d9532q3cnkt80_ndtxh0000gn/T/trace968769714
2021/03/25 15:16:16 Parsing trace...
2021/03/25 15:16:16 Splitting trace...
2021/03/25 15:16:16 Opening browser. Trace viewer is listening on http://127.0.0.1:53704
```

## 7、cpu信息
和原生的差不多
```go
Profiling CPU now, will take 30 secs...
Profile dump saved to: /var/folders/69/m2hx0d9532q3cnkt80_ndtxh0000gn/T/cpu_profile353858022
Binary file saved to: /var/folders/69/m2hx0d9532q3cnkt80_ndtxh0000gn/T/binary686968845
File: binary686968845
Type: cpu
Time: Mar 25, 2021 at 3:18pm (CST)
Duration: 30s, Total samples = 0 
No samples were found with the default sample value type.
Try "sample_index" command to analyze different sample values.
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
top
Showing nodes accounting for 0, 0% of 0 total
flat  flat%   sum%        cum   cum%
```

## 8、堆信息
```go
gops pprof-heap 79106
Profile dump saved to: /var/folders/69/m2hx0d9532q3cnkt80_ndtxh0000gn/T/heap_profile961243197
Binary file saved to: /var/folders/69/m2hx0d9532q3cnkt80_ndtxh0000gn/T/binary963711096
File: binary963711096
Type: inuse_space
Time: Mar 25, 2021 at 3:19pm (CST)
No samples were found with the default sample value type.
Try "sample_index" command to analyze different sample values.
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 0, 0% of 0 total
flat  flat%   sum%        cum   cum%
(pprof) traces
File: binary963711096
Type: inuse_space
Time: Mar 25, 2021 at 3:19pm (CST)
-----------+-------------------------------------------------------
bytes:  1MB
0   runtime/pprof.writeGoroutineStacks
runtime/pprof.writeGoroutine
runtime/pprof.(*Profile).WriteTo
github.com/google/gops/agent.handle
github.com/google/gops/agent.listen
-----------+-------------------------------------------------------
(pprof)

```
## 9、远程端口
```go
gops pprof-heap 127.0.0.1:8848
Profile dump saved to: /var/folders/69/m2hx0d9532q3cnkt80_ndtxh0000gn/T/heap_profile200987266
Binary file saved to: /var/folders/69/m2hx0d9532q3cnkt80_ndtxh0000gn/T/binary663790073
File: binary663790073
Type: inuse_space
Time: Mar 25, 2021 at 3:21pm (CST)
No samples were found with the default sample value type.
Try "sample_index" command to analyze different sample values.
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) 
```





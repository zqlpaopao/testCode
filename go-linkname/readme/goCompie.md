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

### 1、优缺点

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

### 2、案例

#### 使用内联

```go
package main


func appStr(w string)string{
	return w + " word"
}

func main(){
	appStr("hello")
}
```



编译

```go
GOOS=linux GOARCH=386 go tool compile -S main.go > main.S
                           
.
├── main.S
├── main.go
└── main.o

0 directories, 3 files
```



截取部分信息

```go
0x001c 00028 (main.go:13)	MOVL	"".w+32(SP), AX
0x0020 00032 (main.go:13)	PCDATA	$0, $0
0x0020 00032 (main.go:13)	MOVL	AX, 4(SP)
0x0024 00036 (main.go:13)	PCDATA	$1, $1
0x0024 00036 (main.go:13)	MOVL	"".w+36(SP), AX
0x0028 00040 (main.go:13)	MOVL	AX, 8(SP)
0x002c 00044 (main.go:13)	PCDATA	$0, $1
0x002c 00044 (main.go:13)	LEAL	go.string." word"(SB), AX
0x0032 00050 (main.go:13)	PCDATA	$0, $0
0x0032 00050 (main.go:13)	MOVL	AX, 12(SP)
0x0036 00054 (main.go:13)	MOVL	$5, 16(SP)
0x003e 00062 (main.go:13)	CALL	runtime.concatstring2(SB)
```

可以看出，并没有直接调用addStr函数，而是直接调用的go.string." word"，此时就是内联了

#### 没有内联

```go
package main

//go:noinline
func appStr(w string)string{
	return w + " word"
}

func main(){
	appStr("hello")
}
```

编译

```go
 GOOS=linux GOARCH=386 go tool compile -S main.go > main.S

0x0015 00021 (main.go:17)	LEAL	go.string."hello"(SB), AX
0x001b 00027 (main.go:17)	PCDATA	$0, $0
0x001b 00027 (main.go:17)	MOVL	AX, (SP)
0x001e 00030 (main.go:17)	MOVL	$5, 4(SP)
0x0026 00038 (main.go:17)	CALL	"".appStr(SB)
0x002b 00043 (main.go:18)	ADDL	$16, SP
0x002e 00046 (main.go:18)	RET
0x002f 00047 (main.go:18)	NOP
0x002f 00047 (main.go:16)	PCDATA	$1, $-1
0x002f 00047 (main.go:16)	PCDATA	$0, $-2
0x002f 00047 (main.go:16)	CALL	runtime.morestack_noctxt(SB)
```

可以看到调用了0x0026 00038 (main.go:17)	CALL	"".appStr(SB)函数来执行



## 2、go:nosplit-堆溢出检测

nosplite的作用是：==跳过栈溢出检测==

### 栈溢出

> goroutine的初始化是2k，当不够用的时候会增长占空间，借用G0的栈空间，那么不能让其无限增长，所以就要有一个机制去检测
>
> golang默认是开启栈溢出检测的
>
> 不执行栈溢出检测，可以提高性能，同时也可能发生stack overflow导致编译失败



### 用法

```go
//go:nosplit
func sayHi()string{
	return "k"
}
```



## 3、go:noescape-禁止逃逸

`noescape` 的作用是：<font color=red size=5x>**禁止逃逸，而且它必须指示一个只有声明没有主体的函数。**</font>



### 逃逸

> 自动的将超出自己生命周期的变量从函数栈转移到堆中
>
> 优势
>
> ​	GC压力变小，栈空间会在结束的时候进行回收
>
> 劣势
>
> ​	造成不可控因素



## 4、go：norace

`norace` 的作用是：**跳过竞态检测**

```go
go run -race main.go 利用 -race 来使编译器报告数据竞争问题。
```

> 使用 `norace` 除了减少编译时间，我想不到有其他的优点了。但缺点却很明显，那就是数据竞争会导致程序的不确定性。



## 5、go:linkname 转译实现函数

要在func的方法下添加会变文件标示

在看源码的时候看到

```go
// Implemented in runtime.
func runtime_registerPoolCleanup(cleanup func())
func runtime_procPin() int
func runtime_procUnpin()
```

只有方法名、参数和返回参数，没有方法体，于是自己也写了一个，既然报错，为什么呢？



研究后发现，这是golang的一种外部实现写法，针对函数和变量（公有私有都有效）常量无效的写法

话不多说伤上代码

目录结构

```go.
├── a
│   ├── a.go
│   └── a.s
├── b
│   └── b.go
├── main.go
```

a.go

```go
package a

import (
	_ "unsafe"
	_"mysql/go-linkname/b"
)

//变量是可以的
var Aa = ""

//常来是不可以的
const Ab = ""

//只能通过其他正常函数调用才生效，直接调用没反应
//即使是小写的也可以在外部在重新实现
func hello()string

func Greet() string {
	return hello()
}
```



b.go

```go
package b

import (
	_ "unsafe"
)

//go:linkname sayHi  mysql/go-linkname/a.hello
//go:nosplit
func sayHi()string{
	return "k"
}


//go:linkname aaaStr mysql/go-linkname/a.Aa
var aaaStr = "111"

//go:linkname aBStr mysql/go-linkname/a.Ab
var aBStr = "11s1"
```



Main.go

```go
package main

import (
	"fmt"
	"test/go-linkname/a"
)

func main() {
	fmt.Println(a.Greet())


	fmt.Println(a.Aa)
	fmt.Println(a.Ab)
}
```

输出

```go
k
111


```

变量和函数是可以的，常量不可以，但是直接调用函数是不可以的



## 6、go:generate

当运行go generate时，它将扫描与当前包相关的源代码文件，找出所有包含"//go:generate"的特殊注释，提取并执行该特殊注释后面的命令，命令为可执行程序，形同shell下面执行。

```go
package main

import (
	"fmt"
	"test/go-linkname/a"
)
//go:generate  echo hello word
func main() {
	fmt.Println(a.Greet())


	fmt.Println(a.Aa)
	fmt.Println(a.Ab)
	a.Adds()
}
//go:generate  go version
func add(){

}
```

```go
 go generate
hello word
go version go1.14.3 darwin/amd64
```


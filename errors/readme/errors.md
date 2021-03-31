# 1、error 打印调用堆栈信息

当前使用的版本是

```go
go version    
go version go1.14.3 darwin/amd64
```



```go
package main

import (
	"github.com/pkg/errors"
	"fmt"
	//"errors"
)

func f() error {
	return errors.New("is f")
}

func errorWith()error{
	return  f()
	//return errors.WithStack(errors.New("ggg"))
}

func main() {
	// 详细输出（打印出调用堆栈）
	//fmt.Printf("%+v\n", f())
	fmt.Printf("%+v\n", errorWith())

}
```



```go
go run main.go
is f
main.f
        /Users/zhangsan/Documents/GitHub/testCode/errors/main.go:10
main.errorWith
        /Users/zhangsan/Documents/GitHub/testCode/errors/main.go:14
main.main
        /Users/zhangsan/Documents/GitHub/testCode/errors/main.go:21
runtime.main
        /usr/local/go/src/runtime/proc.go:203
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1373
```

在开发的时候就不用一层层的log.Write了，只需要记录一个原始的错误信息，就可以找到调用链

来看下源码

```go
// New使用提供的消息返回一个错误。
// New还记录调用时的堆栈跟踪。
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// 基本错误是有消息和堆栈，但没有调用者的错误
type fundamental struct {
	msg string
	*stack
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func callers(skip int, pcbuf []uintptr) int {
	sp := getcallersp()
	pc := getcallerpc()
	gp := getg()
	var n int
	systemstack(func() {
		n = gentraceback(pc, sp, 0, gp, skip, &pcbuf[0], len(pcbuf), nil, nil, 0)
	})
	return n
}

```

获取寄存器sp 

程序计数器pc

运行的goroutine信息 gp

# 2、error嵌套




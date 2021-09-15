当前有这样的代码
```go
package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "math/rand"
    "strconv"
    "strings"
)

func foo(n int) string {
    var buf bytes.Buffer
    for i := 0; i < 100000; i++ {
        buf.WriteString(strconv.Itoa(n))
    }
    sum := sha256.Sum256(buf.Bytes())

    var b []byte
    for i := 0; i < int(sum[0]); i++ {
        x := sum[(i*7+1)%len(sum)] ^ sum[(i*5+3)%len(sum)]
        c := strings.Repeat("abcdefghijklmnopqrstuvwxyz", 10)[x]
        b = append(b, c)
    }
    return string(b)
}

func main() {
    // ensure function output is accurate
    if foo(12345) == "aajmtxaattdzsxnukawxwhmfotnm" {
        fmt.Println("Test PASS")
    } else {
        fmt.Println("Test FAIL")
    }

    for i := 0; i < 100; i++ {
        foo(rand.Int())
    }
}
```

我们要加下CPU分析的代码，改动如下
```go
package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "math/rand"
    "os"
    "runtime/pprof"
    "strconv"
    "strings"
)

func foo(n int) string {
    var buf bytes.Buffer
    for i := 0; i < 100000; i++ {
        buf.WriteString(strconv.Itoa(n))
    }
    sum := sha256.Sum256(buf.Bytes())

    var b []byte
    for i := 0; i < int(sum[0]); i++ {
        x := sum[(i*7+1)%len(sum)] ^ sum[(i*5+3)%len(sum)]
        c := strings.Repeat("abcdefghijklmnopqrstuvwxyz", 10)[x]
        b = append(b, c)
    }
    return string(b)
}

func main() {
    cpufile, err := os.Create("cpu.pprof")
    if err != nil {
        panic(err)
    }
    err = pprof.StartCPUProfile(cpufile)
    if err != nil {
        panic(err)
    }
    defer cpufile.Close()
    defer pprof.StopCPUProfile()

    // ensure function output is accurate
    if foo(12345) == "aajmtxaattdzsxnukawxwhmfotnm" {
        fmt.Println("Test PASS")
    } else {
        fmt.Println("Test FAIL")
    }

    for i := 0; i < 100; i++ {
        foo(rand.Int())
    }
}
```

编译并运行此程序后，配置文件将会写入./cpu.pprof. 我们可以使用go tool pprof以下命令读取此文件：

$ go tool pprof cpu.pprof
我们现在在 pprof 交互工具中。我们可以通过top10( top1, top2, top99, ...,topn也可以工作)看到我们的程序大部分时间都在做什么。top10展示如下：

```go
 go tool pprof cpu.pprof
Type: cpu
Time: Sep 15, 2021 at 9:46am (CST)
Duration: 1.30s, Total samples = 1.14s (87.40%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1070ms, 93.86% of 1140ms total
Showing top 10 nodes out of 64
      flat  flat%   sum%        cum   cum%
     520ms 45.61% 45.61%      520ms 45.61%  crypto/sha256.block
     150ms 13.16% 58.77%      350ms 30.70%  strconv.formatBits
     120ms 10.53% 69.30%      180ms 15.79%  runtime.mallocgc
      70ms  6.14% 75.44%       70ms  6.14%  runtime.memmove
      70ms  6.14% 81.58%       70ms  6.14%  runtime.pthread_cond_wait
      60ms  5.26% 86.84%       60ms  5.26%  runtime.kevent
      20ms  1.75% 88.60%       20ms  1.75%  runtime.memclrNoHeapPointers
      20ms  1.75% 90.35%       20ms  1.75%  runtime.nextFreeFast
      20ms  1.75% 92.11%      200ms 17.54%  runtime.slicebytetostring
      20ms  1.75% 93.86%      370ms 32.46%  strconv.FormatInt
(pprof)
```
[ 注意：本文中使用“分配”指的是 堆分配[2] 。栈上分配也是分配，但在性能方面，它们并不是没有那么昂贵或重要。]

看起来我们在使用 sha256、strconv、内存分配和垃圾回收方面花费了大量时间。现在我们知道需要改进什么。由于我们没有进行任何类型的复杂计算（可能除了 sha256），我们的大多数性能问题似乎都是由堆分配引起的。我们可以通过替换来精确地验证一下内存分配



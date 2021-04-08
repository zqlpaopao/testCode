看了曹大的一个线程回收的博文，如是自己实践了一下

```go
package main


/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
void output(char *str) {
    usleep(1000000);
    printf("%s\n", str);
}
*/
import "C"
import "unsafe"


import "net/http"
import _ "net/http/pprof"


func init() {
	go http.ListenAndServe(":9999", nil)
}


func main() {
	for i := 0;i < 1000;i++ {
		go func(){
			str := "hello cgo"
			//change to char*
			cstr := C.CString(str)
			C.output(cstr)
			C.free(unsafe.Pointer(cstr))


		}()
	}
	select{}
}
```

![image-20210408133131018](thread.assets/image-20210408133131018.png)

可以看到有1004个线程没有被回收

可见 Goroutine 退出了，历史上创建的线程也是不会退出的。

```go
runtime.LockOSThread()
```

加了这个侯发现线程是8个,因为我是8核的

# 程序调用查看
```go
package main


/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
void output(char *str) {
    usleep(1000000);
    printf("%s\n", str);
}
*/
import "C"
import (
	"log"
	"runtime"
	"sync"
	"unsafe"
)


import "net/http"
import _ "net/http/pprof"


func init() {
	go http.ListenAndServe(":9999", nil)
}


func main() {
	for i := 0;i < 1000;i++ {
		go func(){
			str := "hello cgo"
			//change to char*
			cstr := C.CString(str)
			C.output(cstr)
			C.free(unsafe.Pointer(cstr))
			//runtime.LockOSThread()
			//fmt.Println(str)
			//time.Sleep(10000*time.Second)


		}()
	}

	killThreadService()
	select{}
}

func killThreadService() {
	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":10003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func sayhello(wr http.ResponseWriter, r *http.Request) {
	KillOne()
}

// KillOne kills a thread
func KillOne() {
	var wg sync.WaitGroup
	wg.Add(1)


	go func() {
		defer wg.Done()
		runtime.LockOSThread()
		return
	}()


	wg.Wait()
}
```

curl localhost:10003 每次请求一次线程就会相对应的减少一个

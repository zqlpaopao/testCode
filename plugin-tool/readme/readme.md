# plugin 使用

## 1、构建文件内容
创建文件 main.go

```go
package main

import (
	"fmt"
	"log"
	"plugin"
)

var In int = 65

func Func()string{
	return "zhangSan"
}
```


生成文件
```go
 go build -buildmode=plugin -o plugin.so main.go
```
- -buildmode 固定参数
- -o 输出的文件名称
- main.go 要转换的文件

## 2、使用
```go
p ,err := plugin.Open("plugin.so")
	if err != nil{
		log.Panicln(err)
	}

	v ,err := p.Lookup("In")
	fmt.Println(*v.(*int))

	f , err := p.Lookup("Func")

	fn := f.(func()string)
	fmt.Println(fn())
```

  结果
```go
65
zhangSan
```

# cshared示例
除了上面提到的buildmode=plugin外，还有一种用法就是 buildmode=c-shared ，使用该参数时会生成出来两个文件，一个.so文件，一个.h头文件 ，使用起来就和使用c 生成的库文件和模块文件一样使用。具体使用如下：

## 编写文件
awesome.go
```go
package main
import "C"
import (
	"fmt"
	"math"
	"sort"
	"sync"
)
var count int
var mtx sync.Mutex
//export Add
func Add(a, b int) int {
	return a + b
}
//export Cosine
func Cosine(x float64) float64 {
	return math.Cos(x)
}
//export Sort
func Sort(vals []int) {
	sort.Ints(vals)
}
//export Log
func Log(msg string) int {
	mtx.Lock()
	defer mtx.Unlock()
	fmt.Println(msg)
	count++
	return count
}
func main() {}
```

## 2、生成文件
```go
go build -o awesome.so -buildmode=c-shared awesome.go

.
├── awesome.go
├── awesome.h
└── awesome.so

0 directories, 3 files

```

编译后生成如下两个文件：
awesome.h awesome.h

以下命令查看文件内容
```go
file awesome.so
nm awesome.so | grep -e "T Add" -e "T Cosine" -e "T Sort" -e "T Log"
```
其会输出为shared object文件，并exported 相关对像的symbols 。

## 3、c语言动态链接调用
```c
#include <stdio.h>
#include "awesome.h"
int main() {
    //Call Add() - passing integer params, interger result
    GoInt a = 12;
    GoInt b = 99;
    printf("awesome.Add(12,99) = %d\n", Add(a, b));
    //Call Cosine() - passing float param, float returned
    printf("awesome.Cosine(1) = %f\n", (float)(Cosine(1.0)));
    //Call Sort() - passing an array pointer
    GoInt data[6] = {77, 12, 5, 99, 28, 23};
    GoSlice nums = {data, 6, 6};
    Sort(nums);
    printf("awesome.Sort(77,12,5,99,28,23): ");
    for (int i = 0; i < 6; i++){
        printf("%d,", ((GoInt *)nums.data)[i]);
    }
    printf("\n");
    //Call Log() - passing string value
    GoString msg = {"Hello from C!", 13};
    Log(msg);
}
```
编译并调用
```c
$> gcc -o client client1.c ./awesome.so
$> ./client
awesome.Add(12,99) = 111
awesome.Cosine(1) = 0.540302
awesome.Sort(77,12,5,99,28,23): 5,12,23,28,77,99,
Hello from C!
```

## 4、 c语言动态加载
```c
#include <stdlib.h>
#include <stdio.h>
#include <dlfcn.h>
// define types needed
typedef long long go_int;
typedef double go_float64;
typedef struct{void *arr; go_int len; go_int cap;} go_slice;
typedef struct{const char *p; go_int len;} go_str;
int main(int argc, char **argv) {
    void *handle;
    char *error;
    // use dlopen to load shared object
    handle = dlopen ("./awesome.so", RTLD_LAZY);
    if (!handle) {
        fputs (dlerror(), stderr);
        exit(1);
    }
    // resolve Add symbol and assign to fn ptr
    go_int (*add)(go_int, go_int)  = dlsym(handle, "Add");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }
    // call Add()
    go_int sum = (*add)(12, 99);
    printf("awesome.Add(12, 99) = %d\n", sum);
    // resolve Cosine symbol
    go_float64 (*cosine)(go_float64) = dlsym(handle, "Cosine");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }
    // Call Cosine
    go_float64 cos = (*cosine)(1.0);
    printf("awesome.Cosine(1) = %f\n", cos);
    // resolve Sort symbol
    void (*sort)(go_slice) = dlsym(handle, "Sort");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }
    // call Sort
    go_int data[5] = {44,23,7,66,2};
    go_slice nums = {data, 5, 5};
    sort(nums);
    printf("awesome.Sort(44,23,7,66,2): ");
    for (int i = 0; i < 5; i++){
        printf("%d,", ((go_int *)data)[i]);
    }
    printf("\n");
    // resolve Log symbol
    go_int (*log)(go_str) = dlsym(handle, "Log");
    if ((error = dlerror()) != NULL)  {
        fputs(error, stderr);
        exit(1);
    }
    // call Log
    go_str msg = {"Hello from C!", 13};
    log(msg);
    // close file handle when done
    dlclose(handle);
}
```
编译执行：
```c
$> gcc -o client client2.c -ldl
$> ./client
awesome.Add(12, 99) = 111
awesome.Cosine(1) = 0.540302
awesome.Sort(44,23,7,66,2): 2,7,23,44,66,
Hello from C!
```

## 5、python ctypes 方式调用
```python
from ctypes import *
lib = cdll.LoadLibrary("./awesome.so")
# describe and invoke Add()
lib.Add.argtypes = [c_longlong, c_longlong]
lib.Add.restype = c_longlong
print "awesome.Add(12,99) = %d" % lib.Add(12,99)
# describe and invoke Cosine()
lib.Cosine.argtypes = [c_double]
lib.Cosine.restype = c_double
print "awesome.Cosine(1) = %f" % lib.Cosine(1)
# define class GoSlice to map to:
# C type struct { void *data; GoInt len; GoInt cap; }
class GoSlice(Structure):
    _fields_ = [("data", POINTER(c_void_p)), ("len", c_longlong), ("cap", c_longlong)]
nums = GoSlice((c_void_p * 5)(74, 4, 122, 9, 12), 5, 5)
# call Sort
lib.Sort.argtypes = [GoSlice]
lib.Sort.restype = None
lib.Sort(nums)
print "awesome.Sort(74,4,122,9,12) = [",
for i in range(nums.len):
    print "%d "% nums.data[i],
print "]"
# define class GoString to map:
# C type struct { const char *p; GoInt n; }
class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]
# describe and call Log()
lib.Log.argtypes = [GoString]
lib.Log.restype = c_longlong
msg = GoString(b"Hello Python!", 13)
print "log id %d"% lib.Log(msg)
```
执行结果：
```python
$> python client.py
awesome.Add(12,99) = 111
awesome.Cosine(1) = 0.540302
awesome.Sort(74,4,122,9,12) = [ 4 9 12 74 122 ]
Hello Python!
```

## 6、Python CFFI方式调用
```python
from __future__ import print_function
import sys
from cffi import FFI
is_64b = sys.maxsize > 2**32
ffi = FFI()
if is_64b: ffi.cdef("typedef long GoInt;\n")
else:      ffi.cdef("typedef int GoInt;\n")
ffi.cdef("""
typedef struct {
    void* data;
    GoInt len;
    GoInt cap;
} GoSlice;
typedef struct {
    const char *data;
    GoInt len;
} GoString;
GoInt Add(GoInt a, GoInt b);
double Cosine(double v);
void Sort(GoSlice values);
GoInt Log(GoString str);
""")
lib = ffi.dlopen("./awesome.so")
print("awesome.Add(12,99) = %d" % lib.Add(12,99))
print("awesome.Cosine(1) = %f" % lib.Cosine(1))
data = ffi.new("GoInt[]", [74,4,122,9,12])
nums = ffi.new("GoSlice*", {'data':data, 'len':5, 'cap':5})
lib.Sort(nums[0])
print("awesome.Sort(74,4,122,9,12) = %s" % [
    ffi.cast("GoInt*", nums.data)[i]
    for i in range(nums.len)])
data = ffi.new("char[]", b"Hello Python!")
msg = ffi.new("GoString*", {'data':data, 'len':13})
print("log id %d" % lib.Log(msg[0]))
```
## 7、java调用
```java
import com.sun.jna.*;
import java.util.*;
import java.lang.Long;
public class Client {
   public interface Awesome extends Library {
        // GoSlice class maps to:
        // C type struct { void *data; GoInt len; GoInt cap; }
        public class GoSlice extends Structure {
            public static class ByValue extends GoSlice implements Structure.ByValue {}
            public Pointer data;
            public long len;
            public long cap;
            protected List getFieldOrder(){
                return Arrays.asList(new String[]{"data","len","cap"});
            }
        }
        // GoString class maps to:
        // C type struct { const char *p; GoInt n; }
        public class GoString extends Structure {
            public static class ByValue extends GoString implements Structure.ByValue {}
            public String p;
            public long n;
            protected List getFieldOrder(){
                return Arrays.asList(new String[]{"p","n"});
            }
        }
        // Foreign functions
        public long Add(long a, long b);
        public double Cosine(double val);
        public void Sort(GoSlice.ByValue vals);
        public long Log(GoString.ByValue str);
    }
   static public void main(String argv[]) {
        Awesome awesome = (Awesome) Native.loadLibrary(
            "./awesome.so", Awesome.class);
        System.out.printf("awesome.Add(12, 99) = %s\n", awesome.Add(12, 99));
        System.out.printf("awesome.Cosine(1.0) = %s\n", awesome.Cosine(1.0));
        // Call Sort
        // First, prepare data array
        long[] nums = new long[]{53,11,5,2,88};
        Memory arr = new Memory(nums.length * Native.getNativeSize(Long.TYPE));
        arr.write(0, nums, 0, nums.length);
        // fill in the GoSlice class for type mapping
        Awesome.GoSlice.ByValue slice = new Awesome.GoSlice.ByValue();
        slice.data = arr;
        slice.len = nums.length;
        slice.cap = nums.length;
        awesome.Sort(slice);
        System.out.print("awesome.Sort(53,11,5,2,88) = [");
        long[] sorted = slice.data.getLongArray(0,nums.length);
        for(int i = 0; i < sorted.length; i++){
            System.out.print(sorted[i] + " ");
        }
        System.out.println("]");
        // Call Log
        Awesome.GoString.ByValue str = new Awesome.GoString.ByValue();
        str.p = "Hello Java!";
        str.n = str.p.length();
        System.out.printf("msgid %d\n", awesome.Log(str));
    }
}
```
更多示例可以参考：https://github.com/vladimirvivien/go-cshared-examples

三、总结

除了上面提到的示例外，c-shared模式的golang模块还支持nodejs、lua、ruby、Julia等语言的调用。个人理解是大部分语言都是用C开发的，由于golang自身与 c 的亲缘性，所以其生成的模块都是支持其他语言去调用的。

## 1、创建awesome.go
```go
package main
import "C"
import (
	"fmt"
	"math"
	"sort"
	"sync"
)
var count int
var mtx sync.Mutex
//export Add
func Add(a, b int) int {
	return a + b
}
//export Cosine
func Cosine(x float64) float64 {
	return math.Cos(x)
}
//export Sort
func Sort(vals []int) {
	sort.Ints(vals)
}
//export Log
func Log(msg string) int {
	mtx.Lock()
	defer mtx.Unlock()
	fmt.Println(msg)
	count++
	return count
}
func main() {}
```

编译
```go
go build -o awesome.so -buildmode=c-shared awesome.go
```

go 调用
```go
package main

/*
#include <stdio.h>
#include "awesome.h"
#cgo linux CFLAGS: -L./ -I./
#cgo linux LDFLAGS: -L./ -I./ -lcshared
*/
import "C"

import (
	"fmt"
)

func main() {

		//s := "hello"
		var i int
	i = C.Add(1,2)
	//C.Test()
	fmt.Println(i)
}
```










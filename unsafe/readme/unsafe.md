都听说过unsafe和unsafe.Pointer可以相互转换，其实他们也可以跟指针类型转换

![image-20210317181650227](unsafe.assets/image-20210317181650227.png)

# 1、unsafe包
此包只有三个函数
```go
func Sizeof(x ArbitraryType) uintptr

func Offsetof(x ArbitraryType) uintptr

func Alignof(x ArbitraryType) uintptr
```
- Sizeof返回x所占的字节数,==但并不包含`x`所指向的内容的大小==

  ```go
  func main(){
  	var s int64
  	s = 10
  	a := "jkhgyaaaaaaaaaa567890"
  	fmt.Println(unsafe.Sizeof(s))
  	fmt.Println(unsafe.Sizeof(&s))
  	fmt.Println(unsafe.Sizeof(&a))
  }
  ```

  64系统是8字节

  ```go
  8
  8
  8
  ```
  
- Offetof 主要作用是返回结构体成员在内存中的位置距离结构体起始位置处（结构体的第一个字段的偏移量都是0）的字节数，即偏移量
  
```go
  /**
   * @Author: zhangsan
   * @Description:
   * @File:  main
   * @Version: 1.0.0
   * @Date: 2021/3/17 下午5:57
   */
  
  package main
  
  import (
  	"fmt"
  	"unsafe"
  )
  type c struct {
  
  }
  type b struct {
  	namea string
  	//aa int64
  	bb string
  }
  
  type a struct {
  	c
  	b
  	name string
  	age int64
  }
  
  func main(){
  	var s a
  	fmt.Println(unsafe.Offsetof(s.namea))
  	fmt.Println(unsafe.Offsetof(s.name))
  	fmt.Println(unsafe.Offsetof(s.age))
  
  
  }
  
  ```
  
  空结构体是不占用字节的
  
  ```go
  0
  32
  48
  ```
  
  


/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/29 下午6:42
 */

package main

import (
	"fmt"
	"sync"
)

type Replacer struct {
	once   sync.Once // guards buildOnce method
	oldnew []string
}

func main() {
	oldnew := []string{"aa","bb"}
	a := &Replacer{oldnew: append([]string(nil), oldnew...)}
	fmt.Printf("%#v",a)

	fmt.Println()

	as := append([]string(nil), oldnew...)
	fmt.Printf("%#v",as)


}
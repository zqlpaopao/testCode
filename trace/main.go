/**
 * @Author: zhangSan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/5/11 下午4:26
 */

package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
)

func Sum(w *sync.WaitGroup, i int) {
	defer w.Done()
	var sum, n int64
	for ; n < 1000000000; n++ {
		sum += n
	}
	fmt.Println(i, sum)
}

func main() {
	runtime.GOMAXPROCS(1)

	f, err := os.Create("./trace.o")
	fmt.Println(err)
	defer f.Close()

	//开启trace编程
	_ = trace.Start(f)
	defer trace.Stop()

	var w sync.WaitGroup
	for i := 0; i < 10; i++ {
		w.Add(1)
		go Sum(&w, i)
	}
	w.Wait()
}
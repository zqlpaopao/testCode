package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	var f bool
	go func() {
		time.Sleep(time.Second * 5)
		cond.L.Lock()
		f = true
		cond.L.Unlock()

		cond.Signal()

	}()

	fmt.Println("waiting")

	cond.L.Lock()
	for !f {
		cond.Wait()
	}
	cond.L.Unlock()

	fmt.Println("done")
}
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sourcegraph/conc/iter"
	"github.com/sourcegraph/conc/pool"
	"github.com/sourcegraph/conc/stream"
	"sync/atomic"
	"time"
)

//https://mp.weixin.qq.com/s/AUSse5z1YES9wtKCiMdyWA

func main() {
	//type callBack chan func()
	//
	//var f chan callBack = make(chan callBack)
	//
	//go func() {
	//	for v := range f {
	//		v1 := <-v
	//		v1()
	//	}
	//}()
	//
	//for _, v := range []string{
	//	"第1",
	//	"第2",
	//	"第3",
	//} {
	//	//fmt.Println(k, v)
	//	//fmt.Println(v)
	//	vl := v
	//	var ck = make(callBack)
	//
	//	f <- ck
	//	ck <- func() {
	//		fmt.Println(vl)
	//	}
	//}
	//
	//select {}

	foreach()

	//ExampleStream()
	//wg := conc.NewWaitGroup()
	//for i := 0; i < 10; i++ {
	//	wg.Go(doSomething)
	//}
	//wg.WaitAndRecover()
}

func doSomething() {

	fmt.Println("test")
	panic("panic")
}

func ExampleContextPoolWithCancelOnError() {
	p := pool.New().
		WithMaxGoroutines(4).
		WithContext(context.Background()).
		WithCancelOnError()
	for i := 0; i < 3; i++ {
		i := i
		p.Go(func(ctx context.Context) error {
			if i == 2 {
				return errors.New("I will cancel all other tasks!")
			}
			<-ctx.Done()
			return nil
		})
	}
	err := p.Wait()
	fmt.Println(err)
	// Output:
	// I will cancel all other tasks!
}

func ExampleStream() {
	times := []int{20, 52, 16, 45, 4, 80}

	streams := stream.New()
	for _, millis := range times {
		dur := time.Duration(millis) * time.Millisecond
		streams.Go(func() stream.Callback {
			time.Sleep(dur)
			// This will print in the order the tasks were submitted
			return func() { fmt.Println(dur) }
		})
	}
	streams.Wait()

	// Output:
	// 20ms
	// 52ms
	// 16ms
	// 45ms
	// 4ms
	// 80ms
}

func foreach() {
	input := []int{1, 2, 3, 4}
	iterator := iter.Iterator[int]{
		MaxGoroutines: len(input) * 2,
	}

	iterator.ForEach(input, func(v *int) {
		if *v%2 != 0 {
			*v = -1
		}
	})

	fmt.Println(input)
}

// ////////////////////////
func forModel() {
	var idx atomic.Int64
	// Create the task outside the loop to avoid extra closure allocations.
	i := int(idx.Add(1) - 1)
	fmt.Println(i, 99)
	for ; i < 3; i = int(idx.Add(1) - 1) {
		fmt.Println(i)
	}
}

func SlToMap() {
	input := []int{1, 2, 3, 4}
	mapper := iter.Mapper[int, bool]{
		MaxGoroutines: len(input) / 2,
	}

	results := mapper.Map(input, func(v *int) bool { return *v%2 == 0 })
	fmt.Println(results)
	// Output:
	// [false true false true]
}

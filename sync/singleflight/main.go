package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {

	ch := make(chan int)
	go func() {
		for {
			time.Sleep(10*time.Second)

			ch <- 1
			time.Sleep(1*time.Second)
		}
	}()


	go func() {
		select {
			case v := <-ch:
				fmt.Println(v)
		}
		fmt.Println("end")
	}()

	select {

}
}


// doChan 此函数会在执行完毕后删除map中key信息
func doChan(){
	var singleSetCache singleflight.Group
	getAndSetCache := func(requestID int, cacheKey string) (string, error) {
		log.Printf("request %v start to get and set cache...", requestID)
		retChan := singleSetCache.DoChan(cacheKey, func() (ret interface{}, err error) {
			log.Printf("request %v is setting cache...", requestID)
			time.Sleep(3 * time.Second)
			log.Printf("request %v set cache success!", requestID)
			return "VALUE", nil
		})

		var ret singleflight.Result

		timeout := time.After(5 * time.Second)

		select { //加入了超时机制
		case <-timeout:
			log.Printf("time out!")
			return "", errors.New("time out")
		case ret = <-retChan: //从chan中取出结果
			fmt.Println( ret.Val.(string), ret.Err)
		}
		return "", nil
	}

	cacheKey := "cacheKey"
	for i := 1; i < 10; i++ {
		go func(requestID int) {
			value, _ := getAndSetCache(requestID, cacheKey)
			log.Printf("request %v get value: %v", requestID, value)
		}(i)
	}
	time.Sleep(20 * time.Second)
}




// Do 此函数会在执行完毕后删除map中key信息
func Do(){
	var singleSetCache singleflight.Group

	getAndSetCache := func(requestID int, cacheKey string) (string, error) {


		fmt.Printf("request %v start to get and set cache...\n", requestID)
		value, err, shard := singleSetCache.Do(cacheKey, func() (ret interface{}, err error) { //do的入参key，可以直接使用缓存的key，这样同一个缓存，只有一个协程会去读DB
			fmt.Printf("request %v is setting cache...\n", requestID)

			time.Sleep(3 * time.Second)
			fmt.Printf("request %v set cache success!---------------\n", requestID)
			fmt.Printf("%#v\n",singleSetCache)

			return "VALUE", nil
		})
		fmt.Println(err)
		fmt.Println(shard)
		fmt.Println(value)
		return value.(string), nil
	}

	cacheKey := "cacheKey"
	for i := 1; i < 10; i++ { //模拟多个协程同时请求
		go func(requestID int) {
			value, _ := getAndSetCache(requestID, cacheKey)
			fmt.Printf("request %v get value: %v", requestID, value)
		}(i)
	}
	time.Sleep(20 * time.Second)

	fmt.Printf("%#v\n",singleSetCache)

	for i := 1; i < 10; i++ { //模拟多个协程同时请求
		go func(requestID int) {
			value, _ := getAndSetCache(requestID, cacheKey)
			fmt.Printf("request %v get value: %v", requestID, value)
		}(i)
	}

	time.Sleep(20 * time.Second)
}

func singleflightFn(barrier *singleflight.Group){
	round := 10
	var wg sync.WaitGroup
	wg.Add(round)
	for i := 0; i < round; i++ {
		go func() {
			defer wg.Done()
			// 启用10个协程模拟获取缓存操作
			val, err, b := barrier.Do("get_rand_int", func() (interface{}, error) {
				time.Sleep(time.Second)
				return rand.Int(), nil
			})
			fmt.Println(b)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(val)
			}
		}()
	}
	wg.Wait()
}

func singleflightFnChan(barrier *singleflight.Group){
	round := 10
	var wg sync.WaitGroup
	wg.Add(round)
	for i := 0; i < round; i++ {
		go func() {
			defer wg.Done()
			// 启用10个协程模拟获取缓存操作
			cha := barrier.DoChan("get_rand_int", func() (interface{}, error) {
				time.Sleep(time.Second)
				return rand.Int(), nil
			})
			fmt.Println(<-cha)
			//if err != nil {
			//	fmt.Println(err)
			//} else {
			//	fmt.Println(val)
			//}
		}()
	}
	wg.Wait()
}

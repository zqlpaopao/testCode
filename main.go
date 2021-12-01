package main

import (
	"fmt"
	"sync"
	"time"
)

var Rw map[string]*sync.Mutex



func main(){
	Rw = map[string]*sync.Mutex{
		"1":new(sync.Mutex),
	}
	getsession()
	getsession()

}

func getsession(){

	Rw["1"].Lock()
	fmt.Println(2222)
	time.Sleep(5*time.Second)
	defer  Rw["1"].Unlock()


}
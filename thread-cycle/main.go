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
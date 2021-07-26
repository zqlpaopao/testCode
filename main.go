package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/hello",helloServer)
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	go http.ListenAndServe(":18080",hystrixStreamHandler)
	for {
		time.Sleep(10*time.Second)
	}
	log.Fatal()
}


func helloServer(w http.ResponseWriter, req *http.Request) {
	log.Println("request remote addr",req.RemoteAddr)
	io.WriteString(w, "hello,world")
}
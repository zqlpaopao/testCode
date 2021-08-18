package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test/router/routerApi"
	"time"
)

func main(){
	Shutdown()
	//endless.DefaultReadTimeOut =  5 * time.Second
	//endless.DefaultWriteTimeOut = 5* time.Second
	//endless.DefaultMaxHeaderBytes = 1 << 20
	//endPoint := fmt.Sprintf(":%d",8088)
	//
	//
	//r := &routerApi.Router{}
	//r.Route("GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println(r.Body)
	//	fmt.Println(r.URL)
	//	fmt.Println(r)
	//	fmt.Printf("%#v",r)
	//	w.Write([]byte("Hello，Chongchong!"))
	//})
	//r.Route("GET", `/hello/(?P<Message>\w+)`, func(w http.ResponseWriter, r *http.Request) {
	//	message := routerApi.URLParam(r, "Message")
	//	w.Write([]byte("Hello " + message))
	//})
	//server := endless.NewServer(endPoint,r)
	//
	//server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGUSR1] = append(
	//	server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGUSR1],
	//	preSigUsr1)
	//server.SignalHooks[endless.POST_SIGNAL][syscall.SIGUSR1] = append(
	//	server.SignalHooks[endless.POST_SIGNAL][syscall.SIGUSR1],
	//	postSigUsr1)
	//
	//server.BeforeBegin = func(add string){
	//	log.Printf("Actual pid is %d",syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil{
	//	log.Printf("Server err : %v",err)
	//}

}
func preSigUsr1() {
	log.Println("pre SIGUSR1")
}

func postSigUsr1() {
	log.Println("post SIGUSR1")
}



func Shutdown(){
	r := &routerApi.Router{}
	r.Route("GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
		fmt.Println(r.URL)
		fmt.Println(r)
		fmt.Printf("%#v",r)
		time.Sleep(10 * time.Second)
		fmt.Println(12345)
		w.Write([]byte("Hello，Chongchong!"))
	})
	r.Route("GET", `/hello/(?P<Message>\w+)`, func(w http.ResponseWriter, r *http.Request) {
		message := routerApi.URLParam(r, "Message")
		w.Write([]byte("Hello " + message))
	})
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d",8089),
		Handler:           r,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}


	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit,os.Interrupt)

		<- quit
		log .Println("shutdown Server ...")
		ctx ,cancel := context.WithTimeout(context.Background(),5 *time.Second)
		defer cancel()
		if err := server.Shutdown(ctx);nil != err{
			log.Fatal("Server Shutdown:" ,err)
		}
		log .Println("Server exiting")
	}()

	if err := server.ListenAndServe();nil != err{
		log.Println(err)
	}
}
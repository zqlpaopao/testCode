/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/5/10 下午5:04
 */

package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	servers "test/grpcurl/server"
	server "test/grpcurl/src"
)

func main(){
	grpcServe := grpc.NewServer()

	//注册grpcurl服务，否则报错
	reflection.Register(grpcServe)

	server.RegisterServerServer(grpcServe,servers.NewHelloServices())

	listen,e := net.Listen("tcp","127.0.0.1:8080")
	if e != nil{
		log.Fatal(e)
	}

	grpcServe.Serve(listen)
}

/**
 * @Author: zhangsan
 * @Description:
 * @File:  server.go
 * @Version: 1.0.0
 * @Date: 2021/5/10 下午5:14
 */

package servers

import (
	"fmt"
	"golang.org/x/net/context"
	server "test/grpcurl/src"
)

type HelloServices struct {

}

func NewHelloServices()*HelloServices{
	return &HelloServices{}
}

func(h *HelloServices)Hello(ctx context.Context,req *server.HelloReq)(resp *server.HelloResp,err error){
	fmt.Println(req.Name)
	resp = &server.HelloResp{
		Return:               "hello",
	}
	return
}

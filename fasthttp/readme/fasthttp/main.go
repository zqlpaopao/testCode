/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/5/19 上午10:51
 */

package main

import (
"fmt"
"log"
"net/http"
"runtime"
"time"

_ "expvar"

_ "net/http/pprof"

"github.com/valyala/fasthttp"
)

type HelloGoHandler struct {
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintln(ctx, "Hello, Go!")
}

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	go func() {
		for {
			log.Println("当前routine数量:", runtime.NumGoroutine())
			time.Sleep(time.Second)
		}
	}()

	s := &fasthttp.Server{
		Handler: fastHTTPHandler,
	}
	s.ListenAndServe(":8081")
}
/*

+-----------------------------------------------+
|         				Length（24）			    |
+---------------+--------------+----------------+
|   Type(8)     |  Flags(8)    |
+-+-------------+--------------+---------------------------------+
|R|                         Stream Identifier(32)                |
+=+==============================================================+
|                           Frame Payload(0...)                ...
+----------------------------------------------------------------+

*
 */


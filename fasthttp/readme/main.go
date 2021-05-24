/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/5/19 上午10:49
 */

package main
import (
	_ "expvar"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	go func() {
		for {
			log.Println("当前routine数量:", runtime.NumGoroutine())
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Go!"))
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
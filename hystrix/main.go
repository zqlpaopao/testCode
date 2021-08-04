package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	httpURL := "http://localhost:18080/hello"

	//初始化hystrix配置
	hystrix.ConfigureCommand("svc1", hystrix.CommandConfig{
		// 执行 command 的超时时间
		Timeout: 10,

		// 最大并发量
		MaxConcurrentRequests: 5,

		// 一个统计窗口 10 秒内请求数量
		// 达到这个请求数量后才去判断是否要开启熔断
		RequestVolumeThreshold: 20,

		// 熔断器被打开后
		// SleepWindow 的时间就是控制过多久后去尝试服务是否可用了
		// 单位为毫秒
		SleepWindow: 5000,

		// 错误百分比
		// 请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
		ErrorPercentThreshold: 30,
	})


	for i := 0 ; i < 11 ; i++ {
		for j := 0 ; j< 10 ; j++ {
			go func() {
				req,err := httpDo(httpURL)
				if err != nil {
					log.Println("http request error :",err.Error())
				}else {
					log.Println("http request success:",string(req))
				}
			}()
		}
		time.Sleep(time.Second*1)
	}

}

//http请求
func httpDo(url string) (respBytes []byte, err error) {
	var resp *http.Response
	err = hystrix.Do("svc1",func() error {
		resp,err = http.Get(url)
		if err != nil {
			return err
		}
		respBytes,err = ioutil.ReadAll(resp.Body)
		if err !=nil {
			return err
		}

		return nil
	}, func(error) error {
		return errors.New("http request error")
	})

	return
}
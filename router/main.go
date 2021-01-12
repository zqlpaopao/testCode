/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/12 上午11:06
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"test/router/routerApi"
)
//https://www.toutiao.com/i6877469517241877000/?app=news_article&group_id=6877469517241877000&is_new_connect=0&is_new_user=0&req_id=202101120909520101300361650A2C5759&share_token=35B2A523-16C1-40AF-A2C8-359817722554&timestamp=1610413793&tt_from=weixin&use_new_style=1&utm_campaign=client_share&utm_medium=toutiao_ios&utm_source=weixin&wxshare_count=1

func main(){
	var (
		err error
	)
	//如果是get的换可以做参数的转换存储
	r := &routerApi.Router{}
	r.Route("GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
		fmt.Println(r.URL)
		fmt.Println(r)
		fmt.Printf("%#v",r)
		w.Write([]byte("Hello，Chongchong!"))
	})
	r.Route("GET", `/hello/(?P<Message>\w+)`, func(w http.ResponseWriter, r *http.Request) {
		message := routerApi.URLParam(r, "Message")
		w.Write([]byte("Hello " + message))
	})
	if err = http.ListenAndServe(":8080",r);nil!=err{
		log.Fatal(err)
	}
	fmt.Println("listen 8080")







}


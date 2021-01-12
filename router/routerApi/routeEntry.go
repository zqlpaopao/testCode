/**
 * @Author: zhangsan
 * @Description:
 * @File:  routeEntry
 * @Version: 1.0.0
 * @Date: 2021/1/12 上午11:22
 */

package routerApi

import (
	"log"
	"net/http"
)

/**
* 定义请求实体
* 路径 方法
* 回调函数
*/
type RouteEntry struct {
	Path ,Method string
	Handler http.HandlerFunc
}

/**
* 定义router集合
* 实现了http的 ServeHTTP
*/
type Router struct {
	routes []RouteEntry
}

func (ro *Router) ServeHTTP(w http.ResponseWriter,r * http.Request){
	//防止panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR:", r)
			http.Error(w, "发生错误…", http.StatusInternalServerError)
		}
	}()

	for _ ,e := range ro.routes{
		match := e.Match(r)
		if !match{
			continue
		}
		e.Handler.ServeHTTP(w,r)
		return
	}
	http.NotFound(w,r)
}

/**
* 将每个请求添加到routeEntry实体中
*/
func (rte *Router) Route(method,path string, handFunc http.HandlerFunc){
	e := RouteEntry{
		Path:    path,
		Method:  method,
		Handler: handFunc,
	}
	rte.routes = append(rte.routes,e)
}


/**
* 每个实体的匹配
* 请求方法 请求路径
*/
func (rt *RouteEntry) Match( r *http.Request) bool{
	if r.Method != rt.Method{
		return false
	}
	if r.URL.Path != rt.Path{
		return  false
	}
	return true
}

func URLParam(r *http.Request, name string) string {
	ctx := r.Context()
	params := ctx.Value("params").(map[string]string)
	return params[name]
}
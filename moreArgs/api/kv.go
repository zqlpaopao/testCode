/**
 * @Author: zhangsan
 * @Description:
 * @File:  kv
 * @Version: 1.0.0
 * @Date: 2021/1/11 下午5:21
 */

package api

import (
	"fmt"
	"time"
)

type kv struct {
	name string
}

type PutResp struct {
	name string
}

func NewKv() kv {
	return kv{name: "lisi"}
}

func (k *kv) Put(op...Options) PutResp {
	return  k.Do(op...)
}

//此用可用一个type 区分具体类型，然后分开处理
func (k *kv)Do(op...Options) PutResp {
	args := OpPut(op...)
	time.Sleep(100*time.Second)
	fmt.Println(args)
	fmt.Println("处理业务请求--返回resp")

	return PutResp{name: "zhangSan"}
}
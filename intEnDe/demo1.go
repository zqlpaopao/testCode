/**
 * @Author: zhangsan
 * @Description:
 * @File:  demo1
 * @Version: 1.0.0
 * @Date: 2021/1/11 上午9:59
 */

package main

import (
	"fmt"
	"github.com/speps/go-hashids"
)

// 加密
func Encrypt(salt string, minLength int, params []int) string {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := hashids.NewWithData(hd)
	if err == nil {
		e, err := h.Encode(params)
		if err == nil {
			return e
		}
	}
	return ""
}

// 解密
func Decrypt(salt string, minLength int, hash string) []int {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := hashids.NewWithData(hd)
	if err == nil {
		e, err := h.DecodeWithError(hash)
		if err == nil {
			return e
		}
	}
	return []int{}
}

func main(){

	/*
		加密的key  加密后的长度  加密的参数
	 */
	enCode := Encrypt("this is my salt",30,[]int{1001})
	fmt.Println(enCode)
	//解密
	deCode := Decrypt("this is my salt",30,enCode)
	fmt.Println(deCode)
}
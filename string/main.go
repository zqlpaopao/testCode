/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/2/23 下午2:00
 */

package main

import (
	"fmt"
	"test/string/src"
)

func main(){
	str := "abcdefg"

	//普通str[1:4]
	a := src.StrSplit(1,4,str)
	fmt.Println(a)//bcd

	//中文的截取 []rune转换
	str = "中文的截取abcdefghjk"
	a = src.ChineseSplit(1,4,str)
	fmt.Println(a)//文的截

}
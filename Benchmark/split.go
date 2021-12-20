/**
 * @Author: zhangsan
 * @Description:
 * @File:  split
 * @Version: 1.0.0
 * @Date: 2021/1/11 上午11:19
 */

package Benchmark


import "strings"

func Split(s, sep string) (result []string) {
	result = make([]string, 0, strings.Count(s, sep)+1)
	i := strings.Index(s, sep)
	for i > -1 {
		result = append(result, s[:i])
		s = s[i+len(sep):] // 这里使用len(sep)获取sep的长度
		i = strings.Index(s, sep)
	}
	result = append(result, s)
	return
}


type Addr struct {
	name string
	age int
}

var r  = new(Addr)

func Tr()*Addr{
	cy := *r
	cy.name = "Addr"
	cy.age = 18
	return &cy
}

func Tr1()*Addr{
	return   &Addr{
		name: "addr1",
		age:  18,
	}

}
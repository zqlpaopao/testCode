/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/29 下午6:42
 */

package main

import (
	"fmt"
)

type IdPid struct {
	id,Pid int64
}

func main() {
	mp := map[int64]int64{
		1:0,
		2:0,
		3:2,
		4:1,
		5:0,
		6:2,
		7:3,
	}

	pids := []IdPid{}
	pars := make(map[int64][]int64)

	for i ,v := range mp{
		if v == 0{
			pars[i] = []int64{}
			continue
		}
		if !check(v,pids){
			pids = append(pids,IdPid{
				id:  i,
				Pid: v,
			})
		}
	}

	fmt.Printf("%+v\n",pids)
	fmt.Printf("%+v\n",pars)


	//for i ,v := range pars{
	//	if check(i,pids){
	//
	//	}
	//}

	//maps := map[int64][]int64{
	//	1 : []int64{4},
	//	2 : []int64{3,6},
	//	3 : []int64{7},
	//	4 :[]int64{},
	//	5 :[]int64{},
	//	6 :[]int64{},
	//	7 :[]int64{},
	//}
}

func check(i int64,sli []IdPid)bool{
	for _, v:= range sli{
		if v.Pid == i{
			return true
		}
	}
	return false
}
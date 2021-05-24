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


{
	"qscId":1234567,
	"account":"zhangsan",
	"type":1,
	"payStartTime":1620785243,
	"touchType":1,
	"touchTime":1620785243,
	"attType":1,
	"isPolicy":[{"PN897665432":1620785243},{"PNhhyu765543":1620785243}],
	"falsePolicy":[{"PN8976654321":1620785243},{"PNhhyu7655431":1620785243}],
	"isOrder":[{"OR89766543":1620785243}],
	"newOrder":[{"NO89766543":1620785243}]
}
"type":1,//业绩类型：1-服务类；2-新单；3-其他
"payStartTime":678987651,//支付时间 时间戳
"touchType":1,////触达类型 1-电话，2-连接
"touchTime":1620785243,//触达时间 时间戳
"attType":1,//归属绑定条件 1->=3在保保单，2-在保加单
"isPolicy":[//归属在保保单详情，如无为空，
	{
		"PN897665432":1620785243//保单号：支付时间 时间戳
	},
	{
		"PNhhyu765543":1620785243
	}
	]
"falsePolicy":[//非在保的保单详情
	{
		"PN8976654321":1620785243//保单号：支付时间 时间戳
	},
	{
		"PNhhyu7655431":1620785243
	}
	],
"order":[//归属订单详情
	{
		"OR89766543":1620785243//订单号：支付时间 时间戳
	}
	]
"newOrder":[//加保保单
	{
		"NO89766543":1620785243
	}
	]

1、超过36个月直接置为到期
2、超过12个月，是否存在>=3个在保保单 OR 在保加保
3、坐席离职直接置为归属关系为离职状态
4、不超过12个月的，是否存在>=1个在保保单
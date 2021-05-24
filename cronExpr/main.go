/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/15 上午11:10
 */

package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"os"
	"reflect"
	"time"
)


func getStructValue(expr *cronexpr.Expression)(int64){
	var (
		val reflect.Value
		secondSli []int64
		minuteSli []int64
		hourSli []int64
		monthSli []int64
		weekSli []int64
		yearSli []int64
	)
	val = reflect.ValueOf(*expr).FieldByName("secondList")
	for i := 0; i <val.Len() ;i ++{
		f := val.Index(i)
		secondSli = append(secondSli,f.Int())
	}
	if secondSli != nil && secondSli[1] - secondSli[0] > 1{
		return secondSli[1]- secondSli[0]
	}

	//minute
	val = reflect.ValueOf(*expr).FieldByName("minuteList")
	for i := 0; i <val.Len() ;i ++{
		f := val.Index(i)
		minuteSli = append(minuteSli,f.Int())
	}
	if minuteSli != nil && minuteSli[1]-minuteSli[0] > 1{
		return 60
	}

	//hour
	val = reflect.ValueOf(*expr).FieldByName("hourList")
	for i := 0; i <val.Len() ;i ++{
		f := val.Index(i)
		hourSli = append(hourSli,f.Int())
	}
	if hourSli != nil && hourSli[1] - hourSli[0] > 1{
		return 60
	}

	//month
	val = reflect.ValueOf(*expr).FieldByName("monthList")
	for i := 0; i <val.Len() ;i ++{
		f := val.Index(i)
		monthSli = append(monthSli,f.Int())
	}
	if monthSli != nil && monthSli[1] - monthSli[0] > 1{
		return 60
	}

	//week
	iter := reflect.ValueOf(*expr).FieldByName("daysOfWeek").MapRange()
	for iter.Next() {
		k := iter.Key()
		 //v := iter.Value()
		weekSli = append(weekSli,k.Int())
	}
	if weekSli != nil && weekSli[1] - weekSli[0] > 1{
		return 60
	}

	//year
	val = reflect.ValueOf(*expr).FieldByName("yearList")
	for i := 0; i <val.Len() ;i ++{
		f := val.Index(i)
		yearSli = append(yearSli,f.Int())
	}
	if yearSli != nil && yearSli[1] - yearSli[0] > 1{
		return 60
	}
	return 0
}

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)
	// 秒 分  时  天  月   星期 年 支持到秒 7位
	//第一种   支持7位，只写5位就是分钟了
	if expr, err = cronexpr.Parse("1  * 0 * * * *"); nil != err {
		fmt.Println(err)
	}
	//fmt.Printf("%#v",expr)
	fmt.Println(expr.Next(time.Now()))
	a := CheckTime(expr.Next(time.Now()),time.Now())
	fmt.Println(a)
	os.Exit(3)
	getStructValue(expr)


	os.Exit(3)
	//第二种MustParse() 返回*Expression 没有错误，认为表达式正确

	//获取当前时间
	now = time.Now()
	//传入当前时间，返回下次调度时间
	// 2020-05-02 22:25:00 +0800 CST
	nextTime = expr.Next(now)

	//定时器超时触发
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了")
	})

	//防止主函数退出
	time.Sleep(10 * time.Second)
	fmt.Println(nextTime)

}

func CheckTime(t1,t2 time.Time)bool{
	if t1.Year() != t2.Year(){
		return true
	}
	if t1.Month() != t2.Month(){
		return true
	}
	if t1.Day() != t2.Day(){
		return true
	}
	if t1.Hour() != t2.Hour(){
		return true
	}
	if t1.Minute() != t2.Minute(){
		return true
	}
	return false
}

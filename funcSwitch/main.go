/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/5/10 下午3:36
 */

package main

import "fmt"

/*
函数选项模式
*/


type (
	//逻辑处理函数
	funcS struct{
		opt options
	}

	//函数选项模式的参数实体
	options struct {
		name string
		age uint8
		sex uint8
	}

	//函数选项模式的处理函数
	fn func(*options)
)

/**********************************************实例话处理体********************************/
func Init(fns... fn)*funcS{
	//选项赋值行为
	//initOpt(fns...)

	//返回实体
	return  clone(initOpt(fns...))
}

func initOpt(fns... fn)(opt options){
	for _ ,v := range fns{
		v(&opt)
	}
	return
}

func clone(opt options)*funcS{
	return &funcS{opt: options{
		name:opt.name,
		age:  opt.age,
		sex:  opt.sex,
	}}
}

/**********************************************参数赋值行为********************************/
func initName(name string)fn{
	return func(o *options) {
		o.name = name
	}
}

func initAge(age uint8)fn{
	return func(o *options) {
		o.age = age
	}
}
func initSex(sex uint8)fn{
	return func(o *options) {
		o.sex = sex
	}
}





func main(){
	fS := Init(
		initName("zhangSan"),
		initAge(28),
		initSex(1),
		)

	fmt.Printf("%+v",fS)
	//&{opt:{name:zhangSan age:28 sex:1}}
}
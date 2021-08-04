package main

import "fmt"

/**
	clms
	定义接口，具体不同的业务不同的实现
	软件开发中，我们经常会遇到这样的场景，比如支付，用户支付，可以通过支付宝，微信，或者银联。
	他们最终的目的都是完成付钱的动作。我们可以归纳抽象这些支付渠道为相同的流程；
	基本的支付动作都是相同的，只不过支付实现（策略）不一样。比如调用的api 不一样，参数不一样，签名不一样等。

	我们可以定义基本的支付 interface，然后用不同的实现，完成 provider.DoPay 操作，完成支付。
	这种就是最简单的策略模式。

*/

//定义接口
type Payment interface {
	DoPay(string, string) int
}

//定义支付的结构体，实现空接口
type PaymentStrategy struct {
	Payment Payment
}

//定义策略方法，调用次接口
//相当于接口的入口管理者
func (o *PaymentStrategy) DoPayAction(config, args string) int {
	return o.Payment.DoPay(config, args)
}


/*********************************具体的逻辑实体***************************************/
//支付宝服务
type AlipayStrategy struct{}
func (AlipayStrategy) DoPay(config, args string) int {
	//具体的支付逻辑
	fmt.Println(config)
	fmt.Println(args)
	return 0
}



//微信支付
type WeixinStrategy struct{}
func (WeixinStrategy) DoPay(config, args string) int {
	//写具体的支付逻辑
	fmt.Println(config)
	fmt.Println(args)
	return 0
}



func main(){
	config := "支付宝的配置文件信息"
	args  := "支付宝请求参数"
	p := PaymentStrategy{AlipayStrategy{}}
	p.DoPayAction(config, args) // config 和args是支付相关的参数，比如支付配置信息和订单号等


	config = "微信的配置文件信息"
	args  = "微信请求参数"
	pw := PaymentStrategy{WeixinStrategy{}}
	pw.DoPayAction(config, args) // config 和args是支付相关的参数，比如支付配置信息和订单号等

}




package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//模版模式 ，有些事情要区别对待的，大部分相同，只有少数的不同
// https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247496813&idx=1&sn=fda31b59530deb7f2c05cff20bdc35c2&chksm=fa8325facdf4aceca2c7e4f0e38b416ea74f9b16231c06724b919ef6823de306fd3cf0c4303a&scene=178&cur_album_id=2531498848431669249#rd

type BankBusinessHandler interface {
	// 排队拿号
	TakeRowNumber()
	// 等位
	WaitInHead()
	// 处理具体业务
	HandleBusiness()
	// 对服务作出评价
	Commentate()
	// 钩子方法，
	// 用于在流程里判断是不是VIP， 实现类似VIP不用等的需求
	CheckVipIdentity() bool
}

type BankBusinessExecutor struct {
	handler BankBusinessHandler
}

// 模板方法，处理银行业务
func (b *BankBusinessExecutor) ExecuteBankBusiness() {
	// 适用于与客户端单次交互的流程
	// 如果需要与客户端多次交互才能完成整个流程，
	// 每次交互的操作去调对应模板里定义的方法就好，并不需要一个调用所有方法的模板方法
	b.handler.TakeRowNumber()
	if !b.handler.CheckVipIdentity() {
		b.handler.WaitInHead()
	}
	b.handler.HandleBusiness()
	b.handler.Commentate()
}

type DepositBusinessHandler struct {
	*DefaultBusinessHandler
	userVip bool
}

// 通用的方法还可以抽象到BaseBusinessHandler里，组合到具体实现类里，减少重复代码（实现类似子类继承抽象类的效果）
func (*DepositBusinessHandler) TakeRowNumber() {
	fmt.Println("请拿好您的取件码：" + strconv.Itoa(rand.Intn(100)) +
		" ，注意排队情况，过号后顺延三个安排")
}

func (dh *DepositBusinessHandler) WaitInHead() {
	fmt.Println("排队等号中...")
	time.Sleep(5 * time.Second)
	fmt.Println("请去窗口xxx...")
}

func (*DepositBusinessHandler) HandleBusiness() {
	fmt.Println("账户存储很多万人民币...")
}

func (dh *DepositBusinessHandler) CheckVipIdentity() bool {
	return dh.userVip
}

func (*DepositBusinessHandler) Commentate() {

	fmt.Println("请对我的服务作出评价，满意请按0，满意请按0，(～￣▽￣)～")
}

func NewBankBusinessExecutor(businessHandler BankBusinessHandler) *BankBusinessExecutor {
	return &BankBusinessExecutor{handler: businessHandler}
}

type DefaultBusinessHandler struct {
}

func (*DefaultBusinessHandler) TakeRowNumber() {
	fmt.Println("请拿好您的取件码：" + strconv.Itoa(rand.Intn(100)) +
		" ，注意排队情况，过号后顺延三个安排")
}

func (dbh *DefaultBusinessHandler) WaitInHead() {
	fmt.Println("排队等号中...")
	time.Sleep(5 * time.Second)
	fmt.Println("请去窗口xxx...")
}

func (*DefaultBusinessHandler) Commentate() {

	fmt.Println("请对我的服务作出评价，满意请按0，满意请按0，(～￣▽￣)～")
}

func (*DefaultBusinessHandler) CheckVipIdentity() bool {
	// 留给具体实现类实现
	return false
}

func main() {
	dh := &DepositBusinessHandler{userVip: false}
	bbe := NewBankBusinessExecutor(dh)
	bbe.ExecuteBankBusiness()
}

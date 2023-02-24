package main

import "fmt"

//职责链模式
/*

	是一个递归的过程，适合产品线那种多个流程处理的 最后返回结果，例如数据不同维度的染色
	https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247497039&idx=1&sn=ee5ef2ca2a378e9836564da0f2eae485&chksm=fa8324d8cdf4adce942debfe07b76f656bc3963ec9a70192d0195a9e2b2df9e15589e4630438&scene=178&cur_album_id=2531498848431669249#rd

*/
// 流程中的请求类--患者
type patient struct {
	Name              string
	RegistrationDone  bool
	DoctorCheckUpDone bool
	MedicineDone      bool
	PaymentDone       bool
}

type PatientHandler interface {
	Execute(*patient) error
	SetNext(PatientHandler) PatientHandler
	Do(*patient) error
}

// Next 充当抽象类型，实现公共方法，抽象方法不实现留给实现类自己实现
type Next struct {
	nextHandler PatientHandler
}

func (n *Next) SetNext(handler PatientHandler) PatientHandler {
	n.nextHandler = handler
	return handler
}

func (n *Next) Execute(patient *patient) (err error) {
	// 调用不到外部类型的 Do 方法，所以 Next 不能实现 Do 方法
	if n.nextHandler != nil {
		if err = n.nextHandler.Do(patient); err != nil {
			return
		}

		return n.nextHandler.Execute(patient)
	}

	return
}

// Reception 挂号处处理器
type Reception struct {
	Next
}

func (r *Reception) Do(p *patient) (err error) {
	if p.RegistrationDone {
		fmt.Println("Patient registration already done")
		return
	}
	fmt.Println("Reception registering patient")
	p.RegistrationDone = true
	return
}

// Clinic 诊室处理器--用于医生给病人看病
type Clinic struct {
	Next
}

func (d *Clinic) Do(p *patient) (err error) {
	if p.DoctorCheckUpDone {
		fmt.Println("Doctor checkup already done")
		return
	}
	fmt.Println("Doctor checking patient")
	p.DoctorCheckUpDone = true
	return
}

// Cashier 收费处处理器
type Cashier struct {
	Next
}

func (c *Cashier) Do(p *patient) (err error) {
	if p.PaymentDone {
		fmt.Println("Payment Done")
		return
	}
	fmt.Println("Cashier getting money from patient patient")
	p.PaymentDone = true
	return
}

// Pharmacy 药房处理器
type Pharmacy struct {
	Next
}

func (m *Pharmacy) Do(p *patient) (err error) {
	if p.MedicineDone {
		fmt.Println("Medicine already given to patient")
		return
	}
	fmt.Println("Pharmacy giving medicine to patient")
	p.MedicineDone = true
	return
}

func main() {
	receptionHandler := &Reception{}
	patient := &patient{Name: "abc"}
	// 设置病人看病的链路
	receptionHandler.SetNext(&Clinic{}).SetNext(&Cashier{}).SetNext(&Pharmacy{})
	receptionHandler.Execute(patient)
}

package main

import "fmt"

//观察者接口
type Observer interface {
	Notify(config,args string) error
}
//实体者
type Subject struct {
	Observers []Observer
	Conf      string
}
// 获取实例
func NewSubject(conf string) *Subject {
	return &Subject{
		Conf: conf,
	}
}
// 添加观察者
func (s *Subject) AddObserve(o ...Observer) {
	s.Observers = append(s.Observers, o...)
}
// 通知观察者
func (s *Subject) Notify(config ,args string) (err error) {
	for _, o := range s.Observers {
		err = o.Notify(config ,args)//分别处理err
	}
	return
}
/*****************************不同逻辑实体******************************/
type SubMessage struct {
	name string
	config string
}
func NewSubMessage()*SubMessage{
	return &SubMessage{name: "SubMessage", config: "实体1配置文件"}
}

func(s *SubMessage)Notify(config ,args string)error{
	fmt.Println("公共信息",config,args);fmt.Println(s.config);fmt.Println(s.name);return nil
}

//实体2
type SubMessage1 struct {
	name string
	config string
}

func NewSubMessage1()*SubMessage1{
	return &SubMessage1{name: "SubMessage2", config: "实体2配置文件"}
}

func(s *SubMessage1)Notify(config ,args string)error{
	fmt.Println("公共信息",config,args);fmt.Println(s.config);fmt.Println(s.name);return nil
}
/*****************************公共数据体******************************/
type kafkaMsg struct {
	msg []string
	Subjects *Subject
}

func NewKafkaMsg(m *Subject)*kafkaMsg{
	return &kafkaMsg{msg: []string{}, Subjects: m}
}

func (k *kafkaMsg)getMsg(){
	//获取数据
}

//真正处理
func (k *kafkaMsg)DoTidy()error{
	return k.Subjects.Notify(k.Subjects.Conf,"args")
}

func main(){
	//初始化操作实体并添加观察者
	m := NewSubject("conf")
	m.AddObserve(NewSubMessage(),NewSubMessage1())

	//获取数据
	k := NewKafkaMsg(m)
	k.getMsg()
	if k.DoTidy()!= nil{
		panic("错误")
	}
}
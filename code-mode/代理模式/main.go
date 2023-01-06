package main

import "fmt"

type Vehicle interface {
	Drive()
}

//********************普通版本的 就是定义一个struct 然后实现接口
// https://mp.weixin.qq.com/s/3V-1AcuWq1sWTbIrVdjLSw

type Car struct {
}

func (c *Car) Drive() {
	fmt.Println("car is being driven")
}

// 此时需要增加一个年龄的限制，但是age的限制是针对开车者  而不是car

type Driver struct {
	Age int
}

//包装 Driver 和 Vehicle 类型的包装类型

type CarProxy struct {
	vehicle Vehicle
	driver  *Driver
}

func NewCarProxy(driver *Driver) *CarProxy {
	return &CarProxy{&Car{}, driver}
}

//用包装类型代理vehicle属性的 Drive() 行为时，给它加上驾驶员的年龄限制。

func (c *CarProxy) Drive() {
	if c.driver.Age >= 16 {
		c.vehicle.Drive()
	} else {
		fmt.Println("Driver too young!")
	}
}

//例如数据推送多个kafka es 等等
//这个代理就可以做统一的逻辑处理，然后最后各自的drive 去执行相关逻辑

func main() {
	car := NewCarProxy(&Driver{12})
	car.Drive() // 输出 Driver too young!
	car2 := NewCarProxy(&Driver{22})
	car2.Drive() // 输出 Car is being driven

	c3 := NewCarProxy(&Driver{15})
	c3.vehicle.Drive()
}

package main
/*
https://mp.weixin.qq.com/s/hrnQI8Gx-LeJq9jQBhWsxg
//适合做统一化的kafka这种的有标准化东西的需求，可以参考之前写的kafka的消费的那个qsc
 */
import "fmt"

type mode func (string)string

func (m mode) Print( str string){
	fmt.Println(m(str))
}


func english(name string)string{
	return "hello"  + name
}

func english1(name string)string{
	return "hello" + name
}



func main(){
	greet := mode(english)
	greet.Print("world")
}
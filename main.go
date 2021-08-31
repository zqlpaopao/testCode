package main

import "fmt"

func main() {
	switchFu()
}




func switchFu(){
	var i int = 1

	switch i {
	case 1:
		fmt.Println(1)
		fallthrough
	case 2:
		fmt.Println(2)
	case 3:
		fmt.Println(3)
	default:
		fmt.Println("default")

	}
}

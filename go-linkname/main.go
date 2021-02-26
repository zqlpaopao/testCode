package main

import (
	"fmt"
	"test/go-linkname/a"
)

func main() {
	fmt.Println(a.Greet())
	a.Hello()

	fmt.Println(a.Aa)
	fmt.Println(a.Ab)
}
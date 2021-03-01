package main

import (
	"fmt"
	"test/go-linkname/a"
)
//go:generate  echo hello word
func main() {
	fmt.Println(a.Greet())


	fmt.Println(a.Aa)
	fmt.Println(a.Ab)
	a.Adds()
}
//go:generate  go version
func add(){

}
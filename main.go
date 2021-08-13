package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

func sayhello(wr http.ResponseWriter, r *http.Request) {}

func main() {
	str := "111200.0.0.xlsx"

	//fmt.Println(strings.Split(str,"."))
	//fmt.Println(strings.SplitN(str,".",2))
	//fmt.Println(strings.SplitN(str,".",-1),-1)
	fmt.Println(strings.SplitAfter(str,".")[len(strings.SplitAfter(str,"."))-1])
	//fmt.Println(strings.SplitAfterN(str,".",2))
}
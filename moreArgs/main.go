package main

import (
	"test/moreArgs/api"
)

func main(){
	NewKv := api.NewKv()
	NewKv.Put(api.WithLimit(int64(100)), api.WithGroup(true), api.WithOrder(true))
}


package main

import "test/moreArgs"

func main(){
	NewKv := moreArgs.NewKv()
	NewKv.Put(moreArgs.WithLimit(int64(100)),moreArgs.WithGroup(true),moreArgs.WithOrder(true))
}
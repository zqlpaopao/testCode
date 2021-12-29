package main

import (
	"fmt"
)

type Flags uint

const (
	_  Flags = iota
	FlagUp            // interface is up
	FlagBroadcast                      // interface supports broadcast access capability
	FlagLoopback                       // interface is a loopback interface
	FlagPointToPoint                   // interface belongs to a point-to-point link
	FlagMulticast                      // interface supports multicast access capability
)

var flagNames = []string{
	"up",
	"broadcast",
	"loopback",
	"pointtopoint",
	"multicast",
}

func (f Flags) String() string {
	return  flagNames[f+1]
}

func (f Flags)Int()int{
	return int(f)
}

func (f Flags)New()*Flags{
	return &f
}


func main(){
	fmt.Println(FlagUp)
	fmt.Println(FlagUp.Int())


	fmt.Println(FlagBroadcast)
	fmt.Println(FlagBroadcast.Int())

	fmt.Println(FlagLoopback)
	fmt.Println(FlagLoopback.Int())
	//src.GetLocalIp()


}

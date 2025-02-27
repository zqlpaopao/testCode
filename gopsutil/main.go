package main

import (
"fmt"

"github.com/shirou/gopsutil/v3/mem"
//"github.com/shirou/gopsutil/v3/docker"
// "github.com/shirou/gopsutil/mem"  // to use v2
)

func main() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)


}
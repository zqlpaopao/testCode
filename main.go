package main

import (
	"syscall"
)

func main() {
	var rlimit syscall.Rlimit

	// 限制cpu个数
	rlimit.Cur = 1
	rlimit.Max = 2
	syscall.Setrlimit(syscall.RLIMIT_CPU, &rlimit)
	err := syscall.Getrlimit(syscall.RLIMIT_CPU, &rlimit)
	if err != nil {
		panic(err)
	}

	rlimit.Cur = 100 //以字节为单位
	rlimit.Max = rlimit.Cur + 1024
	err = syscall.Setrlimit(syscall.RLIMIT_CORE, &rlimit)
	err = syscall.Getrlimit(syscall.RLIMIT_CORE, &rlimit)
	if err != nil {
		panic(err)
	}
}


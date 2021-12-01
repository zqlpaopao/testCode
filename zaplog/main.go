package main

import (
	"test/zaplog/src"
	"time"
)

func main() {
	src.InitZapLog()
	for {
		a := []string{"mysql","hello","world"}
		src.Debug("output",a)
		//src.Warn("output",a)
		//src.Error("output",a)
		//src.Debug("output",a)
		time.Sleep( 500 *time.Millisecond)
	}

}
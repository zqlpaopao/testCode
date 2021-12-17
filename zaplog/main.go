package main

import (
	"test/zaplog/src"
	"time"
)

type strucs struct {
	name string
	age int
}
func main() {
	src.InitZapLog()
	for {
		a := []string{"mysql","hello","world"}
		src.Debug("output",a)
		//src.Warn("output",a)
		//src.Error("output",a)
		//src.Debug("output",a)
b := strucs{
	name: "111",
	age:  11,
}
		src.Debug("output",b)

		time.Sleep( 500 *time.Millisecond)

	}

}
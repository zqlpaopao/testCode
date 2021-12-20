package main
//
import "testing"

func BenchmarkAddr(b *testing.B, ) {
	for i := 0; i < 10; i++ {
		Addr()
	}
}
func BenchmarkAddr1(b *testing.B, ) {
	for i := 0; i < 10; i++ {
		Addr1()
	}
}

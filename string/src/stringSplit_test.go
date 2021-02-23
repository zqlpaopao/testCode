package benchmark

import (
	"testing"
	"unicode/utf8"
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
)

var benchmarkStr = "此次进行的是中文的截取的性能测试情况"


func StrSplit(i,j int,str string)string{
	return str[i:j]
}

func BenchmarkStrSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrSplit(3,10, benchmarkStr)
	}
}


func SubStrDecodeRuneInString(s string, length int) string {
	var size, n int
	for i := 0; i < length && n < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[n:])
		n += size
	}

	return s[:n]
}

func BenchmarkSubStrDecodeRuneInString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SubStrDecodeRuneInString(benchmarkStr, 10)
	}
}





func SplitStrRunes(s string, length int) string {
	if utf8.RuneCountInString(s) > length {
		rs := []rune(s)
		return string(rs[:length])
	}

	return s
}

func BenchmarkSubStrRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SplitStrRunes(benchmarkStr, 10)
	}
}


func SubStrRange(s string, length int) string {
	var n, i int
	for i = range s {
		if n == length {
			break
		}

		n++
	}

	return s[:i]
}

func BenchmarkSubStrRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SubStrRange(benchmarkStr, 10)
	}
}


func SubStrRuneIndexInString(s string, length int) string {
	n, _ := exutf8.RuneIndexInString(s, length)
	return s[:n]
}

func BenchmarkSubStrRuneIndexInString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SubStrRuneIndexInString(benchmarkStr, 10)
	}
}
package c

import (
	"test/go-linkname/b"
)

func SayHii(name string) string {
	return b.Hi(name)
}
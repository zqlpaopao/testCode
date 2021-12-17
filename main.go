package main

import (
	"fmt"
	"sync"
)

type NpCopy struct {
	noCopys
	DoNotCopy
}

type noCopys struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopys) Lock()   {}
func (*noCopys) Unlock() {}

type DoNotCopy [0]sync.Mutex
//

func main() {
	var w NpCopy
	var w1 sync.WaitGroup
	No(w, w1)

}

func No(npCopy NpCopy, w1 sync.WaitGroup) {
	fmt.Println(w1)
	fmt.Println(npCopy)
}

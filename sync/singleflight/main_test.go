package main

import (
	"fmt"
	. "golang.org/x/sync/singleflight"
	"testing"
	"time"
)

func process1(g *Group, t *testing.T, ch chan int, key string) {
	for count := 0; count < 10; count++ {
		v, err, shared := g.Do(key, func() (interface{}, error) {
			time.Sleep(1000 * time.Millisecond)
			return "bar", nil
		})
		t.Log("v = ", v, " err = ", err, " shared =", shared, " ch :", ch, "g ")
		if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
			t.Errorf("Do = %v; want %v", got, want)
		}
		if err != nil {
			t.Errorf("Do error = %v", err)
		}
	}
	ch <- 1
}

func TestDo1(t *testing.T) {
	var g Group
	channels := make([]chan int, 10)
	key := "key"
	for i := 0; i < 10; i++ {
		channels[i] = make(chan int)
		go process1(&g, t, channels[i], key)
	}
	for i, ch := range channels {
		<-ch
		fmt.Println("routine ", i, "quit!")
	}
}

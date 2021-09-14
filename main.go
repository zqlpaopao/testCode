package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// for build -ldflags
var (
	// explain for more details: https://colobu.com/2015/10/09/Linux-Signals/
	// Signal selection reference:
	//   1. https://github.com/fvbock/endless/blob/master/endless.go
	//   2. https://blog.csdn.net/chuanglan/article/details/80750119
	hookableSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
		syscall.SIGQUIT,
		syscall.SIGSTOP,
		syscall.SIGKILL,
	}
	defaultHeartbeatTime = 1 * time.Minute
)

func handleSignal(pid int) {
	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, hookableSignals...)
	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		sig := <-sigs
		log.Println("pid[%d], signal: [%v]", pid, sig)
		done <- true
	}()
	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	for {
		log.Println("pid[%d], awaiting signal", pid)

		select {
		case <-done:
			log.Println("exiting")
			return
		case <-time.After(defaultHeartbeatTime):
		}
	}
}

func main(){
	pid := os.Getpid()
	handleSignal(pid)
}
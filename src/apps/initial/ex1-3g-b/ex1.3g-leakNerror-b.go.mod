package main

import (
	"os"
	"time"
	"runtime/trace"
)

func main() {
	trace.Start(os.Stderr)
	defer func() {
		time.Sleep(50 * time.Millisecond)
		trace.Stop()
	}()
	ch := make(chan int, 1)
	go send(ch)
	go recv(ch)
}

func send(ch chan int) {
	ch <- 42
}

func recv(ch chan int) {
	<-ch
}

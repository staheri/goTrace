package main

import (
	"time"
	"runtime/trace"
	"os"
)

func main() {
	trace.Start(os.Stderr)
	defer func() {
		time.Sleep(50 * time.Millisecond)
		trace.Stop()
	}()
	ch := make(chan int)
	go send(ch)
	go recv(ch)
	time.Sleep(100 * time.Millisecond)
	close(ch)
}

func send(ch chan int) {
	ch <- 42
}

func recv(ch chan int) {
	<-ch
}

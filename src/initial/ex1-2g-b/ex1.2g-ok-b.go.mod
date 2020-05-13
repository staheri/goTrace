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
	ch := make(chan int, 1)
	go send(ch)
	<-ch
	time.Sleep(1000 * time.Millisecond)
}

func send(ch chan int) {
	ch <- 42
}

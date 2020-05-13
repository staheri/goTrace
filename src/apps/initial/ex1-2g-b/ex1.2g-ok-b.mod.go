package main

import (
	"os"
	"time"
	"runtime/trace"
)

func main() {
	f,_ := os.Create("trace.out")
	trace.Start(f)
	defer func() {
		time.Sleep(50 * time.Millisecond)
		trace.Stop()
	}()
	ch := make(chan int, 1)
	go send(ch)
	<-ch
}

func send(ch chan int) {
	ch <- 42
}

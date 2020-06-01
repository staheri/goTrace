package main

import (
	"os"
	_ "net"
	"time"
	"runtime/trace"
)

func main() {
	trace.Start(os.Stderr)
	go func() {
		time.Sleep(5 * time.Second)
		trace.Stop()
		os.Exit(0)
	}()
	ch := make(chan int)
	ch <- 1
	<-ch
	close(ch)
}

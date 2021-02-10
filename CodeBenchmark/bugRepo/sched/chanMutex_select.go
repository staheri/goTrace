package main

import (
	"fmt"
	"sync"
	"time"
)

// commit hash: a69a59ffc7e3d028a72d1195c2c1535f447eaa84
// buggy version: 18768fdc2e76ec6c600c8ab57d2d487ee7877794
// https://github.com/moby/moby/pull/27782
// Different timing in select options and conditionals cause deadlock

func main() {
	cv := sync.NewCond(&sync.Mutex{})

	done := make(chan int)
	eventch := make(chan int)
	closec := make(chan int)

	go g0(closec)
	go g1(cv, done, eventch, closec)
	go g2(cv, done, eventch)

	fmt.Println("End of main!")
}

func g0(ch chan int) {
	time.Sleep(time.Microsecond*20)
  //runtime.Gosched()
  // to trigger the bug, it has to happen first
	close(ch)
}

func g1(cv *sync.Cond, done, eventch, closec chan int) {
	select {
	case <-eventch:
	case <-closec:
		cv.L.Lock()
		cv.Wait()
		cv.L.Unlock()
		close(done)
	}
}

func g2(cv *sync.Cond, done, eventch chan int) {
	select {
	case eventch <- 1:
	case <-done:
	}
	cv.Broadcast()
}

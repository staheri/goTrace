package main

// https://github.com/kubernetes/kubernetes/pull/16223
// Buggy version: e755988d5922df4d0e111a0167d9859359113463
// https://github.com/kubernetes/kubernetes/pull/10182
// buggy version: 4b990d128a17eea9058d28a3b3688ab8abafbd94

import (
  "sync"
  "fmt"
	"time"
)

func main() {
  ch := make(chan int)
  m := sync.Mutex{}
	cv := sync.NewCond(&m)
  var m1 sync.Mutex

  // goroutine 1
  go func() {
		time.Sleep(5*time.Millisecond) // simulates computation
    m1.Lock()
    cv.L.Lock()
    cv.Signal()
    cv.L.Unlock()
    m1.Unlock()
  }()

  // goroutine 2
  go func() {
    cv.L.Lock()
    cv.Wait()
    cv.L.Unlock()
    close(ch)
  }()

  // goroutine 3
  go func(){
		time.Sleep(5*time.Millisecond) // simulates computation
    m1.Lock()
    <-ch
    m1.Unlock()
  }()

  fmt.Println("End of main!")
}

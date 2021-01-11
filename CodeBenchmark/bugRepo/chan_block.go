// 58befe3081726ef74ea09198cd9488fb42c51f51
// https://github.com/moby/moby/pull/256
// the bug above happened because the os pipe buffer is full
// and the Copy() function is blocked because it is full and
// wait until the buffer is read out
// Fix: read the buffer before copys to it

// Here we simulate the scenario with two goroutines and a channel
package main

import (
  "fmt"
  _"runtime"
  "time"
)


func main() {
  //runtime.GOMAXPROCS(1)
  ch := make(chan int)
  go func(){
    ch <- 1
  }()
  go func(){
    //runtime.Gosched()
    ch <- 2 // expect to be blocked
  }()
  /*go func(){
    runtime.Gosched()
    <- ch // read
  }()*/
  <- ch
  time.Sleep(time.Millisecond*10)
  fmt.Println("End of main!")
}

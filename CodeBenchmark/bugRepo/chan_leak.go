package main

// d3a6ee1e55a53ee54b91ffb6c53ba674768cf9de
// https://github.com/moby/moby/pull/4395
// Goroutine leak because of undrained channel
// fix: Buffer the channel or drain the channel before return

import (
  "fmt"
  "runtime"
  "time"
)


func main() {
  runtime.GOMAXPROCS(1)

  go func()chan int{
    ch := make(chan int)
    // fix1: ch := make(chan int, 1)
    go func(){
      ch <- 0
      }()
    return ch
    }()
  // fix2: <- ch
  time.Sleep(time.Millisecond*10)
  fmt.Println("End of main!")
}

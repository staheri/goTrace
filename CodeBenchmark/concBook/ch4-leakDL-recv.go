package main

import (
  "fmt"
)

func main(){
  doWork := func(done <-chan int, strings <-chan string) <-chan interface{} {
    completed := make(chan interface{})
    go func() {
      defer fmt.Println("doWork exited.")
      defer close(completed)
      for s := range strings {
        // Do something interesting
        fmt.Println(s)
      }
      }()
    return completed
  }

  done := make(chan int)
  terminated := doWork(done, nil)

  // Perhaps more work is done here
  <-terminated
  fmt.Println("Done.")
}

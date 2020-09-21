package main

import (
  "fmt"
  "time"
)

func main(){
  doWork := func(
    done <-chan int,
    strings <-chan string,
    ) <-chan interface{} {
      terminated := make(chan interface{})
      go func() {
        defer fmt.Println("doWork exited.")
        defer close(terminated)
        for {
          select {
          case s := <-strings:
            // Do something interesting
            fmt.Println(s)
          case <-done:
            return
          }
        }
      }()
      return terminated
    }

  done := make(chan int)
  terminated := doWork(done, nil)

  go func() {
    // Cancel the operation after 1 second.
    time.Sleep(1 * time.Second)
    fmt.Println("Canceling doWork goroutine...")
    close(done)
  }()

  <-terminated
  fmt.Println("Done.")
}

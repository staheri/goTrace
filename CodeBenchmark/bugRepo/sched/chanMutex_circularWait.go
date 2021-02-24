package main
import (
  "sync"
  "fmt"
)

// https://github.com/moby/moby/pull/28462
// https://github.com/moby/moby/issues/28405
// introduce commit: b6c7becb
// Blocking of channel and lock

func main() {
  ch := make(chan int)
  var m sync.Mutex

  // goroutine 1
  go func() {
    m.Lock()
    ch <- 1
    m.Unlock()
  }()

  // goroutine 2
  go func() {
    m.Lock()
    m.Unlock()
    <-ch
  }()

	fmt.Println("End of main!")
}


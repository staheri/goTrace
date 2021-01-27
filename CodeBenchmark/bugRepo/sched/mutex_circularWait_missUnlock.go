package main
import (
  "fmt"
  "sync"
  "time"
)

// https://github.com/moby/moby/pull/17176
// introduce commit: d295dc66
// Lock() but does not drop it if Condition. This will block
// other goroutines trying to acquire the same lock

type sharedObject struct{
  content      bool
  mu           sync.Mutex
}

func main() {
  buggy := decideBuggy()
  so := new(sharedObject)
  so.task(buggy) // call a function that holds a lock
  // go func() {lock and unlock the same lock, send to done}
  go func(so *sharedObject){
    so.mu.Lock()
    time.Sleep(time.Millisecond*2)
    so.mu.Unlock()
    }(so)

  time.Sleep(time.Millisecond*10)
  fmt.Println("End of main!")
}

func (a *sharedObject) task(buggy bool){
  a.mu.Lock()
  if buggy{
    return
  }
  a.mu.Unlock()
}

func decideBuggy() bool{
  ch := make(chan int)
  go func(ch chan int){
    ch <- 1
    }(ch)
  go func(ch chan int){
    //runtime.Gosched()
    ch <- 2
    }(ch)
  r := <- ch
  <- ch
  if r == 1{
    return true
  } else{
    return false
  }
}

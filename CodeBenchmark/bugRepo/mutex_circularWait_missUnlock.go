package main

// https://github.com/moby/moby/pull/17176
// introduce commit: d295dc66
// Lock() but does not drop it if Condition. This will block other goroutines trying to acquire the same lock

import (
  "fmt"
  "sync"
  "runtime"
  "time"
)

type sharedObject struct{
  content      bool
  mu           sync.Mutex
}

func main() {
  runtime.GOMAXPROCS(1)

  buggy := decideBuggy()

  //doneChan := make(chan bool)
  so := new(sharedObject)
  so.task(buggy)

  // call a function that holds a lock
  // go func() {lock and unlock the same lock, send to done}
  go func(so *sharedObject){
    so.mu.Lock()
    time.Sleep(time.Millisecond*2)
    so.mu.Unlock()
    //doneChan <- true
    }(so)

  /*select{
  case <- time.After(time.Second * 1):
    fmt.Println("TO, End of main!")
  case <-doneChan:
    fmt.Println("End of main!")
  }*/

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
    runtime.Gosched()
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

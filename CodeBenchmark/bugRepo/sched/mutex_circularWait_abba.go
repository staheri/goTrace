package main

import (
  "fmt"
  "sync"
  "runtime"
  "time"
)

// https://github.com/moby/moby/pull/4951
// introduce commit : 81f148be

type sharedObject struct{
  content      bool
  mu           sync.Mutex
}

func main() {
  runtime.GOMAXPROCS(1)
  a := new(sharedObject)
  b := new(sharedObject)
  a.content = true
  b.content = false
  go func(){
    b.and(a)
    }()

  go func(){
    a.and(b)
  }()

  time.Sleep(5*time.Millisecond)
  fmt.Println("End of main!")
}

func (a *sharedObject) and(b *sharedObject){
  a.mu.Lock()
  runtime.Gosched()
  b.mu.Lock()
  a.content = a.content && b.content
  b.content = a.content && b.content
  b.mu.Unlock()
  a.mu.Unlock()
}

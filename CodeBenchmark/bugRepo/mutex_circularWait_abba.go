package main

// https://github.com/moby/moby/pull/4951
// introduce commit : 81f148be

import "fmt"
import "sync"
import "runtime"
//import "time"

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


  fmt.Println("End of main!")
}

func (a *sharedObject) and(b *sharedObject){
  //runtime.Gosched()
  a.mu.Lock()
  //runtime.Gosched()
  b.mu.Lock()
  a.content = a.content && b.content
  b.content = a.content && b.content
  b.mu.Unlock()
  a.mu.Unlock()
}

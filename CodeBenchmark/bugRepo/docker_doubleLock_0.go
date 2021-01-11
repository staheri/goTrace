package main

// https://github.com/docker/docker-ce/commit/a19d2a32c2c0b3097097e1de0a46618915b55660
// Introduce commit : f4e25694c15583ed6ed290aff0c29116f7ed361e
// https://github.com/moby/moby/pull/8929
// https://github.com/moby/moby/issues/8909
//introduce commit: e0339d4b

// This is a double lock bug.

import "fmt"
import "sync"

type sharedObject struct{
  content      bool
  mu           sync.Mutex
}

func main() {
  so := new(sharedObject)
  so.content = true
  go func(){
    so.stop()
    }()

  fmt.Println("End of main!")
}

func (so *sharedObject) stop(){
  so.mu.Lock()
  defer so.mu.Unlock()
  so.kill()
}


func (so *sharedObject) kill(){
  so.mu.Lock()
  so.content = false
  so.mu.Unlock()
}

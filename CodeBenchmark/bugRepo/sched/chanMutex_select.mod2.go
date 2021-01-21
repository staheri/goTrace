package main

import (
  "fmt"
  "sync"
  "runtime"
  "time"
  "math/rand"
)

// commit hash: a69a59ffc7e3d028a72d1195c2c1535f447eaa84
// buggy version: 18768fdc2e76ec6c600c8ab57d2d487ee7877794
// https://github.com/moby/moby/pull/27782
// Different timing in select options and conditionals cause deadlock


func main() {
  runtime.GOMAXPROCS(1)
  cv := sync.NewCond(&sync.Mutex{})

  done := make(chan int)
  eventch := make(chan int)
  closec := make(chan int)

  if Reschedule() {runtime.Gosched()}
  go g0(closec)
  if Reschedule() {runtime.Gosched()}
  go g1(cv,done,eventch,closec)
  if Reschedule() {runtime.Gosched()}
  go g2(cv,done,eventch)

  time.Sleep(time.Millisecond*10)
  fmt.Println("End of main!")
}


func g0(ch chan int){
  //time.Sleep(time.Millisecond*10)
  if Reschedule() {runtime.Gosched()}
  // to trigger the bug, it has to happen first
  //if (Cont()){runtime.Gosched()}
  close(ch)
}

func g1(cv *sync.Cond, done,eventch,closec chan int){
  if Reschedule() {runtime.Gosched()}
  select{
  case <- eventch:
  case <- closec:
    if Reschedule() {runtime.Gosched()}
    cv.L.Lock()
    if Reschedule() {runtime.Gosched()}
    cv.Wait()
    if Reschedule() {runtime.Gosched()}
    cv.L.Unlock()
    if Reschedule() {runtime.Gosched()}
    close(done)
  }
}

func g2(cv *sync.Cond, done,eventch chan int){
  if Reschedule() {runtime.Gosched()}
  select {
  case eventch <- 1:
  case <- done:
  }
  if Reschedule() {runtime.Gosched()}
  cv.Broadcast()
}


const depth=2 // depth

type sharedInt struct{
  n    int
  sync.Mutex
}

var cnt sharedInt

func Reschedule() bool {
  rand.Seed(time.Now().UnixNano())
  if rand.Intn(2) == 1 { // coin toss
    cnt.Lock()
    defer cnt.Unlock()
    if cnt.n < depth{
      cnt.n++
      fmt.Println("Reschedule!")
      return true
    }else{
      return false
    }
  }
  return false
}

package main

import "time"

func main(){
  ch := make(chan int,1)
  go send(ch)
  <-ch
  time.Sleep(1000 * time.Millisecond)
}

func send(ch chan int){
  ch <- 42
}

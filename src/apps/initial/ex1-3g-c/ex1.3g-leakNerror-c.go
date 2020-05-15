package main

import "time"

func main(){
  ch := make(chan int)
  go send(ch)
  go recv(ch)
  time.Sleep(100 * time.Millisecond)
  close(ch)
}

func send(ch chan int){
  ch <- 42
}

func recv(ch chan int){
  <-ch
}

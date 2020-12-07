package main

import "fmt"

func main(){
  ch := make(chan int)
  go send(ch)
  fmt.Println(<-ch)
  close(ch)
}

func send(ch chan int){
  ch <- 42
}

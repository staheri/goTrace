package main

import "fmt"

func main(){
  ch := make(chan int)
  go recv(ch)
  ch <- 42
  close(ch)
}

func recv(ch chan int){
  fmt.Println(<- ch)
}

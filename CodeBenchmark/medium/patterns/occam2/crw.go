package main

import "time"
//import "fmt"
//import "runtime"

//import "strconv"

func main() {
  ch1 := make(chan int)
  ch2 := make(chan int)

  //done := make(chan int)
  go p1(ch1,ch2)
  go p2(ch1,ch2)
  time.Sleep(1*time.Second)
}

func p1(ch1, ch2 chan int){
  <- ch1

  ch2 <- 1

}

func p2(ch1, ch2 chan int){
  <- ch2
  ch1 <- 1
}

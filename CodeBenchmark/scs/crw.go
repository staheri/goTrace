package main

import "time"
//import "fmt"
//import "runtime"

//import "strconv"

func main() {
  ch1 := make(chan int)
  ch2 := make(chan int)

  done := make(chan int)
  go p1(ch1,ch2,done)
  go p2(ch1,ch2,done)
  time.Sleep(1*time.Second)
  <-done
}

func p1(ch1, ch2, done chan int){
  ch2 <- 1
  <- ch1
  done <-1
}

func p2(ch1, ch2, done chan int){
  <- ch2
  ch1 <- 1
  done <-1
}

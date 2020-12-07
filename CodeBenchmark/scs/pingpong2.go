package main

import "time"
//import "fmt"

func main() {
    var Ball int
    table := make(chan int)
    go player(table,1)
    go player(table,2)

    table <- Ball
    time.Sleep(1 * time.Second)
    <-table
}

func player(table chan int, i int) {
  for{
    ball := <- table
    ball++
    //fmt.Printf("Player %d: %d\n",i,ball)
    time.Sleep(200 * time.Millisecond)
    table <- ball
  }
}

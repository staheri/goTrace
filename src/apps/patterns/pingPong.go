package main

import "time"
import "fmt"
import "os"
import "strconv"

func main() {
    var Ball int
    table := make(chan int)
    pls,_ := strconv.Atoi(os.Args[1])
    for cnt := 0 ; cnt < pls ; cnt++{
      go player(table,cnt)
    }
    table <- Ball
    time.Sleep(3 * time.Second)
    <-table
}

func player(table chan int, i int) {
//  for ball := range table{
  for{
    ball := <- table
    ball++
    fmt.Printf("\t(%d) Ball: %d\n",i,ball)
    time.Sleep(200 * time.Millisecond)
    table <- ball
  }
}

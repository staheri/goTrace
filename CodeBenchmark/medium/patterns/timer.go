package main

import "time"
import "fmt"

func timer(d time.Duration, i int, c chan int){
    //c := make(chan int)
    go func() {
        time.Sleep(d)
        c <- i
    }()
    //return c
}

func main() {
  c := make(chan int)
  for i := 0; i < 24; i++ {
    timer(100 * time.Millisecond,i,c)
    fmt.Println(<-c)
  }
}

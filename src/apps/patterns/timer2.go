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
  c1 := make(chan int)
  c2 := make(chan int , 1)
  for i := 0; i < 24; i++ {
    timer(100 * time.Millisecond,i,c1)
    timer(150 * time.Millisecond,i,c2)
    fmt.Println(<-c1)
    fmt.Println(<-c2)
  }
}

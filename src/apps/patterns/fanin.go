package main

import (
    "fmt"
    "time"
)

var bound = 100

func producer(ch chan int, d time.Duration) {
    var i int
    for {
        ch <- i
        i++
        if i>bound{
          break
        }
        time.Sleep(d)
    }
}

func reader(out chan int) {
    for x := range out {
        fmt.Println(x)
    }
}

func main() {
    ch := make(chan int)
    out := make(chan int)
    go producer(ch, 100*time.Millisecond)
    go producer(ch, 150*time.Millisecond)
    go reader(out)
    for i := range ch {
        out <- i
    }
}

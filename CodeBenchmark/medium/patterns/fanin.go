package main

import (
    "fmt"
    "time"
    "sync"
)

var bound = 20

func producer(ch chan int, d time.Duration, wg *sync.WaitGroup) {
    var i int
    for {
        ch <- i
        i++
        if i>bound{
          break
        }
        time.Sleep(d)
    }
    defer wg.Done();
}

func reader(out chan int) {
    for x := range out {
        fmt.Println(x)
    }
}

func main() {
    wg  := &sync.WaitGroup{}
    ch := make(chan int)
    out := make(chan int)
    wg.Add(2)
    go producer(ch, 100*time.Millisecond, wg)
    go producer(ch, 150*time.Millisecond, wg)

    go func() {
      wg.Wait()
      close(ch)
    }()
    go reader(out)
    for i := range ch {
        out <- i
    }
    close(out)
}

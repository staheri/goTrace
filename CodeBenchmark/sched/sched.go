package main

import "fmt"
import "runtime"

func main() {
    ch := make(chan int)
    go func(chan int) {
        for _, v := range []int{1, 2,3,4} {
          runtime.Gosched()
          ch <- v
        }
        close(ch)
    }(ch)

    for v := range ch {
        fmt.Println(v)
    }
    fmt.Println("The channel is closed.")
}

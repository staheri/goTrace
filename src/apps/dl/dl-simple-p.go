package main

import (
    "fmt"
)

func main() {
    c1 := make(chan int)
    fmt.Println("push c1: ")
    go func() {
      c1 <- 10
      close(c1)
    }()
    g1 := 0
    g1 = <- c1
    fmt.Println("get g1: ", g1)
}

package main

import "fmt"
//import _ "net"

func test_a(test_channel chan int) {
  test_channel <- 1
  return
}

func test() {
  test_channel := make(chan int)
  for i := 0; i < 10; i++ {
    go test_a(test_channel)
  }
  for {
    fmt.Println(<-test_channel)
  }
}
func main() {
  test()
}

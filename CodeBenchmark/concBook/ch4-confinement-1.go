package main

import (
	"fmt"
)

func main(){
  data := []int{0,1,2,3}
  loopData := func(handleData chan<- int) {
    defer close(handleData)
    for i := range data {
      handleData <- data[i]
    }
  }
  handleData := make(chan int)
  go loopData(handleData)
  for num := range handleData {
    fmt.Println(num)
  }
}

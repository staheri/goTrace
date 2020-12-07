package main

import "time"
//import "fmt"
//import "runtime"
import "math/rand"

func main() {
    //runtime.GOMAXPROCS(1)
    paper := make(chan int) // 0
    matches := make(chan int) // 1
    tobacco := make(chan int) // 2
    order := make(chan int)
    var materials = []chan int{paper, matches, tobacco}
    go agent(materials, order)

    order <- 1
    go smoker1(tobacco,paper,order)
    go smoker2(paper,matches,order)
    go smoker3(matches,tobacco,order)

    time.Sleep(1*time.Second)
    //<- order



}

func smoker1(tobacco,paper,order chan int) {
  <- tobacco
  <- paper
  //fmt.Println("smoker 1")
  order <- 1
}

func smoker2(paper, matches,order chan int) {
  <- paper
  <- matches
  //fmt.Println("smoker 2")
  order <- 1
}
func smoker3(matches, tobacco,order chan int) {
  <- matches
  <- tobacco
  //fmt.Println("smoker 3")
  order <- 1
}

func agent(materials []chan int, order chan int){
  var randx0,randx1 int

  for{
    //fmt.Println("order")
    <- order
    rand.Seed(time.Now().UnixNano())
    randx0 = rand.Intn(len(materials))
    randx1 = rand.Intn(len(materials))
    for randx1 == randx0{
      randx1 = rand.Intn(len(materials))
    }
    //fmt.Printf("Sending to %d and %d\n",randx0,randx1)
    materials[randx0] <- 1
    //fmt.Printf("Sent to %d\n",randx0)
    materials[randx1] <- 1
    //fmt.Printf("Sent to %d\n",randx1)
  }
}

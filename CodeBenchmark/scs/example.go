package main

import "fmt"
import _"net"
//import "time"


func main() {
  ch1 := make(chan string)
  done := make(chan int)

  go send(ch1)
  go recvAndPrnt(ch1,done)
	//time.Sleep(time.Millisecond*10)
	<-done
  fmt.Println("End of Main")
}

func send(ch1 chan string){
  ch1 <- "Hello SCS!"
}

func recvAndPrnt(ch1 chan string, done chan int){
  fmt.Println(<-ch1)
  //done <- 1
}



package main

func main() {
  ch1 := make(chan string)
  done := make(chan int)

  go send(ch1)
  go recv(ch1,done)

  <- done

}

func send(ch1 chan string){
  ch1 <- "Hello SCS!"
}

func recv(ch1 chan string, done chan int){
  <- ch1
  // done <- 1
}

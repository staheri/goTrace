package main

func main() {
  ch1 := make(chan string)

  go send(ch1)
  go recv(ch1)

}

func send(ch1 chan string){
  ch1 <- "Hello SCS!"
}

func recv(ch1 chan string){
  //<- ch1
}

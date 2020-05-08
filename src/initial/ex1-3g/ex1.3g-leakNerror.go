package main

func main(){
  ch := make(chan int)
  go send(ch)
  go recv(ch)
}

func send(ch chan int){
  ch <- 42
}

func recv(ch chan int){
  <-ch
}

package main

func main(){
  ch := make(chan int,1)
  go send(ch)
  <-ch
}

func send(ch chan int){
  ch <- 42
}

package main

func main(){
  ch := make(chan int,1)
  go send(ch)
  go recv(ch)
}

func send(ch chan int){
  for i:=0;i<10;i++{
    ch <- 42
  }
}

func recv(ch chan int){
  for i:=0;i<10;i++{
    <- ch
  }
}

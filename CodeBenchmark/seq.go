package main

import (
	"sync"
  "fmt"
  "runtime"
)

func main(){
  runtime.GOMAXPROCS(1)
  var wg sync.WaitGroup
  for i:=0;i<5;i++{
    wg.Add(1)
    go func(ii int){
      fmt.Println(ii)
      wg.Done()
    }(i)
  }
  wg.Wait()
}

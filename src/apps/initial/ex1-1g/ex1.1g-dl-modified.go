package main

import(
   "runtime/trace"
   "os"
   "log"
   "fmt"
 )

func main(){
  f,err := os.Create("trace_test.out")
  if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
  trace.Start(f)
  ch := make(chan int)
  go func(){
    ch <- 1
  }()
  fmt.Println(<- ch)

  trace.Stop()
  //close(ch)
}

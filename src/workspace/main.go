package main

import (
  "flag"
  "fmt"
  "os"
  //"strings"
  "log"
  //"github.com/staheri/goTrace/trace"
  "trace"
)


func main(){
  flag.Parse()
  args := flag.Args()
  fmt.Println(args[0], " - ", args[1])
  f,err := os.Open(args[0])
  if err != nil{
    log.Fatal(err)
  }
  defer f.Close()
  events, err := trace.Parse(f,args[1])
  if err != nil{
    log.Fatal(err)
  }
  trace.Print(events.Events)
  /*for _,e := range events{

  }*/


}

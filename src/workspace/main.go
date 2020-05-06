package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "log"
  "github.com/staheri/goTrace/trace"
)


func main(){
  flag.Parse()
  args := flag.Args()
  f,err := os.Open(args[0])
  if err != nil{
    log.Fatal(err)
  }
  defer f.Close()
  events, err := trace.Parse(f)
  if err != nil{
    log.Fatal(err)
  }
  trace.Print(events)
  /*for _,e := range events{

  }*/


}

package main

import (
  "flag"
  "fmt"
  _"os"
  _"log"
  _"trace"
  "util"
  "instrument"
  "analyze"
  _"sort"
  _"bytes"
  _"path"
)


func main(){
  flag.Parse()
  args := flag.Args()
  fmt.Println(args[0])
  //f,err := os.Open(args[0])
  //if err != nil{
    //log.Fatal(err)
  //}
  //defer f.Close()
  //events, err := trace.Parse(f,args[1])
  var src instrument.EventSource
  src = instrument.NewNativeRun(args[0])
  events, err := src.Events()
	if err != nil {
		panic(err)
	}
  //trace.Print(events)
  //Procs(events)
  //Grtns(events.Events)
  //Grtns(events)
  objCtg := "grtn"
  context, err := analyze.Convert(events,objCtg,"101010",5)
  if err != nil{
    panic(err)
  }
  util.DispAtrMap(context,objCtg)
  //util.GroupGrtns(events)
}

package main

import (
  "flag"
  "fmt"
  _"os"
  //"strings"
  _"log"
  //"github.com/staheri/goTrace/trace"

  "trace"
  "util"
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
  var src util.EventSource
  src = util.NewNativeRun(args[0])
  events, err := src.Events()
	if err != nil {
		panic(err)
	}
  //trace.Print(events)
  Procs(events)
  //Grtns(events.Events)
  Grtns(events)
}

func Procs(events []*trace.Event) {
  m := make(map[int][]*trace.Event)
  for _,e := range events{
    if _,ok := m[e.P]; ok{
      m[e.P] = append(m[e.P],e)
    } else{
      m[e.P] = append(m[e.P],e)
    }
  }
  util.ToPTable(m)
}

func Grtns(events []*trace.Event) {
  m := make(map[uint64][]*trace.Event)
  for _,e := range events{
    if _,ok := m[e.G]; ok{
      m[e.G] = append(m[e.G],e)
    } else{
      m[e.G] = append(m[e.G],e)
    }
  }
  util.ToGTable(m)

}

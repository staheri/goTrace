package main

import (
  "flag"
  "fmt"
  "os"
  //"strings"
  "log"
  //"github.com/staheri/goTrace/trace"
  "trace"
  "util"
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
  src = NewNativeRun(args[0])
  events, err := src.Events()
	if err != nil {
		panic(err)
	}
  //trace.Print(events.Events)
  //Procs(events.Events)
  Grtns(events.Events)
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
  for k,v := range m{
    fmt.Println(k)
    fmt.Println("***")
    trace.Print(v)
    fmt.Println("***")
  }
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
  for k,v := range m{
    fmt.Println(k)
    fmt.Println("___")
    trace.Print(v)
    fmt.Println("___")
  }
}

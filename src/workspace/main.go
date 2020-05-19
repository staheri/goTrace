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

const dir = "/Users/saeed/goTrace/"
const outpath = dir+"traces/"
const inpath = dir+"/src/apps/"

func main(){
  appPtr := flag.String("app", "initial/ex1-2g/ex1.2g-ok.go", "Target application (*.go)")
  objPtr := flag.String("obj", "grtn", "Object:[grtn,proc,chan]")
  atrPtr := flag.String("atr", "110000", "Attributes: a bitstring showing 1/0 event groups :\n\t\t\"GoRoutine,Channel,Process,GCmem,Syscall,Other\"")
  atrModePtr := flag.Int("atrMode", 0, util.AttributeModesDescription())
  flag.Parse()

  fmt.Println("Analyzing ", inpath+(*appPtr), "...")
  var src instrument.EventSource

  src = instrument.NewNativeRun(inpath+(*appPtr))
  events, err := src.Events()
	if err != nil {
		panic(err)
	}
  //trace.Print(events)
  //Procs(events)
  //Grtns(events.Events)
  //Grtns(events)
  //context, err := analyze.Convert(events,*objPtr,*atrPtr,*atrModePtr)
  //if err != nil{
  //  panic(err)
  //}
  analyze.TestDB()
  //util.DispAtrMap(context,*objPtr)
  //util.WriteContext(outpath+util.AppName(*appPtr), *objPtr , *atrPtr , context, *atrModePtr )
  //util.GroupGrtns(events)
}

package main

import (
  "flag"
  "fmt"
  _"os"
  _"log"
  _"trace"
   "util"
  "instrument"
  _"cl"
  _"sort"
  _"bytes"
  _"path"
  "db"
)

const dir = "/Users/saeed/goTrace/"
const outpath = dir+"traces/"
const inpath = dir+"/CodeBenchmark/"
const datapath = dir+"DataBenchmark/medium/"

func main(){
  appPtr := flag.String("app", "small/initial/ex1-2g/ex1.2g-ok.go", "Target application (*.go)")
  tout   := flag.Int("to", -1, "Timeout for deadlocks")
  //dbName   := flag.String("dbName", "pingPongX9", "Table Name")
  //outName   := flag.String("outName", "medium/test.py", "OutName")
  //filterPtr := flag.String("filter", "CHNL", "FILTERS: CHNL, GCMM, GRTN, MISC, MUTX, PROC, SYSC, WGRP ")

  //objPtr := flag.String("obj", "grtn", "Object:[grtn,proc,chan]")
  //atrPtr := flag.String("atr", "1110000", "Attributes: a bitstring showing 1/0 event groups :\n\t\t\"GoRoutine,Channel,Process,GCmem,Syscall,Other\"")
  //atrModePtr := flag.Int("atrMode", 0, util.AttributeModesDescription())


  // * main block for instrumentation and collecting traces
  flag.Parse()
  fmt.Println("Analyzing ", inpath+(*appPtr), "...")
  var src instrument.EventSource

  src = instrument.NewNativeRun(inpath+(*appPtr),(*tout))
  events, err := src.Events()
	if err != nil {
		panic(err)
	}



  //trace.Print(events)
  //Procs(events)
  //Grtns(events.Events)
  //Grtns(events)
  //
  //context, err := cl.Convert(events,*objPtr,*atrPtr,*atrModePtr)
  //if err != nil{
    //panic(err)
  //}
  //analyze.TestDB()
  //cl.DispAtrMap(context,*objPtr)
  //cl.WriteContext(outpath+util.AppName(*appPtr), *objPtr , *atrPtr , context, *atrModePtr )


  //cl.GroupGrtns(events)
  dbName := db.Store(events,util.AppName(*appPtr))

  /* word2vec block
  filters := []string{"all","CHNL", "GCMM", "GRTN", "MISC", "MUTX", "PROC", "SYSC", "WGRP"}
  for _,filt := range(filters){
    db.WriteData(dbname,datapath,filt,11)
    db.WriteData(dbname,datapath,filt,21)
  }
  */

  //db.WriteData(dbname,datapath,(*filterPtr),11)
  //db.Ops()

  //db.ToFile(dbName)
}

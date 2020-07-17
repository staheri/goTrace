package main

import (
  "flag"
  "fmt"
  "util"
  "instrument"
  "strings"
  "db"
  "os"
)

const WORD_CHUNK_LENGTH = 11
var CLOUTPATH = os.Getenv("GOPATH")+"/traces/clx"
//var validCategories = []string{"CHNL", "GCMM", "GRTN", "MISC", "MUTX", "PROC", "SYSC", "WGRP"}



var (
  flagCmd     string
  flagOut     string
  flagSrc     string
  flagX       string
  flagN       int
  flagBase    string
  flagTO      int
  flagApp     string
  flagArgs    []string
  dbName      string
  validCategories = []string{"CHNL", "GCMM", "GRTN", "MISC", "MUTX", "PROC", "SYSC", "WGRP"}
)




func main(){
  // Read flags
  parseFlags()

  // Obtain dbName
  dbName = dbPointer()

  fmt.Printf("DB Name: %s\n",dbName)

  switch flagCmd {
  case "word":
    for _,arg := range(flag.Args()){
      // For now, only one filter is allowed at a time
      if len(strings.Split(arg,",")) != 1{
        panic("Currently more than one filter is not allowed!")
      }
      db.WordData(dbName,flagOut,arg,WORD_CHUNK_LENGTH)
      //for _,e := range(tl){
        // TODO: Make db.WriteData compatible with combination of filters
      //}
    }

  case "hac":
    for _,arg := range(flag.Args()){
      tl := strings.Split(arg,",")
      db.CLOperations(dbName,CLOUTPATH,flagOut,tl...)
    }

  case "rr":
    for _,arg := range(flag.Args()){
      if len(strings.Split(arg,",")) != 1{
        panic("For rr, only one category is allowed")
      }
      switch arg {
      case "CHNL":
        db.ChannelReport(dbName)
      case "MUTX":
        db.MutexReport(dbName)
        db.RWMutexReport(dbName)
      case "WGRP":
        db.WaitingGroupReport(dbName)
      default:
        panic("Wrong category for rr!")
      }
    }

  case "rg":
    db.ResourceGraph(dbName,flagOut)
  case "diff":
    baseDBName := db.Ops("latest",util.AppName(flagBase),"0")
    for _,arg := range(flag.Args()){
      tl := strings.Split(arg,",")
      db.DIFF(dbName,baseDBName,CLOUTPATH,flagOut,tl...)
    }
  case "dineData":
    db.DineData(dbName, flagOut+"/ch-chid", flagN, true,true) // channel events only + channel ID
    db.DineData(dbName, flagOut+"/ch", flagN, true,false) // channel events only
    db.DineData(dbName, flagOut+"/all-chid", flagN, false,true) // all events + channel ID (for channel events)
    db.DineData(dbName, flagOut+"/all", flagN, false,false) // all events
  }
}


// Parse flags, execute app & store traces (if necessary), return app database handler
func parseFlags() (){
  srcDescription := "native: execute the app and collect from scratch, latest: retrieve data from latest execution, x: retrieve data from specific execution (requires -x option)"
  // Parse flags
  flag.StringVar(&flagCmd,"cmd","","Commands: word, cl, rr, rg, diff")
  flag.StringVar(&flagBase,"base","","Base for \"diff\" command (latest)")
  flag.StringVar(&flagOut,"outdir","","Output directory to write words and/or reports")
  flag.StringVar(&flagSrc,"src","latest",srcDescription)
  flag.StringVar(&flagX,"x","0","Execution version stored in database")
  flag.IntVar(&flagN,"n",0,"Number of philosophers for dineData command")
  flag.StringVar(&flagApp,"app","","Target application (*.go)")
  flag.IntVar(&flagTO,"to",-1,"Timeout for deadlocks")

  flag.Parse()

  // Check cmd
  if flagCmd != "word" && flagCmd != "hac" && flagCmd != "rr" && flagCmd != "rg" && flagCmd != "diff" && flagCmd != "dineData"{
    util.PrintUsage()
    fmt.Printf("flagCMD: %s\n",flagCmd)
    panic("Wrong command")
  }

  // Check Outdir
  if flagOut == "" {
    util.PrintUsage()
    panic("Outdir required")
  }

  // Check src
  if flagSrc != "native" && flagSrc != "latest" && flagSrc != "x"{
    util.PrintUsage()
    panic("Wrong source")
  }

  // Check app
  if flagApp == "" {
    util.PrintUsage()
    panic("App required")
  }

  for _,arg := range(flag.Args()){
    tl := strings.Split(arg,",")
    for _,e := range(tl){
      if ! util.Contains(validCategories,e){
        panic("Invalid category: "+e)
      }
    }
  }

  if flagCmd == "diff" && flagBase == ""{
    util.PrintUsage()
    panic("Undefined base for diff command!")
  }

  if flagCmd == "dineData" && flagN == 0{
    util.PrintUsage()
    panic("Wrong N for dineData!")
  }

  flagArgs = flag.Args()
}

// Find appropriate DB handler According to options
func dbPointer() (dbName string){

  switch flagSrc {

  case "native":
    fmt.Println("Analyzing ", flagApp, "...")
    var src instrument.EventSource
    src = instrument.NewNativeRun(flagApp,flagTO)
    events, err := src.Events()
  	if err != nil {
  		panic(err)
  	}
    dbName = db.Store(events,util.AppName(flagApp))
    return dbName

  case "latest", "x":
    dbName = db.Ops(flagSrc, util.AppName(flagApp), flagX )
    return dbName

  default:
    panic("DbPointer not available!")
  }
}

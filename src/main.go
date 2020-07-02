package main

import (
  "flag"
  "fmt"
  "util"
  "instrument"
)

const WORD_CHUNK_LENGTH = 11
const validCategories =[]string{"CHNL", "GCMM", "GRTN", "MISC", "MUTX", "PROC", "SYSC", "WGRP"}



var (
  flagCmd     string
  flagOut     string
  flagSrc     string
  flagX       string
  flagTO      int
  flagApp     string
  flagArgs    []string
  dbName      string
)




func main(){
  // Read flags
  parseFlags()

  // Obtain dbName
  //dbName = dbPointer()

  /*
  switch flagCmd {
  case "word":
    for _,arg := range(flag.Args()){
      // For now, only one filter is allowed at a time
      if len(strings.Split(arg,",")) != 1{
        panic("Currently more than one filter is not allowed!")
      }
      db.WriteData(dbname,flagOut,filt,WORD_CHUNK_LENGTH)
      //for _,e := range(tl){
        // TODO: Make db.WriteData compatible with combination of filters
      //}
    }

  case "cl":
    for _,arg := range(flag.Args()){
      tl := strings.Split(arg,",")
      db.FormalContext(dbName,outpath,tl)
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
  }*/
}


// Parse flags, execute app & store traces (if necessary), return app database handler
func parseFlags() (){
  // Parse flags
  flag.StringVar(&flagCmd,"cmd","","Commands: word, cl, rr")
  flag.StringVar(&flagOut,"outdir","","Output directory to write words")
  flag.StringVar(&flagSrc,"src","latest",srcDescription)
  flag.StringVar(&flagX,"x","0","Execution version stored in database")
  flag.StringVar(&flagApp,"app","","Target application (*.go)")
  flag.IntVar(&flagTO,"to",-1,"Timeout for deadlocks")

  flag.Parse()

  // Check cmd
  if flagCmd != "word" || flagCmd != "cl" || flagCmd != "rr" {
    printUsage()
    panic("Wrong command")
  }

  // Check Outdir
  if flagOut == "" {
    printUsage()
    panic("Outdir required")
  }

  // Check src
  if flagSrc != "native" || flagSrc != "latest" || flagSrc != "x"{
    printUsage()
    panic("Wrong source")
  }

  // Check app
  if flagApp == "" {
    printUsage()
    panic("App required")
  }

  for _,arg := range(flag.Args()){
    tl := strings.Split(arg,",")
    for _,e := range(tl){
      if !util.Contains(validCategories,e){
        panic("Invalid category: "+e)
      }
    }
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

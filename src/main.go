package main

import (
	"db"
	"flag"
	"fmt"
	"instrument"
	"os"
	"schedtest"
	"strings"
	"util"
)

const WORD_CHUNK_LENGTH = 11

var CLOUTPATH = os.Getenv("GOPATH") + "/traces/clx"

var (
	flagCmd            string
	flagOut            string
	flagSrc            string
	flagX              string
	flagCons           int
	flagAtrMode        int
	flagN              int
	flagBase           string
	flagTO             int
	flagDepth          int
	flagIter           int
	flagApp            string
	flagArgs           []string
	dbName             string
	validCategories    = []string{"CHNL", "GCMM", "GRTN", "MISC", "MUTX", "PROC", "SYSC", "WGCV", "SCHD", "BLCK"}
	validPrimeCmds     = []string{"word", "hac", "rr", "rg", "diff", "dineData", "cleanDB", "dev", "hb", "gtree", "cgraph", "resg"}
	validTestSchedCmds = []string{"test"}
	validSrc           = []string{"native", "x", "latest", "schedTest"}
)

func main() {
	// Read flags
	parseFlags()

	if flagSrc == "schedTest" {
		handleSchedTestCommands()
	} else {
		myapp := instrument.NewAppExec(flagApp, flagSrc, flagX, flagTO)
		dbn, err := myapp.DBPointer()
		if err != nil {
			panic(err)
		}
		myapp.DBName = dbn
		handlePrimaryCommands(myapp.DBName)
		fmt.Println(myapp.ToString())
	}
}

// Parse flags, execute app & store traces (if necessary), return app database handler
func parseFlags() {
	srcDescription := "native: execute the app and collect from scratch, latest: retrieve data from latest execution, x: retrieve data from specific execution (requires -x option)"
	// Parse flags
	flag.StringVar(&flagCmd, "cmd", "", "Commands: word, cl, rr, rg, diff")
	flag.StringVar(&flagBase, "baseX", "0", "Base execution for \"diff\" or \"schedTrace\" command")
	flag.StringVar(&flagOut, "outdir", "", "Output directory to write words and/or reports")
	flag.StringVar(&flagSrc, "src", "latest", srcDescription)
	flag.StringVar(&flagX, "x", "", "Execution version stored in database")
	flag.IntVar(&flagN, "n", 0, "Number of philosophers for dineData command")
	flag.IntVar(&flagCons, "cons", 1, "Number of consecutive elements for HAC & DIFF")
	flag.IntVar(&flagAtrMode, "atrmode", 0, "Modes for HAC & DIFF")
	flag.StringVar(&flagApp, "app", "", "Target application (*.go)")
	flag.IntVar(&flagTO, "to", 0, "Timeout for deadlocks")
	flag.IntVar(&flagDepth, "depth", 0, "Max depth for rescheduling")
	flag.IntVar(&flagIter, "iter", 2, "Testing iteration")

	flag.Parse()

	// Check src validity
	if !util.Contains(validSrc, flagSrc) {
		util.PrintUsage()
		panic("Wrong source")
	}

	// Check prime cmd validity
	if flagSrc != "schedTest" && !util.Contains(validPrimeCmds, flagCmd) {
		util.PrintUsage()
		fmt.Printf("flagCMD: %s\n", flagCmd)
		panic("Wrong prime command")
	}

	// Check prime cmd validity
	if flagSrc == "schedTest" && !util.Contains(validTestSchedCmds, flagCmd) {
		util.PrintUsage()
		fmt.Printf("flagCMD: %s\n", flagCmd)
		panic("Wrong schedTest command")
	}

	// Check Outdir
	if flagOut == "" {
		util.PrintUsage()
		panic("Outdir required")
	}

	// Check app
	if flagApp == "" {
		util.PrintUsage()
		panic("App required")
	}

	// Check validity of categories
	for _, arg := range flagArgs {
		tl := strings.Split(arg, ",")
		for _, e := range tl {
			if !util.Contains(validCategories, e) {
				panic("Invalid category: " + e)
			}
		}
	}

	// diff command needs a base
	if flagCmd == "diff" && flagBase == "" {
		util.PrintUsage()
		panic("Undefined base for diff command!")
	}

	// dineData command needs N
	if flagCmd == "dineData" && flagN == 0 {
		util.PrintUsage()
		panic("Wrong N for dineData!")
	}

	// x command needs X value
	if flagSrc == "x" && flagX == "" {
		util.PrintUsage()
		panic("Needs X value!")
	}

	flagArgs = flag.Args()
}

// handle primary commands
func handlePrimaryCommands(dbName string) {
	switch flagCmd {
	case "word":
		for _, arg := range flagArgs {
			// For now, only one filter is allowed at a time
			if len(strings.Split(arg, ",")) != 1 {
				panic("Currently more than one filter is not allowed!")
			}
			db.WordData(dbName, flagOut, arg, WORD_CHUNK_LENGTH)
			//for _,e := range(tl){
			// TODO: Make db.WriteData compatible with combination of filters
			//}
		}

	case "hac":
		if len(flagArgs) > 0 {
			for _, arg := range flagArgs {
				tl := strings.Split(arg, ",")
				db.HAC(dbName, CLOUTPATH, flagOut, flagCons, flagAtrMode, tl...)
			}
		} else {
			var emptyList []string
			db.HAC(dbName, CLOUTPATH, flagOut, flagCons, flagAtrMode, emptyList...)
		}

	case "rr":
		for _, arg := range flagArgs {
			if len(strings.Split(arg, ",")) != 1 {
				panic("For rr, only one category is allowed")
			}
			switch arg {
			case "CHNL":
				db.ChannelReport(dbName, flagOut)
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
		for _, arg := range flagArgs {
			tl := strings.Split(arg, ",")
			db.SwimLanes(dbName, flagOut, tl...)
		}
	case "diff":
		baseDBName := db.Ops("x", util.AppName(flagBase), "13")
		for _, arg := range flagArgs {
			tl := strings.Split(arg, ",")
			db.JointHAC(dbName, baseDBName, CLOUTPATH, flagOut, flagCons, flagAtrMode, tl...)
		}
	case "dineData":
		db.DineData(dbName, flagOut+"/ch-chid", flagN, true, true)   // channel events only + channel ID
		db.DineData(dbName, flagOut+"/ch", flagN, true, false)       // channel events only
		db.DineData(dbName, flagOut+"/all-chid", flagN, false, true) // all events + channel ID (for channel events)
		db.DineData(dbName, flagOut+"/all", flagN, false, false)     // all events
	case "cleanDB":
		db.Ops("clean all", "", "0")
	case "hb":
		fmt.Println("HB DBNAME:", dbName)
		for _, arg := range flagArgs {
			tl := strings.Split(arg, ",")
			hbtable := db.HBTable(dbName, tl...)
			db.HBLog(dbName, hbtable, flagOut, true)
			fmt.Println("****")
			db.HBLog(dbName, hbtable, flagOut, false)
		}
	case "gtree":
		db.Gtree(dbName, flagOut)
	case "cgraph":
		db.ChannelGraph(dbName, flagOut)
	case "resg":
		db.ResourceGraph(dbName, flagOut)
	case "dev":
		/*for _,arg := range(flagArgs){
		  tl := strings.Split(arg,",")
		  hbtable := db.HBTable(dbName,tl...)
		  db.Dev(dbName,hbtable, flagOut)
		}*/
		db.Checker(dbName)
		//fmt.Println(dbName)
		//db.Gtree(dbName,flagOut)
		//db.Histogram(10,dbName)

		//db.HBLog(dbName,flagOut,true)
		//fmt.Println("****")
		//db.HBLog(dbName,flagOut,false)
	}
}

// handle schedTest commands
func handleSchedTestCommands() {
	mytest := schedtest.SchedTest(flagApp, flagSrc, flagX, flagTO, flagDepth, flagIter)
	fmt.Println(mytest.ToString())
}

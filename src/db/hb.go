package db

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"log"
	"util"
	"os"
)

func aspect2string(aspects ...string) (ret string){
	if len(aspects) != 0{
		ret = ""
		for i,asp := range aspects{
			if i < len(aspects) - 1{
				ret = ret +asp+"_"
			} else{
				ret = ret +asp
			}
		}
	} else{
		ret = "all"
	}
	return ret
}

func asp2int(asp string) (ret int){
	if asp == "CHNL"{
		ret = catCHNL
	}else if asp == "GRTN"{
		ret = catGRTN
	}else if asp == "MUTX"{
		ret = catMUTX
	}else if asp == "SYSC"{
		ret = catSYSC
	}else if asp == "WGCV"{
		ret = catWGCV
	}else if asp == "PROC"{
		ret = catPROC
	}else if asp == "MISC"{
		ret = catMISC
	}else if asp == "GCMM"{
		ret = catGCMM
	}else if asp == "BLCK"{
		ret = catBLCK
	}else if asp == "SCHD"{
		ret = catSCHD
	}else{
		panic("Wrong Aspect")
	}
	return ret
}

func isWhite(event string, aspects ...string)(ret bool){
	fmt.Println(aspects)
	if len(aspects) != 0{
		ret = false
		for _,asp := range aspects{
			 aspID := asp2int(asp)
			 //fmt.Println("Check if "+event+" is in "+asp+ " (aspID:"+strconv.Itoa(aspID)+")")
			 if util.Contains(ctgDescriptions[aspID].Members, event){
				 ret = true
				 break
			 }
		}
	} else{
		ret = true
	}
	return ret
}

// Take sequence of events, create a new DB Schema and insert events into tables
func HBTable(dbName string,aspects ...string) (HBTableName string) {

	if len(aspects) == 0{
		return "Events"
	}

	var err             	      error
	//var res                   sql.Result
	var q, event, _ev   	      string
	var p,eid       		    		int
	var g,_rid	       		   		uint64
	var offset 									int64
	var ts      								int64
	var predG,predClk,linkoff	  sql.NullInt64
	var rclock,rval,reid   	    sql.NullInt64
	var rid,srcLine        	    sql.NullString
	//var buff, output          string


	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	}else{
		log.Println("HBTable: Initial connection established")
	}
	defer db.Close()

	if aspect2string(aspects...) != "all"{
		HBTableName = `Events_`+aspect2string(aspects...)
	} else{
		HBTableName = `Events`
	}


	res,err := db.Query("SHOW TABLES LIKE \""+HBTableName+"\"")
	check(err)
	if res.Next(){
		// table exist
		log.Println("HBTable: Table ", HBTableName ," exists & returns")
		return HBTableName
	}

	stmt := `CREATE TABLE `+HBTableName+` (
					id int NOT NULL AUTO_INCREMENT,
					type varchar(255) NOT NULL,
					vc int NOT NULL,
					ts bigint NOT NULL,
					off int NOT NULL,
					g int NOT NULL,
					p int NOT NULL,
					linkoff bigint,
					predG int,
					predClk int,
					rid varchar(255),
					reid int,
					rval bigint,
					rclock int,
					src varchar(255),
					PRIMARY KEY (id)
					);`
	// create new table
	log.Printf("HBTable: Creating table %v ... \n",HBTableName)
	_,err = db.Exec(stmt)
	if err != nil {
		panic(err)
	}

	insertStmt, err := db.Prepare("INSERT INTO "+HBTableName+" (type, vc ,ts, off, g, p, linkoff, predG, predClk, rid, rval, rclock, src) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	check(err)
	defer insertStmt.Close()


	// Init vector clocks
	msgs          := make(map[msgKey]eventPredecessor) // storing (to be) pred of a send/recv
	links         := make(map[int64]eventPredecessor) // storing (to be) pred of an event
	// Resource clocks
	localClock    := make(map[uint64]uint64) // vc[g]           = local clock
	//chanClock     := make(map[uint64]uint64) // chansClock[cid] = channel clock
	//wgClock       := make(map[uint64]uint64) // wgsClock[cid]   = wg clock
	//mutexClock    := make(map[uint64]uint64) // mutexClock[cid] = mutex clock

	predG    = sql.NullInt64{}
	predClk  = sql.NullInt64{}
	//rid      := sql.NullString{}
	//rval     :=  sql.NullInt32{}
	//rclock   :=  sql.NullInt32{}


	q = "SELECT id,type,ts,offset,g,p,linkoff,rid,reid,rval,rclock,src FROM Events ORDER BY ts;"
	res,err = db.Query(q)
	check(err)
	defer res.Close()
	for res.Next(){
		err = res.Scan(&eid,&_ev,&ts,&offset,&g,&p,&linkoff,&rid,&reid,&rval,&rclock,&srcLine)
		check(err)

		event = _ev[2:] // event is _ev without "Ev*"

		// fresh values for each event
		predG    = sql.NullInt64{}
		predClk  = sql.NullInt64{}

		if !isWhite(_ev,aspects...){
			continue
		}

		// Assign local logical clock
		if _,ok := localClock[g];ok{
			localClock[g] = localClock[g] + 1
		} else{
			localClock[g] = 1
		}

		// Check category of events\
	 if util.Contains(ctgDescriptions[catCHNL].Members, _ev){
			// CHNL event
			// Assign chanClock
			// Assign rid, rval, rclock
			// ChSend? set predG , rval = value
			// ChRecv and MSG[key]? use predG,predClk, else: null,null
			// ChMake/Close? rval = null
			if s, err := strconv.ParseUint(rid.String[1:], 10, 64); err == nil {
				_rid = s
			}else{
				_rid = 0
				panic("_rid")
			}
			if event == "ChRecv"{
				log.Printf("Recv\n")
				//rval = sql.NullInt32{Valid:true, Int32: int32(e.Args[2])} // message val
				if vv,ok := msgs[msgKey{_rid,uint64(reid.Int64),uint64(rval.Int64)}] ; ok{
					// A matching sent is found for the recv
					//fmt.Printf("\tMatching sent is found\n")
					predG    = sql.NullInt64{Valid:true, Int64: int64(vv.g)}
					predClk  = sql.NullInt64{Valid:true, Int64: int64(vv.clock)}
				}else{
					// Receiver without a matching sender
					// might be found later
					//fmt.Printf("\tNo matching (store null for predG,PredCLK)\n\tStore msg[%v,%v,%v] = g:%v localClock:%v",_rid,uint64(reid.Int64),uint64(rval.Int64),g, localClock[g])
					msgs[msgKey{_rid,uint64(reid.Int64),uint64(rval.Int64)}] = eventPredecessor{g, localClock[g]}
					predG = sql.NullInt64{}
					predClk = sql.NullInt64{}
				}
			}else{
				// ChMake, ChSend, ChClose
				if event == "ChSend"{
					// Set Predecessor for a receive (key to the event: {cid, eid, val})
					log.Printf("HBTable: Send\n")
					if vv,ok := msgs[msgKey{_rid,uint64(reid.Int64),uint64(rval.Int64)}] ; ok{
						log.Printf("\tMatching recv is found\n")
						predG    = sql.NullInt64{Valid:true, Int64: int64(vv.g)}
						predClk  = sql.NullInt64{Valid:true, Int64: int64(vv.clock)}
					} else{ // a send for this particular message has been stored before
						msgs[msgKey{_rid,uint64(reid.Int64),uint64(rval.Int64)}] = eventPredecessor{g, localClock[g]}
						log.Printf("\tNo matching (store null for predG,PredCLK)\n\tStore msg[%v,%v,%v] = g:%v localClock:%v",_rid,uint64(reid.Int64),uint64(rval.Int64),g, localClock[g])
						predG = sql.NullInt64{}
						predClk = sql.NullInt64{}
						//panic("Previously stored as sent!")
					}
				}else{ // ChMake. ChClose
				//	rval = sql.NullInt32{}
					log.Printf("HBTable: Make/Close (null predG predClk)\n")
					predG = sql.NullInt64{}
					predClk = sql.NullInt64{}
				}
			}
		} else if linkoff.Valid{
			// Set Predecessor for an event (key to the event: TS)
			//fmt.Printf("Source: %s\n",e)
			//fmt.Printf("LINK: %s\n",e.Link)
			if _,ok := links[linkoff.Int64] ; !ok{
				links[linkoff.Int64] = eventPredecessor{g, localClock[g]}
			} else{ // the link of current event has been linked to another event before
				panic("Previously linked to another event!")
			}
		} else{ // does not fall into any category
			predG     = sql.NullInt64{}
			predClk   = sql.NullInt64{}
		}

		// So far, all predecessor values are set,
		// all resource values obtained from main events table
		// if a recv has found a sender, it is all set
		// Now only check if the current event has a predecessor. If so: set predG, set predClk
		// otherise: everything is null
		if vv,ok := links[offset]; ok{
			// Is there a possibility that this event has resource other than G?
			// No. Events with predecessor links only have G resource
			predG    = sql.NullInt64{Valid:true, Int64: int64(vv.g)}
			predClk  = sql.NullInt64{Valid:true, Int64: int64(vv.clock)}
			//rval     = sql.NullInt32{}
			//rclock   = sql.NullInt32{}
		}
		_,err := insertStmt.Exec(_ev,
															 strconv.Itoa(int(localClock[g])),
															 strconv.Itoa(int(ts)),
															 strconv.Itoa(int(offset)),
															 strconv.FormatUint(g,10),
															 strconv.Itoa(p),
															 linkoff,
															 predG,
															 predClk,
															 rid,
															 rval,
															 rclock,
															 srcLine)

		check(err)
	}
	return HBTableName
}


//func HBLog(dbName, outdir string, resourceView bool, aspects ...string){
func HBLog(dbName, hbtable, outdir string, resourceView bool){
	fmt.Println("HBLog")
	// Variables
	var q, event, _ev         string
	var event1                string
	//var _arg,_val        			string
	var g,logclock,eid   			int
	var predG,predClk,rclock  sql.NullInt32
	var rval 				     			sql.NullInt64
	var rid, srcl 	     			sql.NullString
	var buff, output,srcLine  string


	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("HBLog: Conntected to ",dbName)
	}
	defer db.Close()


	// make sure
	outdir = outdir + "/"

	q = "SELECT id,type,g,vc,predG,predClk,rid,rval,rclock,src FROM "+hbtable+" ORDER BY ts;"
	res, err := db.Query(q)
	check(err)
	defer res.Close()

	if resourceView{
		output = outdir + dbName+hbtable+"_rlog.txt"
		f,err := os.Create(output)
		if err != nil{
			log.Fatal(err)
		}

		buff = "(?<event>.*) [(](?<host>\\S*)[)] (?<clock>{.*})\n"
		f.WriteString(buff)
		buff = "\n\n"
		f.WriteString(buff)

		for res.Next(){
			err = res.Scan(&eid,&_ev,&g,&logclock,&predG,&predClk,&rid,&rval,&rclock,&srcl)
			check(err)

			event = _ev[2:]
			event1 = _ev[2:]
			//fmt.Printf("OUTOPS-EVENT: %v\n",event)
			if event1 == "WgAdd"{
				if rval.Valid && rval.Int64 > 0{
					event = event + "[val:"+strconv.Itoa(int(rval.Int64))+"]"
				}else{
					event = event + "[val:-]"
				}
			}else if event1 == "ChRecv" || event1 =="ChSend"{
				if rval.Valid{
					event = event + "[val:"+strconv.Itoa(int(rval.Int64))+"]"
				}else{
					event = event + "[val:-]"
				}
			}
			if srcl.Valid{
				srcLine = srcl.String
			}else{
				srcLine = ""
			}
			if rid.Valid && rclock.Valid {
				//fmt.Printf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
				buff = fmt.Sprintf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
				f.WriteString(buff)

				//fmt.Printf("%v@%v (%v) {\"G%v\": %v,\"%v\": %v}\n","_"+event,srcLine,rid.String,g,logclock,rid.String,rclock.Int32)
				buff = fmt.Sprintf("%v@%v (%v) {\"G%v\": %v,\"%v\": %v}\n","_"+event,srcLine,rid.String,g,logclock,rid.String,rclock.Int32)
				f.WriteString(buff)
			}else{
				//panic("KIR")

				if predG.Valid {
					if g == int(predG.Int32){
						//happening on same goroutine, just GID is enough
						//fmt.Printf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
						buff = fmt.Sprintf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
						f.WriteString(buff)

					} else{
						//fmt.Printf("%v@%v (G%v) {\"G%v\": %v, \"G%v\": %v }\n",event,srcLine,g,g,logclock,predG.Int32,predClk.Int32)
						buff = fmt.Sprintf("%v@%v (G%v) {\"G%v\": %v, \"G%v\": %v }\n",event,srcLine,g,g,logclock,predG.Int32,predClk.Int32)
						f.WriteString(buff)
					}
				} else{
					//fmt.Printf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
					buff = fmt.Sprintf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
					f.WriteString(buff)
				}
			}
		}
		f.Close()
	}else{
		output = outdir + dbName+hbtable+"_glog.txt"
		f,err := os.Create(output)
		if err != nil{
			log.Fatal(err)
		}
		buff = "(?<event>.*) [(](?<host>\\S*)[)] (?<clock>{.*})\n"
		f.WriteString(buff)
		buff = "\n"
		f.WriteString(buff)
		defer f.Close()

		for res.Next(){
			err = res.Scan(&eid,&_ev,&g,&logclock,&predG,&predClk,&rid,&rval,&rclock,&srcl)
			check(err)

			event = _ev[2:]
			event1 = _ev[2:]
			if event1 == "WgAdd"{
				if rid.Valid{
					event = event + "["+rid.String
				}else{
					event = event + "[-"
				}
				if rval.Valid && rval.Int64 > 0{
					event = event + ",val:"+strconv.Itoa(int(rval.Int64))+"]"
				} else{
					event = event + ",val:-]"
				}
			}else if util.Contains(ctgDescriptions[catCHNL].Members, "Ev"+event1){
				if rid.Valid{
					event = event + "["+rid.String
				}else{
					event = event + "[-"
				}

				if event1 == "ChRecv" || event1=="ChSend"{
					if rval.Valid{
						event = event + ",val:"+strconv.Itoa(int(rval.Int64))+"]"
					}else{
						event = event + "]"
					}
				}else{
					event = event + "]"
				}
			}else if util.Contains(ctgDescriptions[catMUTX].Members, "Ev"+event1) || util.Contains(ctgDescriptions[catWGCV].Members, "Ev"+event1){
				if rid.Valid{
					event = event + " ["+rid.String+"]"
				}else{
					event = event + " [-]"
				}
			}
			if srcl.Valid{
				srcLine = srcl.String
			}else{
				srcLine = ""
			}

			if predG.Valid {
				if g == int(predG.Int32){
					//happening on same goroutine, just GID is enough
					//fmt.Printf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
					buff = fmt.Sprintf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
					f.WriteString(buff)

				} else{
					//fmt.Printf("%v@%v (G%v) {\"G%v\": %v, \"G%v\": %v }\n",event,srcLine,g,g,logclock,predG.Int32,predClk.Int32)
					buff = fmt.Sprintf("%v@%v (G%v) {\"G%v\": %v, \"G%v\": %v }\n",event,srcLine,g,g,logclock,predG.Int32,predClk.Int32)
					f.WriteString(buff)
				}
			} else{
				//fmt.Printf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
				buff = fmt.Sprintf("%v@%v (G%v) {\"G%v\": %v}\n",event,srcLine,g,g,logclock)
				f.WriteString(buff)
			}
		}
	}
}

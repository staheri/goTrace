package db

import (
	"fmt"
	"trace"
	"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"log"
	"util"
	"strings"
)

// Take sequence of events, create a new DB Schema and insert events into tables
func Store(events []*trace.Event, app string) (dbName string) {
	var err error
	var res sql.Result
	var createLoc,createLoc0 string
	// Connecting to mysql driver
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}

	// Creating new database for current experiment
	idx := 0
	dbName = app + "X" + strconv.Itoa(idx)
	fmt.Printf("Attempt to create database: %s\n",dbName)
	_,err = db.Exec("CREATE DATABASE "+dbName + ";")
	for err != nil{
		fmt.Printf("Error: %v\n",err)
		idx = idx + 1
		dbName = app + "X" + strconv.Itoa(idx)
		fmt.Printf("Attempt to create database: %s\n",dbName)
		_,err = db.Exec("CREATE DATABASE "+dbName+ ";")
	}

	// Close conncection to re-establish it again with proper DBname
	db.Close()

	// Re-establish
	//dbName = "dinphilX18"
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	}else{
		fmt.Println("Connection Re-Established")
	}
	defer db.Close()
	db.SetMaxOpenConns(50000)
	db.SetMaxIdleConns(40000)
	db.SetConnMaxLifetime(0)

	// Create the triple tables (events, stackFrames, Args)
	createTables(db)

	// QUERIES
	var eid int64
	// for the events with resources (channels, mutex, WaitingGroup)
	insertEventResourceStmt, err := db.Prepare("INSERT INTO Events (offset, type, vc , ts, g, p, linkoff, predG, predClk, rid, reid, rval, rclock, stkID, src, src0) values (? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? );")
	check(err)
	defer insertEventResourceStmt.Close()

	insertStackStmt, err := db.Prepare("INSERT INTO StackFrames (eventID, stkIDX, pc, func, file, line) values (?, ?, ?, ?, ?, ?)")
	check(err)
	defer insertStackStmt.Close()
	insertArgStmt, err   := db.Prepare("INSERT INTO Args (eventID, arg, value) values (?, ?, ?)")
	check(err)
	defer insertArgStmt.Close()


	grtnInitStmt, err       := db.Prepare("SELECT * FROM Goroutines WHERE gid=?")
	check(err)
	defer grtnInitStmt.Close()
	grtnInsertFullStmt, err := db.Prepare("INSERT INTO Goroutines (gid, parent_id, createLoc, createLoc0, create_eid, crlid, crlid0) VALUES (?, ?, ?, ?, ?, ?, ?)")
	check(err)
	defer grtnInsertFullStmt.Close()
	grtnUpdStartStmt, err   := db.Prepare("UPDATE Goroutines SET startLOC=? , start_eid=? WHERE gid=?")
	check(err)
	defer grtnUpdStartStmt.Close()
	grtnUpdEndStmt, err     := db.Prepare("UPDATE Goroutines SET ended=? WHERE gid=?")
	check(err)
	defer grtnUpdEndStmt.Close()
	grtnInsertStmt, err     := db.Prepare("INSERT INTO Goroutines (gid, parent_id) VALUES (?, ?)")
	check(err)
	defer grtnInsertStmt.Close()

	chnlInsertStmt,err      := db.Prepare("INSERT INTO Channels (cid, make_gid, make_eid) VALUES (?,?,?)")
	check(err)
	defer chnlInsertStmt.Close()
	chnlInitStmt,err        := db.Prepare("SELECT * FROM Channels WHERE cid=?")
	check(err)
	defer chnlInitStmt.Close()
	chnlUpdCloseStmt,err    := db.Prepare("UPDATE Channels SET close_eid=?, close_gid=? WHERE cid=?")
	check(err)
	defer chnlUpdCloseStmt.Close()
	chnlUpdScountStmt,err   := db.Prepare("UPDATE Channels SET cntSends = cntSends + 1 WHERE cid=?")
	check(err)
	defer chnlUpdScountStmt.Close()
	chnlUpdRcountStmt,err   := db.Prepare("UPDATE Channels SET cntRecvs = cntRecvs + 1 WHERE cid=?")
	check(err)
	defer chnlUpdRcountStmt.Close()

	// Init vector clocks
	msgs          := make(map[msgKey]eventPredecessor) // storing (to be) pred of a recv
	links         := make(map[int64]eventPredecessor) // storing (to be) pred of an event
	// Resource clocks
	localClock    := make(map[uint64]uint64) // vc[g]           = local clock
	chanClock     := make(map[uint64]uint64) // chansClock[cid] = channel clock
	wgClock       := make(map[uint64]uint64) // wgsClock[cid]   = wg clock
	mutexClock    := make(map[uint64]uint64) // mutexClock[cid] = mutex clock

	createLocs     := make(map[string]int)
	createLocs0    := make(map[string]int)

	var tkey uint64

	predG    := sql.NullInt64{}
	predClk  := sql.NullInt64{}
	rid      := sql.NullString{}
	rval     := sql.NullInt64{}
	reid     := sql.NullInt64{}
	rclock   := sql.NullInt64{}
	linkoff   := sql.NullInt64{}
	srcLine  := sql.NullString{}
	srcLine0  := sql.NullString{}




	cnt := 0
	var crlid,crlid0 int
	for _,e := range events{
		// Debug info
		cnt+=1
		desc := EventDescriptions[e.Type]
		fmt.Printf("%v: %v\n",cnt,desc.Name)
		//if cnt > TOPX{
		//	break
		//}

		// fresh values for each event
		predG    = sql.NullInt64{}
		predClk  = sql.NullInt64{}
		rid      = sql.NullString{}
		rval     = sql.NullInt64{}
		reid     = sql.NullInt64{}
		rclock   = sql.NullInt64{}
		linkoff   = sql.NullInt64{}
		if len(e.Stk) != 0{
			srcLine   = sql.NullString{Valid:true, String: util.FilterSlash(path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line))}
			srcLine0  = sql.NullString{Valid:true, String: util.FilterSlash(path.Base(e.Stk[0].File)+":"+ e.Stk[0].Fn + ":" + strconv.Itoa(e.Stk[0].Line))}
		}else{
			srcLine  = sql.NullString{}
			srcLine0  = sql.NullString{}
		}

		// Assign local logical clock
		if _,ok := localClock[e.G];ok{
			localClock[e.G] = localClock[e.G] + 1
		} else{
			localClock[e.G] = 1
		}

		// Check category of events\
		// Assign resource clocks (channels, WaitingGroups, mutexes)
		// Assign predG, predClk for ChRecv
		// Assign predG, predClk for Link
		// Assign rid, rval (if any), rclock for all resources

		if util.Contains(ctgDescriptions[catMUTX].Members, "Ev"+desc.Name){
			// MUTX event
			// Assign mutexClock
			// Assign rid, rval=Null, rclock
			// predG, predClk: null

			tkey = e.Args[0] // muid - rwid
			rid =  sql.NullString{Valid:true, String: "M"+strconv.FormatUint(tkey,10)} // muid

			if _,ok := mutexClock[tkey];ok{
				mutexClock[tkey] = mutexClock[tkey] + 1
			} else{
				mutexClock[tkey] = 1
			}

			predG = sql.NullInt64{}
			predClk = sql.NullInt64{}
			rval = sql.NullInt64{}
			reid = sql.NullInt64{}
			rclock = sql.NullInt64{Valid:true, Int64: int64(mutexClock[tkey])}


		} else if util.Contains(ctgDescriptions[catCHNL].Members, "Ev"+desc.Name){
			// CHNL event
			// Assign chanClock
			// Assign rid, rval, rclock
			// ChSend? set predG , rval = value
			// ChRecv and MSG[key]? use predG,predClk, else: null,null
			// ChMake/Close? rval = null

			tkey= e.Args[0] // cid
			rid =  sql.NullString{Valid:true, String: "C"+strconv.FormatUint(tkey,10)} // cid

			if _,ok := chanClock[tkey];ok{
				chanClock[tkey] = chanClock[tkey] + 1
			} else{
				chanClock[tkey] = 1
			}

			rclock = sql.NullInt64{Valid:true, Int64: int64(chanClock[tkey])}

			if desc.Name == "ChRecv"{
				rval = sql.NullInt64{Valid:true, Int64: int64(e.Args[2])} // message val
				reid = sql.NullInt64{Valid:true, Int64: int64(e.Args[1])} // message eid
				if vv,ok := msgs[msgKey{e.Args[0],e.Args[1],e.Args[2]}] ; ok{
					// A matching sent is found for the recv
					predG    = sql.NullInt64{Valid:true, Int64: int64(vv.g)}
					predClk  = sql.NullInt64{Valid:true, Int64: int64(vv.clock)}
				}else{
					// Receiver without a matching sender
					predG = sql.NullInt64{}
					predClk = sql.NullInt64{}
				}
			}else{
				// ChMake, ChSend, ChClose
				if desc.Name == "ChSend"{
					rval = sql.NullInt64{Valid:true, Int64: int64(e.Args[2])} // message val
					reid = sql.NullInt64{Valid:true, Int64: int64(e.Args[1])} // message eid
					// Set Predecessor for a receive (key to the event: {cid, eid, val})
					if _,ok := msgs[msgKey{e.Args[0],e.Args[1],e.Args[2]}] ; !ok{
						msgs[msgKey{e.Args[0],e.Args[1],e.Args[2]}] = eventPredecessor{e.G, localClock[e.G]}
					} /*else{ // a send for this particular message has been stored before
						panic("Previously stored as sent!")
					}*/
				}else{ // ChMake. ChClose
					rval = sql.NullInt64{}
					reid = sql.NullInt64{}

				}
				predG = sql.NullInt64{}
				predClk = sql.NullInt64{}
			}
		} else if util.Contains(ctgDescriptions[catWGCV].Members, "Ev"+desc.Name){
			// WGRP event
			// Assign wgsClock
			// Assign rid, rval=(add? val, else? Null), rclock
			// predG, predClk: null

			tkey= e.Args[0] // wgid/cvid
			if strings.HasPrefix(desc.Name,"Cv"){
				rid =  sql.NullString{Valid:true, String: "CV"+strconv.FormatUint(tkey,10)} // cvid
			} else{
				rid =  sql.NullString{Valid:true, String: "W"+strconv.FormatUint(tkey,10)} // wgid
			}

			if _,ok := wgClock[tkey];ok{
				wgClock[tkey] = wgClock[tkey] + 1
			} else{
				wgClock[tkey] = 1
			}

			predG = sql.NullInt64{}
			predClk = sql.NullInt64{}

			if desc.Name == "WgAdd"{ // it has a val
				rval = sql.NullInt64{Valid:true, Int64: int64(e.Args[1])} // val
			}else{
				rval = sql.NullInt64{}
			}
			rclock = sql.NullInt64{Valid:true, Int64: int64(wgClock[tkey])}
			reid = sql.NullInt64{}
			// All resource events are assigned a logical clock based on their id

		} else if e.Link != nil{
			// Set Predecessor for an event (key to the event: TS)
			fmt.Printf("Source: %s\n",e)
			fmt.Printf("LINK: %s\n",e.Link)
			linkoff = sql.NullInt64{Valid:true, Int64: int64(e.Link.Off)}
			if _,ok := links[int64(e.Link.Off)] ; !ok{
				links[int64(e.Link.Off)] = eventPredecessor{e.G, localClock[e.G]}
			} else{ // the link of current event has been linked to another event before
				panic("Previously linked to another event!")
			}
		} else{ // does not fall into any category
			predG     = sql.NullInt64{}
			predClk   = sql.NullInt64{}
			rid       =  sql.NullString{}
			rval      = sql.NullInt64{}
			reid      = sql.NullInt64{}
			rclock    = sql.NullInt64{}
			linkoff    = sql.NullInt64{}

		}

		// So far, all predecessor values are set,
		// all resource values are set
		// if a recv has found a sender, it is all set
		// Now only check if the current event has a predecessor. If so: set predG, set predClk
		// otherise: everything is null
		if vv,ok := links[int64(e.Off)]; ok{
			// Is there a possibility that this event has resource other than G?
			// No. Events with predecessor links only have G resource
			predG    = sql.NullInt64{Valid:true, Int64: int64(vv.g)}
			predClk  = sql.NullInt64{Valid:true, Int64: int64(vv.clock)}
			rval     = sql.NullInt64{}
			reid     = sql.NullInt64{}
			rclock   = sql.NullInt64{}

			if len(e.Args) > 0{
				// For events that has link (according to Go spec), they might
				// have an argument in Args which is the goroutie ID
				// (e.g GoUnblock has the id of goroutine that it unblocks)
				// So we want to save that under rid in the Events table
				tkey= e.Args[0] // g
				rid =  sql.NullString{Valid:true, String: "G"+strconv.FormatUint(tkey,10)} // g

			} else{
				// an event with link without arg
				rid    = sql.NullString{}
			}
		}
		/*fmt.Printf("INSERT INTO Events (offset=%v, type=%v, vc=%v, ts=%v, g=%v, p=%v, linkoff=%v, predG=%v, predClk=%v, rid=%v, reid=%v, rval=%v, rclock=%v, stkID=%v, src=%v)\n",strconv.Itoa(e.Off),
																					 "Ev"+desc.Name,
																					 strconv.Itoa(int(localClock[e.G])),
																					 strconv.Itoa(int(e.Ts)),
																					 strconv.FormatUint(e.G,10),
																					 strconv.Itoa(e.P),
																					 linkoff,
																					 predG,
																					 predClk,
																					 rid,
																					 reid,
																					 rval,
																					 rclock,
																					 strconv.FormatUint(e.StkID,10),
																					 srcLine,
																					 srcLine0)*/
		res,err = insertEventResourceStmt.Exec(strconv.Itoa(e.Off),
																					 "Ev"+desc.Name,
																					 strconv.Itoa(int(localClock[e.G])),
																					 strconv.Itoa(int(e.Ts)),
																					 strconv.FormatUint(e.G,10),
																					 strconv.Itoa(e.P),
																					 linkoff,
																					 predG,
																					 predClk,
																					 rid,
																					 reid,
																					 rval,
																					 rclock,
																					 strconv.FormatUint(e.StkID,10),
																					 srcLine,
																				 	 srcLine0)

		check(err)
		eid, err = res.LastInsertId()
		check(err)

		//insert stacks
		//insertStackframe(eid, e.StkID, e.Stk, db)
		if len(e.Stk) != 0{
			for _,a := range e.Stk{
				_,err := insertStackStmt.Exec(strconv.FormatInt(eid,10), strconv.FormatUint(e.StkID,10), strconv.FormatUint(a.PC,10), a.Fn, path.Base(a.File), strconv.Itoa(a.Line))
				check(err)
			}
		}

		// insert args
		//insertArgs(eid, e.Args, desc.Args, db)
		if len(e.Args) != 0{
			for i,a := range desc.Args{
				_,err = insertArgStmt.Exec(strconv.FormatInt(eid,10), a, strconv.FormatInt(int64(e.Args[i]),10))
				check(err)
			}
		}

		// insert goroutines
		if desc.Name == "GoCreate" || desc.Name == "GoStart" || desc.Name == "GoEnd"{
			var startLoc string
			//grtnEntry(e, eid, db)
			res, err := grtnInitStmt.Query(strconv.FormatUint(e.G,10))
			check(err)
			if res.Next() {
				// this goroutine already has been added
				// do other stuff with it
				if desc.Name == "GoCreate"{
					// this goroutine has been inserted and it creates another goroutine
					// insert child goroutine with (parent_id of current goroutine) (stack createLOC)
					gid := strconv.FormatInt(int64(e.Args[0]),10) // e.Args[0] for goCreate is "g"
					parent_id := e.G
					//fmt.Printf("Len Stack: %v\n",len(e.Stk))
					if len(e.Stk) != 0{
						createLoc = util.FilterSlash(path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line))
						createLoc0 = util.FilterSlash(path.Base(e.Stk[0].File)+":"+ e.Stk[0].Fn + ":" + strconv.Itoa(e.Stk[0].Line))
					}else{
						createLoc = "Unknown"
						createLoc0 = "Unknown"
					}

					if val,ok := createLocs[createLoc] ; ok{
						crlid = val + 1
					}else{
						crlid = 1
					}
					createLocs[createLoc] = crlid

					if val,ok := createLocs0[createLoc0] ; ok{
						crlid0 = val + 1
					}else{
						crlid0 = 1
					}
					createLocs0[createLoc0] = crlid0


					//q = fmt.Sprintf("INSERT INTO Goroutines (gid, parent_id, createLoc, create_eid) VALUES (%v,%v,\"%s\",%v);",gid,parent_id,createLoc,eid)
					//fmt.Printf(">>> Executing %s...\n",)
					_,err := grtnInsertFullStmt.Exec(gid,parent_id,createLoc,createLoc0,eid,crlid,crlid0)
					check(err)
				} else if desc.Name == "GoStart"{
					// this goroutine has been inserted before (with create) // update its row with startLOC
					gid := e.G
					if len(e.Stk) > 0{
						startLoc = util.FilterSlash(path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line))
					} else {
						startLoc = "XXX"
						//return
					}
					_,err := grtnUpdStartStmt.Exec(startLoc,eid,gid)
					check(err)

				} else if desc.Name == "GoEnd"{
					// this goroutine has been inserted before (with create)
					// Now we need to update its row with GoEnd eventID
					gid := e.G
					//q = fmt.Sprintf("UPDATE Goroutines SET ended=%v WHERE gid=%v;",eid,gid)
					//fmt.Printf(">>> Executing %s...\n",q)
					_,err := grtnUpdEndStmt.Exec(eid,gid)
					check(err)
				}
			}else{
				if desc.Name == "GoCreate"{
					// this goroutine has not been inserted (no create) and it creates another goroutine
					gid := strconv.FormatUint(e.G,10) // current G
					parent_id := -1
					_,err := grtnInsertStmt.Exec(gid,parent_id)
					check(err)

					// insert child goroutine with (parent_id of current goroutine) (stack location of create)
					gid = strconv.FormatInt(int64(e.Args[0]),10) // e.Args[0] for goCreate is "g"
					parent_id = int(e.G)

					if len(e.Stk) != 0{
						createLoc = util.FilterSlash(path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line))
						createLoc0 = util.FilterSlash(path.Base(e.Stk[0].File)+":"+ e.Stk[0].Fn + ":" + strconv.Itoa(e.Stk[0].Line))
					}else{
						createLoc = "Unknown"
						createLoc0 = "Unknown"
					}

					if val,ok := createLocs[createLoc] ; ok{
						crlid = val + 1
					}else{
						crlid = 1
					}
					createLocs[createLoc] = crlid

					if val,ok := createLocs0[createLoc0] ; ok{
						crlid0 = val + 1
					}else{
						crlid0 = 1
					}
					createLocs0[createLoc0] = crlid0

					_,err = grtnInsertFullStmt.Exec(gid,parent_id,createLoc,createLoc0,eid,crlid,crlid0)
					check(err)

				} else{
					// this goroutine has not been inserted before (no create) and started/ended out of nowhere
					panic("GoStart/End before creating...It is not possible!")
				}

			}
			err = res.Close()
			check(err)
		} else if desc.Name == "ChSend" || desc.Name == "ChRecv" || desc.Name == "ChMake" || desc.Name == "ChClose"{
			//chanEntry(e, eid, db)
			// search for channel
			var cid uint64
			cid = e.Args[0]

			res, err := chnlInitStmt.Query(strconv.FormatUint(cid,10))
			check(err)

		  if res.Next(){ // this channel has already been inserted
		    if desc.Name == "ChMake"{ // making a made channel? PANIC!
		      panic("making a made channel? PANIC!")
		    }else{
		      if desc.Name == "ChClose"{
		        // update Channels
		  			//fmt.Printf(">>> Executing %s...\n",q)
		  			_,err := chnlUpdCloseStmt.Exec(eid,e.G,cid)
		  			check(err)
		      } else if desc.Name == "ChSend"{
		        // update Channels

		        //fmt.Printf(">>> Executing %s...\n",q)
		      	_, err := chnlUpdScountStmt.Exec(cid)
		      	check(err)
		      } else if desc.Name == "ChRecv"{
		        // update Channels

		        //fmt.Printf(">>> Executing %s...\n",q)
		      	_, err := chnlUpdRcountStmt.Exec(cid)
		      	check(err)
		      } else{
		        panic("Wrong Place!")
		      }
		    }
		  } else{
		    if desc.Name != "ChMake"{ // Operation on un-made channel? PANIC!
					// panic("Operation on un-made channel? PANIC!")
					// there might be a global channel creation, then what?
					// First insert the uninitiated channel
					_, err := chnlInsertStmt.Exec(cid,-1,-1)
		    	check(err)

					// Then handle current channel op
					if desc.Name == "ChClose"{
		        // update Channels
		        //q = fmt.Sprintf("UPDATE Channels SET close_eid=%v, close_gid=%v WHERE cid=%v;",eid,e.G,cid)
		  			//fmt.Printf(">>> Executing %s...\n",q)
						_,err := chnlUpdCloseStmt.Exec(eid,e.G,cid)
		  			check(err)
		      } else if desc.Name == "ChSend"{
		        // update Channels
						_, err := chnlUpdScountStmt.Exec(cid)
		      	check(err)
		      } else if desc.Name == "ChRecv"{
		        // update Channels
						_, err := chnlUpdRcountStmt.Exec(cid)
		      	check(err)
		      } else{
		        panic("Wrong Place!")
		      }
		    } else{
		      // insert

		      //fmt.Printf(">>> Executing %s...\n",q)
		    	_, err := chnlInsertStmt.Exec(cid,e.G,eid)
		    	check(err)
		    }
		  }
			err=res.Close()
			if err != nil{
				panic(err)
			}
		}
		//insertEvent(e, db)
	}
	return dbName
}

// Create tables for newly created schema db
func createTables(db *sql.DB){
	eventsCreateStmt := `CREATE TABLE Events (
    									id int NOT NULL AUTO_INCREMENT,
    									offset int NOT NULL,
    									type varchar(255) NOT NULL,
											vc int NOT NULL,
    									ts bigint NOT NULL,
    									g int NOT NULL,
    									p int NOT NULL,
											linkoff int,
											predG int,
											predClk int,
											rid varchar(255),
											reid int,
											rval bigint,
											rclock int,
    									stkID int,
											src varchar(255),
											src0 varchar(255),
											PRIMARY KEY (id)
											);`
	stkFrmCreateStmt := `CREATE TABLE StackFrames (
    									id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
							    		eventID int NOT NULL,
							    		stkIDX int NOT NULL,
							    		pc int NOT NULL,
							    		func varchar(255) NOT NULL,
							    		file varchar(255) NOT NULL,
							    		line int NOT NULL
											);`
	argsCreateStmt :=   `CREATE TABLE Args (
											id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    									eventID int NOT NULL,
    									arg varchar(255) NOT NULL,
    									value bigint NOT NULL);`
	grtnCreateStmt  :=  `CREATE TABLE Goroutines (
    									id int NOT NULL AUTO_INCREMENT,
    									gid int NOT NULL,
    									parent_id int NOT NULL,
											ended int DEFAULT -1,
    									createLoc varchar(255),
											createLoc0 varchar(255),
											create_eid int,
											crlID int,
											crlID0 int,
											startLoc varchar(255),
											start_eid int,
    									PRIMARY KEY (id)
											);`
  chanCreateStmt  :=  `CREATE TABLE Channels (
    									id int NOT NULL AUTO_INCREMENT,
    									cid int NOT NULL,
    									make_gid int NOT NULL DEFAULT -1,
                      make_eid int NOT NULL DEFAULT -1,
                      close_gid int NOT NULL DEFAULT -1,
                      close_eid int NOT NULL DEFAULT -1,
											cntSends int DEFAULT 0,
                      cntRecvs int DEFAULT 0,
    									PRIMARY KEY (id)
											);`
	/*msgCreateStmt   :=  `CREATE TABLE Messages (
    									id int NOT NULL AUTO_INCREMENT,
                      message_id int NOT NULL,
                      channel_id int NOT NULL,
    									sender_gid int NOT NULL DEFAULT -1,
                      receiver_gid int NOT NULL DEFAULT -1,
                      PRIMARY KEY (id)
											);`*/

	createTable(eventsCreateStmt,"Events",db)
	createTable(stkFrmCreateStmt,"StackFrames",db)
	createTable(argsCreateStmt,"Args",db)
	createTable(grtnCreateStmt,"Goroutines",db)
  createTable(chanCreateStmt,"Channels",db)
  //createTable(msgCreateStmt,"Messages",db)
}

// Create individual tables for schema db
func createTable(stmt , name string, db *sql.DB) () {
	fmt.Printf("Creating table %v ... \n",name)
	_,err := db.Exec(stmt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v Created! \n",name)
}

// Take an event and insert to tables
func insertEvent(e *trace.Event, db *sql.DB){
	var q string
	var eid int64
	desc := EventDescriptions[e.Type]
	q = "INSERT INTO Events (offset, type, seq , ts, g, p, stkID, hasStk, hasArgs)"
	q = q + " VALUES ("
	q = q + strconv.Itoa(e.Off) + ", "
	//q = q + desc.Name + ", "
	q = q + "\"Ev"+desc.Name + "\", "
	q = q + strconv.Itoa(int(e.Seq)) + ", "
	q = q + strconv.Itoa(int(e.Ts)) + ", "
	q = q + strconv.FormatUint(e.G,10) + ", "
	q = q + strconv.Itoa(e.P) + ", "
	q = q + strconv.FormatUint(e.StkID,10) + ", "
	q = q + strconv.FormatBool(len(e.Stk) != 0) + ", "
	q = q + strconv.FormatBool(len(e.Args) != 0)
	q = q + ");"
	fmt.Printf("> Executing %s...\n",q)
	res,err := db.Exec(q)
	if err != nil {
		panic(err)
	} else{
		eid, err = res.LastInsertId()
	}


	// insert stacks
	//if len(e.Stk) != 0{
	//	insertStackframe(eid, e.StkID, e.Stk, db)
	//}

	// insert args
	if len(e.Args) != 0{
		insertArgs(eid, e.Args, desc.Args, db)
	}

	// insert goroutines
	if desc.Name == "GoCreate" || desc.Name == "GoStart" || desc.Name == "GoEnd"{
		grtnEntry(e, eid, db)
	} else if desc.Name == "ChSend" || desc.Name == "ChRecv" || desc.Name == "ChMake" || desc.Name == "ChClose"{
		chanEntry(e, eid, db)
	}
}

// Insert stack frames
func insertStackframe(eventID int64, stkIDX uint64, frames []*trace.Frame, db *sql.DB) {
	var s string
	for _,a := range frames{
		s = "INSERT INTO StackFrames (eventID, stkIDX, pc, func, file, line)"
		s = s + " VALUES ("
		s = s + strconv.FormatInt(eventID,10) + ", "
		s = s + strconv.FormatUint(stkIDX,10) + ", "
		s = s + strconv.FormatUint(a.PC,10) + ", "
		s = s + "\""+a.Fn + "\", "
		s = s + "\""+path.Base(a.File) + "\", "
		s = s + strconv.Itoa(a.Line)
		s = s +");"
		fmt.Printf("> Executing %s...\n",s)
		_,err := db.Exec(s)
		if err != nil{
			panic(err)
		}
	}
}

// Insert args
func insertArgs(eventID int64, args [4]uint64, descArgs []string, db *sql.DB) {
	var s string

	for i,a := range descArgs{
		s = "INSERT INTO Args (eventID, arg, value)"
		s = s + " VALUES ("
		s = s + strconv.FormatInt(eventID,10) + ", "
		s = s + "\""+ a + "\", "
		s = s + strconv.FormatInt(int64(args[i]),10)
		s = s +");"
		fmt.Printf("> Executing %s ...\n",s)
		_,err := db.Exec(s)
		if err != nil{
			panic(err)
		}
	}
}

// insert/update Goroutine table
func grtnEntry(e *trace.Event, eid int64, db *sql.DB){
	desc := EventDescriptions[e.Type]
	var q string
	var startLoc string
	//search for event gid in the table
	q = "SELECT * FROM Goroutines WHERE gid="+strconv.FormatUint(e.G,10)+";"
	fmt.Printf(">>> Executing %s...\n",q)
	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}

	if res.Next(){
		// this goroutine already has been added
		// do other stuff with it
		if desc.Name == "GoCreate"{
			// this goroutine has been inserted and it creates another goroutine

			// insert child goroutine with (parent_id of current goroutine) (stack createLOC)
			gid := strconv.FormatInt(int64(e.Args[0]),10) // e.Args[0] for goCreate is "g"
			parent_id := e.G
			createLoc := util.FilterSlash(path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line))
			q = fmt.Sprintf("INSERT INTO Goroutines (gid, parent_id, createLoc, create_eid) VALUES (%v,%v,\"%s\",%v);",gid,parent_id,createLoc,eid)
			fmt.Printf(">>> Executing %s...\n",q)
			_,err := db.Exec(q)
			if err != nil{
				panic(err)
			}


		} else if desc.Name == "GoStart"{
			// this goroutine has been inserted before (with create)
			// update its row with startLOC
			gid := e.G
			if len(e.Stk) > 0{
				startLoc = util.FilterSlash(path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line))
			} else {
				//startLoc = "NIL"
				return
			}

			q = fmt.Sprintf("UPDATE Goroutines SET startLOC=\"%s\", start_eid=%v WHERE gid=%v;",startLoc,eid,gid)
			fmt.Printf(">>> Executing %s...\n",q)
			_,err := db.Exec(q)
			if err != nil{
				panic(err)
			}

		} else if desc.Name == "GoEnd"{
			// this goroutine has been inserted before (with create)
			// Now we need to update its row with GoEnd eventID
			// select * from
			// update
			gid := e.G
			q = fmt.Sprintf("UPDATE Goroutines SET ended=%v WHERE gid=%v;",eid,gid)
			fmt.Printf(">>> Executing %s...\n",q)
			_,err := db.Exec(q)
			if err != nil{
				panic(err)
			}
		}

	} else{
		if desc.Name == "GoCreate"{
			// this goroutine has not been inserted (no create) and it creates another goroutine

			// insert current goroutine
			gid := strconv.FormatUint(e.G,10) // current G
			parent_id := -1
			q = fmt.Sprintf("INSERT INTO Goroutines (gid, parent_id) VALUES (%s,%v);",gid,parent_id)
			fmt.Printf(">>> Executing %s...\n",q)
			res,err := db.Exec(q)
			if err != nil{
				panic(err)
			}
			tmp,_ := res.LastInsertId()
			fmt.Printf(">>> LAST ID: %v \n",tmp)

			// insert child goroutine with (parent_id of current goroutine) (stack location of create)
			gid = strconv.FormatInt(int64(e.Args[0]),10) // e.Args[0] for goCreate is "g"
			parent_id = int(e.G)
			createLoc := util.FilterSlash(path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line))
			q = fmt.Sprintf("INSERT INTO Goroutines (gid, parent_id, createLoc, create_eid) VALUES (%v,%v,\"%s\",%v);",gid,parent_id,createLoc,eid)
			fmt.Printf(">>> Executing %s...\n",q)
			_,err = db.Exec(q)
			if err != nil{
				panic(err)
			}

		} else{
			// this goroutine has not been inserted before (no create) and started/ended out of nowhere
			panic("GoStart/End before creating...It is not possible!")
		}
	}
	err=res.Close()
	if err != nil{
		panic(err)
	}
}

// insert/update channel/message tables
func chanEntry(e *trace.Event, eid int64, db *sql.DB){
	desc := EventDescriptions[e.Type]
	var q string
  var cid uint64

  // search for channel
  if desc.Name == "ChMake" || desc.Name == "ChClose"{
    cid = e.Args[0]
  } else{
    cid = e.Args[1]
  }
  q = "SELECT * FROM Channels WHERE cid="+strconv.FormatUint(cid,10)+";"
  fmt.Printf(">>> Executing %s...\n",q)
	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}

  if res.Next(){ // this channel has already been inserted
    if desc.Name == "ChMake"{ // making a made channel? PANIC!
      panic("making a made channel? PANIC!")
    }else{
      if desc.Name == "ChClose"{
        // update Channels
        q = fmt.Sprintf("UPDATE Channels SET close_eid=%v, close_gid=%v WHERE cid=%v;",eid,e.G,cid)
  			fmt.Printf(">>> Executing %s...\n",q)
  			_,err := db.Exec(q)
  			if err != nil{
  				panic(err)
  			}
      } else if desc.Name == "ChSend"{
        // update Channels
        q = fmt.Sprintf("UPDATE Channels SET cntSends = cntSends + 1 WHERE cid=%v;",cid)
        fmt.Printf(">>> Executing %s...\n",q)
      	_, err := db.Exec(q)
      	if err != nil {
      		panic(err)
      	}
      } else if desc.Name == "ChRecv"{
        // update Channels
        q = fmt.Sprintf("UPDATE Channels SET cntRecvs = cntRecvs + 1 WHERE cid=%v;",cid)
        fmt.Printf(">>> Executing %s...\n",q)
      	_, err := db.Exec(q)
      	if err != nil {
      		panic(err)
      	}
      } else{
        panic("Wrong Place!")
      }
    }
  } else{
    if desc.Name != "ChMake"{ // Operation on un-made channel? PANIC!
			// panic("Operation on un-made channel? PANIC!")
			// there might be a global channel creation, then what?
			// First insert the uninitiated channel
      q = fmt.Sprintf("INSERT INTO Channels (cid, make_gid, make_eid) VALUES (%v,%v,%v);",cid,-1,-1)
      fmt.Printf(">>> Executing %s...\n",q)
    	_, err := db.Exec(q)
    	if err != nil {
    		panic(err)
    	}
			// Then handle current channel op
			if desc.Name == "ChClose"{
        // update Channels
        q = fmt.Sprintf("UPDATE Channels SET close_eid=%v, close_gid=%v WHERE cid=%v;",eid,e.G,cid)
  			fmt.Printf(">>> Executing %s...\n",q)
  			_,err := db.Exec(q)
  			if err != nil{
  				panic(err)
  			}
      } else if desc.Name == "ChSend"{
        // update Channels
        q = fmt.Sprintf("UPDATE Channels SET cntSends = cntSends + 1 WHERE cid=%v;",cid)
        fmt.Printf(">>> Executing %s...\n",q)
      	_, err := db.Exec(q)
      	if err != nil {
      		panic(err)
      	}
      } else if desc.Name == "ChRecv"{
        // update Channels
        q = fmt.Sprintf("UPDATE Channels SET cntRecvs = cntRecvs + 1 WHERE cid=%v;",cid)
        fmt.Printf(">>> Executing %s...\n",q)
      	_, err := db.Exec(q)
      	if err != nil {
      		panic(err)
      	}
      } else{
        panic("Wrong Place!")
      }
    } else{
      // insert
      q = fmt.Sprintf("INSERT INTO Channels (cid, make_gid, make_eid) VALUES (%v,%v,%v);",cid,e.G,eid)
      fmt.Printf(">>> Executing %s...\n",q)
    	_, err := db.Exec(q)
    	if err != nil {
    		panic(err)
    	}
    }
  }
	err=res.Close()
	if err != nil{
		panic(err)
	}
}

// Operations on db
func Ops(command, appName, X string) (dbName string){
	// Vars
	var dbs,q string
	var xx    int
	// Connecting to mysql driver
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	defer db.Close()

	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	//fmt.Println("Command: "+command)
	if command == "clean all"{
		res,err := db.Query("SHOW DATABASES;")
		if err != nil {
			log.Fatal(err)
		}
		for res.Next(){
			err := res.Scan(&dbs)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("DB: %s \n",dbs)
			if dbs[len(dbs)-1] >= '0' && dbs[len(dbs)-1] <= '9'{
				q = "DROP DATABASE "+dbs+";"
				//fmt.Printf(">>> Executing %s...\n",q)
				_,err2 := db.Exec(q)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
		}
		err=res.Close()
		if err != nil{
			panic(err)
		}
		return ""
	}else if command == "x"{
		//fmt.Println("SHOW DATABASES LIKE \""+appName+"X"+X+"\";")
		res,err := db.Query("SHOW DATABASES LIKE \""+appName+"X"+X+"\";")
		if err != nil {
			log.Fatal(err)
		}
		if res.Next(){
			err := res.Scan(&dbs)
			if err != nil {
				log.Fatal(err)
			}
			return dbs
		}else{
			panic("Database "+appName+"X"+X+" does not exist!")
		}
		err=res.Close()
		if err != nil{
			panic(err)
		}
	}else if command == "latest"{
		xx = 0
		for {
			res,err := db.Query("SHOW DATABASES LIKE \""+appName+"X"+strconv.Itoa(xx)+"\";")
			if err != nil {
				log.Fatal(err)
			}
			if !res.Next(){
			return appName+"X"+strconv.Itoa(xx-1)
			} else{
				xx += 1
				continue
			}
			err=res.Close()
			if err != nil{
				panic(err)
			}
		}
	}else{
		panic("Ops command unknown!")
	}
	return ""
}

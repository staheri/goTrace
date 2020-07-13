package db

import (
	"fmt"
	"trace"
	"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"log"
)


// Take sequence of events, create a new DB Schema and insert events into tables
func Store(events []*trace.Event, app string) (dbName string) {
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
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	}else{
		fmt.Println("Connection Re-Established")
	}
	defer db.Close()

	// Create the triple tables (events, stackFrames, Args)
	createTables(db)

	for _,e := range events{
		if e.Link != nil{
			fmt.Printf(" > > > %v (g%v) -> %v (g%v)\n",EventDescriptions[e.Type].Name,e.G,EventDescriptions[e.Link.Type].Name,e.Link.G)
		}
		insertEvent(e, db)
	}
	return dbName
}

// Create tables for newly created schema db
func createTables(db *sql.DB){
	eventsCreateStmt := `CREATE TABLE Events (
    									id int NOT NULL AUTO_INCREMENT,
    									offset int NOT NULL,
    									type varchar(255) NOT NULL,
											seq int NOT NULL,
    									ts bigint NOT NULL,
    									g int NOT NULL,
    									p int NOT NULL,
    									stkID int,
											hasSTK bool,
											hasArgs bool,
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
											create_eid int,
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
  msgCreateStmt   :=  `CREATE TABLE Messages (
    									id int NOT NULL AUTO_INCREMENT,
                      message_id int NOT NULL,
                      channel_id int NOT NULL,
    									sender_gid int NOT NULL DEFAULT -1,
                      receiver_gid int NOT NULL DEFAULT -1,
                      PRIMARY KEY (id)
											);`

	createTable(eventsCreateStmt,"Events",db)
	createTable(stkFrmCreateStmt,"StackFrames",db)
	createTable(argsCreateStmt,"Args",db)
	createTable(grtnCreateStmt,"Goroutines",db)
  createTable(chanCreateStmt,"Channels",db)
  createTable(msgCreateStmt,"Messages",db)
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
	if len(e.Stk) != 0{
		insertStackframe(eid, e.StkID, e.Stk, db)
	}

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
func insertArgs(eventID int64, args [3]uint64, descArgs []string, db *sql.DB) {
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
			createLoc := path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line)
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
				startLoc = path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line)
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
			createLoc := path.Base(e.Stk[len(e.Stk)-1].File)+":"+ e.Stk[len(e.Stk)-1].Fn + ":" + strconv.Itoa(e.Stk[len(e.Stk)-1].Line)
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
      	_, err := db.Query(q)
      	if err != nil {
      		panic(err)
      	}
      } else if desc.Name == "ChRecv"{
        // update Channels
        q = fmt.Sprintf("UPDATE Channels SET cntRecvs = cntRecvs + 1 WHERE cid=%v;",cid)
        fmt.Printf(">>> Executing %s...\n",q)
      	_, err := db.Query(q)
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
    	_, err := db.Query(q)
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
      	_, err := db.Query(q)
      	if err != nil {
      		panic(err)
      	}
      } else if desc.Name == "ChRecv"{
        // update Channels
        q = fmt.Sprintf("UPDATE Channels SET cntRecvs = cntRecvs + 1 WHERE cid=%v;",cid)
        fmt.Printf(">>> Executing %s...\n",q)
      	_, err := db.Query(q)
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
    	_, err := db.Query(q)
    	if err != nil {
    		panic(err)
    	}
    }
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
		}
	}else{
		panic("Ops command unknown!")
	}
}

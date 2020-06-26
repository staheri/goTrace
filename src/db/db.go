package db

import (
	"fmt"
	"trace"
	"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"log"
	"os"
	"os/exec"
	"strings"
	"bytes"
	"github.com/jedib0t/go-pretty/table"
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
		fmt.Println(err)
	}else{
		fmt.Println("Connection Re-Established")
	}
	defer db.Close()

	// Create the triple tables (events, stackFrames, Args)
	createTables(db)

	for _,e := range events{
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
    									value int NOT NULL);`
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
      q = fmt.Sprintf("INSERT INTO Channels (cid, make_gid, make_eid) VALUES (%v,%v,%v);",e.Args[1],-1,-1)
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
      q = fmt.Sprintf("INSERT INTO Channels (cid, make_gid, make_eid) VALUES (%v,%v,%v);",e.Args[1],e.G,eid)
      fmt.Printf(">>> Executing %s...\n",q)
    	_, err := db.Query(q)
    	if err != nil {
    		panic(err)
    	}
    }
  }
}

func Ops(command string){
	// Connecting to mysql driver
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}

	res,err := db.Query("SHOW DATABASES;")
	if err != nil {
		log.Fatal(err)
	}

	var dbs,q string
	for res.Next(){
		err := res.Scan(&dbs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("DB: %s \n",dbs)
		if dbs[len(dbs)-1] >= '0' && dbs[len(dbs)-1] <= '9'{
			q = "DROP DATABASE "+dbs+";"
			fmt.Printf(">>> Executing %s...\n",q)
			_,err2 := db.Exec(q)

			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}

	if command == "CLEAN"{

	}
}

func WriteData(dbName, datapath, filter string, chunkSize int){
	// Re-establish
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()



	var q,output, dbs,tmps string
	var chunk  int
	appgs := appGoroutineFinder(db)

	// SeqAll
	chunk = 0
	dbs = ""
	tmps = ""
	output = datapath + dbName+"_l"+strconv.Itoa(chunkSize)+"_seqALL_"+filter+".py"
	f,err := os.Create(output)
	if err != nil{
		log.Fatal(err)
	}
	if filter != "all"{
		q = "SELECT t1.type FROM Events t1 INNER JOIN global.cat"+filter+" t2 ON t1.type = t2.eventName;"
	} else{
		q = "SELECT type FROM Events;"
	}

	fmt.Printf(">>> Executing %s...\n",q)
	res,err := db.Query(q)
	if err != nil{
		log.Fatal(err)
	}
	tmps = "data = [\n\t["
	for res.Next(){
		err := res.Scan(&dbs)
		if err != nil {
			log.Fatal(err)
		}
		tmps = tmps + "\""
		tmps = tmps + dbs
		//fmt.Printf("DB: %s \n",dbs)
		tmps = tmps + "\","
		chunk = chunk + 1
		if chunk % chunkSize == 0{
			chunk = 0
			tmps = tmps + "],"
			f.WriteString(tmps)
			tmps = "\n\t["
		}
	}
	tmps = tmps + "]\n]"
	f.WriteString(tmps)
	f.Close()

	// SeqApp
	chunk = 0
	dbs = ""
	tmps = ""
	output = datapath + dbName+"_l"+strconv.Itoa(chunkSize)+"_seqAPP_"+filter+".py"

	f,err = os.Create(output)
	if err != nil{
		log.Fatal(err)
	}

	if filter != "all"{
		q = "SELECT t1.type FROM Events t1 INNER JOIN global.cat"+filter+" t2 ON t1.type = t2.eventName "
	} else{
		q = "SELECT t1.type FROM Events t1 "
	}

	// extend the query for selecting apps
	fmt.Println(len(appgs))
	q = q + "WHERE "
	for i,g := range(appgs){
		q = q + "t1.g=" + strconv.Itoa(g)
		if i  !=  len(appgs) - 1 {
			q = q + " OR "
		}else{
			q = q + ";"
		}
	}
	// Executing Query
	fmt.Printf(">>> Executing %s...\n",q)
	res,err = db.Query(q)
	if err != nil{
		log.Fatal(err)
	}
	tmps = "data = [\n\t["
	for res.Next(){
		err := res.Scan(&dbs)
		if err != nil {
			log.Fatal(err)
		}
		tmps = tmps + "\""
		tmps = tmps + dbs
		//fmt.Printf("DB: %s \n",dbs)
		tmps = tmps + "\","
		chunk = chunk + 1
		if chunk % chunkSize == 0{
			chunk = 0
			tmps = tmps + "],"
			f.WriteString(tmps)
			tmps = "\n\t["
		}
	}
	tmps = tmps + "]\n]\n"
	f.WriteString(tmps)
	f.Close()


	// Goroutines
	chunk = 0
	dbs = ""
	tmps = ""
	output = datapath + dbName+"_l"+strconv.Itoa(chunkSize)+"_grtnAPP_"+filter+".py"

	f,err = os.Create(output)
	if err != nil{
		log.Fatal(err)
	}

	for i,g := range(appgs){
		chunk = 0
	  if filter != "all"{
	    q = "SELECT t1.type FROM Events t1 INNER JOIN global.cat"+filter+" t2 ON t1.type = t2.eventName WHERE t1.g="+ strconv.Itoa(g)+";"
	  } else{
	    q = "SELECT t1.type FROM Events t1 WHERE t1.g="+ strconv.Itoa(g)+";"
	  }

		// Executing Query
		fmt.Printf(">>> Executing %s...\n",q)
		res,err := db.Query(q)
		if err != nil{
			log.Fatal(err)
		}
		tmps = "data_g"+strconv.Itoa(i)+" = [\n\t["
		for res.Next(){
			err := res.Scan(&dbs)
			if err != nil {
				log.Fatal(err)
			}
			tmps = tmps + "\""
			tmps = tmps + dbs
			//fmt.Printf("DB: %s \n",dbs)
			tmps = tmps + "\","
			chunk = chunk + 1
			if chunk % chunkSize == 0{
				chunk = 0
				tmps = tmps + "],"
				f.WriteString(tmps)
				tmps = "\n\t["
			}
		}
		tmps = tmps + "]\n]\n"
		f.WriteString(tmps)
	}
	f.Close()
}

func appGoroutineFinder(db *sql.DB) (appGss []int){
	var q string
	q = "SELECT id,gid,startLoc FROM Goroutines;"
	fmt.Printf(">>> Executing %s...\n",q)
	res,err := db.Query(q)
	if err != nil{
		log.Fatal(err)
	}
	var dbs sql.NullString
	var id,gid int

	for res.Next(){
		err := res.Scan(&id,&gid,&dbs)
		if err != nil {
			log.Fatal(err)
		}
		if dbs.Valid{
			if strings.Split(strings.Split(dbs.String,":")[1],".")[0] == "main" {
				// this is the app goroutine
				// add it to slice
				appGss = append(appGss,gid)
			}
		}
	}
	return appGss
}

func FormalContext(dbName, outpath string, aspects ...string ){
	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()

	var q,subq,event string
	var id int

	data := make(map[int][]string)


	q = `SELECT t1.id, t2.type
	     FROM Goroutines t1
			 INNER JOIN Events t2 ON t1.gid=t2.g `

	subq = ""
	if len(aspects) != 0{
		for i,asp := range aspects{
			 subq = subq + "SELECT * FROM global.cat"+asp
			 if i < len(aspects) - 1{
				 subq = subq + " UNION "
			 }
		}
		q = q + "INNER JOIN ("+subq+") t3 ON t3.eventName=t2.type ORDER BY t2.ts;"
	} else{
		q = q + " ORDER BY t2.ts;"
	}
	// query the database
	fmt.Printf(">>> Executing %s...\n",q)
	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}


	// create directory
	outdir := outpath + "/" +dbName + "/"
	filts := ""
	if len(aspects) != 0{
		for i,asp := range aspects{
			filts = filts + asp
			if i < len(aspects) - 1{
				 filts = filts + "_"
			}
		}
	} else{
		filts = filts + "all"
	}
	outdir = outdir + filts
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
    os.MkdirAll(outdir, 0755)
	}


 	for res.Next(){
		err = res.Scan(&id,&event)
		if err != nil{
			panic(err)
		}
		//if val,ok := data[id];ok{
		data[id] = append(data[id],event)
		//}else{}
	}

	// store files in the outpath folder
	for k,v := range data{
		output := outdir+"/g"+strconv.Itoa(k)+".txt"
		f,err := os.Create(output)
		if err != nil{
			log.Fatal(err)
		}
		fmt.Printf("\ndata[%v]:\n\t",k)
		for _,e := range v{
			fmt.Printf("%v\n\t",e)
			f.WriteString(fmt.Sprintf("%v\n",e))
		}
		f.Close()
	}

	// Execute C++ cl on outdir
	_cmd := CLPATH + "/cltrace -m 1 -p "+outdir
	cmd := exec.Command(CLPATH + "/cltrace","-m","1","-p",outdir)
	fmt.Printf(">>> Executing %s...\n",_cmd)
	err = cmd.Run()
	if err != nil{
		log.Fatal(err)
	}

	// Execute python hac on outdir/cl
	_cmd = "python "+ HACPATH + "/main.py " + outdir+"/cl/"+dbName+"_"+filts+".dot "+RESPATH+"/"+dbName+"_"+filts

	cmd = exec.Command("python",HACPATH + "/main.py",outdir+"/cl/"+dbName+"_"+filts+".dot",RESPATH+"/"+dbName+"_"+filts)
	fmt.Printf(">>> Executing %s...\n",_cmd)
	//err = cmd.Run()
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    return
	}
	fmt.Println("Result: " + out.String())
}

func ChannelReport(dbName string){

	// Variables
	var q, event             string
	var report, tmp               string
	var file, funct          string
	var id, cid, ts, gid     int
	var make_eid, make_gid   int
	var close_eid, close_gid int

	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()

	// Query channels
	q = `SELECT id,cid,make_eid,make_gid,close_eid,close_gid
	     FROM Channels;`
	//fmt.Printf("Executing: %v\n",q)

	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}

	// Generate report for each channel
	for res.Next(){
		err = res.Scan(&id,&cid,&make_eid,&make_gid,&close_eid,&close_gid)
		if err != nil{
			panic(err)
		}
		report = "Channel global ID: "+strconv.Itoa(cid)+"\n"
		report = report + "Owner: "
		// Now generate reports
		// Channel ID:
		// Owner: Goroutine ID + file + func
		// Closed:
		// Generate the tables of sends/recvs



		if make_eid != -1{
			// Query to find location of channel make
			q = `SELECT t2.file,t2.func
			     FROM Channels t1
					 INNER JOIN StackFrames t2
					 ON t1.make_eid=t2.eventID
					 WHERE t1.cid=`+strconv.Itoa(cid)+";"

			//fmt.Printf("Executing: %v\n",q)
			res1, err1 := db.Query(q)
			if err1 != nil {
				panic(err1)
			}
			for res1.Next(){
				err1 = res1.Scan(&file,&funct)
				if err1 != nil {
					panic(err1)
				}
				//report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
			}
			report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
		} else{ // global declaration of channel
			report = report + "N/A (e.g., created globaly)\n"
		}

		report = report + "Closed? "

		if close_eid != -1{
			// Query to find location of channel make
			q = `SELECT t2.file,t2.func
			     FROM Channels t1
					 INNER JOIN StackFrames t2
					 ON t1.close_eid=t2.eventID
					 WHERE t1.cid=`+strconv.Itoa(cid)+";"

			fmt.Printf("Executing: %v\n",q)
			res1, err1 := db.Query(q)
			if err1 != nil {
				panic(err1)
			}
			for res1.Next(){
				err1 = res1.Scan(&file,&funct)
				if err1 != nil {
					panic(err1)
				}
				//report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
			}
			report = report + "Yes, G"+strconv.Itoa(close_gid)+": "+file+" >> "+funct+"\n"
		} else{ // global declaration of channel
			report = report + "No\n"
		}

		// now generate table

		// query to obtain send/recv for channelID=cid
		q = `SELECT t1.id, t1.type, t1.ts, t1.g
		     FROM Events t1
				 INNER JOIN global.catCHNL t2 ON t1.type=t2.eventName
				 INNER JOIN Args t3 ON t1.id=t3.eventID
				 WHERE t3.arg="cid" AND t3.value=`+strconv.Itoa(cid)+`
				 ORDER BY t1.ts;`
		//fmt.Printf("Executing: %v\n",q)
		res1, err1 := db.Query(q)
		if err1 != nil {
			panic(err1)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"TS","Send", "Recv"})

		for res1.Next(){
			err1 = res1.Scan(&id,&event,&ts,&gid)
			if err1 != nil{
				panic(err1)
			}
			// now find stack entry for current row
			q = `SELECT file,func
			     FROM StackFrames
					 WHERE eventID=`+strconv.Itoa(id)+" ORDER BY id;"
			//fmt.Printf("Executing: %v\n",q)
		 	res2, err2 := db.Query(q)
		 	if err2 != nil {
		 		panic(err2)
		 	}
			for res2.Next(){
				err2 = res2.Scan(&file,&funct)
				if err2 != nil{
					panic(err2)
				}
			}
			var row []interface{}
			row = append(row,ts)
			tmp = "G"+strconv.Itoa(gid)+": "+file+" >> "+funct+"\n"
			if event == "EvChSend"{
				row = append(row,tmp)
				row = append(row,"-")
			}else{
				row = append(row,"-")
				row = append(row,tmp)
			}
			t.AppendRow(row)
		}
		fmt.Printf("%v\n",report)
		t.Render()
	}
}


func MutexReport(dbName string){

	// Variables
	var q, event             string
	var report, tmp               string
	var file, funct          string
	var muid,id, ts, gid     int

	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()

	// Query events to find mutex IDs
	q = `SELECT DISTINCT(t3.value)
 	     FROM Events t1
 			 INNER JOIN global.catMUTX t2 ON t1.type=t2.eventName
 			 INNER JOIN args t3 ON t1.id=t3.eventID
 			 WHERE t3.arg="muid";`
	//fmt.Printf("Executing: %v\n",q)

	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	var lockIDs []int
	for res.Next(){
		err = res.Scan(&muid)
		if err != nil {
			panic(err)
		}
		lockIDs=append(lockIDs,muid)
	}

	for _,muid := range lockIDs{
		report = "Mutex global ID: "+strconv.Itoa(muid)+"\n"

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"TS","Lock", "Unlock","RWLock","RWUnlock"})

		// query to obtain send/recv for mutexID=muid
		q = `SELECT t1.id, t1.type, t1.ts, t1.g
		     FROM Events t1
				 INNER JOIN global.catMUTX t2 ON t1.type=t2.eventName
				 INNER JOIN Args t3 ON t1.id=t3.eventID
				 WHERE t3.arg="muid" AND t3.value=`+strconv.Itoa(muid)+`
				 ORDER BY t1.ts;`
		//fmt.Printf("Executing: %v\n",q)
		res1, err1 := db.Query(q)
		if err1 != nil {
			panic(err1)
		}

		for res1.Next(){
			err1 = res1.Scan(&id,&event,&ts,&gid)
			if err1 != nil{
				panic(err1)
			}
			// now find stack entry for current row
			q = `SELECT file,func
			     FROM StackFrames
					 WHERE eventID=`+strconv.Itoa(id)+`
					 ORDER BY id`
			//fmt.Printf("Executing: %v\n",q)
		 	res2, err2 := db.Query(q)
		 	if err2 != nil {
		 		panic(err2)
		 	}
			for res2.Next(){
				err2 = res2.Scan(&file,&funct)
				if err2 != nil{
					panic(err2)
				}
			}
			var row []interface{}
			row = append(row,ts)
			tmp = "G"+strconv.Itoa(gid)+": "+file+" >> "+funct+"\n"
			if event == "EvMuLock"{
				row = append(row,tmp)
				row = append(row,"-")
				row = append(row,"-")
				row = append(row,"-")
			}else if event == "EvMuUnlock"{
				row = append(row,"-")
				row = append(row,tmp)
				row = append(row,"-")
				row = append(row,"-")
			} else if event == "EvRWMLock"{
				row = append(row,"-")
				row = append(row,"-")
				row = append(row,tmp)
				row = append(row,"-")
			} else{
				row = append(row,"-")
				row = append(row,"-")
				row = append(row,"-")
				row = append(row,tmp)
			}
			t.AppendRow(row)
		}
		fmt.Printf("%v\n",report)
		t.Render()
	}
}





func WaitingGroupReport(dbName string){

	// Variables
	var q, event             string
	var report, tmp               string
	var file, funct          string
	var wid,id, ts, gid, val    int

	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()

	// Query events to find mutex IDs
	q = `SELECT DISTINCT(t3.value)
 	     FROM Events t1
 			 INNER JOIN global.catWGRP t2 ON t1.type=t2.eventName
 			 INNER JOIN args t3 ON t1.id=t3.eventID
 			 WHERE t3.arg="wid";`
	//fmt.Printf("Executing: %v\n",q)

	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	var wIDs []int
	for res.Next(){
		err = res.Scan(&wid)
		if err != nil {
			panic(err)
		}
		wIDs=append(wIDs,wid)
	}

	for _,wid := range wIDs{
		report = "WaitingGroup global ID: "+strconv.Itoa(wid)+"\n"

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"TS","ADD(value+LOC)", "DONE","WAIT"})

		// query to obtain send/recv for mutexID=muid
		q = `SELECT t1.id, t1.type, t1.ts, t1.g
		     FROM Events t1
				 INNER JOIN global.catWGRP t2 ON t1.type=t2.eventName
				 INNER JOIN Args t3 ON t1.id=t3.eventID
				 WHERE t3.arg="wid" AND t3.value=`+strconv.Itoa(wid)+`
				 ORDER BY t1.ts;`
		//fmt.Printf("Executing: %v\n",q)
		res1, err1 := db.Query(q)
		if err1 != nil {
			panic(err1)
		}

		for res1.Next(){
			err1 = res1.Scan(&id,&event,&ts,&gid)
			if err1 != nil{
				panic(err1)
			}
			// now find stack entry for current row
			q = `SELECT file,func
			     FROM StackFrames
					 WHERE eventID=`+strconv.Itoa(id)+`
					 ORDER BY id`
			//fmt.Printf("Executing: %v\n",q)
		 	res2, err2 := db.Query(q)
		 	if err2 != nil {
		 		panic(err2)
		 	}
			for res2.Next(){
				err2 = res2.Scan(&file,&funct)
				if err2 != nil{
					panic(err2)
				}
			}
			var row []interface{}
			row = append(row,ts)
			tmp = "G"+strconv.Itoa(gid)+": "+file+" >> "+funct+"\n"
			if event == "EvWgAdd"{
				// find value
				q =  `SELECT value FROM args WHERE arg="val" and eventID=`+strconv.Itoa(id)+`;`
				//fmt.Printf("Executing: %v\n",q)
			 	res3, err3 := db.Query(q)
			 	if err3 != nil {
			 		panic(err3)
			 	}
				if res3.Next(){
					err3 = res3.Scan(&val)
					if err3 != nil{
						panic(err3)
					}
				}
				if val < 0{
					continue
				} else{
					row = append(row,"Value: "+strconv.Itoa(val)+" * "+tmp)
					row = append(row,"-")
					row = append(row,"-")
				}
			}else if event == "EvWgDone"{
				row = append(row,"-")
				row = append(row,tmp)
				row = append(row,"-")
			} else if event == "EvWgWait"{
				row = append(row,"-")
				row = append(row,"-")
				row = append(row,tmp)
			}
			t.AppendRow(row)
		}
		fmt.Printf("%v\n",report)
		t.Render()
	}
}

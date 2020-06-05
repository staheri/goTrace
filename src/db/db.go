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
func Store(events []*trace.Event, app string) () {
	// Connecting to mysql driver
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}

	// Creating new database for current experiment
	idx := 0
	dbName := app + "_" + strconv.Itoa(idx)
	fmt.Printf("Attempt to create database: %s\n",dbName)
	_,err = db.Exec("CREATE DATABASE "+dbName + ";")
	for err != nil{
		fmt.Printf("Error: %v\n",err)
		idx = idx + 1
		dbName = app + "_" + strconv.Itoa(idx)
		fmt.Printf("Attempt to create database: %s\n",dbName)
		_,err = db.Exec("CREATE DATABASE "+dbName+ ";")
	}


	_,err = db.Exec("USE "+dbName+";")
	if err!=nil{
		panic(err)
	}

	// Create the triple tables (events, stackFrames, Args)
	createTables(db)

	for _,e := range events{
		insertEvent(e, db)
	}
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
	grtnCreateStmt  :=   `CREATE TABLE Goroutines (
    									id int NOT NULL AUTO_INCREMENT,
    									gid int NOT NULL,
    									parent_id int NOT NULL,
											ended bool,
    									createLoc varchar(255),
											startLoc varchar(255),
    									PRIMARY KEY (id)
											);`

	createTable(eventsCreateStmt,"Events",db)
	createTable(stkFrmCreateStmt,"StackFrames",db)
	createTable(argsCreateStmt,"Args",db)
	createTable(grtnCreateStmt,"Goroutines",db)
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
	stmt, err := db.Prepare("INSERT INTO Events (offset, type, seq , ts, g, p, stkID, hasStk, hasArgs) VALUES(?, ?, ?, ?, ?, ?, ?, ?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.
	if res, err := stmt.Exec(strconv.Itoa(e.Off),
	                       "\"Ev"+desc.Name + "\"",
													strconv.Itoa(int(e.Seq)),
													strconv.Itoa(int(e.Ts)),
													strconv.FormatUint(e.G,10),
													strconv.Itoa(e.P),
													strconv.FormatUint(e.StkID,10),
													strconv.FormatBool(len(e.Stk) != 0),
													strconv.FormatBool(len(e.Args) != 0)
													); err != nil {

			log.Fatal(err)
		}else{
			eid, err = res.LastInsertId()
		}
	/*q = "INSERT INTO Events (offset, type, seq , ts, g, p, stkID, hasStk, hasArgs)"
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
	}*/

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
		grtnEntry(e,desc)
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

func Ops(){
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

	var dbs string
	for res.Next(){
		err := res.Scan(&dbs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("DB: %s \n",dbs)
	}
}

func grtnEntry(e *trace.Event){
	desc := EventDescriptions[e.Type]
	//search for event gid in the table
	res, err := db.Query("SELECT * FROM Goroutines WHERE gid="+strconv.FormatUint(e.G,10),+";")
	if err != nil {
		log.Fatal(err)
	}

	if res.Next(){
		// this goroutine already has been added
		// do other stuff with it
		if desc.Name == "GoCreate"{
			// this goroutine has been inserted and it creates another goroutine
			q = "INSERT INTO Goroutines (gid, parent_id, createLoc) VALUES (?,?,?)"
			gid := strconv.FormatInt(int64(e.Args[0]),10) // e.Args[0] for goCreate is "g"
			parent_id := res.idOfReturnedRow
			createLoc := path.Base(e.Stk[-1].File)+":"+ e.Stk[-1].Fn + ":" + strconv.Itoa(e.Stk[-1].Line)
			// insert child goroutine with (parent_id of current goroutine) (stack createLOC)

		} else if desc.Name == "GoStart"{
			// this goroutine has been inserted before (with create)
			// Now we need to update its row with startLOC
			// select * from
			// update
		} else if desc.Name == "GoEnd"{
			// this goroutine has been inserted before (with create)
			// Now we need to update its row with GoEnd eventID
			// select * from
			// update
		}

	} else{
		if desc.Name == "GoCreate"{
			// this goroutine has not been inserted (no create) and it creates another goroutine
			// insert current goroutine
			// insert child goroutine with (parent_id of current goroutine) (stack location of create)
		} else if desc.Name == "GoStart"{
			// this goroutine has not been inserted before (no create) and started out of nowhere
			// assert false
		} else if desc.Name == "GoEnd"{
			// this goroutine has not been inserted before (no create) and ended out of nowhere
			// assert false
		}
		insStmt, err := db.Prepare("INSERT INTO Goroutines* FROM Goroutines WHERE gid=?")
		if err != nil {
			log.Fatal(err)
		}
		defer chkStmt.Close()
		res,err := chkStmt.Exec(e.g)

	}
	if err != nil {
		log.Fatal(err)
	}

	var dbs string
	for res.Next(){
		err := res.Scan(&dbs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("DB: %s \n",dbs)
	}

}

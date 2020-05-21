package main

import (
	"fmt"
	"trace"
	_"util"

	_"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"


)

var(
	exID
)
func main(){
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}

	name := "test1"
	// for each execution (e.g., source code)
	//   - create a new DB (if not exist)
	//   - create tables (events, stackFrames, args) -> func CreateTables()
	//   - for each event in events[]
	//     - generate insert query for:
	//            func insertDB(event):
	//               -
	//            - events
	//            - StackFrms (if any)
	//            - Args (if any)
	//   UP TO HERE, the data is stored in
	//            - DB(exec).Table(events)
	//            - DB(exec).Table(StackFrms)
	//            - DB(exec).Table(Args)
	//   Now we can generate queries and check
	_,err = db.Exec("CREATE DATABASE IF NOT EXISTS "+name)
	if err != nil {
			panic(err)
	}

	fmt.Println("Database Created")

	sqlStatement := "USE "+name+";"
	res,err := db.Exec(sqlStatement)
	if err == nil {
 		fmt.Printf("%v\n",res)
		lastId, err1 := res.LastInsertId()
		if err1 != nil {
			log.Fatal(err1)
		}
		rowCnt, err1 := res.RowsAffected()
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
 	}else{
 		fmt.Println("ERRRRRR")
 		panic(err)
 	}

	sqlStatement = `SHOW TABLES;`
	res,err = db.Exec(sqlStatement)
	if err == nil {
 		fmt.Printf("%v\n",res)
		lastId, err1 := res.LastInsertId()
		if err1 != nil {
			log.Fatal(err1)
		}
		rowCnt, err1 := res.RowsAffected()
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
 	}else{
 		fmt.Println("ERRRRRR")
 		panic(err)
 	}


func CreateTables(){
	eventsCreateStmt := `CREATE TABLE Events (
    									id int NOT NULL AUTO_INCREMENT,
    									offset int NOT NULL,
    									type varchar(255) NOT NULL,
    									timestamp int NOT NULL,
    									goroutine int NOT NULL,
    									process int NOT NULL,
    									stkID int,
    									stkFrmID int,
    									argsID int,
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

	CreateTable(eventsCreateStmt,"Events")
	CreateTable(stkCreateStmt,"StackFrames")
	CreateTable(argsCreateStmt,"Args")
}

func createTable(stmt , name string) () {
	fmt.Printf("Creating table %v ... \n",name)
	_,err := db.Exec(trans)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v Created! \n",name)
}


func Store(events []*trace.Event, app string){

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
	dbName = app + strconv.Itoa(idx)
	_,err := db.Exec("CREATE DATABASE "+dbName)
	for err != nil{
		idx = idx + 1
		dbName = app + strconv.Itoa(idx)
		_,err = db.Exec("CREATE DATABASE "+dbName)
	}

	// Create the triple tables (events, stackFrames, Args)
	CreateTables()

	var query string
	for i,e := range events{
		// generateQuery of e
		// insert into catGRTN (eventName) VALUES ("EvGoCreate");
		query = genInsertEventQuery(i,e)
		_,err := db.Exec(query)
		if err != nil{
			panic(err)
		}
	}
}


func genInsertEventQuery(i int, e *trace.Event){
	desc := EventDescriptions[e.Type]
	offset := e.Off
	typ    := "Ev"+desc.Name
	seq    := e.Seq
	ts     := e.Ts
	p      := e.P
	g      := e.G
	stkid  := e.StkId
	// inject stacks - get stkFrame id and store it
	// insert args - get argsID and store it
	
}
	 /*
   _,err = db.Exec("USE "+name)
   if err != nil {
       panic(err)
   }

  /* _,err = db.Exec("CREATE TABLE example ( id integer, data varchar(32) )")
   if err != nil {
       panic(err)
   }
	 */

	// sqlStatement := `SHOW TABLES`
	 //res,err := db.Exec(sqlStatement)
 	//var table string
 	//if err != nil {
 		/*for res.Next() {
 			res.Scan(&table)
 			fmt.Println(table)
 		}*/
 		//fmt.Printf("%v",res)
 	//}else{
 		//fmt.Println("ERRRRRR")
 		//panic(err)
 	//}

}

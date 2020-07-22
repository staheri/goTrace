package db

import (
	"fmt"
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

var(
	GOPATH    string
	CLPATH    string
	HACPATH   string
)

func WordData(dbName, outdir, filter string, chunkSize int){
	// make sure
	outdir = outdir + "/"

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
	output = outdir + dbName+"_l"+strconv.Itoa(chunkSize)+"_seqALL_"+filter+".py"
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
		if (chunk + 1) % chunkSize == 0{
			tmps = tmps + "\""
		}else{
			tmps = tmps + "\","
		}

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
	output = outdir + dbName+"_l"+strconv.Itoa(chunkSize)+"_seqAPP_"+filter+".py"

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
		if (chunk + 1) % chunkSize == 0{
			tmps = tmps + "\""
		}else{
			tmps = tmps + "\","
		}
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
	output = outdir + dbName+"_l"+strconv.Itoa(chunkSize)+"_grtnAPP_"+filter+".py"

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
			if (chunk + 1) % chunkSize == 0{
				tmps = tmps + "\""
			}else{
				tmps = tmps + "\","
			}
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

func CLOperations(dbName, cloutpath,resultpath string, aspects ...string ){
	// Paths
	setPaths()
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


	q = `SELECT (t1.id)-1, t2.type
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
	cloutdir := cloutpath + "/" +dbName + "/"
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
	cloutdir = cloutdir + filts
	if _, err := os.Stat(cloutdir); os.IsNotExist(err) {
    os.MkdirAll(cloutdir, 0755)
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
		output := cloutdir+"/g"+strconv.Itoa(k)+".txt"
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

	// Execute C++ cl on cloutdir
	_cmd := CLPATH + "/cltrace -m 1 -p "+cloutdir
	cmd := exec.Command(CLPATH + "/cltrace","-m","1","-p",cloutdir)
	fmt.Printf(">>> Executing %s...\n",_cmd)
	err = cmd.Run()
	if err != nil{
		log.Fatal(err)
	}

	// Execute python hac on cloutdir/cl
	_cmd = "python "+ HACPATH + "/main.py " + cloutdir+"/cl/"+dbName+"_"+filts+".dot "+resultpath+"/"+dbName+"_"+filts

	cmd = exec.Command("python",HACPATH + "/main.py",cloutdir+"/cl/"+dbName+"_"+filts+".dot",resultpath+"/"+dbName+"_"+filts)
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
	var report, tmp          string
	var file, funct          string
	var id, cid, ts, gid     int
	var make_eid, make_gid   int
	var close_eid, close_gid int
	var line                 int
	var val, pos             int

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


	// PREPARED STATEMENTS:
	// Query to find location of channel make
	chmakeLocStmt,err := db.Prepare("SELECT t2.file,t2.func,t2.line FROM Channels t1 INNER JOIN StackFrames t2 ON t1.make_eid=t2.eventID WHERE t1.cid=?")
	check(err)
	defer chmakeLocStmt.Close()

	chcloseLocStmt,err := db.Prepare("SELECT t2.file,t2.func,t2.line FROM Channels t1 INNER JOIN StackFrames t2 ON t1.close_eid=t2.eventID WHERE t1.cid=?")
	check(err)
	defer chcloseLocStmt.Close()

	// query to obtain send/recv for channelID=cid
	chsendrecvStmt,err := db.Prepare(`SELECT t1.id, t1.type, t1.ts, t1.g FROM Events t1 INNER JOIN global.catCHNL t2 ON t1.type=t2.eventName INNER JOIN Args t3 ON t1.id=t3.eventID WHERE t3.arg="cid" AND t3.value=? AND (t1.type="EvChSend" OR t1.type="EvChRecv") ORDER BY t1.ts`)
	check(err)
	defer chsendrecvStmt.Close()

	valStmt,err := db.Prepare("SELECT value from args where eventID=? and arg=\"val\"")
	check(err)
	defer valStmt.Close()

	posStmt,err := db.Prepare("SELECT value from args where eventID=? and arg=\"pos\"")
	check(err)
	defer posStmt.Close()

	// now find stack entry for current row
	stackEntryStmt,err := db.Prepare("SELECT file,func,line FROM StackFrames WHERE eventID=? ORDER BY id")
	check(err)
	defer stackEntryStmt.Close()

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
			res1, err1 := chmakeLocStmt.Query(cid)
			check(err1)
			for res1.Next(){
				err1 = res1.Scan(&file,&funct,&line)
				check(err1)
				//report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
			}
			res1.Close()
			report = report + "G"+strconv.Itoa(make_gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"\n"
		} else{ // global declaration of channel
			report = report + "N/A (e.g., created globaly)\n"
		}

		report = report + "Closed? "

		if close_eid != -1{
			res1, err1 := chcloseLocStmt.Query(cid)
			check(err1)
			for res1.Next(){
				err1 = res1.Scan(&file,&funct,&line)
				check(err1)
				//report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
			}
			res1.Close()
			report = report + "Yes, G"+strconv.Itoa(close_gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"\n"
		} else{ // global declaration of channel
			report = report + "No\n"
		}

		// now generate table
		//fmt.Printf("Executing: %v\n",q)
		res1, err1 := chsendrecvStmt.Query(cid)
		check(err1)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"TS","Send", "Recv"})

		for res1.Next(){
			err1 = res1.Scan(&id,&event,&ts,&gid)
			check(err1)

			// now find stack entry for current row
		 	res2, err2 := stackEntryStmt.Query(id)
		 	check(err2)
			for res2.Next(){
				err2 = res2.Scan(&file,&funct,&line)
				if err2 != nil{
					panic(err2)
				}
			}
			res2.Close()
			var row []interface{}
			row = append(row,ts)


			res4,err := valStmt.Query(id)
			check(err)
			if res4.Next(){
				err := res4.Scan(&val)
				check(err)
			}

			res5,err := posStmt.Query(id)
			check(err)
			if res5.Next(){
				err := res5.Scan(&pos)
				check(err)
			}

			/*
			if val == 1{
				tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"-NOPE\n"
			}else if val == 2{
				tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"-FORK\n"
			}else if val == 3{
				tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"-FREE?\n"
			}else if val == 4{
				tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"-REL\n"
			}else{
				tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"-XX\n"
			}
			*/
			tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+">"+strconv.Itoa(val)+"@"+strconv.Itoa(pos)+"\n"
			res4.Close()
			res5.Close()

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
		res1.Close()
		//fmt.Printf("%v\n",report)
		//t.RenderMarkdown()
	}
	err=res.Close()
	check(err)
}

func MutexReport(dbName string){

	// Variables
	var q, event             string
	var report, tmp          string
	var file, funct          string
	var muid,id,ts,gid,line  int

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
			q = `SELECT file,func,line
			     FROM StackFrames
					 WHERE eventID=`+strconv.Itoa(id)+`
					 ORDER BY id`
			//fmt.Printf("Executing: %v\n",q)
		 	res2, err2 := db.Query(q)
		 	if err2 != nil {
		 		panic(err2)
		 	}
			for res2.Next(){
				err2 = res2.Scan(&file,&funct,&line)
				if err2 != nil{
					panic(err2)
				}
			}
			var row []interface{}
			row = append(row,ts)
			tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"\n"
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

		fmt.Printf("%v\n",report)
		t.RenderMarkdown()
	}
}

func RWMutexReport(dbName string){

	// Variables
	var q, event             string
	var report, tmp          string
	var file, funct          string
	var rwid,id,ts,gid,line  int

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
 			 WHERE t3.arg="rwid";`
	//fmt.Printf("Executing: %v\n",q)

	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	var lockIDs []int
	for res.Next(){
		err = res.Scan(&rwid)
		if err != nil {
			panic(err)
		}
		lockIDs=append(lockIDs,rwid)
	}

	for _,rwid := range lockIDs{
		report = "RWMutex global ID: "+strconv.Itoa(rwid)+"\n"

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"TS","RWMLock", "RWMUnlock","RWMrLock","RWMrUnlock"})

		// query to obtain send/recv for mutexID=muid
		q = `SELECT t1.id, t1.type, t1.ts, t1.g
		     FROM Events t1
				 INNER JOIN global.catMUTX t2 ON t1.type=t2.eventName
				 INNER JOIN Args t3 ON t1.id=t3.eventID
				 WHERE t3.arg="rwid" AND t3.value=`+strconv.Itoa(rwid)+`
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
			q = `SELECT file,func,line
			     FROM StackFrames
					 WHERE eventID=`+strconv.Itoa(id)+`
					 ORDER BY id`
			//fmt.Printf("Executing: %v\n",q)
		 	res2, err2 := db.Query(q)
		 	if err2 != nil {
		 		panic(err2)
		 	}
			for res2.Next(){
				err2 = res2.Scan(&file,&funct,&line)
				if err2 != nil{
					panic(err2)
				}
			}
			var row []interface{}
			row = append(row,ts)
			tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"\n"
			if event == "EvRWMLock"{
				row = append(row,tmp)
				row = append(row,"-")
				row = append(row,"-")
				row = append(row,"-")
			}else if event == "EvRWMUnlock"{
				row = append(row,"-")
				row = append(row,tmp)
				row = append(row,"-")
				row = append(row,"-")
			} else if event == "EvRWMrLock"{
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

		fmt.Printf("%v\n",report)
		t.RenderMarkdown()
	}
}

func WaitingGroupReport(dbName string){

	// Variables
	var q, event                  string
	var report, tmp               string
	var file, funct               string
	var wid,id,ts,gid,val,line    int

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
			q = `SELECT file,func,line
			     FROM StackFrames
					 WHERE eventID=`+strconv.Itoa(id)+`
					 ORDER BY id`
			//fmt.Printf("Executing: %v\n",q)
		 	res2, err2 := db.Query(q)
		 	if err2 != nil {
		 		panic(err2)
		 	}
			for res2.Next(){
				err2 = res2.Scan(&file,&funct,&line)
				if err2 != nil{
					panic(err2)
				}
			}
			var row []interface{}
			row = append(row,ts)
			tmp = "G"+strconv.Itoa(gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"\n"
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

		fmt.Printf("%v\n",report)
		t.RenderMarkdown()
	}
}

func setPaths(){
	GOPATH = os.Getenv("GOPATH")
	CLPATH = GOPATH +"/cl"
	HACPATH = GOPATH + "/scripts/hac"
}

func check(err error){
	if err != nil{
		panic(err)
	}
}

func ResourceGraph(dbName, resultpath string, categories ...string ){
//func ResourceGraph(dbName, resultpath string){
	// Variables
	var subq, q, event,arg   string
	var _arg                 sql.NullString
	var _value               sql.NullInt32
	var value                int
	//var file, fuct,arg       string
	var sevent               string
	var id, eid, gid         int
	//var idx, jdx             int
	gmap := make(map[int]int)
	var gmat [][]string
	var row  []string

	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()

	q = `SELECT id,gid FROM Goroutines;`
	res,err := db.Query(q)
	check(err)
	for res.Next(){
		err = res.Scan(&id,&gid)
		check(err)
		gmap[gid]=id-1
		//fmt.Printf("eventG: %v - tableG: %v\n",gid,id-1)
	}
	q = `SELECT t2.id, t2.g, t2.type, t4.arg, t4.value
			 FROM Events t2 `

	subq = ""
 	if len(categories) != 0{
 		for i,cat := range categories{
 			 subq = subq + "SELECT * FROM global.cat"+cat
 			 if i < len(categories) - 1{
 				 subq = subq + " UNION "
 			 }
 		}
 		q = q + "INNER JOIN ("+subq+") t3 ON t3.eventName=t2.type "
 	}
	q = q + `LEFT JOIN Args t4
		 			 ON t2.id=t4.eventID AND (t4.arg="g"
			 		 		OR t4.arg="muid"
			 				OR t4.arg="cid"
			 				OR t4.arg="rwid"
			 				OR t4.arg="wid")
						ORDER BY t2.ts;`
	res,err = db.Query(q)
	check(err)
	for res.Next(){
		err = res.Scan(&eid,&gid,&event,&_arg,&_value)
		check(err)
		sevent = strings.Split(event,"Ev")[1]
		if _arg.Valid{
			arg = _arg.String
		}else{
			arg = ""
		}
		if _value.Valid{
			value = int(_value.Int32)
		}else{
			value = -1
		}
		if arg=="g"{
			//fmt.Printf("%v-%v\n\tVal:%v gmap[%v]:%v\n",gid,sevent,value,value,gmap[value])
			value = gmap[value]
			//fmt.Printf("%v-%v-%v\n",gid,sevent,value)
		}
		if value >= 0{
			sevent = sevent + "-" + strconv.Itoa(value)
		}
		if strings.Contains(sevent,"WgAdd"){
			q = `SELECT value FROM Args WHERE eventID=`+strconv.Itoa(eid)+` and arg="val"`
			fmt.Printf("Executing %v\n",q)
			res1,err1 := db.Query(q)
			check(err1)
			if res1.Next(){
				err2:=res1.Scan(&value)
				check(err2)
				sevent = sevent + "-(" + strconv.Itoa(value)+")"
			}
		}
		if strings.Contains(sevent,"ChSend") || strings.Contains(sevent,"ChRecv"){
			q = `SELECT value FROM Args WHERE eventID=`+strconv.Itoa(eid)+` and arg="eid"`
			res1,err1 := db.Query(q)
			check(err1)
			if res1.Next(){
				err2 := res1.Scan(&value)
				check(err2)
				sevent = sevent + "-" + strconv.Itoa(value)
			}
		}

		//fmt.Println(gmap[gid],"-",sevent)

		row = nil
		for i:=0;i<gmap[gid];i++{
			row = append(row,"-")
		}
		row = append(row,sevent)
		for i:=gmap[gid]+1;i<len(gmap);i++{
			row = append(row,"-")
		}
		gmat = append(gmat,row)
		//fmt.Println(row)
		//fmt.Println(gmat)
	}
	/*fmt.Println(gmat)
	for _,r := range(gmat){
		for _,s := range(r){
			fmt.Printf("%v,",s)
		}
		fmt.Printf("\n")
	}*/

	// Write dot
	outdot := resultpath+"/"+dbName+"_rg.dot"
	outpdf := resultpath+"/"+dbName+"_rg.pdf"
	f,err := os.Create(outdot)
	if err != nil{
		log.Fatal(err)
	}
	f.WriteString(mat2dot(gmat))
	f.Close()

	// Create pdf
	_cmd := "dot -Tpdf "+ outdot + " -o " + outpdf

	cmd := exec.Command("dot","-Tpdf",outdot,"-o",outpdf)
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

func mat2dot(mat [][]string) string{
	if len(mat) < 1{
		panic("Mat is empty")
	}
	if len(mat[0]) < 1{
		panic("Mat row is empty")
	}

	tmp := ""
	dot := "digraph G{\n\trankdir=TB"

	//subgraph G labels (-1)
	tmp = "\n\tsubgraph{"
	tmp = tmp + "\n\t\tnode [margin=0 fontsize=8 width=0.75 shape=box style=dashed]"
	tmp = tmp + "\n\t\trank=same;"
	tmp = tmp + "\n\t\trankdir=LR"
	for i,_ := range(mat[0]){
		tmp=tmp+"\n\t\t\"-1,"+strconv.Itoa(i)+"\" [label=\"G"+strconv.Itoa(i)+"\"]"
	}
	tmp = tmp + "\n\n\t\tedge [dir=none, style=invis]"

	for i:=0;i<len(mat[0])-1;i++{
		tmp = tmp + "\n\t\t\"-1,"+strconv.Itoa(i)+"\" -> \"-1,"+strconv.Itoa(i+1)+"\""
	}
	tmp = tmp+"\t}"
	dot = dot + tmp
	dot = dot + "\n"
	// For loop for all the subgraphs
	for i,row := range(mat){
		tmp = "\n\tsubgraph{"
		tmp = tmp + "\n\t\tnode [margin=0 fontsize=8 width=0.75 shape=box style=invis]"
		tmp = tmp + "\n\t\trank=same;"
		tmp = tmp + "\n\t\trankdir=LR\n"
		for j,el := range(row){
			tmp = tmp + "\n\t\t\""+strconv.Itoa(i)+","+strconv.Itoa(j)+"\" "
			if el != "-"{
				if strings.Contains(el,"Mu") || strings.Contains(el,"RWM"){
					tmp = tmp + "[label=\""+el+"\",style=\"dotted,filled\", fillcolor=green]"
				}else if strings.Contains(el,"Wg"){
					tmp = tmp + "[label=\""+el+"\",style=\"dashed,filled\", fillcolor=gold]"
				}else if strings.Contains(el,"ChSend"){
					tmp = tmp + "[label=\""+el+"\",style=\"bold,filled\", fillcolor=cyan]"
				}else if strings.Contains(el,"ChRecv"){
					tmp = tmp + "[label=\""+el+"\",style=\"bold,filled\", fillcolor=violet]"
				}else{
					tmp = tmp + "[label=\""+el+"\",style=filled]"
				}
			}
		}

		tmp = tmp + "\n\n\t\tedge [dir=none, style=invis]"

		for j:=0;j<len(row)-1;j++{
			tmp = tmp + "\n\t\t\""+strconv.Itoa(i)+","+strconv.Itoa(j)+"\" -> \""+strconv.Itoa(i)+","+strconv.Itoa(j+1)+"\""
		}
		tmp = tmp+"\t}"
		dot = dot + tmp
		dot = dot + "\n"
	}


	//subgraph X
	tmp = "\n\tsubgraph{"
	tmp = tmp + "\n\t\tnode [margin=0 fontsize=8 width=0.75 shape=box style=invis]"
	tmp = tmp + "\n\t\trank=same;"
	tmp = tmp + "\n\t\trankdir=LR"
	for i,_ := range(mat[0]){
		tmp=tmp+"\n\t\t\"x,"+strconv.Itoa(i)+"\""
	}
	tmp = tmp + "\n\n\t\tedge [dir=none, style=invis]"

	for i:=0;i<len(mat[0])-1;i++{
		tmp = tmp + "\n\t\t\"x,"+strconv.Itoa(i)+"\" -> \"x,"+strconv.Itoa(i+1)+"\""
	}
	tmp = tmp+"\t}"
	dot = dot + tmp
	dot = dot + "\n"
	// Edges
	dot = dot + "\n\tedge [dir=none, color=gray88]"
	for j := 0; j<len(mat[0]) ; j++{
		for i:= -1; i<len(mat) ; i++{
			if i == len(mat)-1{
				dot = dot + "\n\t\""+strconv.Itoa(i)+","+strconv.Itoa(j)+"\" -> \"x,"+strconv.Itoa(j)+"\""
			}else{
				dot = dot + "\n\t\""+strconv.Itoa(i)+","+strconv.Itoa(j)+"\" -> \""+strconv.Itoa(i+1)+","+strconv.Itoa(j)+"\""
			}
			dot = dot + "\n"
		}
	}
	dot = dot + "\n}"


	return dot
}

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
			q = `SELECT t2.file,t2.func,t2.line
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
				err1 = res1.Scan(&file,&funct,&line)
				if err1 != nil {
					panic(err1)
				}
				//report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
			}
			report = report + "G"+strconv.Itoa(make_gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"\n"
		} else{ // global declaration of channel
			report = report + "N/A (e.g., created globaly)\n"
		}

		report = report + "Closed? "

		if close_eid != -1{
			// Query to find location of channel make
			q = `SELECT t2.file,t2.func,t2.line
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
				err1 = res1.Scan(&file,&funct,&line)
				if err1 != nil {
					panic(err1)
				}
				//report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
			}
			report = report + "Yes, G"+strconv.Itoa(close_gid)+": "+file+">"+funct+":"+strconv.Itoa(line)+"\n"
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
				 AND (t1.type="EvChSend" OR t1.type="EvChRecv")
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
			q = `SELECT file,func,line
			     FROM StackFrames
					 WHERE eventID=`+strconv.Itoa(id)+" ORDER BY id;"
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

		fmt.Printf("%v\n",report)
		t.RenderMarkdown()
	}
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

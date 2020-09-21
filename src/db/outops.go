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
	"util"
)

func Gtree(dbName, outdir string){
	// Variables
	var q,event            string
	//var report, tmp          string
	//var file, funct          string
	var g,parent,ended,_g     int
	var stkn,stk0     sql.NullString
	//var make_eid, make_gid   int
	//var close_eid, close_gid int
	//var line                 int
	//var val, pos, eid        int*/
	var label   string
	//var q        string
	//var line                string
	//var _arg,_val        			string
	nodes := make(map[int]string) //key: id, val: label
	edges := make(map[int][]int) //key: parent_id, val: child_id
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}

	// Trace size
	q = "SELECT id,gid,parent_id,ended,createLoc,createLoc0 FROM goroutines;"
	res, err := db.Query(q)
	check(err)
	for res.Next(){
		err = res.Scan(&_g,&g,&parent,&ended,&stkn,&stk0)
		check(err)
		label = "[ label = \"{"
		label = label + strconv.Itoa(_g-1)
		label = label + " | "
		if stkn.Valid && stk0.Valid {
			label = label + "bot_stack: "+stkn.String+" \\l top_stack:"+stk0.String+"\\l"
		}else{
			label = label + "bot_stack: - \\l top_stack:-\\l"
		}

		label = label + " | "
		// now add event histogram to the goroutine

		//data := make(map[int]string)  // key: id val: event
		freq := make([]int,num_of_ctgs) //[catX freqs] len:8

		q = "SELECT type FROM events WHERE g="+strconv.Itoa(g)+";"
		res1, err1 := db.Query(q)
		check(err1)
		for res1.Next(){
			err = res1.Scan(&event)
			check(err)

			for k := 0 ; k < num_of_ctgs ; k++{
				if util.Contains(ctgDescriptions[k].Members,event){
					freq[k]++
				}
			}
		}
		res1.Close()
		fmt.Printf("G %v\n",_g-1)
		for k,item := range freq{
			s := fmt.Sprintf("%v:  %v \\l ",ctgDescriptions[k].Category,item)
			fmt.Println(s)
			label = label + s
		}

		// rest
		label = label + "}\""
		if ended != -1{
			label = label + " style=bold ]"
		} else{
			label = label + " style=dashed]"
		}
		nodes[g] = label
		edges[parent] = append(edges[parent],g)
	}
	res.Close()
	out := "digraph{\n\tnode[shape=record,style=filled,fillcolor=gray95]\n\n\t"
	for k,v := range nodes{
		out = out +strconv.Itoa(k) + " " + v + "\n\t"
	}
	out = out + "\n\n\t"
	for k,v := range edges{
		if k != -1{
			for _,vv := range v{
				out = out + strconv.Itoa(k) + " -> " + strconv.Itoa(vv) + "\n\t"
			}
		}
	}
	out = out + "}"
	fdot,err := os.Create(outdir+"/"+dbName+"_gtree.dot")
	check(err)
	fdot.WriteString(out)
	fdot.Close()

	// Create pdf
	_cmd := "dot -Tpdf "+ outdir+"/"+dbName+"_gtree.dot" + " -o " + outdir+"/"+dbName+"_gtree.pdf"

	cmd := exec.Command("dot","-Tpdf",outdir+"/"+dbName+"_gtree.dot","-o",outdir+"/"+dbName+"_gtree.pdf")
	fmt.Printf(">>> Executing %s...\n",_cmd)
	//err = cmd.Run()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    return
	}
	fmt.Println("Result: " + stdout.String())

	fmt.Println(out)
}


func Histogram(t float64, dbName string){
	// Variables
	var q, event             string
	//var report, tmp          string
	//var file, funct          string
	//var g,logclock     int
	//var predG,predClk  sql.NullInt32
	//var make_eid, make_gid   int
	//var close_eid, close_gid int
	//var line                 int
	//var val, pos, eid        int*/
	//var q        string
	//var line                string
	//var _arg,_val        			string
	var length   			int
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}

	// Trace size
	q = "SELECT COUNT(*) FROM events;"
	res, err := db.Query(q)
	check(err)
	if res.Next(){
		err = res.Scan(&length)
		check(err)
	}
	res.Close()

	data := make(map[int]string) // key: id val: event
	freq := make(map[int][]int) // key: chunck# val: [catX freqs] len:8
	q = "SELECT type FROM events;"
	res, err = db.Query(q)
	check(err)
	for res.Next(){
		err = res.Scan(&event)
		check(err)
		data[len(data)+1]=event
	}
	res.Close()
	chunkSize := len(data)/20
	start := 1
	end   := start + chunkSize
	for i:=0 ; i<20 ; i++{
		freq[i] = make([]int,num_of_ctgs)
		start = i * chunkSize + 1
		end   = start + chunkSize
		for j := start ; j < end ; j++{
			for k := 0 ; k < num_of_ctgs ; k++{
				if util.Contains(ctgDescriptions[k].Members, data[j]){
					freq[i][k]++
				}
			}
		}
		fmt.Printf("chunk# %v\n",i)
		for k,item := range freq[i]{
			fmt.Printf("\t%v:  %v\n",ctgDescriptions[k].Category,item)
		}
	}

	// for i=0 ... binSize:
	//		calculate start & end
	//    initate the freq data structure
	//    count numbers
	//    for map[binIDX(start-end)]=[vector of counts per category]
	//    format for visualization
	db.Close()
}

//func Dev(dbName,hbtable, outdir string){
func Dev(){
	// Variables
	//var q, event             string
	//var report, tmp          string
	//var file, funct          string
	//var g,logclock     int
	//var predG,predClk  sql.NullInt32
	//var make_eid, make_gid   int
	//var close_eid, close_gid int
	//var line                 int
	//var val, pos, eid        int*/
	var q        string
	var line                string
	//var _arg,_val        			string
	var cnt   			int

	for i:=0 ; i < 24 ; i++{
		dbName := "fft2X"+strconv.Itoa(i)
		line = ""
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
		if err != nil {
			fmt.Println(err)
		}

		// Trace size
		q = "SELECT COUNT(*) FROM events;"
		res, err := db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()

		// PROC
		q = "SELECT COUNT(*) FROM events inner join global.catPROC on type=eventName;"
		res, err = db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()

		// GCMM
		q = "SELECT COUNT(*) FROM events inner join global.catGCMM on type=eventName;"
		res, err = db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()

		// WGRP
		q = "SELECT COUNT(*) FROM events inner join global.catWGRP on type=eventName;"
		res, err = db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()

		// WGRP - wait
		q = "SELECT COUNT(*) FROM events where type =\"EvWgWait\";"
		res, err = db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()

		// WGRP - add
		q = "SELECT COUNT(*) FROM events where type =\"EvWgAdd\";"
		res, err = db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()

		// WGRP - done
		q = "SELECT COUNT(*) FROM events where type =\"EvWgDone\";"
		res, err = db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()

		// CH-send
		q = "SELECT COUNT(*) FROM events where type =\"EvChSend\";"
		res, err = db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()

		// CH-Recv
		q = "SELECT COUNT(*) FROM events where type =\"EvChRecv\";"
		res, err = db.Query(q)
		check(err)
		if res.Next(){
			err = res.Scan(&cnt)
			check(err)
			line = line + strconv.Itoa(cnt) + ","
		}
		res.Close()
		fmt.Println(dbName+","+line)
		db.Close()
	}
}

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

func ChannelReport(dbName, outdir string){

	// Variables
	var q, event             string
	var tmp          string
	var file, funct          string
	var createDesc,closeDesc string
	var id, cid, ts, gid     int
	var make_eid, make_gid   int
	var close_eid, close_gid int
	var line                 int
	var val, pos, eid        int


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
	chsendrecvStmt,err := db.Prepare(`SELECT t1.id, t1.type, t1.ts, t1.g FROM Events t1 WHERE t1.rid=? AND (t1.type="EvChSend" OR t1.type="EvChRecv") ORDER BY t1.ts`)
	check(err)
	defer chsendrecvStmt.Close()

	valStmt,err := db.Prepare("SELECT value from args where eventID=? and arg=\"val\"")
	check(err)
	defer valStmt.Close()

	eidStmt,err := db.Prepare("SELECT value from args where eventID=? and arg=\"cheid\"")
	check(err)
	defer eidStmt.Close()


	posStmt,err := db.Prepare("SELECT value from args where eventID=? and arg=\"pos\"")
	check(err)
	defer posStmt.Close()

	// now find stack entry for current row
	stackEntryStmt,err := db.Prepare("SELECT file,func,line FROM StackFrames WHERE eventID=? ORDER BY id")
	check(err)
	defer stackEntryStmt.Close()

	mdtab := ""

	// Generate report for each channel
	for res.Next(){
		err = res.Scan(&id,&cid,&make_eid,&make_gid,&close_eid,&close_gid)
		if err != nil{
			panic(err)
		}
		commTypes := make(map[int][]int) // commTypes[gid] = []10 categories of messages

		if make_eid != -1{
			res1, err1 := chmakeLocStmt.Query(cid)
			check(err1)
			for res1.Next(){
				err1 = res1.Scan(&file,&funct,&line)
				check(err1)
				//report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
			}
			res1.Close()
			createDesc = "G"+strconv.Itoa(make_gid)+"<br>"+file+"<br>"+funct+":"+strconv.Itoa(line)
		}

		if close_eid != -1{
			res1, err1 := chcloseLocStmt.Query(cid)
			check(err1)
			for res1.Next(){
				err1 = res1.Scan(&file,&funct,&line)
				check(err1)

				//report = report + "G"+strconv.Itoa(make_gid)+": "+file+" >> "+funct+"\n"
			}
			res1.Close()
			closeDesc = "G"+strconv.Itoa(close_gid)+"<br>"+file+"<br>"+funct+":"+strconv.Itoa(line)
		}

		// now generate table
		//fmt.Printf("Executing: %v\n",q)
		res1, err1 := chsendrecvStmt.Query("C"+strconv.Itoa(cid))
		check(err1)

		rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"TS","Send", "Recv"})

		detail_table := table.NewWriter()
		detail_table.SetOutputMirror(os.Stdout)

    detail_table.AppendHeader(table.Row{"Channel "+strconv.Itoa(cid),"Creates","Send","Send","Send","Send","TOT Send","Recv","Recv","Recv","Recv","Recv","TOT Recv","Close","Total"}, rowConfigAutoMerge)
    detail_table.AppendHeader(table.Row{"","","vacant","blocked","recv-ready","select","","onClose","direct","blocked","send-ready","select","","",""})

		mdtab = mdtab + "|Channel "+strconv.Itoa(cid)+"|Creates|Send|Send|Send|Send|TOT Send|Recv|Recv|Recv|Recv|Recv|TOT Recv|Close|Total|\n"
		mdtab = mdtab + "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|\n"
		mdtab = mdtab + "|||vacant|blocked|recv-ready|select||onClose|direct|blocked|send-ready|select||||\n"

		mdtab2 := "|TS|Send|Recv|\n"
		mdtab2 = mdtab2 +  "|---|---|---|\n"

		for res1.Next(){
			err1 = res1.Scan(&id,&event,&ts,&gid)
			check(err1)

			// add g to commTypes map
			if _,ok := commTypes[gid]; !ok{
				commTypes[gid] = []int{0,0,0,0,0,0,0,0,0}
			}

			// find stack entry for current row
		 	res2, err2 := stackEntryStmt.Query(id)
		 	check(err2)
			for res2.Next(){
				err2 = res2.Scan(&file,&funct,&line)
				if err2 != nil{
					panic(err2)
				}
			}
			res2.Close()


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

			res6,err := eidStmt.Query(id)
			check(err)
			if res6.Next(){
				err := res6.Scan(&eid)
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
			//tmp = "G"+strconv.Itoa(gid)+": "+file+":"+funct+":"+strconv.Itoa(line)+":"+strconv.Itoa(val)+"#"+strconv.Itoa(eid)+"("+strconv.Itoa(pos)+")"
			tmp = "G"+strconv.Itoa(gid)+": "+file+":"+funct+":"+strconv.Itoa(line)+"("+descOf(pos,event)+")"
			res4.Close()
			res5.Close()
			res6.Close()

			fmt.Println(event,pos,indexOf(pos,event))
			commTypes[gid][indexOf(pos,event)]++
			mdtab2 = mdtab2 + "|"+strconv.Itoa(ts)+"|"

			var row []interface{}
			row = append(row,ts)

			if event == "EvChSend"{
				row = append(row,tmp)
				row = append(row,"-")
				mdtab2 = mdtab2 +tmp+"|-|\n"
			}else{
				row = append(row,"-")
				row = append(row,tmp)
				mdtab2 = mdtab2 +"-|"+tmp+"|\n"
			}
			t.AppendRow(row)
		}


		rowTotSend := make(map[int]int) // rowtot[g] = total of row g
		rowTotRecv := make(map[int]int) // rowtot[g] = total of row g
		colTot := make(map[int]int) // rowtot[g] = total of row g

		for idx := 0 ; idx < 12 ; idx++{
			colTot[idx]=0
		}
		for k,v := range commTypes{
			// clear row
			var row []interface{}
			//row=row[:0] // clear row

			// init rowtot
			rowTotSend[k] = 0
			rowTotRecv[k] = 0

			//G
			row = append(row,"G"+strconv.Itoa(k))
			mdtab = mdtab + "|" + "G"+strconv.Itoa(k)

			// Make
			if k == make_gid {
				row = append(row,createDesc)
				mdtab = mdtab + "|" + createDesc
			} else{
				row = append(row,"-")
				mdtab = mdtab + "|-"
			}

			// vacant
			row = append(row,v[0])
			mdtab = mdtab + "|" + strconv.Itoa(v[0])
			rowTotSend[k] += v[0]
			colTot[0] += v[0]

			// s-blocked
			row = append(row,v[1])
			mdtab = mdtab + "|" + strconv.Itoa(v[1])
			rowTotSend[k] += v[1]
			colTot[1] += v[1]

			// recv-ready
			row = append(row,v[2])
			mdtab = mdtab + "|" + strconv.Itoa(v[2])
			rowTotSend[k] += v[2]
			colTot[2] += v[2]

			// s-select
			row = append(row,v[3])
			mdtab = mdtab + "|" + strconv.Itoa(v[3])
			rowTotSend[k] += v[3]
			colTot[3] += v[3]

			// total send
			row = append(row,rowTotSend[k])
			mdtab = mdtab + "|" + strconv.Itoa(rowTotSend[k])
			colTot[4] += rowTotSend[k]

			// onClose
			row = append(row,v[4])
			mdtab = mdtab + "|" + strconv.Itoa(v[4])
			rowTotRecv[k] += v[4]
			colTot[5] += v[4]

			// direct
			row = append(row,v[5])
			mdtab = mdtab + "|" + strconv.Itoa(v[5])
			rowTotRecv[k] += v[5]
			colTot[6] += v[5]

			// r-blocked
			row = append(row,v[6])
			mdtab = mdtab + "|" + strconv.Itoa(v[6])
			rowTotRecv[k] += v[6]
			colTot[7] += v[6]

			// send-ready
			row = append(row,v[7])
			mdtab = mdtab + "|" + strconv.Itoa(v[7])
			rowTotRecv[k] += v[7]
			colTot[8] += v[7]

			// select
			row = append(row,v[8])
			mdtab = mdtab + "|" + strconv.Itoa(v[8])
			rowTotRecv[k] += v[8]
			colTot[9] += v[8]

			// total recv
			row = append(row,rowTotRecv[k])
			mdtab = mdtab + "|" + strconv.Itoa(rowTotRecv[k])
			colTot[10] += rowTotRecv[k]

			// Close
			if k == close_gid {
				row = append(row,closeDesc)
				mdtab = mdtab + "|" +  closeDesc
			} else{
				row = append(row,"-")
				mdtab = mdtab + "|-"
			}

			// total
			row = append(row,rowTotRecv[k]+rowTotSend[k])
			mdtab = mdtab + "|" +  strconv.Itoa(rowTotRecv[k]+rowTotSend[k])
			colTot[11] += rowTotRecv[k]+rowTotSend[k]

			mdtab = mdtab + "|\n"
			detail_table.AppendRow(row)
			//detail_table.AppendSeparator()
		}

		//row=row[:0] // clear row
		var row []interface{}

		row = append(row,"Total")
		row = append(row,"-")
		mdtab = mdtab + "|Total|-|"
		for idx := 0 ; idx < 11 ; idx++{
			row = append(row,colTot[idx])
			mdtab = mdtab + strconv.Itoa(colTot[idx]) + "|"
		}
		mdtab = mdtab + "-|"+strconv.Itoa(colTot[11])+"|\n\n"
		row = append(row,"-")
		row = append(row,colTot[11])

		detail_table.AppendFooter(row)
		mdtab = mdtab + "\n\n" + mdtab2 + "\n\n"

		fmt.Println(mdtab + "\n\n" + mdtab2 + "\n\n")

		t.Render()
		detail_table.Render()
		res1.Close()

	}

	err=res.Close()
	check(err)

	f, err := os.Create(outdir+"/"+dbName+"_chans.md")
	check(err)
	f.WriteString(mdtab)
	f.Close()
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

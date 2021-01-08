package db

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	_"log"
	"os"
	_"os/exec"
	"strings"
	_"bytes"
	"github.com/jedib0t/go-pretty/table"
	_"util"
	_"text/tabwriter"
)

func Checker(dbName string) {
	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()
	// END DB

	//fmt.Println(appGoroutineFinder(db))
	var gs []int
	var g int
	var event,rid string

	// last events stroe every last event of goroutines
	lastEvents := make(map[int]string)

	// Query mutexes
	q := `select gid from goroutines;`
	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	for res.Next(){
		err = res.Scan(&g)
		check(err)
		gs = append(gs,g) // append g to gs
	}
	res.Close()

	lastEventStmt,err := db.Prepare("SELECT type FROM Events WHERE g=? ORDER BY id DESC LIMIT 1")
	check(err)
	defer lastEventStmt.Close()

	resStmt,err := db.Prepare("SELECT type,rid FROM Events WHERE rid IS NOT NULL AND g=?")
	check(err)
	defer resStmt.Close()


	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Goroutine","Last Event","Resources","Goroutines"})

	for _,gi := range(gs){
		// New row
		var row []interface{}
		row = append(row,gi)

		// Last event
		res,err = lastEventStmt.Query(gi)
		check(err)
		for res.Next(){
			err = res.Scan(&event)
			check(err)
			lastEvents[gi]=event
			//gs = append(gs,g) // append g to gs
			row = append(row,event)
		}

		// Resources
		resMap := make(map[string]int)
		var resources []interface{}
		var otherg []interface{}
		res,err = resStmt.Query(gi)
		check(err)
		for res.Next(){
			err = res.Scan(&event,&rid)
			check(err)
			if _,ok := resMap[rid]; !ok{
				resMap[rid] = 1
				if strings.HasPrefix(rid,"G"){
					otherg = append(otherg,rid)
				}else{
					resources = append(resources,rid)
				}

			}
		}
		row = append(row,resources)
		row = append(row,otherg)

		t.AppendRow(row)
		res.Close()
	}

	t.Render()
	//fmt.Println(lastEvents)
	textReport(lastEvents)
}

func textReport(lastEvents map[int]string){
	//writer := tabwriter.NewWriter(os.Stdout,0 , 16, 1, '\t', tabwriter.AlignRight)
	totalG := len(lastEvents)
	var suspicious []int
	var   isGlobalDL   bool
	var   numWaiting   int
	var   numBlock     int
	var   numEnd       int

	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	for k,v := range(lastEvents){
		if k == 1 && v != "EvGoSched"{
			isGlobalDL = true
			continue
		}
		switch v {
		case "EvGoWaiting":
			numWaiting++
		case "EvGoEnd":
			numEnd++
		case "EvGoBlock":
			numBlock++
		default:
			if k != 0 && k != 1{
				// the goroutine is in app goroutine which has not ended!
				suspicious = append(suspicious,k)
			}
		}
	}

	fmt.Println("Total # goroutines:",totalG)
	fmt.Println("Total # runtime goroutines (0,1,tracing):",numWaiting+numBlock+2)
	if isGlobalDL{
		fmt.Println("Global Deadlock:",string(colorRed),"TRUE",string(colorReset))
	} else{
		fmt.Println("Global Deadlock:",string(colorGreen),"FALSE",string(colorReset))
	}

	if len(suspicious) != 0{
		temp := ""
		for _,i := range(suspicious){
			temp = temp + strconv.Itoa(i) + " "
		}
		fmt.Println("Leaked Goroutines:",string(colorRed),temp,string(colorReset))
	} else{
		fmt.Println("Leaked Goroutines:",string(colorGreen),"NONE",string(colorReset))
	}
}
/*func ResourceGraph(dbName, outdir string){

	// Variables
	var q1,q2, event             string
	var rid          string
	var  g     int
	var pos int

	edges := make(map[string][]string)
	waitingEdges := make(map[string][]string)
	nodes := make(map[string]int)




	res, err = db.Query(q2)
	if err != nil {
		panic(err)
	}
	for res.Next(){
		err = res.Scan(&event,&g,&rid,&pos)
		if err != nil{
			panic(err)
		}
		if rid != "M3"{ // trace lock, ignore it
			if event == "EvChSend"{
				if pos != 0{
					edges["G"+strconv.Itoa(g)] = append(edges["G"+strconv.Itoa(g)],rid)
				}else{
					waitingEdges["G"+strconv.Itoa(g)] = append(waitingEdges["G"+strconv.Itoa(g)],rid)
				}
			} else{
				if pos != 0{
					edges[rid] = append(edges[rid],"G"+strconv.Itoa(g))
				}else{
					waitingEdges[rid] = append(waitingEdges[rid],"G"+strconv.Itoa(g))
				}
			}
			if _,ok := nodes["G"+strconv.Itoa(g)] ; !ok{
				nodes["G"+strconv.Itoa(g)] = 1
			}
			if _,ok := nodes[rid] ; !ok{
				nodes[rid] = 1
			}
		}

	}
	res.Close()


	nodes_st := ""
	for k,_ := range nodes{
		if strings.HasPrefix(k,"G"){
			nodes_st = nodes_st + k + " [label = \""+k+"\" shape=circle]\n\t"
		}else if strings.HasPrefix(k,"C"){
			nodes_st = nodes_st + k + " [label = \""+k+"\" shape=diamond style=bold]\n\t"
		} else{
			nodes_st = nodes_st + k + " [label = \""+k+"\" shape=invtriangle style=bold]\n\t"
		}
	}

	edges_st := ""
	for k,v := range edges{
		freq := make(map[string]int)
		for _,d := range v{
			freq[d]++
		}
		for kk,vv := range freq{
			edges_st = edges_st + k + " -> " + kk + " [label=\""+strconv.Itoa(vv)+"\"]\n\t"
		}
	}

	edges_st = edges_st + "\n\n"

	for k,v := range waitingEdges{
		freq := make(map[string]int)
		for _,d := range v{
			freq[d]++
		}
		for kk,vv := range freq{
			edges_st = edges_st + k + " -> " + kk + " [label=\""+strconv.Itoa(vv)+"\" style=dashed]\n\t"
		}
	}

	fmt.Println(nodes_st)
	fmt.Println(edges_st)


	f, err := os.Create(outdir+"/"+dbName+"_resg.dot")
	check(err)
	f.WriteString("digraph {\n\t"+nodes_st+"\n\n\t"+edges_st+"\n}")
	f.Close()

	// Create pdf
	_cmd := "dot -Tpdf "+ outdir+"/"+dbName+"_resg.dot" + " -o " + outdir+"/"+dbName+"_resg.pdf"

	cmd := exec.Command("dot","-Tpdf",outdir+"/"+dbName+"_resg.dot","-o",outdir+"/"+dbName+"_resg.pdf")
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


	_cmd = "open" + outdir+"/"+dbName+"_resg.pdf"

	cmd = exec.Command("open",outdir+"/"+dbName+"_resg.pdf")
	fmt.Printf(">>> Executing %s...\n",_cmd)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    return
	}
	fmt.Println("Result: " + stdout.String())
	//fmt.Println(out)
}*/

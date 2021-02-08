package db

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"log"
	"os"
	"strings"
	"github.com/jedib0t/go-pretty/table"

)

func longLeakReport(dbName string) {
	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	}else{
		log.Println("Cheker(long): Connected to ",dbName)
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

func Checker(dbName string, long bool){
	if long {
		longLeakReport(dbName)
	} else{
		shortLeakReport(dbName)
	}
}

func shortLeakReport(dbName string) {
	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	}else{
		log.Println("Cheker(short): Connected to ",dbName)
	}
	defer db.Close()
	// END DB

	//fmt.Println(appGoroutineFinder(db))
	var gs []int
	var g int
	var event string

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

	for _,gi := range(gs){
		// Last event
		res,err = lastEventStmt.Query(gi)
		check(err)
		for res.Next(){
			err = res.Scan(&event)
			check(err)
			lastEvents[gi]=event
		}
		res.Close()
	}

	// ****************

	var suspicious []int
	var   isGlobalDL   bool
	var   numIgnore   int

	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"


	for k,v := range(lastEvents){
		if k == 1 && v != "EvGoSched"{
			isGlobalDL = true
			continue
		}
		switch v {
		case "EvGoWaiting","EvGoEnd","EvGoBlock":
			numIgnore++
		default:
			if k != 0 && k != 1{
				// the goroutine is in app goroutine which has not ended!
				suspicious = append(suspicious,k)
			}
		}
	}

	if isGlobalDL{
		fmt.Println(string(colorRed),"Fail (global deadlock)",string(colorReset))
	} else if len(suspicious) != 0{
		fmt.Println(string(colorRed),"Fail (partial deadlock - leak)",string(colorReset))
	} else{
		fmt.Println(string(colorGreen),"Pass",string(colorReset))
	}
}

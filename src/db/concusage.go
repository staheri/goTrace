package db

import (
	_"bytes"
	"database/sql"
	_"fmt"
	_ "github.com/go-sql-driver/mysql"
	_"github.com/jedib0t/go-pretty/table"
	"log"
	_"os"
	_"os/exec"
	"strconv"
	"strings"
	"util"
	_"sort"
)

/*type ConUse struct{
	g              int
	event          string
	badSelect      bool
	blocking       bool
	unblocking     bool
	callStack      []*CallSite
	modSite        string
	schedSite      string
}

type CallSite struct{
	file     string
	func     string
	line     int
}
*/

type ConUse struct{
	File         string
	Funct        string
	Rid          string
	Line         string
	Event        string
	G            uint64
}

func ConcUsage(dbName string) map[string]int {

	// Variables
	var id                          int
	var g                           uint64
	var event,file,funct,line,rid   string
	var _rid                        sql.NullString
	var ignored                     []string

	concUsage := make(map[string]int)
	fullstack := make(map[int][]string) // key: eventid, value: slice of stacks
	blacklist := make(map[string]int)

	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	} else {
		log.Println("ConcUsage: Connected to ",dbName)
	}
	defer db.Close()
	// END DB

	// Query catSCHD
	q := `SELECT t1.id,t1.type,t1.g,t1.rid,t3.file,t3.func,t3.line FROM Events t1
				INNER JOIN global.catSCHD t2
				ON t1.type=t2.eventName
				INNER JOIN Stackframes t3 on t1.id=t3.eventid;`

	res, err := db.Query(q)
	check(err)

	// to ignore
	r2ignoreStmt,err := db.Prepare("select file,func from stackframes where eventid=?;")
	check(err)
	defer r2ignoreStmt.Close()

	// store fullstack first
	for res.Next() {
		err = res.Scan(&id, &event, &g, &_rid, &file, &funct, &line)
		check(err)
		if _rid.Valid{
			rid = _rid.String
		}else{
			rid = ""
		}

		//fmt.Println(ignored)
		//fmt.Println(id,event)
		if !util.Contains(ignored,rid){
			if toIgnore,isIgnore := ridToIgnore(r2ignoreStmt,rid,id);isIgnore && event!="EvGoCreate"{
				ignored = append(ignored,toIgnore)
				continue
			}
			if !isBadSelect(db,event,id){
				fullstack[id] = append(fullstack[id], file+":"+funct+":"+line)
			}else{
				y,_ := strconv.Atoi(line)
				fullstack[id] = append(fullstack[id], file+":"+funct+":"+strconv.Itoa(y-1))
				blacklist[file+":"+line] = 1
			}
		}
	}
	res.Close()

	// now iterate over full stack to find the last in-source location
	for _, v := range fullstack {
		t := strings.Split(v[len(v)-1], ":")
		src := t[0]
		// start from end, once source changes, break!
		for i := len(v) - 2; i >= 0; i-- {
			t1 := strings.Split(v[i], ":")
			if src != t1[0] {
				break
			}
			t = t1
			src = t[0]
		}

		// add loc to concusage
		// check if it is not in the black list
		loc := t[0] + ":" + t[2]
		if _, ok := concUsage[loc]; !ok {
			if _,ok2 := blacklist[loc]; !ok2{
				if rid != "" {
					concUsage[loc] = 1
				} else {
					concUsage[loc] = 0
				}
			}
		}
	}
	return concUsage
}

func ConcUsageStruct(dbName string)  []*ConUse {
	// Variables
	var id                          int
	var g,gc                        uint64
	var event,file,funct,line,rid   string
	//var mainStartLocation           string
	var _rid                        sql.NullString
	var ret                         []*ConUse
	var ignored                     []string
	var whiteList                     []string

	concUsage    := make(map[string]int)
	fullstack    := make(map[int][]string) // key: eventid, value: slice of stacks
	blacklist    := make(map[string]int)
	conuseStruct := make(map[int]*ConUse)


	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	} else {
		log.Println("ConcUsage: Connected to ",dbName)
	}
	defer db.Close()
	// END DB

	// Query catSCHD
	q := `SELECT t1.id,t1.type,t1.g,t1.rid,t3.file,t3.func,t3.line FROM Events t1
				INNER JOIN Stackframes t3 on t1.id=t3.eventid
				INNER JOIN global.catSCHD t2
				ON t1.type=t2.eventName;`

	res, err := db.Query(q)
	check(err)


	// to ignore

	//defer


	//defer

	// store fullstack first
	for res.Next() {
		err = res.Scan(&id, &event, &g, &_rid, &file, &funct, &line)
		check(err)
		if _rid.Valid{
			rid = _rid.String
		}else{
			if event == "EvGoCreate"{
				goCreateGStmt,err1 := db.Prepare("select value from args where eventid=? and arg=\"g\";")
				check(err1)
				res1,err1 := goCreateGStmt.Query(id)
				check(err1)
				if res1.Next(){
					err1 = res1.Scan(&gc)
					check(err1)
					rid = "G"+strconv.Itoa(int(gc))
				}
				res1.Close()
				goCreateGStmt.Close()
			} else if event == "EvSelect"{
				rid = "Se"
			}

		}

		//fmt.Println(ignored)
		//fmt.Println(id,event)
		if !util.Contains(ignored,rid){
			if !util.Contains(whiteList,rid){
				r2ignoreStmt,err1 := db.Prepare("select file,func from stackframes where eventid=?;")
				check(err1)
				if toIgnore,isIgnore := ridToIgnore(r2ignoreStmt,rid,id);isIgnore && event!="EvGoCreate"{
					ignored = append(ignored,toIgnore)
					r2ignoreStmt.Close()
					continue
				} else{
					whiteList = append(whiteList,rid)
				}
				r2ignoreStmt.Close()
			}

			if !isBadSelect(db,event,id){
				fullstack[id] = append(fullstack[id], file+":"+funct+":"+line)
			}else{
				y,_ := strconv.Atoi(line)
				fullstack[id] = append(fullstack[id], file+":"+funct+":"+strconv.Itoa(y-1))
				blacklist[file+":"+line] = 1
				//cu := ConUse{file:file,funct:funct,line:strconv.Itoa(y-1),g:g,event:event,rid:rid}
				//conuseStruct[id] = append(conuseStruct[id],&cu)
			}
			if _,ok := conuseStruct[id]; !ok{
				//fmt.Println(">>>>> ", g,event,rid)
				conuseStruct[id] = &ConUse{G:g,Event:event,Rid:rid}
			}
		}
	}
	res.Close()
	//goCreateGStmt.Close()
	//r2ignoreStmt.Close()

	// now iterate over full stack to find the last in-source location
	for k, v := range fullstack {
		t := strings.Split(v[len(v)-1], ":")
		src := t[0]
		// start from end, once source changes, break!
		for i := len(v) - 2; i >= 0; i-- {
			t1 := strings.Split(v[i], ":")
			if src != t1[0] {
				break
			}
			t = t1
			src = t[0]
		}

		// add loc to concusage
		// check if it is not in the black list
		loc := t[0] + ":" + t[2]

		if _, ok := concUsage[loc]; !ok {
			if _,ok2 := blacklist[loc]; !ok2{
				conuseStruct[k].File = t[0]
				conuseStruct[k].Funct = t[1]
				conuseStruct[k].Line = t[2]
				ret = append(ret,conuseStruct[k])
				if rid != "G"{
					concUsage[loc] = 1
				} else {
					concUsage[loc] = 0
				}
			}
		}
	}
	return ret
}

package db

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
  "os"
  "log"
)


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

func mat2dot(mat [][]string, header []string) string{

	width := "2"
	fontsize := "11"

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
	tmp = tmp + "\n\t\tnode [margin=0 fontsize="+fontsize+" width="+width+" shape=box style=dashed]"
	tmp = tmp + "\n\t\trank=same;"
	tmp = tmp + "\n\t\trankdir=LR"
	for i,g := range(header){
		tmp=tmp+"\n\t\t\"-1,"+strconv.Itoa(i)+"\" [label=\""+g+"\"]"
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
		tmp = tmp + "\n\t\tnode [margin=0 fontsize="+fontsize+" width="+width+" shape=box style=invis]"
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
	tmp = tmp + "\n\t\tnode [margin=0 fontsize="+fontsize+" width="+width+" shape=box style=invis]"
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
			fmt.Println(dbs.String)
			if dbs.String != "XXX" && strings.Split(strings.Split(dbs.String,":")[1],".")[0] == "main" {
				// this is the app goroutine
				// add it to slice
				appGss = append(appGss,gid)
			}
		}
	}
	res.Close()

	return appGss
}

func indexOf(pos int, event string) int {
	switch event {
	case "EvChSend":
		return pos-1
	case "EvChRecv":
		if pos == 4 || pos == 5 {
			return 7
		}else if  pos == 6 || pos == 7 {
			return 8
		}else{
			return pos+3
		}
	default:
		return 0
	}
}

func descOf(pos int, event string) string {
	switch event {
	case "EvChSend":
		switch pos {
		case 1:
			return "vacant"
		case 2:
			return "blocked"
		case 3:
			return "recv-ready"
		case 4:
			return "select"
		default:
			return "unknown"
		}
	case "EvChRecv":
		switch pos {
		case 1:
			return "onClose"
		case 2:
			return "direct"
		case 3:
			return "blocked"
		case 4:
			return "send-ready"
		case 5:
			return "send-ready"
		case 6:
			return "select"
		case 7:
			return "select"
		default:
			return "unknown"
		}
	default:
		return "unknown"
	}
}

// check if event,id is bad select (pos=1,2,3)
func isBadSelect(db *sql.DB, event string, id int) (bool){
	if event != "EvSelect"{
		return false
	}

	value := 0
	q := `SELECT value FROM args WHERE arg="pos" and eventid=`+strconv.Itoa(id)+`;`
	res,err := db.Query(q)
	check(err)
	defer res.Close()
	if res.Next(){
		err = res.Scan(&value)
		if value != 0 {
			return true
		}
	}
	res.Close()
	return false
}

// check if the rid is one of print, trace, rand, reschedule()
func ridToIgnore(prep *sql.Stmt,rid string,id int) (string,bool){
	var file,funct string
	//q := `select file,func from stackframes where eventid=?;`
	//fmt.Println("ridToIgnore > ",rid)

	if strings.HasPrefix(rid,"G0"){ // if rid is G0, do not ignore
		return "",false
	}
	res,err := prep.Query(id)
	check(err)
	for res.Next(){
		err = res.Scan(&file,&funct)
		if file == "rand.go" && funct == "math/rand.(*Rand).Intn"{ // random lock
			//fmt.Println("******* rand ")
			return rid,true
		}
		if file == "print.go" && strings.Split(funct,".")[0] == "fmt"{ // print lock
			//fmt.Println("******* print ")
			return rid,true
		}
		if funct == "main.Reschedule"{ // reschedule lock
			//fmt.Println("******* reschedule ")
			return rid,true
		}
		if funct == "runtime/trace.Stop"{ // trace lock
			return rid,true
		}
		if strings.HasPrefix(funct,"runtime/trace"){ // trace lock
			//fmt.Println("******* trace ")
			return rid,true
		}
	}
	res.Close()
	//fmt.Println(">>>>>>>>> PASSED ")
	return "",false
}

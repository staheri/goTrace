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

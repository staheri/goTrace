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

func Checker(dbName string) {
	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()

	// Query mutexes
	q = `SELECT type,g,rid
	     FROM Events
			 WHERE type="EvMuUnlock" OR type="EvMuLock"`

	res, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	for res.Next(){
		err = res.Scan(&event,&g,&rid)
		if err != nil{
			panic(err)
		}
		if rid != "M3"{ // trace lock, ignore it
			if event == "EvMuLock"{
				edges["G"+strconv.Itoa(g)] = append(edges["G"+strconv.Itoa(g)],rid)
			} else{
				edges[rid] = append(edges[rid],"G"+strconv.Itoa(g))
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



}
func ResourceGraph(dbName, outdir string){

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
}

package db

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"log"
	"os"
	"os/exec"
	"bytes"
	"strings"

)

func DIFF(dbName, baseDBName, cloutpath,resultpath string, consec, rid int, aspects ...string ){
	var cnt int
	var tmp string
	var _rid sql.NullString
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

	// Establish connection to baseDB
	dbBase, errBase := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+baseDBName)
	if err != nil {
		fmt.Println(errBase)
	}else{
		fmt.Println("Connection Established")
	}
	defer dbBase.Close()

	var q,subq,event string
	var id int

	data     := make(map[int][]string)
	dataBase := make(map[int][]string)


	q = `SELECT (t1.id)-1, t2.type, t2.rid
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

	resBase, errBase := dbBase.Query(q)
	if errBase != nil {
		panic(errBase)
	}

	// create directory
	cloutdir := cloutpath + "/diff_" +dbName + "/"
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
	optionsDirName := "_c"+strconv.Itoa(consec)+"_a"+strconv.Itoa(rid)
	cloutdir = cloutdir +filts + optionsDirName
	if _, err := os.Stat(cloutdir); os.IsNotExist(err) {
    os.MkdirAll(cloutdir, 0755)
	}

	// Parse results
 	for res.Next(){
		err = res.Scan(&id,&event,&_rid)
		if err != nil{
			panic(err)
		}
		//if val,ok := data[id];ok{
		if rid > 0 && _rid.Valid{
			if !strings.HasPrefix(_rid.String, "G"){
				data[id] = append(data[id],_rid.String+":"+event)
			}else{
				data[id] = append(data[id],event)
			}
		} else{
			data[id] = append(data[id],event)
		}
	}

	// Parse base results
	for resBase.Next(){
		errBase = resBase.Scan(&id,&event,&_rid)
		if errBase != nil{
			panic(errBase)
		}
		//if val,ok := data[id];ok{
		//dataBase[id] = append(dataBase[id],event)
		if rid > 0 && _rid.Valid{
			if !strings.HasPrefix(_rid.String, "G"){
				dataBase[id] = append(dataBase[id],_rid.String+":"+event)
			}else{
				dataBase[id] = append(dataBase[id],event)
			}
		} else{
			dataBase[id] = append(dataBase[id],event)
		}
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
		cnt = 0
		tmp = ""
		for _,e := range v{
			fmt.Printf("%v\n\t",e)
			tmp = tmp + e + "-"
			cnt = cnt + 1
			//f.WriteString(fmt.Sprintf("%v\n",e))
			if cnt % consec == 0{
				//fmt.Printf("%v\n\t",tmp)
				f.WriteString(fmt.Sprintf("%v\n",tmp))
				cnt = 0
				tmp = ""
			}
		}
		if tmp != ""{
			f.WriteString(fmt.Sprintf("%v\n",tmp))
		}
		f.Close()
	}

	// store base files in the outpath folder
	for k,v := range dataBase{
		output := cloutdir+"/base_g"+strconv.Itoa(k)+".txt"
		f,err := os.Create(output)
		if err != nil{
			log.Fatal(err)
		}
		fmt.Printf("\ndata[%v]:\n\t",k)
		cnt = 0
		tmp = ""
		for _,e := range v{
			fmt.Printf("%v\n\t",e)
			tmp = tmp + e + "-"
			cnt = cnt + 1
			//f.WriteString(fmt.Sprintf("%v\n",e))
			if cnt % consec == 0{
				//fmt.Printf("%v\n\t",tmp)
				f.WriteString(fmt.Sprintf("%v\n",tmp))
				cnt = 0
				tmp = ""
			}
		}
		if tmp != ""{
			f.WriteString(fmt.Sprintf("%v\n",tmp))
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
	_cmd = "python "+ HACPATH + "/main.py " + cloutdir+"/cl/diff_"+dbName+"_"+filts+optionsDirName+".dot "+resultpath+"/"+dbName+"_"+filts+optionsDirName

	cmd = exec.Command("python",HACPATH + "/main.py",cloutdir+"/cl/diff_"+dbName+"_"+filts+optionsDirName+".dot",resultpath+"/"+dbName+"_"+filts+optionsDirName)
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

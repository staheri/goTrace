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
	"math"

)

var(
	GOPATH    string
	CLPATH    string
	HACPATH   string
)

func HAC(dbName, cloutpath,resultpath string, consec, atrmode int, aspects ...string ){
	// Paths
	setPaths()

	// create directory
	cloutdir := cloutpath + "/hac_" +dbName + "/"
	filts := aspect2string(aspects...)
	optionsDirName := "_c"+strconv.Itoa(consec)+"_a"+strconv.Itoa(atrmode)
	cloutdir = cloutdir +filts + optionsDirName
	if _, err := os.Stat(cloutdir); os.IsNotExist(err) {
    os.MkdirAll(cloutdir, 0755)
	}


	// genClContext extracts events by querying database, create folders and store object files
	genClContext(dbName,cloutdir,"N", consec, atrmode, aspects...)

	// Execute C++ cl on cloutdir
	_cmd := CLPATH + "/cltrace -m 1 -p "+cloutdir
	cmd := exec.Command(CLPATH + "/cltrace","-m","1","-p",cloutdir)
	fmt.Printf(">>> Executing %s...\n",_cmd)
	err := cmd.Run()
	if err != nil{
		log.Fatal(err)
	}

	outname := "hac_"+dbName+"_"+filts+optionsDirName

	_cmd = "python "+ HACPATH + "/main.py " + cloutdir+"/cl/"+outname+".dot "+resultpath+"/"+outname

	cmd = exec.Command("python",HACPATH + "/main.py",cloutdir+"/cl/"+outname+".dot",resultpath+"/"+outname)
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

func JointHAC(dbName,baseDBName, cloutpath,resultpath string, consec, atrmode int, aspects ...string ){
	// Paths
	setPaths()

	// create directory
	cloutdir := cloutpath + "/diff_" +dbName + "/"
	filts := aspect2string(aspects...)
	optionsDirName := "_c"+strconv.Itoa(consec)+"_a"+strconv.Itoa(atrmode)
	cloutdir = cloutdir +filts + optionsDirName
	if _, err := os.Stat(cloutdir); os.IsNotExist(err) {
    os.MkdirAll(cloutdir, 0755)
	}

	// extracting
	genClContext(baseDBName,cloutdir,"G", consec, atrmode, aspects...)
	genClContext(dbName,cloutdir,"B", consec, atrmode, aspects...)

	// Execute C++ cl on cloutdir
	_cmd := CLPATH + "/cltrace -m 1 -p "+cloutdir
	cmd := exec.Command(CLPATH + "/cltrace","-m","1","-p",cloutdir)
	fmt.Printf(">>> Executing %s...\n",_cmd)
	err := cmd.Run()
	if err != nil{
		log.Fatal(err)
	}

	outname := "diff_"+dbName+"_"+filts+optionsDirName

	// Execute python hac on cloutdir/cl
	_cmd = "python "+ HACPATH + "/main.py " + cloutdir+"/cl/"+outname+".dot "+resultpath+"/"+outname

	cmd = exec.Command("python",HACPATH + "/main.py",cloutdir+"/cl/"+outname+".dot",resultpath+"/"+outname)
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


func genClContext(dbName, cloutdir,prefix string, consec, atrmode int, aspects ...string ){
	// Paths
	setPaths()

	// Vars
	var q,subq,event,tmp  string
	var cnt            int
	var _rid,crl          sql.NullString
	var crlid             sql.NullInt32

	// hold query results
	data := make(map[string][]string)

	// Establish connection to DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}
	defer db.Close()


	q = `SELECT t1.createLoc,t1.crlid,t2.type,t2.rid
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

 	for res.Next(){
		err = res.Scan(&crl,&crlid,&event,&_rid)
		if err != nil{
			panic(err)
		}
		//if val,ok := data[id];ok{
		//data[id] = append(data[id],event)
		//}else{}
		ids := "root"
		if crl.Valid{
			ids = crl.String+"("+strconv.Itoa(int(crlid.Int32))+")"
		}


		if atrmode == 1 && _rid.Valid{
			if !strings.HasPrefix(_rid.String, "G"){
				data[ids] = append(data[ids],_rid.String+":"+event)
			}else{
				data[ids] = append(data[ids],event)
			}
		} else{
			data[ids] = append(data[ids],event)
		}
	}

	// store files in the outpath folder
	//2nd pass
	for k,v := range data{
		output := cloutdir+"/"+prefix+"_"+k+".txt"
		f,err := os.Create(output)
		if err != nil{
			log.Fatal(err)
		}
		//fmt.Printf("\ndata[%v]:\n\t",k)
		cnt = 0
		tmp = ""
		data2 := make(map[string]float64)
		for _,e := range v{
			//fmt.Printf("%v\n\t",e)
			tmp = tmp + e + "-"
			cnt = cnt + 1
			//f.WriteString(fmt.Sprintf("%v\n",e))
			if cnt % consec == 0{
				//fmt.Printf("%v\n\t",tmp)
				//f.WriteString(fmt.Sprintf("%v\n",tmp))
				if val,ok := data2[tmp];ok{
					data2[tmp] = val + 1
				}else{
					data2[tmp] = 1
				}
				cnt = 0
				tmp = ""
			}
		}
		if tmp != ""{
			//f.WriteString(fmt.Sprintf("%v\n",tmp))
			if val,ok := data2[tmp];ok{
				data2[tmp] = val + 1
			}else{
				data2[tmp] = 1
			}
		}
		fmt.Printf("write > %v\n",output)
		for kk,vv := range data2{
			fmt.Printf("data2[%v]:%v\n",kk,freq(vv,atrmode))
			f.WriteString(fmt.Sprintf("%v:%v\n",kk,int(freq(vv,atrmode))))
		}
		f.Close()
	}
}

func freq(val float64, mode int) (float64){
	if mode == 2{ // exact freq
		return float64(val)
	} else if mode == 3{ // log10
		if val == 1{
			return float64(val)
		}else{
			return math.Log10(val)
		}
	}else if mode == 4{ //log2
		if val == 1{
			return float64(val)
		}else{
			return math.Log2(val)
		}
	}else{
		return 0
	}
}

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

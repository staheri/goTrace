package db

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"os"
)

func DineData(dbName, outdir string, N int, chanOnly, chanID bool){
	// make sure
	outdir = outdir + "/"

	if _, err := os.Stat(outdir); os.IsNotExist(err) {
    os.MkdirAll(outdir, 0755)
	}

	// Establish
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection with DB Established")
	}
	defer db.Close()

	// find phil and fork goroutines
	phils := make(map[int]int) // map[phil_id] = gid
	rphils := make(map[int]int) // map[gid] = phil_id
	forks := make(map[int]int) // map[fork_id] = gid
	rforks := make(map[int]int) // map[gid] = fork_id
	var tmp []int
	var id,gid int

	res,err := db.Query("Select id,gid from goroutines order by id")
	check(err)
	for res.Next(){
		err = res.Scan(&id,&gid)
		check(err)
		tmp = append(tmp,gid)
	}
	res.Close()

	for i:= 0 ; i < N ; i++{
		forks[N-i-1] = tmp[len(tmp)-i-1]
		rforks[tmp[len(tmp)-i-1]] = N-i-1
	}

	for i:= 0 ; i < N ; i++{
		phils[N-i-1] = tmp[len(tmp)-N-i-1]
		rphils[tmp[len(tmp)-N-i-1]] = N-i-1
	}
	//fmt.Println(phils)
	//fmt.Println(forks)

	sepAllStmt, err := db.Prepare("Select type FROM events WHERE g=? order by ts")
	check(err)
	defer sepAllStmt.Close()

	sepChStmt, err := db.Prepare("select events.type from events inner join global.catCHNL t1 on events.type=t1.eventName where g=? order by ts")
	check(err)
	defer sepChStmt.Close()

	sepChWIDStmt, err := db.Prepare("Select events.type, args.value from events left join args on args.eventID=events.id where g=? and args.arg=\"cid\" order by events.ts")
	check(err)
	defer sepChWIDStmt.Close()



	globalAllStmt, err := db.Prepare("select g,type from events order by ts")
	check(err)
	defer globalAllStmt.Close()

	globalChStmt, err := db.Prepare("select g,type from events inner join global.catCHNL t1 on events.type=t1.eventName order by ts")
	check(err)
	defer globalChStmt.Close()

	globalChWIDStmt, err := db.Prepare("select g,type,args.value from events left join args on args.eventID=events.id where args.arg=\"cid\" order by events.ts")
	check(err)
	defer globalChWIDStmt.Close()

	// Variables
	var event string
	var chid  int
	// Separates
	//     Phils
	for pid,gid := range phils{
		output := outdir + "Phil-"+strconv.Itoa(pid)+".txt"
		f,err := os.Create(output)
		check(err)
		if chanOnly{
			if chanID{
				// iterate over rows of sepChWIDStmt
				res,err := sepChWIDStmt.Query(gid)
				check(err)
				for res.Next(){
					err := res.Scan(&event,&chid)
					check(err)
					f.WriteString(event+"-"+strconv.Itoa(chid)+"\n")
				}
				f.Close()
			}else{
				// iterate over rows of sepChStmt
				res,err := sepChStmt.Query(gid)
				check(err)
				for res.Next(){
					err := res.Scan(&event)
					check(err)
					f.WriteString(event+"\n")
				}
				f.Close()
			}
		}else{
			// iterate over rows of sepAllStmt
			res,err := sepAllStmt.Query(gid)
			check(err)
			for res.Next(){
				err := res.Scan(&event)
				check(err)
				f.WriteString(event+"\n")
			}
			f.Close()
		}
	}

	//     Forks
	for pid,gid := range forks{
		output := outdir + "Fork-"+strconv.Itoa(pid)+".txt"
		f,err := os.Create(output)
		check(err)
		if chanOnly{
			if chanID{
				// iterate over rows of sepChWIDStmt
				res,err := sepChWIDStmt.Query(gid)
				check(err)
				for res.Next(){
					err := res.Scan(&event,&chid)
					check(err)
					f.WriteString(event+"-"+strconv.Itoa(chid)+"\n")
				}
				f.Close()
			}else{
				// iterate over rows of sepChStmt
				res,err := sepChStmt.Query(gid)
				check(err)
				for res.Next(){
					err := res.Scan(&event)
					check(err)
					f.WriteString(event+"\n")
				}
				f.Close()
			}
		}else{
			// iterate over rows of sepAllStmt
			res,err := sepAllStmt.Query(gid)
			check(err)
			for res.Next(){
				err := res.Scan(&event)
				check(err)
				f.WriteString(event+"\n")
			}
			f.Close()
		}
	}

	// All
	output := outdir + "global.txt"
	f,err := os.Create(output)
	check(err)
	if chanOnly{
		if chanID{
			// iterate over rows of globalChWIDStmt
			res,err := globalChWIDStmt.Query()
			check(err)
			for res.Next(){
				err := res.Scan(&gid,&event,&chid)
				check(err)
				if vv,ok := rphils[gid];ok{
					f.WriteString("Phil"+strconv.Itoa(vv)+":"+event+"-"+strconv.Itoa(chid)+"\n")
				} else if vv,ok := rforks[gid];ok{
					f.WriteString("Fork"+strconv.Itoa(vv)+":"+event+"-"+strconv.Itoa(chid)+"\n")
				} else{
					continue
				}
			}
			f.Close()
		}else{
			// iterate over rows of globalChStmt
			res,err := globalChStmt.Query()
			check(err)
			for res.Next(){
				err := res.Scan(&gid,&event)
				check(err)
				if vv,ok := rphils[gid];ok{
					f.WriteString("Phil"+strconv.Itoa(vv)+":"+event+"\n")
				} else if vv,ok := rforks[gid];ok{
					f.WriteString("Fork"+strconv.Itoa(vv)+":"+event+"\n")
				} else{
					continue
				}
			}
			f.Close()
		}
	}else{
		// iterate over rows of globalAllStmt
		res,err := globalAllStmt.Query()
		check(err)
		for res.Next(){
			err := res.Scan(&gid,&event)
			check(err)
			if vv,ok := rphils[gid];ok{
				f.WriteString("Phil"+strconv.Itoa(vv)+":"+event+"\n")
			} else if vv,ok := rforks[gid];ok{
				f.WriteString("Fork"+strconv.Itoa(vv)+":"+event+"\n")
			} else{
				continue
			}
		}
		f.Close()
	}
}

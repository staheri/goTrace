package main

import (
	"fmt"
	_"trace"
	_"util"

	_"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"


)

func main(){
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}

	name := "test1"

	_,err = db.Exec("CREATE DATABASE IF NOT EXISTS "+name)
	if err != nil {
			panic(err)
	}

	fmt.Println("Database Created")

	sqlStatement := "USE "+name+";"
	res,err := db.Exec(sqlStatement)
	if err == nil {
 		fmt.Printf("%v\n",res)
		lastId, err1 := res.LastInsertId()
		if err1 != nil {
			log.Fatal(err1)
		}
		rowCnt, err1 := res.RowsAffected()
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
 	}else{
 		fmt.Println("ERRRRRR")
 		panic(err)
 	}

	sqlStatement = `SHOW TABLES;`
	res,err = db.Exec(sqlStatement)
	if err == nil {
 		fmt.Printf("%v\n",res)
		lastId, err1 := res.LastInsertId()
		if err1 != nil {
			log.Fatal(err1)
		}
		rowCnt, err1 := res.RowsAffected()
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
 	}else{
 		fmt.Println("ERRRRRR")
 		panic(err)
 	}


	 /*
   _,err = db.Exec("USE "+name)
   if err != nil {
       panic(err)
   }

  /* _,err = db.Exec("CREATE TABLE example ( id integer, data varchar(32) )")
   if err != nil {
       panic(err)
   }
	 */

	// sqlStatement := `SHOW TABLES`
	 //res,err := db.Exec(sqlStatement)
 	//var table string
 	//if err != nil {
 		/*for res.Next() {
 			res.Scan(&table)
 			fmt.Println(table)
 		}*/
 		//fmt.Printf("%v",res)
 	//}else{
 		//fmt.Println("ERRRRRR")
 		//panic(err)
 	//}

}

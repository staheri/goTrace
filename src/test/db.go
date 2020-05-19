package main

import (
	"fmt"
	_"trace"
	_"util"

	_"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"


)

func main(){
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection Established")
	}

	/*name := "test1"

	_,err = db.Exec("CREATE DATABASE "+name)
   if err != nil {
       panic(err)
   }

   _,err = db.Exec("USE "+name)
   if err != nil {
       panic(err)
   }

   _,err = db.Exec("CREATE TABLE example ( id integer, data varchar(32) )")
   if err != nil {
       panic(err)
   }
	 */

	 sqlStatement := `SHOW TABLES`
	 res,err := db.Exec(sqlStatement)
 	//var table string
 	if err != nil {
 		/*for res.Next() {
 			res.Scan(&table)
 			fmt.Println(table)
 		}*/
 		fmt.Printf("%v",res)
 	}else{
 		fmt.Println("ERRRRRR")
 		panic(err)
 	}

}

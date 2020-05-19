package analyze

import (
	"fmt"
  _"trace"
	_"util"

	_"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"


)

func TestDB(){
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
 	if err != nil {
 	fmt.Println(err)
 	}else{
 		fmt.Println("Connection Established")
 	}
 	defer db.Close()
}

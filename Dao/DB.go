package Dao

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql-master"
)

var SqlDB *sql.DB

func init() {
	var err error
	SqlDB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/skydrive?parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = SqlDB.Ping()

	if err != nil {
		log.Fatal(err.Error())
	}
}

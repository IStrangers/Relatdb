package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/dbname")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(`select * from USER`)
	if err != nil {
		return
	}
	println(rows)
}

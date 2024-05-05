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
	rows, err := db.Query(`
		create table myBase.User(
		    id INT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    		email VARCHAR(50) UNIQUE COMMENT '邮箱',
    		age INT UNSIGNED DEFAULT 1 COMMENT '年龄',
    		noId INT(6) ZEROFILL
		);
	`)
	if err != nil {
		return
	}
	println(rows)
}

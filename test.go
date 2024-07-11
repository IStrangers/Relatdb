package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/default")
	if err != nil {
		panic(err)
	}
	res, err := db.Exec(`
		create table if not exists User(
		    id INT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    		email VARCHAR(50) UNIQUE COMMENT '邮箱',
    		age INT UNSIGNED DEFAULT 1 COMMENT '年龄',
    		noId INT(6) ZEROFILL
		);
		insert into User VALUES (1,'1@qq.com',1,1);
		insert into User(id,noId,email,age) VALUES (1,1,'1@qq.com',1);
	`)
	if err != nil {
		panic(err)
		return
	}
	println(res)
}

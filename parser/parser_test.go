package parser

import "testing"

func TestParser(t *testing.T) {
	parser := CreateParser(0, `
		Create Database if not exists myBase;
		create table myBase.User(
		    
		);
		CREATE INDEX idx_name on myBase.User(name);
		CREATE UNIQUE INDEX idx_name on myBase.User(name);
		CREATE SPATIAL INDEX idx_name on myBase.User(name);
		CREATE FULLTEXT INDEX idx_name on myBase.User(name);
		DROP DATABASE myBase;
		DROP TABLE myBase.User;
		DROP INDEX index_name ON myBase.User;
		INSERT INTO myBase.User (name,age,address) VALUES 
		("名称1",1,'地址1'),
		("名称2",2,'地址2'),
		("名称3",3,'地址3');
		DELETE FROM myBase.User WHERE name = '名称' or age > 20 and addres != '地址' ORDER BY age DESC LIMIT 0,10;
		UPDATE myBase.User SET name = '更新名称',age = 1 WHERE name = '名称' ORDER BY age DESC LIMIT 0,10;
		SELECT CONNECTION_ID();
	`, true, true)
	statements := parser.Parse()
	println(statements)
}

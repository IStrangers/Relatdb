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
	`, true, true)
	statements := parser.Parse()
	println(statements)
}

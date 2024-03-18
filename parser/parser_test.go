package parser

import "testing"

func TestParser(t *testing.T) {
	parser := CreateParser(0, `
		Create Database if not exists myBase;
	`, true, true)
	node := parser.Parse()
	println(node)
}

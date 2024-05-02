package main

import (
	"Relatdb/server"
	"Relatdb/store"
	"Relatdb/utils"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	store := store.CreateStore(&store.Options{
		Path: utils.ConcatFilePaths(wd, "data"),
	})
	server := server.CreateServer(&server.Options{
		BindIp:   "localhost",
		BindPort: 3306,
	}, store)
	server.Start()
}

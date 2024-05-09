package main

import (
	"Relatdb/server"
	"Relatdb/store/icna"
	"Relatdb/utils"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	store := icna.NewIcnaStore(&icna.Options{
		Path: utils.ConcatFilePaths(wd, "data"),
	})
	server := server.NewServer(&server.Options{
		BindIp:   "localhost",
		BindPort: 3306,
	}, store)
	server.Start()
}

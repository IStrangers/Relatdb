package main

import "Relatdb/server"

func main() {
	relatdb := server.CreateRelatdb(&server.Options{
		BindIp:   "localhost",
		BindPort: 3306,
	})
	relatdb.Start()
}

package main

import "Relatdb/server"

func main() {
	server := server.CreateServer(&server.Options{
		BindIp:   "localhost",
		BindPort: 3306,
	})
	server.Start()
}

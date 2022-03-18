package main

import (
	"intra/api"
	"intra/cmd"
	"intra/db"
)

func main() {
	cmd.Init()
	db.Connect()
	defer db.Close()

	var server = api.CreateRouter()
	server.Run()

}

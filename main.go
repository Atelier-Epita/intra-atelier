package main

import (
	"intra/api"
	"intra/cmd"
	"intra/db"
)

// @title L'Atelier Intranet Backend
// @version 0.1

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
func main() {
	cmd.Init()
	db.Connect()
	defer db.Close()

	/*
		args := os.Args[1:]
		if len(args) > 0 && (args[0] == "-m" || args[0] == "--migrate") {
			db.DB.AutoMigrate(&models.File{})
			db.DB.AutoMigrate(&models.Participant{})
			db.DB.AutoMigrate(&models.Group{})
			db.DB.AutoMigrate(&models.Inventory{})
			db.DB.AutoMigrate(&models.User{})
			db.DB.AutoMigrate(&models.Event{})
			db.DB.AutoMigrate(&models.Equipment{})
		}
	*/

	var server = api.CreateRouter()
	server.Run()
}

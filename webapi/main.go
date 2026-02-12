package main

import (
	"log"
	"services/webapi/server"
)

func main() {
	app := server.NewApp()
	if err := app.Run("7000"); err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
}

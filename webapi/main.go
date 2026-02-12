package main

import (
	"log"
	_ "services/webapi/docs"
	"services/webapi/server"
)

// @title objects-service
// @version 1.0
// @description objects service: CRUD words and categories

//-- @host localhost:8080
// @BasePath /

//-- @securityDefinitions.apikey ApiKeyAuth
//-- @in header
//-- @name Authorization

func main() {
	app := server.NewApp()
	if err := app.Run("7000"); err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
}

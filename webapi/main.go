package main

import (
	"log"
	"os"
	_ "webapi/docs"
	"webapi/server"

	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := server.NewApp()
	if err := app.Run(os.Getenv("port")); err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
}

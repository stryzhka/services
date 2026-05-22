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

// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer {token}"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	app := server.NewApp()
	if err := app.Run(os.Getenv("port")); err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
	log.Println("Server started on port " + os.Getenv("port"))
}

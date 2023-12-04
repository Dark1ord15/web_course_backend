package main

import (
	app "Road_services/internal/api"
	"log"
)

// @title Платные дороги
// @version 1.0
// @description Road application
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log.Println("Application start!")

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	application.StartServer()

	log.Println("Application terminated!")
}

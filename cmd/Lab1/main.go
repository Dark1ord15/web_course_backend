package main

import (
	"Lab1/internal/api"
	"log"
)

func main() {
	log.Println("Applic start!")
	api.StartServer()
	log.Println("Application termited!")
}

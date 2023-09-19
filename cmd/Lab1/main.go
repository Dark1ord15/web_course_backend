package main

import (
	"Lab1/internal/api"
	"log"
)

func main() {
	log.Println("Applica start!")
	api.StartServer()
	log.Println("Application termited!")
}

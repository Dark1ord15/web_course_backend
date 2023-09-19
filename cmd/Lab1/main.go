package main

import (
	"Lab1/internal/api"
	"log"
)

func main() {
	log.Println("Appl start!")
	api.StartServer()
	log.Println("Application termited!")
}

package main

import (
	"log"

	"rga/backend/router"
)

func main() {
	log.Println("Starting server...")
	router.StartServer()
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func loadLocal() {
	err := godotenv.Load()

	if err == nil {
		log.Panic(err)
	}

	o
}

func main() {
	r := mux.NewRouter()
	log.Println("Server is listening on")
}

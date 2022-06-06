package router

import (
	"log"
	"net/http"
	"rga/backend/controllers"
	"time"

	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", controllers.Home).Methods("GET")
	r.HandleFunc("/admin/table", controllers.ResetUserTable).Methods("POST")
	r.HandleFunc("/users/{id}", controllers.DeleteAccount).Methods("DELETE")
	r.HandleFunc("/users/{id}", controllers.EditInfo).Methods("PATCH")
	r.HandleFunc("/auth/register", controllers.Register).Methods("POST")
	r.HandleFunc("/auth/login", controllers.Login).Methods("POST")

	// Start server
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:3001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

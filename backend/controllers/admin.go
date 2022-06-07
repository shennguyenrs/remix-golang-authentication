package controllers

import (
	"context"
	"log"
	"net/http"
	"rga/backend/config"
	"rga/backend/models"
)

// Create reset users table
func ResetUserTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	db := config.InitializeDB()

	if err := db.ResetModel(ctx, (*models.User)(nil)); err != nil {
		log.Panic(err)
	}

	w.WriteHeader(http.StatusCreated)
}

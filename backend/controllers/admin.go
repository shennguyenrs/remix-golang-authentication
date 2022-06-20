package controllers

import (
	"context"
	"net/http"
	"rga/backend/config"
	"rga/backend/models"
)

// Create reset users table
func ResetUserTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	db := config.InitializeDB()

	if err := db.ResetModel(ctx, (*models.User)(nil)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset user table"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Reset user table successfully"))
	return
}

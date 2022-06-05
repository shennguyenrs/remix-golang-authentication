package controllers

import (
	"context"
	"log"
	"net/http"
	"rga/backend/config"
	"rga/backend/models"

	"github.com/uptrace/bun/extra/bundebug"
)

// Create reset users table
func ResetUserTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	db := config.InitializeDB()

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := db.ResetModel(ctx, (*models.User)(nil)); err != nil {
		log.Panic(err)
	}

	w.WriteHeader(http.StatusCreated)
}

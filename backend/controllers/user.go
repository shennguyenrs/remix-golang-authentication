package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"rga/backend/config"
	"rga/backend/models"
	"rga/backend/utils"
)

func GetAccount(w http.ResponseWriter, r *http.Request) {
	userId, token := utils.ExtractIdnToken(w, r)

	if isValidToken := utils.ValidateToken(w, token, userId); isValidToken {
		// Return account using user id
		ctx := context.Background()
		db := config.InitializeDB()

		user := new(models.User)

		if err := db.NewSelect().Model(user).Where("id = ?", userId).Scan(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get user"))
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
		}

		return
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userId, token := utils.ExtractIdnToken(w, r)

	if isValidToken := utils.ValidateToken(w, token, userId); isValidToken {
		// Delete account using user id
		ctx := context.Background()
		db := config.InitializeDB()

		if _, err := db.NewDelete().Model((*models.User)(nil)).Where("id = ?", userId).Exec(ctx); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Failed to delete user record"))
		} else {
			w.WriteHeader(http.StatusOK)
			http.Redirect(w, r, "localhost:300", http.StatusOK)
		}

		return
	}
}

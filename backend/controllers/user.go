package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
		var user models.User

		if _, err := db.NewSelect().Where("id = ?", userId).Model(&user).Exec(ctx); err != nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get user"))
		}
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userId, token := utils.ExtractIdnToken(w, r)

	if isValidToken := utils.ValidateToken(w, token, userId); isValidToken {
		// Delete account using user id
		ctx := context.Background()
		db := config.InitializeDB()

		if _, err := db.NewDelete().Where("id = ?", userId).Exec(ctx); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Failed to delete user record"))
		} else {
			w.WriteHeader(http.StatusOK)
			http.Redirect(w, r, "localhost:300", http.StatusOK)
		}
	}
}

func EditInfo(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing new user information"))
	}

	userId, token := utils.ExtractIdnToken(w, r)

	if isValidToken := utils.ValidateToken(w, token, userId); isValidToken {
		// Update user information based on te request body
		ctx := context.Background()
		db := config.InitializeDB()

		var updateUser models.User

		json.Unmarshal(reqBody, &updateUser)

		if _, err := db.NewUpdate().Model(updateUser).WherePK().Exec(ctx); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Failed to update user record"))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

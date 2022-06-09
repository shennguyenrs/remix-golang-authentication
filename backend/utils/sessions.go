package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Generate user session
func GenerateSession(w http.ResponseWriter, r *http.Request, userEmail string, userId int) {
	tokenString, err := GenerateToken(userEmail)
	if err != nil {
		log.Panic("Failed to get generated token string")
	}

	// Write token to return header
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(`{"token": "Bearer ` + tokenString + `", "userId": "` + strconv.Itoa(userId) + `"}`)
	return
}

package utils

import (
	"log"
	"net/http"
)

// Generate user session
func GenerateSession(w http.ResponseWriter, r *http.Request, userEmail string) {
	tokenString, err := GenerateToken(userEmail)
	var headerToken string = "Bearer " + tokenString

	if err != nil {
		log.Panic("Failed to get generated token string")
	}

	// Write token to return header
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Authorization", headerToken)
	http.Redirect(w, r, "localhost:3000", http.StatusOK)
}

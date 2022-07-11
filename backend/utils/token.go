package utils

import (
	"context"
	"net/http"
	"rga/backend/config"
	"rga/backend/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type Claims struct {
	Email              string `json:"email" validate:"required,email"`
	jwt.StandardClaims `                    validate:"required"`
}

func GenerateToken(email string) (tokenString string, err error) {
	secretByte := config.GetJwtSecret()
	expirationTimes := time.Now().AddDate(0, 0, 7)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimes.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretByte)
	return
}

func ValidateToken(w http.ResponseWriter, tokenString string, userId int) bool {
	secretByte := config.GetJwtSecret()

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretByte, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parsing token error"))
		return false
	}

	if claims, ok := tkn.Claims.(*Claims); !ok || !tkn.Valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parsing claims error"))
		return false
	} else {

		// Find the user email in the database and compare with user id
		// If both are the same return true
		// Else return unauthorized request
		ctx := context.Background()
		db := config.InitializeDB()
		user := new(models.User)

		if err := db.NewSelect().Model(user).Where("email = ?", claims.Email).Scan(ctx); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}

		if user.ID != userId {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}

		return true
	}
}

func ExtractIdnToken(w http.ResponseWriter, r *http.Request) (userId int, token string) {
	uIdString := mux.Vars(r)["id"]
	tokenString := r.Header.Get("Authorization")

	userId, err := strconv.Atoi(uIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Failed to parse id"))
	}

	token = tokenString[7:]
	return
}

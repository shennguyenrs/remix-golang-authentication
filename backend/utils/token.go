package utils

import (
	"log"
	"rga/backend/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email              string `json:"email" validate:"required,email"`
	jwt.StandardClaims `                    validate:"required"`
}

func GenerateToken(email string) (tokenString string, err error) {
	localEnv, err := config.GetEnvMap()
	if err != nil {
		log.Panic("Failed to load .env file")
	}

	secretString := localEnv["JWT_SECRET"]
	secretByte := []byte(secretString)

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

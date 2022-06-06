package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Get local map secrets
func getEnvMap() (localEnv map[string]string) {
	localEnv, err := godotenv.Read()
	if err != nil {
		log.Panic("Failed to load .env file")
	}

	return
}

func GetJwtSecret() (secretByte []byte) {
	localEnv := getEnvMap()
	secretString := localEnv["JWT_SECRET"]
	secretByte = []byte(secretString)
	return
}

func GetDBString() (dsn string) {
	localEnv := getEnvMap()
	dsn = localEnv["DB_POSTGRES"]
	return
}

package config

import (
	"github.com/joho/godotenv"
)

// Get local map secrets
func GetEnvMap() (localEnv map[string]string, err error) {
	localEnv, err = godotenv.Read()
	return
}

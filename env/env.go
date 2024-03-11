package env

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error : %v", err.Error())
		log.Fatal("Error loading .env file")
	}
}

func GetAllowedOrigins() []string{
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	log.Printf("Allowed Origins: %v\n", allowedOrigins)
	if allowedOrigins == "" {
		return []string{"*"}
	}

	return strings.Split(allowedOrigins, ",")
}
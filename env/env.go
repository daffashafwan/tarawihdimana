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

	if allowedOrigins == "" {
		log.Print("using all origins")
		return []string{"*"}
	}

	return strings.Split(allowedOrigins, ",")
}
package env

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error : %v", err.Error())
		log.Fatal("Error loading .env file")
	}
}
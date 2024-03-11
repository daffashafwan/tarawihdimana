package main

import (
	"log"
	"net/http"
	"tarawihdimana/routes"

	appHandler "tarawihdimana/handlers"
	"tarawihdimana/env"
	"tarawihdimana/middleware"
)

func main() {

	env.LoadEnv()

	appHandler.InitAPIKEY()

	router := routes.NewRouter()

	rateLimitedRouter := middleware.RateLimitMiddleware(router)

	corsHandler := middleware.CorsMiddleware(rateLimitedRouter)

	// Start the server with the CORS and rate-limited handler
	log.Print("Server starting\n")
	err := http.ListenAndServe(":9999", corsHandler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
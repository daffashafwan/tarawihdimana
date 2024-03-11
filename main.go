package main

import (
	"log"
	"net/http"
	"tarawihdimana/routes"

	"github.com/gorilla/handlers"
	appHandler "tarawihdimana/handlers"
	"tarawihdimana/env"
)

func main() {

	env.LoadEnv()

	appHandler.InitAPIKEY()

	router := routes.NewRouter()

	// Define CORS options
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins(env.GetAllowedOrigins())
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Print("Server starting\n")
	// Wrap the router with the CORS middleware
	err := http.ListenAndServe(":9999", handlers.CORS(headersOk, originsOk, methodsOk)(router))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
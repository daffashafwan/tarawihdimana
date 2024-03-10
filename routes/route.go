package routes

import (
	"github.com/gorilla/mux"
	"tarawihdimana/handlers"
)

// NewRouter creates and returns a new router.
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// Handle the "/getRandomNearestMosque" route with the corresponding handler
	router.HandleFunc("/random-nearest-mosque", handlers.GetRandomNearestMosqueHandler).Methods("GET")

	return router
}
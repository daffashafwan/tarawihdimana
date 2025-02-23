package routes

import (
	"github.com/gorilla/mux"
	"tarawihdimana/handlers"
)

// NewRouter creates and returns a new router.
func NewRouter() *mux.Router {
    router := mux.NewRouter()

    // Add a prefix to the routes
    apiPrefix := "/tarawihdimana"
    apiRouter := router.PathPrefix(apiPrefix).Subrouter()

    // Handle the "/getRandomNearestMosque" route with the corresponding handler
    apiRouter.HandleFunc("/random-nearest-mosque", handlers.GetRandomNearestMosqueHandler).Methods("POST")
    apiRouter.HandleFunc("/prayer-times", handlers.GetPrayerTimesHandler).Methods("GET")

    return router
}
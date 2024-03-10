package tarawihdimana

import (
	"fmt"
	"net/http"
)

func main() {
	// Define the endpoint for finding the nearest place
	http.HandleFunc("/nearest-place", getNearestPlaceHandler)

	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":9999", nil)
}
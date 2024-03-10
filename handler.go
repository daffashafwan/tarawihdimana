package tarawihdimana

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Address struct {
	Label        string   `json:"label"`
	CountryCode  string   `json:"countryCode"`
	CountryName  string   `json:"countryName"`
	CountyCode   string   `json:"countyCode"`
	County       string   `json:"county"`
	City         string   `json:"city"`
	District     string   `json:"district"`
	Subdistrict  string   `json:"subdistrict"`
	Street       string   `json:"street"`
	PostalCode   string   `json:"postalCode"`
}

type Category struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
}

type Item struct {
	Title      string     `json:"title"`
	Address    Address    `json:"address"`
	Position   Location   `json:"position"`
	Access     []Location `json:"access"`
	Distance   int        `json:"distance"`
	Categories []Category `json:"categories"`
}

type Response struct {
	Items []Item `json:"items"`
}


func getNearestPlaceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse latitude, longitude, and radius from query parameters
	latitude, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	radius, err := strconv.Atoi(r.URL.Query().Get("radius"))
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("HERE_API_KEY")
	if apiKey == "" {
		http.Error(w, "HERE API key not provided", http.StatusInternalServerError)
		return
	}

	// Replace this URL with your actual API endpoint
	apiURL := fmt.Sprintf("https://discover.search.hereapi.com/v1/discover?q=masjid&in=circle:%f,%f;r=%d&limit=50&apikey=%s", latitude, longitude, radius, apiKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}


	// Respond with the nearest place information in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
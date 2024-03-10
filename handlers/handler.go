package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Address struct {
	Label       string `json:"label"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	CountyCode  string `json:"countyCode"`
	County      string `json:"county"`
	City        string `json:"city"`
	District    string `json:"district"`
	Subdistrict string `json:"subdistrict"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
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

type RequestData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    int     `json:"radius"`
	Limit     int     `json:"limit"`
	UseCache  bool    `json:"useCache"`
}

var cacheAPIResponse map[string][]Item
var ApiKey string
var cacheMutex sync.RWMutex

func InitAPIKEY() {
	apiKey := os.Getenv("HERE_API_KEY")
	if apiKey == "" {
		panic("HERE API key not provided")
	}
	ApiKey = apiKey
}

func GetRandomNearestMosqueHandler(w http.ResponseWriter, r *http.Request) {
	var request RequestData
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate parsed data
	if request.Latitude < -90 || request.Latitude > 90 {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	if request.Longitude < -180 || request.Longitude > 180 {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	if request.Radius <= 0 {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	if request.Limit <= 0 || request.Limit > 50 {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	arrayRes := make([]Item, 0)
	key := fmt.Sprintf("%.6f,%.6f,%d,%d", request.Latitude, request.Longitude, request.Radius, request.Limit)

	cacheMutex.RLock()
	res, exist := cacheAPIResponse[key]
	cacheMutex.RUnlock()

	if !request.UseCache || !exist {
		// Fetch and process the API response as usual
		response, err := callOutbound(request, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Ensure cacheAPIResponse is initialized
		cacheMutex.Lock()
		if cacheAPIResponse == nil {
			cacheAPIResponse = make(map[string][]Item)
		}

		// Update or add the response to the cache
		cacheAPIResponse[key] = response.Items
		cacheMutex.Unlock()

		// Use the response items
		arrayRes = response.Items
	}else {
		arrayRes = res
	}

	randomIndex := rand.Intn(len(arrayRes))
	randomMosque := arrayRes[randomIndex]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomMosque)
}

func callOutbound(request RequestData, w http.ResponseWriter) (Response, error) {
	apiURL := fmt.Sprintf("https://discover.search.hereapi.com/v1/discover?q=masjid&in=circle:%f,%f;r=%d&limit=%d&apikey=%s", request.Latitude, request.Longitude, request.Radius, request.Limit, ApiKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return Response{}, errors.Wrap(err, "Outbound Call")
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return Response{}, errors.Wrap(err, "Parse Resp")
	}
	return response, nil
}

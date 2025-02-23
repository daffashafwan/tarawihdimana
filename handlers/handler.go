package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

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

type CityData struct {
	ID     string `json:"id"`
	Lokasi string `json:"lokasi"`
}

type CityResponse struct {
	Status  bool    `json:"status"`
	Data    []CityData  `json:"data"`
}

type Jadwal struct {
	Tanggal string `json:"tanggal"`
	Imsak   string `json:"imsak"`
	Subuh   string `json:"subuh"`
	Terbit  string `json:"terbit"`
	Dhuha   string `json:"dhuha"`
	Dzuhur  string `json:"dzuhur"`
	Ashar   string `json:"ashar"`
	Maghrib string `json:"maghrib"`
	Isya    string `json:"isya"`
	Date    string `json:"date"`
}

type Data struct {
	ID     int    `json:"id"`
	Lokasi string `json:"lokasi"`
	Daerah string `json:"daerah"`
	Jadwal Jadwal `json:"jadwal"`
}

type PrayerTimeResponss struct {
	Data    Data    `json:"data"`
}

var cacheAPIResponse map[string][]Item
var cityDatas []CityData 
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
	cacheToggle, err := strconv.ParseBool(os.Getenv("SEARCH_RESPONSE_USE_CACHE"))
	if err != nil {
		cacheToggle = false
	}
	err = json.NewDecoder(r.Body).Decode(&request)
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


		if cacheToggle {
			// Ensure cacheAPIResponse is initialized
			cacheMutex.Lock()
			if cacheAPIResponse == nil {
				cacheAPIResponse = make(map[string][]Item)
			}

			// Update or add the response to the cache
			cacheAPIResponse[key] = response.Items
			cacheMutex.Unlock()
		}

		// Use the response items
		arrayRes = response.Items
	}else {
		arrayRes = res
	}

	if len(arrayRes) == 0 {
		err = errors.New("No mosque found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	randomIndex := rand.Intn(len(arrayRes))
	randomMosque := arrayRes[randomIndex]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomMosque)
}

func GetPrayerTimesHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
    if city == "" {
        http.Error(w, "Query parameter 'city' is required", http.StatusBadRequest)
        return
    }

	dataCity := cityDatas
	if len(dataCity) == 0 {
		cityResponse, err := callCityAPI(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dataCity = cityResponse.Data
	}

	cityID, found := containsCity(city, dataCity)
	if !found {
		http.Error(w, "City not found", http.StatusNotFound)
		return
	}

	prayerTimeResponse, err := callPrayerTimesAPI(w, cityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prayerTimeResponse)
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

func containsCity(requestCity string, cityList []CityData) (string, bool) {
	for _, city := range cityList {
		if isSameCity(requestCity, city.Lokasi) {
			return city.ID, true
		}

	}
	return "", false
}

func isSameCity(city1, city2 string) bool {
	return normalizeCityName(city1) == normalizeCityName(city2)
}

func normalizeCityName(name string) string {
	var ignoredPrefixes = []string{"KAB. ", "KOTA "}
	var specialCities = map[string]string {
		"JAKARTA BARAT": "Jakarta",
		"JAKARTA UTARA": "Jakarta",
		"JAKARTA SELATAN": "Jakarta",
		"JAKARTA TIMUR": "Jakarta",
		"JAKARTA PUSAT": "Jakarta",
	}
	if specialName, found := specialCities[strings.ToTitle(name)]; found {
		name = specialName
	}
	name = strings.ToUpper(name)
	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, prefix)
			break
		}
	}
	return name
}

func callPrayerTimesAPI(w http.ResponseWriter, cityID string) (PrayerTimeResponss, error) {
	apiURL := fmt.Sprintf("https://api.myquran.com/v2/sholat/jadwal/%s/%s", cityID, time.Now().Format("2006-01-02"))
	resp, err := http.Get(apiURL)
	if err != nil {
		return PrayerTimeResponss{}, errors.Wrap(err, "Outbound Call")
	}
	defer resp.Body.Close()

	var response PrayerTimeResponss
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return PrayerTimeResponss{}, errors.Wrap(err, "Parse Resp")
	}
	return response, nil
}

func callCityAPI(w http.ResponseWriter) (CityResponse, error) {
	apiURL := fmt.Sprintf("https://api.myquran.com/v2/sholat/kota/semua")
	resp, err := http.Get(apiURL)
	if err != nil {
		return CityResponse{}, errors.Wrap(err, "Outbound Call")
	}
	defer resp.Body.Close()

	var response CityResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return CityResponse{}, errors.Wrap(err, "Parse Resp")
	}
	
	cacheMutex.Lock()
	cityDatas = response.Data
	cacheMutex.Unlock()
	return response, nil
}

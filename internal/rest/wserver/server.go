package wserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/abkhan/simple-weather/domain"
	"github.com/abkhan/simple-weather/internal/repository/weatherapi"
	"github.com/gorilla/mux"
)

// Server represents a weather server interface
type Server interface {
	Start()
	WeatherHandler(w http.ResponseWriter, r *http.Request)
}

type server struct {
	weatherApi weatherapi.WeatherApi
}

var _ Server = &server{}

func New(w weatherapi.WeatherApi) Server {
	return &server{weatherApi: w}
}

// Start starts the weather server
func (s *server) Start() {
	router := mux.NewRouter()

	// some extras paths
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/about", AboutHandler).Methods("GET")

	router.HandleFunc("/weather", s.WeatherHandler).Methods("POST")

	// Start the HTTP server with the router on port 8080
	http.Handle("/", router)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/ is not implemented !")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "simple weather api")
}

func (s *server) WeatherHandler(w http.ResponseWriter, r *http.Request) {

	// Read the JSON data from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req := &domain.WeatherRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// now call the openweather api to get the weather
	resp, e := s.weatherApi.GetByLatLon(req.Lat, req.Lon)
	if e != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// reponse as JSON
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type header to indicate JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(js)
}

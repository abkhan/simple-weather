package domain

// WeatherRequest is the data object for SimpleWeather POST request
type WeatherRequest struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

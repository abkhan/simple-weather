package weatherapi

import (
	"fmt"
	"io"
	"net/http"
)

// WeatherApi is the interface for making HTTP requests
// to openweather api
type WeatherApi interface {
	GetByLatLon(lan, lon string) (*WeatherResp, error)
}

type weatherApi struct {
	url string
	api string
}

var _ WeatherApi = &weatherApi{}

func New(url, apiKey string) WeatherApi {
	return &weatherApi{url: url, api: apiKey}
}

// GetByLocation makes api calls to openweather api
func (s *weatherApi) GetByLatLon(lat, lon string) (*WeatherResp, error) {

	url := s.url + "?lat=" + lat + "&lon=" + lon + "&appid=" + s.api
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return condition(body)
}

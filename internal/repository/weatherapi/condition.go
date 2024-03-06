package weatherapi

import "encoding/json"

// WeatherResp struct is sent back to caller
type WeatherResp struct {
	Condition   string `json:"condition"`
	Temperature string `json:"temperature"`
}

// condition function takes the openweather api response and returns the
// data for the client of this api
func condition(body []byte) (*WeatherResp, error) {

	apiresp := &apiResponse{}
	resp := &WeatherResp{}

	err := json.Unmarshal(body, apiresp)
	if err != nil {
		return resp, err
	}

	resp.Condition = apiresp.Weather[0].Main
	resp.Temperature = "hot"
	if apiresp.Main.Temp < 300 {
		resp.Temperature = "moderate"
	}
	if apiresp.Main.Temp < 200 {
		resp.Temperature = "cold"
	}

	return resp, nil
}

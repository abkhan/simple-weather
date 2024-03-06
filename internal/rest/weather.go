package rest

import "github.com/labstack/echo/v4"

type WeatherApiService interface {
	Weather()
}

type WeatherApiHandler struct {
	Service WeatherApiService
}

// NewWeatherHandler will initialize the weather resources endpoint
func NewWeatherApiHandler(e *echo.Echo, svc WeatherApiService) {
	//handler := &ArticleHandler{
	//}
	//e.POST("/weather", handler.WeatherApiCall)
}

package clients

import (
	"fmt"
	"lab-goexpert-2/service-b/dtos"
)

func GetTemperatureFromWeatherApiResponse(weatherApiResponse dtos.WeatherApiResponse) (dtos.TemperatureResponse, error) {

	tempC := weatherApiResponse.Current.TempC
	tempF := (tempC * 1.8) + 32
	tempK := tempC + 273

	return dtos.TemperatureResponse{
		City:  weatherApiResponse.Location.Name,
		TempC: fmt.Sprintf("%.1f", tempC),
		TempF: fmt.Sprintf("%.1f", tempF),
		TempK: fmt.Sprintf("%.1f", tempK),
	}, nil
}

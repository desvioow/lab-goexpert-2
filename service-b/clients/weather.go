package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"lab-goexpert-2/service-b/dtos"
	"net/http"
)

func FetchWeather(cityName string) (dtos.WeatherApiResponse, error) {
	const weatherApiKey = "4d0046deef4e4342bd9192050251307"
	const baseUrl = "http://api.weatherapi.com/v1"

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		baseUrl,
		weatherApiKey,
		cityName,
	)

	resp, err := http.Get(url)
	if err != nil {
		return dtos.WeatherApiResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dtos.WeatherApiResponse{}, err
	}

	var weatherApiResponse dtos.WeatherApiResponse
	err = json.Unmarshal(body, &weatherApiResponse)
	if err != nil {
		return dtos.WeatherApiResponse{}, err
	}

	return weatherApiResponse, nil
}

func IsWeatherApiResponseEmpty(response dtos.WeatherApiResponse) bool {
	empty := dtos.WeatherApiResponse{}
	return response == empty
}

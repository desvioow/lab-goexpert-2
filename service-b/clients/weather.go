package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"lab-goexpert-2/service-b/models"
	"net/http"
)

func FetchWeather(cityName string) (models.WeatherApiResponse, error) {
	const weatherApiKey = "4d0046deef4e4342bd9192050251307"
	const baseUrl = "http://api.weatherapi.com/v1"

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		baseUrl,
		weatherApiKey,
		cityName,
	)

	resp, err := http.Get(url)
	if err != nil {
		return models.WeatherApiResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.WeatherApiResponse{}, err
	}

	var weatherApiResponse models.WeatherApiResponse
	err = json.Unmarshal(body, &weatherApiResponse)
	if err != nil {
		return models.WeatherApiResponse{}, err
	}

	return weatherApiResponse, nil
}

func IsWeatherApiResponseEmpty(response models.WeatherApiResponse) bool {
	empty := models.WeatherApiResponse{}
	return response == empty
}

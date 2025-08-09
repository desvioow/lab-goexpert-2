package main

import (
	"encoding/json"
	"fmt"
	"lab-goexpert-2/service-b/clients"
	"lab-goexpert-2/service-b/dtos"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Post("/", cepHandler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World B"))
	})

	http.ListenAndServe(":8081", r)
}

func cepHandler(w http.ResponseWriter, r *http.Request) {

	var dto dtos.CepDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
		return
	}

	regex := regexp.MustCompile(`^\d{8}$`)

	if regex.MatchString(dto.Cep) == false {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	viaCepResponse, err := clients.FetchCep(dto.Cep)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if clients.IsViaCepResponseEmpty(viaCepResponse) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}

	weatherApiResponse, err := clients.FetchWeather(viaCepResponse.Localidade)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find weather for the city"))
		return
	}
	if clients.IsWeatherApiResponseEmpty(weatherApiResponse) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find city"))
		return
	}

	temperature, err := clients.GetTemperatureFromWeatherApiResponse(weatherApiResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temperature)
	return

}

/*
Requisitos - Serviço B (responsável pela orquestração):

    O sistema deve receber um CEP válido de 8 digitos
    O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.
    O sistema deve responder adequadamente nos seguintes cenários:
        Em caso de sucesso:
            Código HTTP: 200
            Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
        Em caso de falha, caso o CEP não seja válido (com formato correto):
            Código HTTP: 422
            Mensagem: invalid zipcode
       Em caso de falha, caso o CEP não seja encontrado:
            Código HTTP: 404
            Mensagem: can not find zipcode

*/

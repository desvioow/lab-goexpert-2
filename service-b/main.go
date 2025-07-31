package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"lab-goexpert-2/service-b/clients"
	"net/http"
	"regexp"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{cep}", cepHandler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", r)
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	regex := regexp.MustCompile(`^\d{8}$`)

	if regex.MatchString(cep) == false {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	viaCepResponse, err := clients.FetchCep(cep)
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
	}
	json.NewEncoder(w).Encode(temperature)
	w.WriteHeader(http.StatusOK)
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

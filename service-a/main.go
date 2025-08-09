package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

type CepDTO struct {
	Cep string `json:"cep"`
}

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World A"))
	})
	r.Post("/", cepHandler)

	http.ListenAndServe(":8080", r)
}

func cepHandler(w http.ResponseWriter, r *http.Request) {

	var dto CepDTO
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

	jsonData, err := json.Marshal(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error serializing data"))
		return
	}

	serviceBRequestBody := bytes.NewBuffer(jsonData)

	// Send POST request
	serviceBResponse, err := http.Post("http://service-b:8081", "application/json", serviceBRequestBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error connecting to service B"))
		return
	}
	defer serviceBResponse.Body.Close()

	responseBody, err := io.ReadAll(serviceBResponse.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(serviceBResponse.StatusCode)
	w.Write(responseBody)
}

/*
Requisitos - Serviço A (responsável pelo input):

O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
Caso seja válido, será encaminhado para o Serviço B via HTTP
Caso não seja válido, deve retornar:
Código HTTP: 422
Mensagem: invalid zipcode
*/

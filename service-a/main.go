package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
)

type CepDTO struct {
	Cep string `json:"cep"`
}

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
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

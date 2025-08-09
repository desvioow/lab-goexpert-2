package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"lab-goexpert-2/service-b/dtos"
	"net/http"
)

func FetchCep(cep string) (dtos.ViaCepResponse, error) {

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		return dtos.ViaCepResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dtos.ViaCepResponse{}, err
	}
	var viaCepResponse dtos.ViaCepResponse
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		return dtos.ViaCepResponse{}, err
	}

	return viaCepResponse, nil
}

func IsViaCepResponseEmpty(response dtos.ViaCepResponse) bool {
	empty := dtos.ViaCepResponse{}
	return response == empty
}

package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"lab-goexpert-2/service-b/models"
	"net/http"
)

func FetchCep(cep string) (models.ViaCepResponse, error) {

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		return models.ViaCepResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ViaCepResponse{}, err
	}
	var viaCepResponse models.ViaCepResponse
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		return models.ViaCepResponse{}, err
	}

	return viaCepResponse, nil
}

func IsViaCepResponseEmpty(response models.ViaCepResponse) bool {
	empty := models.ViaCepResponse{}
	return response == empty
}

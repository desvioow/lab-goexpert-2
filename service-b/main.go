package main

import (
	"encoding/json"
	"fmt"
	"lab-goexpert-2/service-b/clients"
	"lab-goexpert-2/service-b/dtos"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/openzipkin/zipkin-go"
	zipkinMiddleware "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinReporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var tracer *zipkin.Tracer

func main() {
	reporter := zipkinReporter.NewReporter("http://zipkin:9411/api/v2/spans")
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint("service-b", "service-b:8081")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	tracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	r := chi.NewRouter()

	r.Use(zipkinMiddleware.NewServerMiddleware(tracer, zipkinMiddleware.TagResponseSize(true)))

	r.Post("/", cepHandler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am service B"))
	})

	log.Println("Service B starting on :8081")
	http.ListenAndServe(":8081", r)
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	span := zipkin.SpanFromContext(r.Context())

	var dto dtos.CepDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		span.Tag("error", "invalid json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
		return
	}

	span.Tag("cep", dto.Cep)

	regex := regexp.MustCompile(`^\d{8}$`)

	if regex.MatchString(dto.Cep) == false {
		span.Tag("error", "invalid zipcode")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	viaCepSpan := tracer.StartSpan("fetch-cep on service B", zipkin.Parent(span.Context()))
	viaCepSpan.Tag("api.name", "viacep")
	viaCepSpan.Tag("cep", dto.Cep)

	viaCepResponse, err := clients.FetchCep(dto.Cep)
	if err != nil {
		viaCepSpan.Tag("error", fmt.Sprintf("Error: %v", err))
		viaCepSpan.Finish()
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	viaCepSpan.Finish()

	if clients.IsViaCepResponseEmpty(viaCepResponse) {
		span.Tag("error", "can not find zipcode")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}

	weatherSpan := tracer.StartSpan("fetch-weather", zipkin.Parent(span.Context()))
	weatherSpan.Tag("api.name", "weather")
	weatherSpan.Tag("city", viaCepResponse.Localidade)

	weatherApiResponse, err := clients.FetchWeather(viaCepResponse.Localidade)
	if err != nil {
		weatherSpan.Tag("error", "can not find weather for the city")
		weatherSpan.Finish()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find weather for the city"))
		return
	}
	weatherSpan.Finish()

	if clients.IsWeatherApiResponseEmpty(weatherApiResponse) {
		span.Tag("error", "can not find city")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find city"))
		return
	}

	temperature, err := clients.GetTemperatureFromWeatherApiResponse(weatherApiResponse)
	if err != nil {
		span.Tag("error", "temperature conversion error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	span.Tag("http.status_code", strconv.Itoa(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temperature)
	return
}

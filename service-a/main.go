package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/openzipkin/zipkin-go"
	zipkinMiddleware "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinReporter "github.com/openzipkin/zipkin-go/reporter/http"
)

type CepDTO struct {
	Cep string `json:"cep"`
}

var tracer *zipkin.Tracer

func main() {
	reporter := zipkinReporter.NewReporter("http://zipkin:9411/api/v2/spans")
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint("service-a", "service-a:8080")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	tracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	r := chi.NewRouter()

	r.Use(zipkinMiddleware.NewServerMiddleware(tracer, zipkinMiddleware.TagResponseSize(true)))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am service A"))
	})
	r.Post("/", cepHandler)

	log.Println("Service A starting on :8080")
	http.ListenAndServe(":8080", r)
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpan("cep-processing on service A")
	defer span.Finish()

	ctx := zipkin.NewContext(r.Context(), span)

	var dto CepDTO
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

	jsonData, err := json.Marshal(dto)
	if err != nil {
		span.Tag("error", "error serializing data")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error serializing data"))
		return
	}

	serviceBRequestBody := bytes.NewBuffer(jsonData)

	serviceBSpan := tracer.StartSpan("call-service-b", zipkin.Parent(span.Context()))
	serviceBSpan.Tag("service.name", "service-b")
	serviceBSpan.Tag("http.method", "POST")
	serviceBSpan.Tag("http.url", "http://service-b:8081")

	client, err := zipkinMiddleware.NewClient(tracer, zipkinMiddleware.ClientTrace(true))
	if err != nil {
		serviceBSpan.Tag("error", "failed to create traced client")
		serviceBSpan.Finish()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error creating traced client"))
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://service-b:8081", serviceBRequestBody)
	if err != nil {
		serviceBSpan.Tag("error", "failed to create request")
		serviceBSpan.Finish()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error creating request"))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	serviceBResponse, err := client.Do(req)
	if err != nil {
		serviceBSpan.Tag("error", "error connecting to service B")
		serviceBSpan.Finish()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error connecting to service B"))
		return
	}
	defer serviceBResponse.Body.Close()

	serviceBSpan.Tag("http.status_code", strconv.Itoa(serviceBResponse.StatusCode))
	serviceBSpan.Finish()

	responseBody, err := io.ReadAll(serviceBResponse.Body)
	if err != nil {
		span.Tag("error", "error reading response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	span.Tag("http.status_code", strconv.Itoa(serviceBResponse.StatusCode))
	w.WriteHeader(serviceBResponse.StatusCode)
	w.Write(responseBody)
}

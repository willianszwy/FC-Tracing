package handlers

import (
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"regexp"
	"willianszwy/FC-Cloud-Run/internal/temperature"
	"willianszwy/FC-Cloud-Run/internal/viacep"
	"willianszwy/FC-Cloud-Run/internal/weather"
)

type TemperatureHandler struct {
	viaCepClient  *viacep.ViaCep
	weatherClient *weather.Weather
	tr            trace.Tracer
}

type RequestBody struct {
	Zipcode string `json:"zipcode"`
}

func New(viacepClient *viacep.ViaCep, weatherClient *weather.Weather, tr trace.Tracer) *TemperatureHandler {
	return &TemperatureHandler{
		viaCepClient:  viacepClient,
		weatherClient: weatherClient,
		tr:            tr,
	}
}

func (t *TemperatureHandler) Handler(writer http.ResponseWriter, request *http.Request) {
	log.Println(request.Header)
	carrier := propagation.HeaderCarrier(request.Header)
	ctx := request.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := t.tr.Start(ctx, "Service B")
	defer span.End()
	log.Println("starting request")

	var req RequestBody
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(fmt.Sprintf("[zipcode:%s]", req.Zipcode))

	regex := regexp.MustCompile("^[0-9]{8}$")
	if !regex.MatchString(req.Zipcode) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		http.Error(writer, "invalid zipCode", http.StatusUnprocessableEntity)
		return
	}

	city, err := t.viaCepClient.FindByZipCode(ctx, req.Zipcode)
	if err != nil {
		log.Println("error", err.Error())
		writer.WriteHeader(http.StatusNotFound)
		http.Error(writer, "can not find zipcode", http.StatusNotFound)
		return
	}

	tempByCity, err := t.weatherClient.FindTempByCity(ctx, city.Name)
	if err != nil {
		log.Println("error", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := temperature.New(city.Name, tempByCity.Current.TempC, tempByCity.Current.TempF)
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

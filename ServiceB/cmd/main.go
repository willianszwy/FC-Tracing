package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"willianszwy/FC-Cloud-Run/configs"
	"willianszwy/FC-Cloud-Run/internal/handlers"
	"willianszwy/FC-Cloud-Run/internal/viacep"
	"willianszwy/FC-Cloud-Run/internal/weather"
)

func main() {

	config, err := configs.LoadConfig("")
	if err != nil {
		panic(err)
	}
	log.Println("Load config...")
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	viaCepClient := viacep.New(http.DefaultClient)
	weatherClient := weather.New(http.DefaultClient, config.WeatherAPIKey)
	temperatureHandler := handlers.New(viaCepClient, weatherClient)

	r.Post("/temperature", temperatureHandler.Handler)

	http.ListenAndServe(":8080", r)
}

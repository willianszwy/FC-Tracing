package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"regexp"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/temperature", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("starting request")

		zipcode := request.URL.Query().Get("zipcode")
		log.Println(fmt.Sprintf("[zipcode:%s]", zipcode))

		regex := regexp.MustCompile("^[0-9]{8}$")
		if !regex.MatchString(zipcode) {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			http.Error(writer, "invalid zipCode", http.StatusUnprocessableEntity)
			return
		}

		http.Post("localhost:8080/", "application/json", nil)
	})

	http.ListenAndServe(":8080", r)
}

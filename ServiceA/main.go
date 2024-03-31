package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
	"regexp"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("starting request service A")

		zipcode := request.URL.Query().Get("zipcode")
		log.Println(fmt.Sprintf("[zipcode:%s]", zipcode))

		regex := regexp.MustCompile("^[0-9]{8}$")
		if !regex.MatchString(zipcode) {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			http.Error(writer, "invalid zipCode", http.StatusUnprocessableEntity)
			return
		}

		endpoint := "http://service-b:8080/temperature"
		body, _ := json.Marshal(map[string]string{
			"zipcode": zipcode,
		})

		response, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Println("error: ", err)
			http.Error(writer, "error calling service B", http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		defer response.Body.Close()
		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("impossible to read all body of response: %s", err)
		}
		log.Printf("body: %s", string(resBody))
		writer.Write(resBody)

	})

	http.ListenAndServe(":8081", r)
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
)

var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)

func initTracer(url string) (func(context.Context) error, error) {
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("service-a"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown, nil
}

func main() {
	log.Println("Start service A...")
	url := flag.String("zipkin", "http://zipkin:9411/api/v2/spans", "zipkin url")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initTracer(*url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tr := otel.GetTracerProvider().Tracer("component-main")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		carrier := propagation.HeaderCarrier(request.Header)
		ctx := request.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
		ctx, span := tr.Start(ctx, "zipcode service")
		defer span.End()

		log.Println("starting request service A")

		var reqBody RequestBody
		err := json.NewDecoder(request.Body).Decode(&reqBody)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(fmt.Sprintf("[zipcode:%s]", reqBody.Zipcode))

		regex := regexp.MustCompile("^[0-9]{8}$")
		if !regex.MatchString(reqBody.Zipcode) {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			http.Error(writer, "invalid zipCode", http.StatusUnprocessableEntity)
			return
		}

		endpoint := "http://service-b:8080/temperature"
		body, _ := json.Marshal(map[string]string{
			"zipcode": reqBody.Zipcode,
		})

		req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(body))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
		response, err := http.DefaultClient.Do(req)
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

type RequestBody struct {
	Zipcode string `json:"zipcode"`
}

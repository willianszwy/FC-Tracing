package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"willianszwy/FC-Cloud-Run/internal/interfaces"
)

type City struct {
	Name string `json:"localidade"`
}

type ViaCep struct {
	client interfaces.HTTPClient
	tr     trace.Tracer
}

func New(client interfaces.HTTPClient, tr trace.Tracer) *ViaCep {
	return &ViaCep{
		client: client,
		tr:     tr,
	}
}

func (vc *ViaCep) FindByZipCode(ctx context.Context, zipCode string) (City, error) {
	ctx, span := vc.tr.Start(ctx, "Viacep")
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://viacep.com.br/ws/"+zipCode+"/json", nil)
	if err != nil {
		return City{}, fmt.Errorf("error creating request %w", err)
	}
	resp, err := vc.client.Do(req)
	if err != nil {
		return City{}, fmt.Errorf("error doing request %w", err)
	}
	defer resp.Body.Close()
	var city City
	err = json.NewDecoder(resp.Body).Decode(&city)
	if err != nil {
		return City{}, fmt.Errorf("error deconding request %w", err)
	}
	if city.Name == "" {
		return City{}, fmt.Errorf("error city notfound")
	}
	return city, nil
}

package handlers

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"willianszwy/FC-Cloud-Run/internal/viacep"
	"willianszwy/FC-Cloud-Run/internal/weather"
)

type ClientMock struct {
	Res *http.Response
	Err error
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return c.Res, c.Err
}

func TestTemperatureHandler_Handler(t *testing.T) {

	json := ` {
      "localidade": "São Paulo"
    }`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}

	json = ` {
    "current": {
        "temp_c": 18.0,
        "temp_f": 64.4
      }
    }`
	body2 := io.NopCloser(bytes.NewReader([]byte(json)))
	client2 := ClientMock{
		Res: &http.Response{Body: body2, StatusCode: 200},
	}

	viaCepClient := viacep.New(&client)
	weatherClient := weather.New(&client2, "")
	temperatureHandler := New(viaCepClient, weatherClient)

	req := httptest.NewRequest("GET", "http://example.com/foo?zipcode=00000000", nil)
	w := httptest.NewRecorder()
	temperatureHandler.Handler(w, req)

	resp := w.Result()

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestTemperatureHandler_Handler_InvalidZipcode(t *testing.T) {

	json := ` {
      "localidade": "São Paulo"
    }`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}

	json = ` {
    "current": {
        "temp_c": 18.0,
        "temp_f": 64.4
      }
    }`
	body2 := io.NopCloser(bytes.NewReader([]byte(json)))
	client2 := ClientMock{
		Res: &http.Response{Body: body2, StatusCode: 200},
	}

	viaCepClient := viacep.New(&client)
	weatherClient := weather.New(&client2, "")
	temperatureHandler := New(viaCepClient, weatherClient)

	req := httptest.NewRequest("GET", "http://example.com/foo?zipcode=invalidcep", nil)
	w := httptest.NewRecorder()
	temperatureHandler.Handler(w, req)

	resp := w.Result()

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

}

func TestTemperatureHandler_Handler_ZipcodeNotFound(t *testing.T) {

	client := ClientMock{
		Res: nil,
		Err: errors.New("zipcode notfound"),
	}

	json := ` {
    "current": {
        "temp_c": 18.0,
        "temp_f": 64.4
      }
    }`
	body2 := io.NopCloser(bytes.NewReader([]byte(json)))
	client2 := ClientMock{
		Res: &http.Response{Body: body2, StatusCode: 200},
	}

	viaCepClient := viacep.New(&client)
	weatherClient := weather.New(&client2, "")
	temperatureHandler := New(viaCepClient, weatherClient)

	req := httptest.NewRequest("GET", "http://example.com/foo?zipcode=00000000", nil)
	w := httptest.NewRecorder()
	temperatureHandler.Handler(w, req)

	resp := w.Result()

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

}

func TestTemperatureHandler_Handler_WeatherError(t *testing.T) {

	json := ` {
      "localidade": "São Paulo"
    }`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}

	client2 := ClientMock{
		Res: nil,
		Err: errors.New("error weather"),
	}

	viaCepClient := viacep.New(&client)
	weatherClient := weather.New(&client2, "")
	temperatureHandler := New(viaCepClient, weatherClient)

	req := httptest.NewRequest("GET", "http://example.com/foo?zipcode=00000000", nil)
	w := httptest.NewRecorder()
	temperatureHandler.Handler(w, req)

	resp := w.Result()

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

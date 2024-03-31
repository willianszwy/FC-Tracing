package weather

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type ClientMock struct {
	Res *http.Response
	Err error
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return c.Res, c.Err
}

func TestFindTempByCity(t *testing.T) {
	json := ` {
    "location": {
        "name": "Embu",
        "region": "Eastern",
        "country": "Kenya",
        "lat": -0.53,
        "lon": 37.45,
        "tz_id": "Africa/Nairobi",
        "localtime_epoch": 1711228128,
        "localtime": "2024-03-24 0:08"
    },
    "current": {
        "last_updated_epoch": 1711227600,
        "last_updated": "2024-03-24 00:00",
        "temp_c": 18.0,
        "temp_f": 64.4,
        "is_day": 0,
        "condition": {
            "text": "Partly cloudy",
            "icon": "//cdn.weatherapi.com/weather/64x64/night/116.png",
            "code": 1003
        },
        "wind_mph": 4.3,
        "wind_kph": 6.8,
        "wind_degree": 60,
        "wind_dir": "ENE",
        "pressure_mb": 1011.0,
        "pressure_in": 29.86,
        "precip_mm": 0.28,
        "precip_in": 0.01,
        "humidity": 88,
        "cloud": 75,
        "feelslike_c": 18.0,
        "feelslike_f": 64.4,
        "vis_km": 10.0,
        "vis_miles": 6.0,
        "uv": 1.0,
        "gust_mph": 4.8,
        "gust_kph": 7.7
    }
}`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}
	weatherApi := New(&client, "asdfasfasf")
	assert.NotNil(t, weatherApi)

	temp, err := weatherApi.FindTempByCity(context.TODO(), "Cidade")

	assert.NotNil(t, temp)
	assert.Nil(t, err)
	assert.Equal(t, 18.0, temp.Current.TempC)
}

func TestFindTempByCity_NewRequestError(t *testing.T) {
	const expectedError = "FindTempByCity : error creating request net/http: nil Context"
	json := ` {
    "current": {
        "temp_c": 18.0,
        "temp_f": 64.4        
    }
}`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}
	weatherApi := New(&client, "asdfasdfasd")
	assert.NotNil(t, weatherApi)

	temp, err := weatherApi.FindTempByCity(nil, "")

	assert.Equal(t, Response{}, temp)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func TestFindTempByCity_DoError(t *testing.T) {
	client := ClientMock{
		Res: nil,
		Err: errors.New("error"),
	}
	weatherApi := New(&client, "asdfasdfasd")
	assert.NotNil(t, weatherApi)

	temp, err := weatherApi.FindTempByCity(context.TODO(), "")

	assert.Equal(t, Response{}, temp)
	assert.NotNil(t, err)
	assert.Equal(t, "FindTempByCity: error doing request error", err.Error())
}

func TestFindTempByCity_UnMarshallError(t *testing.T) {
	const expectedError = "FindTempByCity: error deconding request json: cannot unmarshal string into Go struct field .current.temp_c of type float64"
	json := ` {
    "current": {
        "temp_c": "teste",
        "temp_f": "teste"        
    }
}`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}
	weatherApi := New(&client, "asdfasdfasd")
	assert.NotNil(t, weatherApi)

	temp, err := weatherApi.FindTempByCity(context.TODO(), "")

	assert.Equal(t, Response{}, temp)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

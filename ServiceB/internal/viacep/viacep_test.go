package viacep

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

func TestFindByZipCode(t *testing.T) {
	json := ` {
      "cep": "01001-000",
      "logradouro": "Praça da Sé",
      "complemento": "lado ímpar",
      "bairro": "Sé",
      "localidade": "São Paulo",
      "uf": "SP",
      "ibge": "3550308",
      "gia": "1004",
      "ddd": "11",
      "siafi": "7107"
    }`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}
	viaCep := New(&client)
	assert.NotNil(t, viaCep)

	city, err := viaCep.FindByZipCode(context.TODO(), "00000-000")

	assert.NotNil(t, city)
	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", city.Name)
}

func TestFindByZipCode_NewRequestError(t *testing.T) {
	const expectedError = "error creating request parse \"https://viacep.com.br/ws/$%ˆ&$%/json\": invalid URL escape \"%ˆ\""

	json := ` {
      "cep": "01001-000",
      "logradouro": "Praça da Sé",
      "complemento": "lado ímpar",
      "bairro": "Sé",
      "localidade": "São Paulo",
      "uf": "SP",
      "ibge": "3550308",
      "gia": "1004",
      "ddd": "11",
      "siafi": "7107"
    }`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}
	viaCep := New(&client)
	assert.NotNil(t, viaCep)

	city, err := viaCep.FindByZipCode(context.TODO(), "$%ˆ&$%")

	assert.Equal(t, City{}, city)
	assert.Equal(t, "", city.Name)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func TestFindByZipCode_DoError(t *testing.T) {
	client := ClientMock{
		Res: nil,
		Err: errors.New("error"),
	}
	viaCep := New(&client)
	assert.NotNil(t, viaCep)

	city, err := viaCep.FindByZipCode(context.TODO(), "00000-000")

	assert.Equal(t, City{}, city)
	assert.Equal(t, "", city.Name)
	assert.NotNil(t, err)
	assert.Equal(t, "error doing request error", err.Error())
}

func TestFindByZipCode_UnMarshallError(t *testing.T) {
	const expectedError = "error deconding request json: cannot unmarshal number into Go struct field City.localidade of type string"
	json := ` {
      "localidade": 1231312
    }`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}
	viaCep := New(&client)
	assert.NotNil(t, viaCep)

	city, err := viaCep.FindByZipCode(context.TODO(), "00000-000")

	assert.Equal(t, City{}, city)
	assert.Equal(t, "", city.Name)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func TestFindByZipCode_EmptyCityError(t *testing.T) {
	const expectedError = "error city notfound"
	json := ` {
      "localidade": ""
    }`
	body := io.NopCloser(bytes.NewReader([]byte(json)))
	client := ClientMock{
		Res: &http.Response{Body: body, StatusCode: 200},
	}
	viaCep := New(&client)
	assert.NotNil(t, viaCep)

	city, err := viaCep.FindByZipCode(context.TODO(), "00000-000")

	assert.Equal(t, City{}, city)
	assert.Equal(t, "", city.Name)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

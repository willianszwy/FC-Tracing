package temperature

type Temperature struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
	Kelvin     float64 `json:"kelvin"`
}

func New(city string, celsius, fahrenheit float64) *Temperature {
	return &Temperature{
		City:       city,
		Celsius:    celsius,
		Fahrenheit: fahrenheit,
		Kelvin:     celsius + 273,
	}
}

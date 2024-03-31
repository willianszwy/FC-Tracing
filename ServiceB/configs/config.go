package configs

import "github.com/spf13/viper"

var cfg *Config

type Config struct {
	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	viper.SetConfigName("configs")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

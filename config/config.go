package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port        int    `env:"PORT"`
	PostgresURL string `env:"POSTGRES_URL"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

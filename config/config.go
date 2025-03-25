package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBURL: os.Getenv("DB_URL"),
	}, nil
}

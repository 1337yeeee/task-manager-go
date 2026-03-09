package config

import (
	"os"
)

type Config struct {
	APIPort string

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	SSLMode    string

	JWTSecret string
}

func Load() (Config, error) {
	return Config{
		APIPort: os.Getenv("API_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		SSLMode:    os.Getenv("SSL_MODE"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}, nil
}

func (c Config) Get(key string) string {
	return os.Getenv(key)
}

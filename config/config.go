package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config struct hold all the configuration of the server.
type Config struct {
	ServerAddr string
}

// NewConfig function will load environment variables and return them as [Config] struct.
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
	}, nil
}

// getEnv will return environment variable with a key.
// If the variable is not found it will return the fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

// getEnvInt will return environment variable parsed as int with a key.
// If the variable is not found or not valid int it will return the fallback.
func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return valueInt
	}

	return fallback
}

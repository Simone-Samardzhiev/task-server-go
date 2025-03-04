package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

// DatabaseConfig struct holds database configuration.
type DatabaseConfig struct {
	// Url is the url of the database.
	Url string
	// MaxOpenConnection specify how many connections can be open.
	MaxOpenConnections int
	// MaxIdleConnections specify how many idle connections can be open.
	MaxIdleConnections int
}

// Config struct holds all the configuration of the server.
type Config struct {
	// ServerAddr is the port of the server.
	ServerAddr string
	// DatabaseConfig is the database configuration.
	DatabaseConfig DatabaseConfig
	AuthConfig     AuthConfig
}

// AuthConfig struct holds authentication configuration.
type AuthConfig struct {
	// JwtSecret used to sign tokens
	JwtSecret []byte
	// JwtIssuer used to set the issuer of the tokens.
	JwtIssuer string
}

// NewConfig function will load environment variables and return them as [Config] struct.
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		DatabaseConfig: DatabaseConfig{
			Url:                getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
			MaxOpenConnections: getEnvInt("MAX_OPEN_CONNECTIONS", 10),
			MaxIdleConnections: getEnvInt("MAX_IDLE_CONNECTIONS", 10),
		},
		AuthConfig: AuthConfig{
			JwtSecret: []byte(getEnv("JWT_SECRET", "secret_for_jwt")),
			JwtIssuer: getEnv("JWT_ISSUER", "com.localhost"),
		},
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

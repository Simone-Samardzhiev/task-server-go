package config

import "os"

// Config will hold all environment variables.
type Config struct {
	// The port of the server.
	Port string
	// The database user.
	DBUser string
	// The database password.
	DBPass string
	// The database address.
	DBHost string
	// The database name.
	DBName string
	// Secret used to sign tokens.
	JWTSecret string
}

// Envs is a variable holding the loaded configuration.
var Envs = initConfig()

// initConfig will load the configuration.
func initConfig() Config {
	err := LoadEnv()
	if err != nil {
		panic(err)
	}

	return Config{
		Port:      getEnv("PORT", "8080"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "postgres"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBName:    getEnv("DB_NAME", "postgres"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
	}
}

// getEnv will try to get environment variable by key or return the fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

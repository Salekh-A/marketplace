package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseDSN   string
}

func Load() Config {
	serverAddress := flag.String(
		"a",
		getEnv("SERVER_ADDRESS", ":8082"),
		"server address",
	)

	databaseDSN := flag.String(
		"d",
		getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5434/products?sslmode=disable"),
		"database dsn",
	)

	flag.Parse()

	return Config{
		ServerAddress: *serverAddress,
		DatabaseDSN:   *databaseDSN,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}

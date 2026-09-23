package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseDSN   string
	RedisAddress  string
}

func Load() Config {
	serverAddress := flag.String(
		"a",
		getEnv("SERVER_ADDRESS", ":8083"),
		"server address",
	)

	databaseDSN := flag.String(
		"d",
		getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5435/analytics?sslmode=disable"),
		"database dsn",
	)

	redisAddress := flag.String(
		"r",
		getEnv("REDIS_ADDRESS", "localhost:6379"),
		"redis address",
	)

	flag.Parse()

	return Config{
		ServerAddress: *serverAddress,
		DatabaseDSN:   *databaseDSN,
		RedisAddress:  *redisAddress,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return defaultValue
}

package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress    string
	DatabaseDSN      string
	AnalyticsAddress string
	ProductAddress   string
}

func Load() Config {
	serverAddress := flag.String(
		"a",
		getEnv("SERVER_ADDRESS", ":8081"),
		"server address",
	)

	databaseDSN := flag.String(
		"d",
		getEnv(
			"DATABASE_DSN",
			"postgres://postgres:postgres@localhost:5433/orders?sslmode=disable",
		),
		"database dsn",
	)

	analyticsAddress := flag.String(
		"analytics",
		getEnv("ANALYTICS_ADDRESS", "localhost:9093"),
		"analytics grpc address",
	)

	productAddress := flag.String(
		"product",
		getEnv("PRODUCT_ADDRESS", "localhost:9092"),
		"product grpc address",
	)

	flag.Parse()

	return Config{
		ServerAddress:    *serverAddress,
		DatabaseDSN:      *databaseDSN,
		AnalyticsAddress: *analyticsAddress,
		ProductAddress:   *productAddress,
	}

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}

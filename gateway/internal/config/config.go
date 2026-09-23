package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress    string
	OrderAddress     string
	ProductAddress   string
	AnalyticsAddress string
}

func Load() Config {
	serverAddress := flag.String(
		"a",
		getEnv("SERVER_ADDRESS", ":8080"),
		"server address",
	)

	orderAddress := flag.String(
		"o",
		getEnv("ORDER_ADDRESS", "localhost:8081"),
		"order service address",
	)

	productAddress := flag.String(
		"p",
		getEnv("PRODUCT_ADDRESS", "localhost:8082"),
		"product service address",
	)

	analyticsAddress := flag.String(
		"n",
		getEnv("ANALYTICS_ADDRESS", "localhost:8083"),
		"analytics service address",
	)

	flag.Parse()

	return Config{
		ServerAddress:    *serverAddress,
		OrderAddress:     *orderAddress,
		ProductAddress:   *productAddress,
		AnalyticsAddress: *analyticsAddress,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return defaultValue
}

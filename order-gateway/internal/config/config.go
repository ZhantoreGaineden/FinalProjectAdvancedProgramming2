package config

import "os"

type Config struct {
	HTTPPort         string
	OrderServiceAddr string
}

func Load() Config {
	return Config{
		HTTPPort:         getEnv("ORDER_GATEWAY_HTTP_PORT", "8083"),
		OrderServiceAddr: getEnv("ORDER_SERVICE_ADDR", "localhost:50053"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

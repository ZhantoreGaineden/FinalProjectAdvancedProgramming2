package config

import "os"

type Config struct {
	GRPCPort    string
	DatabaseURL string
	NATSURL     string
}

func Load() Config {
	return Config{
		GRPCPort:    getEnv("ORDER_SERVICE_GRPC_PORT", "50053"),
		DatabaseURL: getEnv("ORDER_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/order_db?sslmode=disable"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

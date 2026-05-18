package config

import "os"

type Config struct {
	GRPCPort    string
	DatabaseURL string
	NATSURL     string
}

func Load() Config {
	return Config{
		GRPCPort:    getEnv("USER_SERVICE_GRPC_PORT", "50052"),
		DatabaseURL: getEnv("USER_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/user_db?sslmode=disable"),
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

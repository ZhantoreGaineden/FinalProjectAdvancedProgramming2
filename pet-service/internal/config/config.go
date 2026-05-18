package config

import "os"

type Config struct {
	GRPCPort    string
	DatabaseURL string
	RedisAddr   string
}

func Load() Config {
	return Config{
		GRPCPort:    getEnv("PET_SERVICE_GRPC_PORT", "50051"),
		DatabaseURL: getEnv("PET_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/pet_db?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

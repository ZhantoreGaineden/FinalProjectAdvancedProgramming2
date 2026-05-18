package config

import "os"

type Config struct {
	HTTPPort       string
	PetServiceAddr string
}

func Load() Config {
	return Config{
		HTTPPort:       getEnv("PET_GATEWAY_HTTP_PORT", "8081"),
		PetServiceAddr: getEnv("PET_SERVICE_ADDR", "localhost:50051"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

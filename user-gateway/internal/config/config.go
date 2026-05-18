package config

import "os"

type Config struct {
	HTTPPort        string
	UserServiceAddr string
}

func Load() Config {
	return Config{
		HTTPPort:        getEnv("USER_GATEWAY_HTTP_PORT", "8082"),
		UserServiceAddr: getEnv("USER_SERVICE_ADDR", "localhost:50052"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

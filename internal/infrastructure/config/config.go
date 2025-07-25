package config

import "os"

type Config struct {
	DatabasePath string `json:"database_path"`
	Port         string `json:"port"`
}

func Load() *Config {
	return &Config{
		DatabasePath: getEnv("DATABASE_PATH", "ecommerce.db"),
		Port:         getEnv("PORT", "3000"),
	}
}

func getEnv(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

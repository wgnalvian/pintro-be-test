package config

import (
	"os"
)

type Config struct {
	MONGO_URL     string
	DATABASE_NAME string
	APP_DEBUG     bool
	JWT_SECRET    string
	MIDTRANS_KEY  string
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() *Config {
	return &Config{
		APP_DEBUG:     GetEnv("APP_DEBUG", "true") == "true",
		MONGO_URL:     GetEnv("MONGO_URL", "mongodb://localhost:27017"),
		DATABASE_NAME: GetEnv("DATABASE_NAME", "payment"),
		JWT_SECRET:    GetEnv("JWT_SECRET", "secret"),
		MIDTRANS_KEY:  GetEnv("MIDTRANS_KEY", ""),
	}
}

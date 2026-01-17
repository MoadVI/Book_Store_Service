package config

import "os"

type Config struct {
	DBPath     string
	ServerPort string
}

func LoadConfig() *Config {
	return &Config{
		DBPath:     getEnv("DB_PATH", "internal/db/database.json"),
		ServerPort: getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

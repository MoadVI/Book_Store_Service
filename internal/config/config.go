package config

import (
	"time"
)

type Config struct {
	DBPath                string
	ServerPort            string
	ReportInterval        time.Duration
	ReportOutputDirectory string
}

func LoadConfig() *Config {
	return &Config{
		DBPath:                "database.json",
		ServerPort:            "8080",
		ReportInterval:        24 * time.Hour,
		ReportOutputDirectory: "output-reports",
	}
}

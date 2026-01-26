package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN       string
	DB        string
	GinConfig string
	Port      string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DSN:       getEnv("DB_URI", ""),
		DB:        getEnv("MONGO_DB", ""),
		GinConfig: getEnv("GIN_MODE", "debug"),
		Port:      getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, df string) string {
	res := os.Getenv(key)
	if res == "" {
		return df
	}
	return res
}

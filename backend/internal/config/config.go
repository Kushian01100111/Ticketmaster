package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/session"
	redisDB "github.com/Kushian01100111/Tickermaster/internal/storage/redisdb"
	"github.com/joho/godotenv"
)

var (
	ErrJWTSecretEmpty = errors.New("jwt secret not loaded")
	ErrTTLEmpty       = errors.New("TTL secret not loaded")
	ErrSkewEmpty      = errors.New("clockskew secret not loaded")
	ErrLoading        = errors.New("error loading")
)

type Config struct {
	DSN          string
	DB           string
	GinConfig    string
	Port         string
	ResendAPIKey string
	EmailFrom    string
	JWTSecrets   session.JWTConfig
	RDBSecrets   redisDB.RedisSecrets
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	jwt := getEnv("JWT_SECRET", "")
	if jwt == "" {
		return nil, ErrJWTSecretEmpty
	}

	issuer := getEnv("HTTP_ISSUER", "")
	if issuer == "" {
		return nil, ErrLoading
	}

	audience := getEnv("HTTP_AUDIENCE", "")
	if audience == "" {
		return nil, ErrLoading
	}

	ttl := getEnv("HTTP_TTL", "")
	if ttl == "" {
		return nil, ErrTTLEmpty
	}

	TTL, err := strconv.Atoi(ttl)
	if err != nil {
		return nil, ErrLoading
	}

	skew := getEnv("HTTP_SKEW", "")
	if skew == "" {
		return nil, ErrSkewEmpty
	}

	Skew, err := strconv.Atoi(skew)
	if err != nil {
		return nil, ErrLoading
	}

	secrets := session.JWTConfig{
		Secret:    jwt,
		Issuer:    issuer,
		Audience:  audience,
		AccessTTL: time.Duration(TTL) * time.Minute,
		ClockSkew: time.Duration(Skew) * time.Minute,
	}

	rdbAddr := getEnv("REDIS_ADDR", "")
	if rdbAddr == "" {
		return nil, ErrLoading
	}

	rdbUser := getEnv("REDIS_USER", "")
	if rdbUser == "" {
		return nil, ErrLoading
	}

	rdbPassword := getEnv("REDIS_PASSWORD", "")
	if rdbPassword == "" {
		return nil, ErrLoading
	}

	rdb := redisDB.RedisSecrets{
		Addr:     rdbAddr,
		User:     rdbUser,
		Password: rdbPassword,
	}

	return &Config{
		DSN:          getEnv("DB_URI", ""),
		DB:           getEnv("MONGO_DB", ""),
		GinConfig:    getEnv("GIN_MODE", "debug"),
		Port:         getEnv("PORT", "8080"),
		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		EmailFrom:    getEnv("RESEND_EMAIL", ""),
		JWTSecrets:   secrets,
		RDBSecrets:   rdb,
	}, nil
}

func getEnv(key, df string) string {
	res := os.Getenv(key)
	if res == "" {
		return df
	}
	return res
}

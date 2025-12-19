package config

import (
	"os"
	"strconv"
	"strings"
	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	JWTSecret string
	RedisURL  string
	AllowedOrigins []string
	RateLimitLogin int
	RateLimitGeneral int
}
	

func LoadConfig() *Config {
	godotenv.Load()
	return &Config{
		DBURL:     getEnv("DATABASE_URL", "postgres://user:secret@localhost:5432/userdb?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", ""), // Load from ENV
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ","),
		
		// Parse Integers with defaults
		RateLimitLogin:   getEnvAsInt("RATE_LIMIT_LOGIN", 5),
		RateLimitGeneral: getEnvAsInt("RATE_LIMIT_GENERAL", 100),
	}
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
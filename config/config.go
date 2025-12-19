package config

import (
	"log"
	"os"
)

type Config struct {
	DBURL     string
	JWTSecret string
	RedisURL  string
}
	

func LoadConfig() *Config {
	redisUrl := os.Getenv("REDIS_URL")
    log.Println("🔍 DEBUG: Loaded REDIS_URL:", redisUrl)
	return &Config{
		DBURL:     getEnv("DATABASE_URL", "postgres://user:secret@localhost:5432/userdb?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", ""), // Load from ENV
		RedisURL: "redis://localhost:6379",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
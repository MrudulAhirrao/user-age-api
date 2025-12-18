package config

import "os"

type Config struct {
	DBURL     string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		DBURL:     getEnv("DATABASE_URL", "postgres://user:secret@localhost:5432/userdb?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", ""), // Load from ENV
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string
	ReadTimeout time.Duration
	WriteTimeot time.Duration
	IdleTimeout time.Duration

	DatabaseURL      string
	DatabaseMaxConns int
	DatabaseMinConns int

	BcryptCost int

	JWTSecret string
}

func Load() *Config {

	dbURL := getEnv("URLL", "")

	if dbURL == "" {
		host := getEnv("HOST", "localhost")
		port := getEnv("POSTGRES_PORT", "")
		user := getEnv("POSTGRES_USER", "")
		password := getEnv("POSTGRES_PASSWORD", "")
		dbName := getEnv("POSTGRES_DB", "")
		sslmode := getEnv("POSTGRES_SSLMODE", "disable")

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbName, sslmode)
	}
	fmt.Println(dbURL)
	return &Config{
		Port:             getEnv("PORT", "8081"),
		ReadTimeout:      getEnvAsDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeot:      getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:      getEnvAsDuration("IDLE_TIMEOUT", 60*time.Second),
		DatabaseURL:      dbURL,
		DatabaseMaxConns: getEnvAsInt("DB_MAX_CONNS", 20),
		DatabaseMinConns: getEnvAsInt("DB_MIN_CONNS", 5),
		BcryptCost:       getEnvAsInt("BCRYPT_COST", 12),
		JWTSecret:        getEnv("JWT_SECRET", "default-secret"),
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return value
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return value
	}
	return fallback
}

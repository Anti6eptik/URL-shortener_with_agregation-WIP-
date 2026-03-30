package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	BcryptCost  int
}

func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		host := os.Getenv("HOST")
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv("POSTGRES_PORT")

		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbName := os.Getenv("POSTGRES_DB")

		sslmode := os.Getenv("POSTGRES_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbName, sslmode)
	}

	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		cost = 12
	}
	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		BcryptCost:  cost,
	}
}

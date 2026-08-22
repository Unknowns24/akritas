package config

import "os"

type Config struct {
	HTTPAddr         string
	PostgresDSN      string
	PaginationSecret string
}

func Load() Config {
	return Config{
		HTTPAddr:         envOr("AKRITAS_HTTP_ADDR", ":8080"),
		PostgresDSN:      envOr("AKRITAS_POSTGRES_DSN", "postgres://akritas:akritas@127.0.0.1:5432/akritas?sslmode=disable"),
		PaginationSecret: envOr("AKRITAS_PAGINATION_SECRET", "akritas-dev-pagination-secret"),
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

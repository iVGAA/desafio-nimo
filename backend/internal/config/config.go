package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config agrega variáveis de ambiente usadas pelo servidor e pela conexão com o Postgres.
type Config struct {
	Port         string
	DatabaseURL  string
	AllowOrigins string
}

// Load lê variáveis de ambiente com defaults adequados ao desenvolvimento local.
func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://combustivel:combustivel_dev@localhost:5432/combustivel_db?sslmode=disable"
	}

	origins := os.Getenv("CORS_ALLOW_ORIGINS")
	if origins == "" {
		origins = "http://localhost:3000"
	}

	return Config{
		Port:         port,
		DatabaseURL:  dsn,
		AllowOrigins: origins,
	}, nil
}

// MustPort retorna a porta como inteiro ou erro se inválida.
func (c Config) MustPort() (int, error) {
	p, err := strconv.Atoi(c.Port)
	if err != nil {
		return 0, fmt.Errorf("PORT inválida %q: %w", c.Port, err)
	}
	return p, nil
}

package config

import (
	"fmt"
	"os"
)

// Config описывает конфигурацию приложения.
type Config struct {
	HTTPPort string
	DSN      string
}

// Load читает конфигурацию из переменных окружения.
func Load() (*Config, error) {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=org_structure port=5432 sslmode=disable"
	}
	return &Config{
		HTTPPort: port,
		DSN:      dsn,
	}, nil
}

// Addr возвращает адрес HTTP-сервера.
func (c *Config) Addr() string {
	return fmt.Sprintf(":%s", c.HTTPPort)
}

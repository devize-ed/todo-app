package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	DB       DBConfig
	Server   ServerConfig
	LogLevel string
}

type DBConfig struct {
	DB_URI string
}

type postgresConfig struct {
	user     string `env:"POSTGRES_USER"`
	password string `env:"POSTGRES_PASSWORD"`
	db       string `env:"POSTGRES_DB"`
}

type ServerConfig struct {
	Addr    string        `env:"SERVER_ADDRESS"`
	Timeout time.Duration `env:"SERVER_TIMEOUT"`
}

func (c *Config) Load() error {
	filename := flag.String("config", "/cfg/config.yaml", "config file")
	flag.Parse()

	if _, err := os.Stat(*filename); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %v", err)
	}

	cfg, err := os.ReadFile(*filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(cfg, &c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	

	return nil
}

func (pc *postgresConfig) buildURI() (string, error) {
	if err := env.Parse(&pc); err != nil {
		return "", fmt.Errorf("failed to parse config: %w", err)
	}

	return fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", pc.user, pc.password, pc.db), nil
}

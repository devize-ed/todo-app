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
	Postgres PostgresConfig `yaml:"postgres"`
	Server   ServerConfig   `yaml:"server"`
	LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

type PostgresConfig struct {
	User     string `yaml:"user" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
	DB       string `yaml:"db" env:"POSTGRES_DB"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT"`
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	SSLMode  string `yaml:"ssl_mode" env:"POSTGRES_SSLMODE"`
}

type ServerConfig struct {
	Addr    string        `yaml:"address" env:"SERVER_ADDRESS"`
	Timeout time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT"`
}

func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Addr: "localhost:8080",
		},
		Postgres: PostgresConfig{
			User:     "postgres",
			Password: "postgres",
			DB:       "postgres",
			Port:     5432,
			Host:     "localhost",
			SSLMode:  "disable",
		},
		LogLevel: "debug",
	}
}

func GetFilePath() string {
	filepath := flag.String("config", "/cfg/config.yaml", "config file")
	flag.Parse()
	return *filepath
}

func loadFromFile(filepath string, cfg *Config) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %v", err)
	}

	cfgFile, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(cfgFile, cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return nil
}

func Load(filepath string) (*Config, error) {
	// default config
	cfg := defaultConfig()
	// load config from file
	if err := loadFromFile(filepath, cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// parse config from environment variables
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func (p *PostgresConfig) BuildURI() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", p.User, p.Password, p.Host, p.Port, p.DB, p.SSLMode)
}

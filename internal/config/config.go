package config

import (
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	Addr  string `env:"ADDR" envDefault:":4000"`
	DB    DBconfig
	Redis RedisConfig
	JWT   JWTConfig
	OAuth OAuthConfig
}

type DBconfig struct {
	DSN          string `env:"DB_DSN,required"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
}

type RedisConfig struct {
	// Redis related configurations
}

type JWTConfig struct {
	// JWT related configurations
}

type OAuthConfig struct {
	// OAuth related configurations
}

// you can add more configuration fields as needed
func LoadConfig() (*Config, error) {
	// Load configuration from file, environment variables, etc.
	_ = godotenv.Load()

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}
	return cfg, nil
}

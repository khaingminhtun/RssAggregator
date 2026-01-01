package config

import (
	"log"
	"time"

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
	// Secret is parsed directly into a byte slice for jwt-go
	Secret string `env:"JWT_SECRET,required"`
	Issuer string `env:"JWT_ISSUER" envDefault:"my-app"`

	// caarlos0/env automatically parses strings like "15m" into time.Duration
	AccessExpiry  time.Duration `env:"JWT_ACCESS_EXPIRY" envDefault:"15m"`
	RefreshExpiry time.Duration `env:"JWT_REFRESH_EXPIRY" envDefault:"168h"`
}

type OAuthConfig struct {
	// OAuth related configurations
}

// you can add more configuration fields as needed
func LoadConfig() (*Config, error) {
	// Load configuration from file, environment variables, etc.
	_ = godotenv.Load()

	cfg := &Config{}

	// Parse environment variables into the struct
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}
	return cfg, nil
}

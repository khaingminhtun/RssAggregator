package config

import (
	"fmt"
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
	User         string `env:"POSTGRES_USER,required"`
	Password     string `env:"POSTGRES_PASSWORD,required"`
	Host         string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port         string `env:"POSTGRES_PORT" envDefault:"5432"`
	Name         string `env:"POSTGRES_DB,required"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
}

func (d *DBconfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

type RedisConfig struct {
	// Redis related configurations
	Host         string        `env:"REDIS_HOST" envDefault:"localhost"`
	Port         string        `env:"REDIS_PORT" envDefault:"6379"`
	Password     string        `env:"REDIS_PASSWORD"`
	DB           int           `env:"REDIS_DB" envDefault:"0"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"3s"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, relying on environment variables")
	}

	cfg := &Config{}

	// Parse environment variables into the struct
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}
	return cfg, nil
}

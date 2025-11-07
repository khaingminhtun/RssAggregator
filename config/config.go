package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaingminhtun/rssagg/internal/db"
)

type Config struct {
	DB *db.Queries
}

func NewConfig() *Config {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set in env")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("Cannot create pgx pool:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Cannot ping database:", err)
	}

	return &Config{
		DB: db.New(pool),
	}
}

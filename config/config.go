package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaingminhtun/rssagg/internal/db"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB    *db.Queries
	Redis *redis.Client
	Ctx   context.Context
}

// NewConfi intializes
func NewConfig() *Config {

	ctx := context.Background()

	// -----------------------------
	// PostgreSQL setup
	// -----------------------------
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set in env")
	}
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("Cannot create pgx pool:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Cannot ping database:", err)
	}

	queries := db.New(pool)

	// -----------------------------
	// Redis setup
	// -----------------------------

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("Redis addr not in env")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password by default
		DB:       0,  // default db
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Cannot connect to Redis:", err)
	}

	// Return conbinded config
	return &Config{
		DB:    queries,
		Redis: rdb,
		Ctx:   ctx,
	}
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaingminhtun/rssagg/internal/config"
	"github.com/khaingminhtun/rssagg/internal/pkg/log"
)

func main() {

	// 1. set up logger first
	log.Init()

	// 2. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("unable to load config", "error", err)
		os.Exit(1)
	}

	// 3. Database
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DB.DSN())
	if err != nil {
		log.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Error("database ping failed", "error", err)
		os.Exit(1)
	}

	// 4. Initialize app struct
	app := &application{
		config: *cfg,
		db:     dbPool,
	}

	// 5. Start background tasks (workers + scheduler)
	app.setupBackgroundTasks(ctx)

	// 6. Set up HTTP routes (from api.go)
	routes := app.routes()

	// 7. Create HTTP server
	srv := createServer(app.config.Addr, routes)

	// 8. Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Server started at %s", app.config.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error:", err)
		}
	}()

	<-stop
	log.Info("Shutdown signal received, stopping server and background tasks...")

	// cancel context → stops workers & scheduler
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown:", err)
	} else {
		log.Info("Server shutdown gracefully")
	}

	log.Info("All background tasks stopped")
	fmt.Println("Bye!")
}

// createServer is a helper to create HTTP server
func createServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

package main

import (
	"context"
	"fmt"

	"os"

	"github.com/jackc/pgx/v5"
	"github.com/khaingminhtun/rssagg/internal/config"
	"github.com/khaingminhtun/rssagg/internal/log"
)

func main() {
	ctx := context.Background()

	// 1. set up logger first
	log.Init()

	//2. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("unable to load config", "error", err)
		os.Exit(1)
	}

	//3. connect to database
	dbConn, err := pgx.Connect(ctx, cfg.DB.DSN)
	if err != nil {
		log.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbConn.Close(ctx)

	//4. Initialize application struct
	app := &application{
		config: *cfg,
		db:     dbConn,
	}

	//5. Set up routes
	routes := app.routes()

	//6. Run the server
	log.Info("starting server", "addr", cfg.Addr)
	err = app.run(routes)
	if err != nil {
		log.Error("unable to start server", "error", err)
		os.Exit(1)
	}

	// Application entry point
	fmt.Println("Hello go")
}

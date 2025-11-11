package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	userRouters "github.com/khaingminhtun/rssagg/api/users/routers"
	authRouters "github.com/khaingminhtun/rssagg/auth/routers"
	"github.com/khaingminhtun/rssagg/config"
	middleware "github.com/khaingminhtun/rssagg/middlewares"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, World!")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in env")
	}
	fmt.Printf("Server will start at port: %s\n", port)

	router := chi.NewRouter()

	// 🔥 CORS setup
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://yourdomain.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//middlewars
	router.Use(middleware.Logger)
	router.Use(middleware.RateLimiter(1, 3))

	// Load config
	cfg := config.NewConfig()

	// Register your routes directly on 'router'
	router.Route("/v1", func(r chi.Router) {
		authRouters.Init(r, cfg)
		userRouters.Init(r, cfg)
	})

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %v", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/khaingminhtun/rssagg/handlers"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, World!")

	godotenvErr := godotenv.Load()
	if godotenvErr != nil {
		log.Fatalf("Error loading .env file: %v", godotenvErr)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in env")
	}
	fmt.Printf("Server wil start at port: %s\n", port)

	router := chi.NewRouter()

	// 🔥 CORS setup
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://yourdomain.com"}, // allowed frontends
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Cache preflight requests for 5 minutes
	}))

	// Get method
	v1Router := chi.NewRouter()

	v1Router.Get("/get", handlers.GetHello)
	v1Router.Get("/error", handlers.HandlerError)
	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server startin on port %v", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

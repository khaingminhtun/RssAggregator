package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/auth"
	"github.com/khaingminhtun/rssagg/internal/config"
	"github.com/khaingminhtun/rssagg/internal/feeds"
)

type application struct {
	config config.Config
	db     *pgx.Conn
}

// routes sets up the application routes and middleware
func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	//initialize jwt
	jwtSVc := auth.NewJWTService(app.config.JWT)

	//authentication
	authService := auth.NewAuthService(repo.New(app.db), jwtSVc)
	authHandler := auth.NewAuthHandler(authService)
	r.Post("/api/v1/register", authHandler.RegisterUser)
	r.Post("/api/v1/login", authHandler.Authenticate)
	r.Post("/api/v1/refresh", authHandler.RefreshToken)
	r.Post("/api/v1/logout", authHandler.Logout)

	// --- Feeds ---
	fetcher := &feeds.FetcherService{
		HttpClient: &http.Client{},
	}
	feedService := feeds.NewFeedService(repo.New(app.db), fetcher)
	feedHandler := feeds.NewFeedHandler(feedService)

	r.Post("/api/v1/createFeed", feedHandler.CreateFeed)
	r.Get("/api/v1/feeds/{feedID}/posts", feedHandler.GetFeedPosts)

	//protected routes example)
	r.Group(func(r chi.Router) {
		r.Use(jwtSVc.JWTMiddleware)
		r.Get("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
			userID, _ := auth.UserIDFromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("your user id: %d", userID)))
		})
	})
	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Printf("server has started at addr %s", app.config.Addr)
	return srv.ListenAndServe()
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/auth"
	"github.com/khaingminhtun/rssagg/internal/config"
	"github.com/khaingminhtun/rssagg/internal/feeds"
	"github.com/khaingminhtun/rssagg/internal/mechanism"
	"github.com/khaingminhtun/rssagg/internal/posts"
)

type application struct {
	config config.Config
	db     *pgxpool.Pool
}

// routes sets up the application routes and middleware
func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// --- Global middleware --
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// --- Health check ---
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	// ------ Auth -----------
	// --- Auth setup ---
	jwtSVc := auth.NewJWTService(app.config.JWT)
	authService := auth.NewAuthService(repo.New(app.db), jwtSVc)
	authHandler := auth.NewAuthHandler(authService)

	// -- Auth routes ---
	r.Post("/api/v1/register", authHandler.RegisterUser)
	r.Post("/api/v1/login", authHandler.Authenticate)
	r.Post("/api/v1/refresh", authHandler.RefreshToken)
	r.Post("/api/v1/logout", authHandler.Logout)

	// ------Feeds --------
	// ---- Feed setup ----
	fetcher := &feeds.FetcherService{
		HttpClient: &http.Client{},
	}
	// repositories
	feedRepo := feeds.NewFeedRepository(app.db)
	postRepo := posts.NewPostRepository(app.db)
	// services and handlers
	feedService := feeds.NewFeedService(feedRepo, postRepo, fetcher)
	feedHandler := feeds.NewFeedHandler(feedService)

	// ----- Posts -----
	// ------ Posts setup ----

	postService := posts.NewPostService(postRepo)
	postHandler := posts.NewPostHandler(postService)

	// Public feed routes
	r.Route("/api/v1", func(r chi.Router) {
		// create feed
		r.Post("/createFeed", feedHandler.CreateFeed)

		// feed posts
		// r.Get("/feeds/{feedID}/posts", feedHandler.GetFeedPosts)

		// get feed by ID
		r.Get("/feeds/{id}", feedHandler.GetFeedByID)

		// get all feeds
		r.Get("/feeds", feedHandler.GetAllFeeds)

		// get feeds by userID
		r.Get("/users/{userID}/feeds", feedHandler.GetFeedsByUserID)

		// update feeds by feedid
		r.Patch("/feeds/{id}", feedHandler.UpdateFeed)

		r.Delete("/feeds/{id}", feedHandler.DeleteFeedByID)
		r.Delete("/feeds/{id}/unused", feedHandler.DeleteFeedIfUnused)
		r.Delete("/users/{user_id}/feeds", feedHandler.DeleteAllFeedsByUserID)

		// get all posts
		r.Get("/posts", postHandler.GetAllPosts)

	})

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

// run background tasks
func (app *application) setupBackgroundTasks() *mechanism.FeedQueue {
	log.Println("[BACKGROUND] Initializing background feed processing")

	// --- Repositories ---
	feedRepo := feeds.NewFeedRepository(app.db)
	postRepo := posts.NewPostRepository(app.db)

	// --- Fetcher service ---
	fetcher := &feeds.FetcherService{
		HttpClient: &http.Client{},
	}

	// --- Feed processor ---
	processor := mechanism.NewFeedProcessorService(fetcher, postRepo, feedRepo)

	// --- Queue ---
	queue := mechanism.NewFeedQueue(100)
	log.Println("[BACKGROUND] Feed queue initialized (size=100)")

	// --- Workers ---
	worker := mechanism.NewFeedWorker(processor, queue)
	ctx := context.Background()
	worker.Start(ctx, 5)
	log.Println("[BACKGROUND] Feed workers started (count=5)")

	// --- Scheduler ---
	scheduler := mechanism.NewFeedScheduler(feedRepo, queue, 5*time.Second)
	go scheduler.Start(ctx)
	log.Println("[BACKGROUND] Feed scheduler started (interval=5s)")

	log.Println("[BACKGROUND] Background feed system is running")

	return queue
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

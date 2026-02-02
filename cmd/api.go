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
	"github.com/khaingminhtun/rssagg/internal/adapters/cache"
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
	redis  *cache.RedisStore
}

// routes sets up the application routes and middleware
func (app *application) routes() *chi.Mux {
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
	jwtSvc := auth.NewJWTService(app.config.JWT)
	authService := auth.NewAuthService(repo.New(app.db), jwtSvc)
	authHandler := auth.NewAuthHandler(authService)

	r.Post("/api/v1/register", authHandler.RegisterUser)
	r.Post("/api/v1/login", authHandler.Authenticate)
	r.Post("/api/v1/refresh", authHandler.RefreshToken)
	r.Post("/api/v1/logout", authHandler.Logout)

	// ------ Feeds & Posts ------
	// 3. Domain caches
// feedCache := cache.NewFeedCache(app.redis, 1*time.Hour)
	fetcher := feeds.NewFetcherService(15*time.Second, 5, 3)
	feedRepo := feeds.NewFeedRepository(app.db)
	postRepo := posts.NewPostRepository(app.db)

	feedService := feeds.NewFeedService(feedRepo, postRepo, fetcher)
	feedHandler := feeds.NewFeedHandler(feedService)

	postService := posts.NewPostService(postRepo)
	postHandler := posts.NewPostHandler(postService)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/createFeed", feedHandler.CreateFeed)
		r.Get("/feeds/{id}", feedHandler.GetFeedByID)
		r.Get("/feeds", feedHandler.GetAllFeeds)
		r.Get("/users/{userID}/feeds", feedHandler.GetFeedsByUserID)
		r.Patch("/feeds/{id}", feedHandler.UpdateFeed)
		r.Delete("/feeds/{id}", feedHandler.DeleteFeedByID)
		r.Delete("/feeds/{id}/unused", feedHandler.DeleteFeedIfUnused)
		r.Delete("/users/{user_id}/feeds", feedHandler.DeleteAllFeedsByUserID)
		r.Get("/posts", postHandler.GetAllPosts)
	})

	// protected routes example
	r.Group(func(r chi.Router) {
		r.Use(jwtSvc.JWTMiddleware)
		r.Get("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
			userID, _ := auth.UserIDFromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("your user id: %d", userID)))
		})
	})

	return r
}

// setupBackgroundTasks starts the feed queue, workers, and scheduler
func (app *application) setupBackgroundTasks(ctx context.Context) *mechanism.FeedQueue {
	log.Println("[BACKGROUND] Initializing background feed processing")

	feedRepo := feeds.NewFeedRepository(app.db)
	postRepo := posts.NewPostRepository(app.db)
	fetcher := feeds.NewFetcherService(15*time.Second, 5, 3)
	processor := mechanism.NewFeedProcessorService(fetcher, postRepo, feedRepo)

	queue := mechanism.NewFeedQueue(100)
	log.Println("[BACKGROUND] Feed queue initialized (size=100)")

	worker := mechanism.NewFeedWorker(processor, queue)
	worker.Start(ctx, 5)
	log.Println("[BACKGROUND] Feed workers started (count=5)")

	scheduler := mechanism.NewFeedScheduler(feedRepo, queue, 5*time.Second)
	go scheduler.Start(ctx)
	log.Println("[BACKGROUND] Feed scheduler started (interval=5s)")

	log.Println("[BACKGROUND] Background feed system is running")
	return queue
}

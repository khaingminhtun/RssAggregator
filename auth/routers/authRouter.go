package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/khaingminhtun/rssagg/auth/handlers"
	"github.com/khaingminhtun/rssagg/auth/services"
	"github.com/khaingminhtun/rssagg/config"
)

// for users
func Init(r chi.Router, cfg *config.Config) {
	// Initialize service & handler inside feature
	service := services.NewAuthService(cfg)
	authHandler := handlers.NewAuthHandler(service)

	// Register routes under this router
	SetupRoutes(r, authHandler)
}

func SetupRoutes(r chi.Router, authHandler *handlers.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Authenicate)
		r.Post("/getrefreshtoken", authHandler.GetRefreshToken)
		r.Post("/logout", authHandler.Logout)

		// Add other endpoints here
	})
}

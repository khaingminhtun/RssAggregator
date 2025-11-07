package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/khaingminhtun/rssagg/config"
)

// for users
func Init(r chi.Router, cfg *config.Config) {
	// Initialize service & handler inside feature
	userService := NewUserService(cfg)
	userHandler := NewUserHandler(userService)

	// Register routes under this router
	SetupRoutes(r, userHandler)
}

func SetupRoutes(r chi.Router, userHandler *UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/create-user", userHandler.CreateUser)
		// Add other endpoints here
	})
}

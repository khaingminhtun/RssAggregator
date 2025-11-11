package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/khaingminhtun/rssagg/api/users/handlers"
	"github.com/khaingminhtun/rssagg/api/users/services"
	"github.com/khaingminhtun/rssagg/config"
)

// for users
func Init(r chi.Router, cfg *config.Config) {
	// Initialize service & handler inside feature
	userService := services.NewUserService(cfg)
	userHandler := handlers.NewUserHandler(userService)

	// Register routes under this router
	SetupRoutes(r, userHandler)
}

func SetupRoutes(r chi.Router, userHandler *handlers.UserHandler) {
	r.Route("/users", func(r chi.Router) {

		r.Get("/{id}", userHandler.GetUserById)
		// Add other endpoints here
	})
}

package users

import (
	"context"

	"github.com/khaingminhtun/rssagg/config"
	"github.com/khaingminhtun/rssagg/internal/db"
)

type UserService struct {
	cfg *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{
		cfg: cfg,
	}
}

// RegisterUser handles user creation logic
func (s *UserService) RegisterUser(ctx context.Context, name, email string) (*db.User, error) {
	user, err := s.cfg.DB.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

package services

import (
	"context"
	"errors"

	"github.com/khaingminhtun/rssagg/auth/dtos"
	"github.com/khaingminhtun/rssagg/config"
	"github.com/khaingminhtun/rssagg/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	config *config.Config
}

func NewAuthService(c *config.Config) *AuthService {
	return &AuthService{
		config: c,
	}
}

// register user
func (a *AuthService) RegisterUser(ctx context.Context, RegisterRequest dtos.RegisterRequest) (*db.User, error) {
	//1. Hash passwword
	hash, err := bcrypt.GenerateFromPassword([]byte(RegisterRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	//2. store in db
	user, err := a.config.DB.CreateUser(ctx, db.CreateUserParams{
		Name:         RegisterRequest.Name,
		Email:        RegisterRequest.Email,
		PasswordHash: string(hash),
	})

	if err != nil {
		return nil, errors.New("email alreay exits or db error")
	}

	return &db.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil

}

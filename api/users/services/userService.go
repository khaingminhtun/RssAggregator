package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khaingminhtun/rssagg/config"
	"github.com/khaingminhtun/rssagg/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	config *config.Config
}

func NewUserService(c *config.Config) *UserService {
	return &UserService{
		config: c,
	}
}

// register user
func (s *UserService) RegisterUser(ctx context.Context, name, email, passwordHash string) (*db.User, error) {
	//1. Hash passwword
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	//2. store in db
	user, err := s.config.DB.CreateUser(ctx, db.CreateUserParams{
		Name:         name,
		Email:        email,
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

// get user one
func (s *UserService) GetUserById(ctx context.Context, id pgtype.UUID) (*db.User, error) {
	user, err := s.config.DB.GetUserById(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil

}

//get all users

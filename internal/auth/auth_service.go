package auth

import (
	"context"

	repo "github.com/khaingminhtun/rssagg/internal/adapters/database/db"
	"github.com/khaingminhtun/rssagg/internal/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/log"
	"github.com/khaingminhtun/rssagg/internal/utils"
)

type AuthService interface {
	RegisterUser(ctx context.Context, request RegisterRequest) error
}

type service struct {
	repo repo.Querier
}

func NewAuthService(repo repo.Querier) AuthService {
	return &service{repo: repo}
}

func (s *service) RegisterUser(ctx context.Context, request RegisterRequest) error {

	//1. input
	if request.Name == "" || request.Email == "" || request.Password == "" {
		return errorHandle.ErrInvalidInput
	}

	// 2. hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return errorHandle.ErrInternal
	}

	//3. create user in db
	user, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: hashedPassword,
		AuthType:     "local",
	})

	if err != nil {
		// Handle DB-level errs
		if errorHandle.IsUniqueViolation(err) {
			return errorHandle.ErrUserAlreadyExists
		}
		return errorHandle.ErrInternal
	}
	log.Info("user registered successfully", "user_id", user.ID)

	return nil
}

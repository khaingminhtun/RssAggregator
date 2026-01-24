package auth

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/pkg/log"
	"github.com/khaingminhtun/rssagg/internal/pkg/utils"
)

type AuthService interface {
	RegisterUser(ctx context.Context, request RegisterRequest) error
	AuthenticateUser(ctx context.Context, request AuthRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
}

type service struct {
	repo repo.Querier
	jwt  *JWTService
}

func NewAuthService(repo repo.Querier, jwtSvc *JWTService) AuthService {
	return &service{repo: repo, jwt: jwtSvc}
}

// 1. register user
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
		AuthType:     string(AuthTypeLocal),
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

// 2. Authenticate user
func (s *service) AuthenticateUser(ctx context.Context, request AuthRequest) (*AuthResponse, error) {
	// 1. input
	if request.Email == "" || request.Password == "" {
		return nil, errorHandle.ErrInvalidInput
	}

	// 2. get user from db
	user, err := s.repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, errorHandle.ErrUserNotFound
	}

	// 3. verify password
	if err := utils.CheckPassword(user.PasswordHash, request.Password); err != nil {
		return nil, errorHandle.ErrInvalidCredentials
	}

	// 4. generate access and refresh tokens
	accessToken, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, errorHandle.ErrInternal
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errorHandle.ErrInternal
	}

	log.Info("user authenticated successfully", "user_id", user.ID)

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	claims, err := s.jwt.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID := claims.UserID

	accessToken, err := s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, errorHandle.ErrInternal
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, errorHandle.ErrInternal
	}

	log.Info("refresh token route", "user_id", userID)

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/khaingminhtun/rssagg/auth/dtos"
	"github.com/khaingminhtun/rssagg/auth/utils"
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

// authenticate user
func (a *AuthService) AuthenticateUser(ctx context.Context, req dtos.AuthRequest) (*dtos.AuthResponse, error) {
	// fetch user from db
	user, err := a.config.DB.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found with email: %s", req.Email)

	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")

	}

	// generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	// 4️⃣ Store refresh token in Redis (important!)
	redisKey := "refresh:" + user.ID.String()

	// Save refresh token with expiry (same as token expiry)
	err = utils.SetKey(a.config, redisKey, refreshToken, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token in redis: %v", err)
	}

	// 5️⃣ Return response
	return &dtos.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// refresh token
func (a *AuthService) RefreshToken(ctx context.Context, oldRefreshToken string) (*dtos.AuthResponse, error) {
	// 1️⃣ Decode the token
	token, err := utils.DecodeToken(oldRefreshToken)
	if err != nil {
		log.Printf("[RefreshToken] DecodeToken error: %v\n", err)
		return nil, errors.New("invalid refresh token")
	}

	// 2️⃣ Validate token type
	if err := utils.ValidateTokenType(token, "refresh"); err != nil {
		log.Printf("[RefreshToken] ValidateTokenType error: %v\n", err)
		return nil, errors.New("invalid refresh token type")
	}

	// 3️⃣ Get claims
	claims, err := utils.GetClaims(oldRefreshToken)
	if err != nil {
		log.Printf("[RefreshToken] GetClaims error: %v\n", err)
		return nil, errors.New("cannot read token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Printf("[RefreshToken] missing 'sub' in token claims: %v\n", claims)
		return nil, errors.New("invalid token claims")
	}

	// 4️⃣ Check refresh token in Redis
	redisKey := "refresh:" + userID
	storedToken, err := utils.GetKey(a.config, redisKey)
	if err != nil {
		log.Printf("[RefreshToken] Redis GetKey error: %v\n", err)
		return nil, errors.New("cannot verify refresh token")
	}

	if storedToken != oldRefreshToken {
		log.Printf("[RefreshToken] refresh token mismatch: expected %s, got %s\n", storedToken, oldRefreshToken)
		return nil, errors.New("refresh token expired or invalid")
	}

	// 5️⃣ Generate new tokens
	accessToken, newRefreshToken, err := utils.GenerateTokens(userID, nil)
	if err != nil {
		log.Printf("[RefreshToken] GenerateTokens error: %v\n", err)
		return nil, errors.New("failed to generate new tokens")
	}

	// 6️⃣ Overwrite refresh token in Redis
	if err := utils.SetKey(a.config, redisKey, newRefreshToken, 7*24*time.Hour); err != nil {
		log.Printf("[RefreshToken] Redis SetKey error: %v\n", err)
		return nil, errors.New("failed to store refresh token")
	}

	// 7️⃣ Return new tokens
	return &dtos.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// LogoutUser handles business logic of logout
func (s *AuthService) LogoutUser(ctx context.Context, refreshToken string) error {
	// Decode refresh token
	_, err := utils.DecodeToken(refreshToken)
	if err != nil {
		return err
	}

	claims, err := utils.GetClaims(refreshToken)
	if err != nil {
		return err
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return err
	}

	// Delete refresh token from Redis
	redisKey := "refresh:" + userID
	return utils.DeleteKey(s.config, redisKey)
}

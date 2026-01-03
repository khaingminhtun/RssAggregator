package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khaingminhtun/rssagg/internal/config"
)

// Define errors for better error handling in your handlers
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims defines the structure of the data stored inside the JWT
type Claims struct {
	UserID int32 `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenPair represents the response sent to the frontend
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// JWTService handles the creation and parsing of tokens
type JWTService struct {
	secret        []byte
	issuer        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTService initializes the service with config values
func NewJWTService(cfg config.JWTConfig) *JWTService {
	return &JWTService{
		secret:        []byte(cfg.Secret),
		issuer:        cfg.Issuer,
		accessExpiry:  cfg.AccessExpiry,
		refreshExpiry: cfg.RefreshExpiry,
	}
}

// GenerateAccessToken creates a short-lived token (e.g., 15 mins) for API access
func (s *JWTService) GenerateAccessToken(userID int32) (string, error) {
	return s.createToken(userID, s.accessExpiry)
}

// GenerateRefreshToken creates a long-lived token (e.g., 7 days)
func (s *JWTService) GenerateRefreshToken(userID int32) (string, error) {
	return s.createToken(userID, s.refreshExpiry)
}

// createToken is a private helper to sign the JWT
func (s *JWTService) createToken(userID int32, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken parses the string and returns the Claims if valid
func (s *JWTService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

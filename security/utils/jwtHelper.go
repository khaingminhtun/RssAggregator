package utils

import (
	"context"
	"errors"
	"time"

	"github.com/khaingminhtun/rssagg/security/jwtauth"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Default durations for access and refresh tokens
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// GenerateAccessToken generates a JWT access token with claims and expiry
func GenerateAccessToken(ja *jwtauth.JWTAuth, userID string, extraClaims map[string]interface{}) (string, error) {
	claims := map[string]interface{}{
		"sub": userID, // subject: user id
		"typ": "access",
	}
	for k, v := range extraClaims {
		claims[k] = v
	}
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, AccessTokenDuration)

	_, tokenStr, err := ja.Encode(claims)
	return tokenStr, err
}

// GenerateRefreshToken generates a JWT refresh token with claims and expiry
func GenerateRefreshToken(ja *jwtauth.JWTAuth, userID string, extraClaims map[string]interface{}) (string, error) {
	claims := map[string]interface{}{
		"sub": userID,
		"typ": "refresh",
	}
	for k, v := range extraClaims {
		claims[k] = v
	}
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, RefreshTokenDuration)

	_, tokenStr, err := ja.Encode(claims)
	return tokenStr, err
}

// GenerateTokens is a helper to create both access & refresh tokens for login
func GenerateTokens(ja *jwtauth.JWTAuth, userID string, extraClaims map[string]interface{}) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(ja, userID, extraClaims)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateRefreshToken(ja, userID, extraClaims)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// ValidateTokenType ensures the token is of the expected type (access/refresh)
func ValidateTokenType(token jwt.Token, expectedType string) error {
	claims, err := token.AsMap(context.Background())
	if err != nil {
		return err
	}
	if typ, ok := claims["typ"].(string); !ok || typ != expectedType {
		return errors.New("token is unauthorized")
	}
	return nil
}

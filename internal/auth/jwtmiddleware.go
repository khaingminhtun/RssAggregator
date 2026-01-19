package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
)

type contextKey string

const userIDKey contextKey = "userID"

// JWTMiddleware validates the access token and injects userID into context
func (s *JWTService) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorHandle.RespondHTTPError(w, ErrInvalidToken)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			errorHandle.RespondHTTPError(w, ErrInvalidToken)
			return
		}

		claims, err := s.ValidateToken(parts[1])
		if err != nil {
			errorHandle.RespondHTTPError(w, err)
			return
		}

		// Inject userID into request context
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext retrieves userID from context
func UserIDFromContext(ctx context.Context) (int32, bool) {
	userID, ok := ctx.Value(userIDKey).(int32)
	return userID, ok
}

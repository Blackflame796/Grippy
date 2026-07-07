package middlewares

import (
	"context"
	"net/http"
	"strings"

	entity "Grippy/internal/domain"
)

type contextKey string

const UserKey contextKey = "user_info"

type TokenParser interface {
	ParseAccessToken(tokenStr string) (*entity.Claims, error)
}

func NewAuthMiddleware(parser TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			accessToken := parts[1]

			claims, err := parser.ParseAccessToken(accessToken)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

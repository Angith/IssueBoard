package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/supabase-community/gotrue-go"
)

type contextKey string

const UserKey contextKey = "user"

func AuthMiddleware(authClient gotrue.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			user, err := authClient.WithToken(token).GetUser()
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

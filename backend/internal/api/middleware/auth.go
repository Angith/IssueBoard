package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/repository/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(supabaseURL string, userRepo *repository.UserRepository) func(http.Handler) http.Handler {
	jwksURL := supabaseURL + "/auth/v1/.well-known/jwks.json"
	kf, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		fmt.Printf("Failed to create JWKS keyfunc: %v\n", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, kf.Keyfunc)

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userIDStr, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Missing user_id in token", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user_id format", http.StatusUnauthorized)
				return
			}

			// Ensure user exists in our local DB for FK constraints
			email, _ := claims["email"].(string)
			user := &models.User{
				ID:    userID,
				Email: email,
			}
			if err := userRepo.CreateOrUpdate(r.Context(), user); err != nil {
				// We don't necessarily want to fail here if it's just a sync issue,
				// but for this project we'll assume it's critical.
				http.Error(w, "Failed to sync user", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userIDStr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	userID, _ := ctx.Value(UserIDKey).(string)
	return userID
}

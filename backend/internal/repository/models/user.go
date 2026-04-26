package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	GitHubID   string    `json:"github_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	OAuthToken string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	ID           uuid.UUID `json:"id"`
	GitHubRepoID int64     `json:"github_repo_id"`
	FullName     string    `json:"full_name"`
	Owner        string    `json:"owner"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
}

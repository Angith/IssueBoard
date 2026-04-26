package models

import (
	"time"

	"github.com/google/uuid"
)

type Issue struct {
	ID            uuid.UUID `json:"id"`
	RepositoryID  uuid.UUID `json:"repository_id"`
	GitHubIssueID int64     `json:"github_issue_id"`
	Number        int       `json:"number"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	State         string    `json:"state"`
	URL           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Labels        []*Label  `json:"labels,omitempty"`
}

type Label struct {
	ID           uuid.UUID `json:"id"`
	RepositoryID uuid.UUID `json:"repository_id"`
	Name         string    `json:"name"`
	Color        string    `json:"color"`
	Description  string    `json:"description"`
}

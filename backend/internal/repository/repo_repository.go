package repository

import (
	"context"

	"github.com/angith/issueboard/internal/repository/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoRepository struct {
	db *pgxpool.Pool
}

func NewRepoRepository(db *pgxpool.Pool) *RepoRepository {
	return &RepoRepository{db: db}
}

func (r *RepoRepository) Create(ctx context.Context, repo *models.Repository) error {
	query := `
		INSERT INTO repositories (user_id, github_repo_id, full_name, owner, name, url)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, github_repo_id) DO UPDATE SET
			full_name = EXCLUDED.full_name,
			owner = EXCLUDED.owner,
			name = EXCLUDED.name,
			url = EXCLUDED.url
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query, repo.UserID, repo.GitHubRepoID, repo.FullName, repo.Owner, repo.Name, repo.URL).
		Scan(&repo.ID, &repo.CreatedAt)
}

func (r *RepoRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Repository, error) {
	query := `SELECT id, user_id, github_repo_id, full_name, owner, name, url, created_at FROM repositories WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repository
	for rows.Next() {
		repo := &models.Repository{}
		err := rows.Scan(&repo.ID, &repo.UserID, &repo.GitHubRepoID, &repo.FullName, &repo.Owner, &repo.Name, &repo.URL, &repo.CreatedAt)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

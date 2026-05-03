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

func (r *RepoRepository) AddToUser(ctx context.Context, userID uuid.UUID, repo *models.Repository) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Upsert the repository globally
	query := `
		INSERT INTO repositories (github_repo_id, full_name, owner, name, url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (github_repo_id) DO UPDATE SET
			full_name = EXCLUDED.full_name,
			owner = EXCLUDED.owner,
			name = EXCLUDED.name,
			url = EXCLUDED.url
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, query, repo.GitHubRepoID, repo.FullName, repo.Owner, repo.Name, repo.URL).
		Scan(&repo.ID, &repo.CreatedAt)
	if err != nil {
		return err
	}

	// 2. Link the repository to the user
	linkQuery := `
		INSERT INTO user_repository (user_id, repository_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, repository_id) DO NOTHING
	`
	_, err = tx.Exec(ctx, linkQuery, userID, repo.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *RepoRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Repository, error) {
	query := `
		SELECT r.id, r.github_repo_id, r.full_name, r.owner, r.name, r.url, r.created_at
		FROM repositories r
		JOIN user_repository ur ON r.id = r.repository_id
		WHERE ur.user_id = $1
	`
	// Wait, typo in join: ur.repository_id
	query = `
		SELECT r.id, r.github_repo_id, r.full_name, r.owner, r.name, r.url, r.created_at
		FROM repositories r
		JOIN user_repository ur ON r.id = ur.repository_id
		WHERE ur.user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repository
	for rows.Next() {
		repo := &models.Repository{}
		err := rows.Scan(&repo.ID, &repo.GitHubRepoID, &repo.FullName, &repo.Owner, &repo.Name, &repo.URL, &repo.CreatedAt)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func (r *RepoRepository) GetByID(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) (*models.Repository, error) {
	query := `
		SELECT r.id, r.github_repo_id, r.full_name, r.owner, r.name, r.url, r.created_at
		FROM repositories r
		JOIN user_repository ur ON r.id = ur.repository_id
		WHERE ur.user_id = $1 AND r.id = $2
	`
	repo := &models.Repository{}
	err := r.db.QueryRow(ctx, query, userID, repoID).
		Scan(&repo.ID, &repo.GitHubRepoID, &repo.FullName, &repo.Owner, &repo.Name, &repo.URL, &repo.CreatedAt)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

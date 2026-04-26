package repository

import (
	"context"

	"github.com/angith/issueboard/internal/repository/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateOrUpdate(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (github_id, username, email, oauth_token, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (github_id) DO UPDATE SET
			username = EXCLUDED.username,
			email = EXCLUDED.email,
			oauth_token = EXCLUDED.oauth_token,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, user.GitHubID, user.Username, user.Email, user.OAuthToken).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByGitHubID(ctx context.Context, githubID string) (*models.User, error) {
	query := `SELECT id, github_id, username, email, oauth_token, created_at, updated_at FROM users WHERE github_id = $1`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, githubID).
		Scan(&user.ID, &user.GitHubID, &user.Username, &user.Email, &user.OAuthToken, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

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
		INSERT INTO users (id, email, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			updated_at = CURRENT_TIMESTAMP
		RETURNING created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, user.ID, user.Email).
		Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, email, created_at, updated_at FROM users WHERE id = $1`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) SaveGitHubToken(ctx context.Context, userID string, encryptedToken []byte) error {
	query := `
		INSERT INTO user_settings (user_id, github_token, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id) DO UPDATE SET
			github_token = EXCLUDED.github_token,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(ctx, query, userID, encryptedToken)
	return err
}

func (r *UserRepository) GetGitHubToken(ctx context.Context, userID string) ([]byte, error) {
	query := `SELECT github_token FROM user_settings WHERE user_id = $1`
	var token []byte
	err := r.db.QueryRow(ctx, query, userID).Scan(&token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *UserRepository) DeleteGitHubToken(ctx context.Context, userID string) error {
	query := `UPDATE user_settings SET github_token = NULL WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

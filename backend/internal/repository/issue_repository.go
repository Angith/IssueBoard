package repository

import (
	"context"

	"github.com/angith/issueboard/internal/repository/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IssueRepository struct {
	db *pgxpool.Pool
}

func NewIssueRepository(db *pgxpool.Pool) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) CreateOrUpdate(ctx context.Context, issue *models.Issue) error {
	query := `
		INSERT INTO issues (repository_id, github_issue_id, number, title, body, state, url, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (repository_id, github_issue_id) DO UPDATE SET
			number = EXCLUDED.number,
			title = EXCLUDED.title,
			body = EXCLUDED.body,
			state = EXCLUDED.state,
			url = EXCLUDED.url,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query, issue.RepositoryID, issue.GitHubIssueID, issue.Number, issue.Title, issue.Body, issue.State, issue.URL, issue.UpdatedAt).
		Scan(&issue.ID, &issue.CreatedAt)
}

func (r *IssueRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.Issue, error) {
	query := `SELECT id, repository_id, github_issue_id, number, title, body, state, url, created_at, updated_at FROM issues WHERE repository_id = $1`
	rows, err := r.db.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []*models.Issue
	for rows.Next() {
		issue := &models.Issue{}
		err := rows.Scan(&issue.ID, &issue.RepositoryID, &issue.GitHubIssueID, &issue.Number, &issue.Title, &issue.Body, &issue.State, &issue.URL, &issue.CreatedAt, &issue.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Load labels for this issue
		labels, err := r.GetLabelsForIssue(ctx, issue.ID)
		if err != nil {
			return nil, err
		}
		issue.Labels = labels

		issues = append(issues, issue)
	}
	return issues, nil
}

func (r *IssueRepository) GetLabelsForIssue(ctx context.Context, issueID uuid.UUID) ([]*models.Label, error) {
	query := `
		SELECT l.id, l.repository_id, l.name, l.color, l.description
		FROM labels l
		JOIN issue_labels il ON l.id = il.label_id
		WHERE il.issue_id = $1
	`
	rows, err := r.db.Query(ctx, query, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []*models.Label
	for rows.Next() {
		l := &models.Label{}
		err := rows.Scan(&l.ID, &l.RepositoryID, &l.Name, &l.Color, &l.Description)
		if err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}

func (r *IssueRepository) SyncLabels(ctx context.Context, issueID uuid.UUID, labels []*models.Label) error {
	// First, ensure labels exist in the labels table for the repository
	for _, l := range labels {
		query := `
			INSERT INTO labels (repository_id, name, color, description)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (repository_id, name) DO UPDATE SET
				color = EXCLUDED.color,
				description = EXCLUDED.description
			RETURNING id
		`
		err := r.db.QueryRow(ctx, query, l.RepositoryID, l.Name, l.Color, l.Description).Scan(&l.ID)
		if err != nil {
			return err
		}
	}

	// Then, clear existing associations and re-add
	_, err := r.db.Exec(ctx, "DELETE FROM issue_labels WHERE issue_id = $1", issueID)
	if err != nil {
		return err
	}

	for _, l := range labels {
		_, err := r.db.Exec(ctx, "INSERT INTO issue_labels (issue_id, label_id) VALUES ($1, $2)", issueID, l.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

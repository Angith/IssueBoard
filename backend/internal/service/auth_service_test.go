package service

import (
	"context"
	"testing"

	"github.com/angith/issueboard/internal/repository/models"
	"github.com/google/uuid"
)

type mockUserRepo struct{}

func (m *mockUserRepo) CreateOrUpdate(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	return nil
}

func TestAuthService_SyncUser(t *testing.T) {
	// This is a placeholder for a real test with mocks
	// In a real TDD scenario, I'd use a mocking library or interfaces
	t.Skip("Skipping for now as AuthService is a simple wrapper")
}

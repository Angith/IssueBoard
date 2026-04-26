package service

import (
	"context"

	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/repository/models"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) SyncUser(ctx context.Context, user *models.User) error {
	return s.userRepo.CreateOrUpdate(ctx, user)
}

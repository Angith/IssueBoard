package service

import (
	"context"

	"github.com/angith/issueboard/internal/github"
	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/repository/models"
	"github.com/google/uuid"
)

type RepoService struct {
	repoRepo      *repository.RepoRepository
	githubService *github.RepoService
}

func NewRepoService(repoRepo *repository.RepoRepository, githubService *github.RepoService) *RepoService {
	return &RepoService{
		repoRepo:      repoRepo,
		githubService: githubService,
	}
}

func (s *RepoService) AddRepository(ctx context.Context, userID uuid.UUID, url string) (*models.Repository, error) {
	gRepo, err := s.githubService.GetRepository(ctx, url)
	if err != nil {
		return nil, err
	}

	repo := &models.Repository{
		GitHubRepoID: gRepo.GetID(),
		FullName:     gRepo.GetFullName(),
		Owner:        gRepo.GetOwner().GetLogin(),
		Name:         gRepo.GetName(),
		URL:          gRepo.GetHTMLURL(),
	}

	if err := s.repoRepo.AddToUser(ctx, userID, repo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *RepoService) ListRepositories(ctx context.Context, userID uuid.UUID) ([]*models.Repository, error) {
	return s.repoRepo.ListByUserID(ctx, userID)
}

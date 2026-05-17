package service

import (
	"context"

	"github.com/angith/issueboard/internal/github"
	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/repository/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	log := logrus.WithFields(logrus.Fields{"user_id": userID, "url": url})
	log.Debug("Fetching repository metadata from GitHub")

	gRepo, err := s.githubService.GetRepository(ctx, url)
	if err != nil {
		log.WithError(err).Error("Failed to fetch repository from GitHub")
		return nil, err
	}

	repo := &models.Repository{
		GitHubRepoID: gRepo.GetID(),
		FullName:     gRepo.GetFullName(),
		Owner:        gRepo.GetOwner().GetLogin(),
		Name:         gRepo.GetName(),
		URL:          gRepo.GetHTMLURL(),
	}

	log = log.WithField("repo", repo.FullName)
	log.Debug("Persisting repository to database")

	if err := s.repoRepo.AddToUser(ctx, userID, repo); err != nil {
		log.WithError(err).Error("Failed to add repository to user")
		return nil, err
	}

	log.Info("Repository added successfully")
	return repo, nil
}

func (s *RepoService) ListRepositories(ctx context.Context, userID uuid.UUID) ([]*models.Repository, error) {
	log := logrus.WithField("user_id", userID)
	log.Debug("Listing repositories for user")

	repos, err := s.repoRepo.ListByUserID(ctx, userID)
	if err != nil {
		log.WithError(err).Error("Failed to list repositories")
		return nil, err
	}

	log.WithField("count", len(repos)).Debug("Repositories listed successfully")
	return repos, nil
}

func (s *RepoService) RemoveRepository(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) error {
	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID})
	log.Debug("Removing repository for user")

	if err := s.repoRepo.RemoveFromUser(ctx, userID, repoID); err != nil {
		log.WithError(err).Error("Failed to remove repository")
		return err
	}

	log.Info("Repository removed successfully")
	return nil
}

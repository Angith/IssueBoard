package service

import (
	"context"
	"fmt"

	"github.com/angith/issueboard/internal/github"
	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/repository/models"
	"github.com/google/uuid"
	gogithub "github.com/google/go-github/v60/github"
	"github.com/sirupsen/logrus"
)

type IssueService struct {
	issueRepo     *repository.IssueRepository
	repoRepo      *repository.RepoRepository
	githubService *github.IssueService
}

func NewIssueService(issueRepo *repository.IssueRepository, repoRepo *repository.RepoRepository, githubService *github.IssueService) *IssueService {
	return &IssueService{
		issueRepo:     issueRepo,
		repoRepo:      repoRepo,
		githubService: githubService,
	}
}

type IssueCategory struct {
	Label  *models.Label   `json:"label"`
	Issues []*models.Issue `json:"issues"`
}

type IssueBoard struct {
	Repository           string           `json:"repository"`
	IsTrackingConfigured bool             `json:"is_tracking_configured"`
	Categories           []*IssueCategory `json:"categories"`
}

func (s *IssueService) GetAvailableLabels(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) ([]*gogithub.Label, error) {
	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID})
	log.Debug("Fetching available labels from GitHub")

	repo, err := s.repoRepo.GetByID(ctx, userID, repoID)
	if err != nil {
		log.WithError(err).Warn("Repository not found or access denied when fetching labels")
		return nil, fmt.Errorf("repository not found or access denied")
	}

	labels, err := s.githubService.GetAvailableLabels(ctx, repo.Owner, repo.Name)
	if err != nil {
		log.WithError(err).Error("Failed to fetch labels from GitHub")
		return nil, err
	}

	log.WithField("count", len(labels)).Debug("Available labels fetched successfully")
	return labels, nil
}

func (s *IssueService) GetTrackedLabels(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) ([]string, error) {
	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID})
	log.Debug("Fetching tracked labels")

	labels, err := s.repoRepo.GetTrackedLabels(ctx, userID, repoID)
	if err != nil {
		log.WithError(err).Error("Failed to get tracked labels")
		return nil, err
	}

	log.WithField("count", len(labels)).Debug("Tracked labels retrieved")
	return labels, nil
}

func (s *IssueService) UpdateTrackedLabels(ctx context.Context, userID uuid.UUID, repoID uuid.UUID, labels []string) error {
	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID, "label_count": len(labels)})
	log.Debug("Updating tracked labels")

	if err := s.repoRepo.UpdateTrackedLabels(ctx, userID, repoID, labels); err != nil {
		log.WithError(err).Error("Failed to update tracked labels")
		return err
	}

	log.Info("Tracked labels updated successfully")
	return nil
}

func (s *IssueService) GetCategorizedIssues(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) (*IssueBoard, error) {
	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID})
	log.Debug("Getting categorized issues")

	// Verify user has access to this repo
	repo, err := s.repoRepo.GetByID(ctx, userID, repoID)
	if err != nil {
		log.WithError(err).Warn("Repository not found or access denied")
		return nil, fmt.Errorf("repository not found or access denied")
	}

	log = log.WithField("repo", repo.FullName)

	trackedLabels, err := s.GetTrackedLabels(ctx, userID, repoID)
	if err != nil {
		return nil, err
	}

	if len(trackedLabels) == 0 {
		log.Info("No tracked labels configured; returning empty board")
		return &IssueBoard{
			Repository:           repo.FullName,
			IsTrackingConfigured: false,
			Categories:           []*IssueCategory{},
		}, nil
	}

	issues, err := s.issueRepo.ListByRepoID(ctx, repoID)
	if err != nil {
		log.WithError(err).Error("Failed to list issues from database")
		return nil, err
	}

	// Filter issues by tracked labels
	trackedMap := make(map[string]bool)
	for _, l := range trackedLabels {
		trackedMap[l] = true
	}

	var filteredIssues []*models.Issue
	for _, issue := range issues {
		hasTrackedLabel := false
		for _, l := range issue.Labels {
			if trackedMap[l.Name] {
				hasTrackedLabel = true
				break
			}
		}
		if hasTrackedLabel {
			filteredIssues = append(filteredIssues, issue)
		}
	}

	// If no issues cached for these labels, trigger a refresh
	if len(filteredIssues) == 0 {
		log.Warn("No cached issues found for tracked labels; triggering auto-refresh from GitHub")
		if err := s.RefreshIssues(ctx, userID, repoID); err != nil {
			return nil, err
		}

		issues, err = s.issueRepo.ListByRepoID(ctx, repoID)
		if err != nil {
			log.WithError(err).Error("Failed to list issues after refresh")
			return nil, err
		}

		filteredIssues = nil
		for _, issue := range issues {
			hasTrackedLabel := false
			for _, l := range issue.Labels {
				if trackedMap[l.Name] {
					hasTrackedLabel = true
					break
				}
			}
			if hasTrackedLabel {
				filteredIssues = append(filteredIssues, issue)
			}
		}
	}

	board := s.groupIssues(filteredIssues, trackedLabels)
	board.Repository = repo.FullName
	board.IsTrackingConfigured = true
	log.WithField("issue_count", len(filteredIssues)).Debug("Returning categorized issue board")
	return board, nil
}

func (s *IssueService) RefreshIssues(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) error {
	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID})
	log.Debug("Starting issue refresh from GitHub")

	// 1. Get repository info from DB and verify access
	repo, err := s.repoRepo.GetByID(ctx, userID, repoID)
	if err != nil {
		log.WithError(err).Warn("Repository not found or access denied during refresh")
		return fmt.Errorf("repository not found or access denied")
	}

	log = log.WithField("repo", repo.FullName)

	// 2. Get tracked labels
	trackedLabels, err := s.GetTrackedLabels(ctx, userID, repoID)
	if err != nil {
		return err
	}

	if len(trackedLabels) == 0 {
		// Nothing to sync
		log.Debug("No tracked labels configured; skipping refresh")
		return nil
	}

	// 3. Fetch from GitHub by labels
	gIssues, err := s.githubService.GetIssuesByLabels(ctx, repo.Owner, repo.Name, trackedLabels)
	if err != nil {
		log.WithError(err).Error("Failed to fetch issues from GitHub")
		return err
	}

	log.WithField("fetched_count", len(gIssues)).Debug("Issues fetched from GitHub; syncing to database")

	// 4. Sync to DB
	for _, gi := range gIssues {
		issue := &models.Issue{
			RepositoryID:  repoID,
			GitHubIssueID: gi.GetID(),
			Number:        gi.GetNumber(),
			Title:         gi.GetTitle(),
			Body:          gi.GetBody(),
			State:         gi.GetState(),
			URL:           gi.GetHTMLURL(),
			UpdatedAt:     gi.GetUpdatedAt().Time,
		}

		if err := s.issueRepo.CreateOrUpdate(ctx, issue); err != nil {
			log.WithError(err).WithField("github_issue_id", gi.GetID()).Error("Failed to upsert issue")
			return err
		}

		// Sync labels
		var labels []*models.Label
		for _, gl := range gi.Labels {
			labels = append(labels, &models.Label{
				RepositoryID: repoID,
				Name:         gl.GetName(),
				Color:        gl.GetColor(),
				Description:  gl.GetDescription(),
			})
		}
		if err := s.issueRepo.SyncLabels(ctx, issue.ID, labels); err != nil {
			log.WithError(err).WithField("github_issue_id", gi.GetID()).Error("Failed to sync labels for issue")
			return err
		}
	}

	log.WithField("synced_count", len(gIssues)).Info("Issues refreshed and synced successfully")
	return nil
}

func (s *IssueService) groupIssues(issues []*models.Issue, trackedLabels []string) *IssueBoard {
	categoryMap := make(map[string]*IssueCategory)
	
	// Initialize empty categories for tracked labels so they always show up
	for _, labelName := range trackedLabels {
		categoryMap[labelName] = &IssueCategory{
			Label: &models.Label{Name: labelName, Color: "dddddd"}, // Default color, will be overwritten if issues exist
			Issues: []*models.Issue{},
		}
	}

	for _, issue := range issues {
		for _, label := range issue.Labels {
			if _, exists := categoryMap[label.Name]; exists {
				categoryMap[label.Name].Label = label // Update with true color/description
				categoryMap[label.Name].Issues = append(categoryMap[label.Name].Issues, issue)
			}
		}
	}

	var categories []*IssueCategory
	for _, cat := range categoryMap {
		categories = append(categories, cat)
	}

	return &IssueBoard{
		Categories: categories,
	}
}

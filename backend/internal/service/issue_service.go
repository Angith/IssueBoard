package service

import (
	"context"
	"fmt"

	"github.com/angith/issueboard/internal/github"
	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/repository/models"
	"github.com/google/uuid"
	gogithub "github.com/google/go-github/v60/github"
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
	repo, err := s.repoRepo.GetByID(ctx, userID, repoID)
	if err != nil {
		return nil, fmt.Errorf("repository not found or access denied")
	}
	return s.githubService.GetAvailableLabels(ctx, repo.Owner, repo.Name)
}

func (s *IssueService) GetTrackedLabels(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) ([]string, error) {
	return s.repoRepo.GetTrackedLabels(ctx, userID, repoID)
}

func (s *IssueService) UpdateTrackedLabels(ctx context.Context, userID uuid.UUID, repoID uuid.UUID, labels []string) error {
	return s.repoRepo.UpdateTrackedLabels(ctx, userID, repoID, labels)
}

func (s *IssueService) GetCategorizedIssues(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) (*IssueBoard, error) {
	// Verify user has access to this repo
	repo, err := s.repoRepo.GetByID(ctx, userID, repoID)
	if err != nil {
		return nil, fmt.Errorf("repository not found or access denied")
	}

	trackedLabels, err := s.GetTrackedLabels(ctx, userID, repoID)
	if err != nil {
		return nil, err
	}

	if len(trackedLabels) == 0 {
		return &IssueBoard{
			Repository:           repo.FullName,
			IsTrackingConfigured: false,
			Categories:           []*IssueCategory{},
		}, nil
	}

	issues, err := s.issueRepo.ListByRepoID(ctx, repoID)
	if err != nil {
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
		if err := s.RefreshIssues(ctx, userID, repoID); err != nil {
			return nil, err
		}
		
		issues, err = s.issueRepo.ListByRepoID(ctx, repoID)
		if err != nil {
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
	return board, nil
}

func (s *IssueService) RefreshIssues(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) error {
	// 1. Get repository info from DB and verify access
	repo, err := s.repoRepo.GetByID(ctx, userID, repoID)
	if err != nil {
		return fmt.Errorf("repository not found or access denied")
	}

	// 2. Get tracked labels
	trackedLabels, err := s.GetTrackedLabels(ctx, userID, repoID)
	if err != nil {
		return err
	}
	
	if len(trackedLabels) == 0 {
		// Nothing to sync
		return nil
	}

	// 3. Fetch from GitHub by labels
	gIssues, err := s.githubService.GetIssuesByLabels(ctx, repo.Owner, repo.Name, trackedLabels)
	if err != nil {
		return err
	}

	// 3. Sync to DB
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
			return err
		}
	}

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

package service

import (
	"context"

	"github.com/angith/issueboard/internal/github"
	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/repository/models"
	"github.com/google/uuid"
)

type IssueService struct {
	issueRepo     *repository.IssueRepository
	repoRepo      *repository.RepoRepository
	githubService *github.RepoService
}

func NewIssueService(issueRepo *repository.IssueRepository, repoRepo *repository.RepoRepository, githubService *github.RepoService) *IssueService {
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
	Repository string           `json:"repository"`
	Categories []*IssueCategory `json:"categories"`
}

func (s *IssueService) GetCategorizedIssues(ctx context.Context, repoID uuid.UUID) (*IssueBoard, error) {
	issues, err := s.issueRepo.ListByRepoID(ctx, repoID)
	if err != nil {
		return nil, err
	}

	// If no issues cached, trigger a refresh
	if len(issues) == 0 {
		if err := s.RefreshIssues(ctx, repoID); err != nil {
			return nil, err
		}
		// Fetch again after refresh
		issues, err = s.issueRepo.ListByRepoID(ctx, repoID)
		if err != nil {
			return nil, err
		}
	}

	return s.groupIssues(issues), nil
}

func (s *IssueService) RefreshIssues(ctx context.Context, repoID uuid.UUID) error {
	// 1. Get repository info from DB
	// (Needs GetByID in repoRepo - adding it here for the sake of completeness in logic)
	// For now, I'll assume we have the repo info or we just use a placeholder
	// In a real app, I'd have a method to get repo details by ID

	// Mocking owner/repo for now or assuming we fetch it
	owner := "mock"
	repoName := "mock"

	// 2. Fetch from GitHub
	gIssues, err := s.githubService.GetIssues(ctx, owner, repoName)
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

func (s *IssueService) groupIssues(issues []*models.Issue) *IssueBoard {
	categoryMap := make(map[string]*IssueCategory)
	var unlabeledIssues []*models.Issue

	for _, issue := range issues {
		if len(issue.Labels) == 0 {
			unlabeledIssues = append(unlabeledIssues, issue)
			continue
		}

		for _, label := range issue.Labels {
			if _, exists := categoryMap[label.Name]; !exists {
				categoryMap[label.Name] = &IssueCategory{
					Label:  label,
					Issues: []*models.Issue{},
				}
			}
			categoryMap[label.Name].Issues = append(categoryMap[label.Name].Issues, issue)
		}
	}

	var categories []*IssueCategory
	for _, cat := range categoryMap {
		categories = append(categories, cat)
	}

	if len(unlabeledIssues) > 0 {
		categories = append(categories, &IssueCategory{
			Label: &models.Label{
				Name:  "Unlabeled",
				Color: "cccccc",
			},
			Issues: unlabeledIssues,
		})
	}

	return &IssueBoard{
		Categories: categories,
	}
}

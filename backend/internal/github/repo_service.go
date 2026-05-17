package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v60/github"
)

type RepoService struct {
	client *Client
}

func NewRepoService(client *Client) *RepoService {
	return &RepoService{client: client}
}

func (s *RepoService) GetRepository(ctx context.Context, url string) (*github.Repository, error) {
	owner, repo, err := parseGitHubURL(url)
	if err != nil {
		return nil, err
	}

	gRepo, _, err := s.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository from GitHub: %w", err)
	}

	return gRepo, nil
}

func parseGitHubURL(url string) (owner, repo string, err error) {
	// Simple parser for https://github.com/owner/repo
	trimmed := strings.TrimPrefix(url, "https://github.com/")
	trimmed = strings.TrimSuffix(trimmed, ".git")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format")
	}
	return parts[0], parts[1], nil
}

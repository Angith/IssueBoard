package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
)

func (s *RepoService) GetIssues(ctx context.Context, owner, repo string) ([]*github.Issue, error) {
	opts := &github.IssueListByRepoOptions{
		State: "open",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allIssues []*github.Issue
	for {
		issues, resp, err := s.client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allIssues, nil
}

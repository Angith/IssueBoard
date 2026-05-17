package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
)

type IssueService struct {
	client *Client
}

func NewIssueService(client *Client) *IssueService {
	return &IssueService{client: client}
}

func (s *IssueService) GetAvailableLabels(ctx context.Context, owner, repo string) ([]*github.Label, error) {
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allLabels []*github.Label
	for {
		labels, resp, err := s.client.Issues.ListLabels(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list labels: %w", err)
		}
		allLabels = append(allLabels, labels...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allLabels, nil
}

func (s *IssueService) GetIssuesByLabels(ctx context.Context, owner, repo string, labels []string) ([]*github.Issue, error) {
	var allIssues []*github.Issue
	seen := make(map[int64]bool)

	for _, label := range labels {
		opts := &github.IssueListByRepoOptions{
			State:  "open",
			Labels: []string{label},
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}

		for {
			issues, resp, err := s.client.Issues.ListByRepo(ctx, owner, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list issues for label %s: %w", label, err)
			}
			
			for _, issue := range issues {
				if !seen[issue.GetID()] {
					allIssues = append(allIssues, issue)
					seen[issue.GetID()] = true
				}
			}

			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
	}

	return allIssues, nil
}

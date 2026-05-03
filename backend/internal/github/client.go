package github

import (
	"context"
	"net/http"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type Client struct {
	*github.Client
}

func NewClient(accessToken string) *Client {
	var httpClient *http.Client
	if accessToken != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	}

	return &Client{
		Client: github.NewClient(httpClient),
	}
}

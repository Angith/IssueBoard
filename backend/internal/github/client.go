package github

import (
	"net/http"

	"github.com/angith/issueboard/internal/api/middleware"
	"github.com/google/go-github/v60/github"
)

type Client struct {
	*github.Client
	fallbackToken string
}

// ctxRoundTripper gets the token from the request context
type ctxRoundTripper struct {
	fallbackToken string
	base          http.RoundTripper
}

func (rt *ctxRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token := middleware.GetGitHubToken(req.Context())
	if token == "" {
		token = rt.fallbackToken
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return rt.base.RoundTrip(req)
}

func NewClient(fallbackToken string) *Client {
	httpClient := &http.Client{
		Transport: &ctxRoundTripper{
			fallbackToken: fallbackToken,
			base:          http.DefaultTransport,
		},
	}

	return &Client{
		Client:        github.NewClient(httpClient),
		fallbackToken: fallbackToken,
	}
}

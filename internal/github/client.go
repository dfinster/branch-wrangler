package github

import (
	"context"
	"time"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	auth   *AuthConfig
	owner  string
	repo   string
}

type PullRequest struct {
	Number int
	Title  string
	State  string
	Draft  bool
	Merged bool
	URL    string
}

func NewClient(owner, repo string) (*Client, error) {
	auth := NewAuthConfig()
	token, err := auth.GetToken()
	if err != nil {
		return nil, err
	}

	if err := auth.ValidateToken(token); err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
		auth:   auth,
		owner:  owner,
		repo:   repo,
	}, nil
}

func (c *Client) GetPullRequestsForBranch(ctx context.Context, branch string) ([]PullRequest, error) {
	opts := &github.PullRequestListOptions{
		Head:        c.owner + ":" + branch,
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allPRs []PullRequest

	for {
		prs, resp, err := c.client.PullRequests.List(ctx, c.owner, c.repo, opts)
		if err != nil {
			return nil, err
		}

		for _, pr := range prs {
			allPRs = append(allPRs, PullRequest{
				Number: pr.GetNumber(),
				Title:  pr.GetTitle(),
				State:  pr.GetState(),
				Draft:  pr.GetDraft(),
				Merged: pr.GetMerged(),
				URL:    pr.GetHTMLURL(),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allPRs, nil
}

func (c *Client) GetRateLimit(ctx context.Context) (*github.RateLimits, error) {
	limits, _, err := c.client.RateLimits(ctx)
	return limits, err
}

func (c *Client) BranchExists(ctx context.Context, branch string) (bool, error) {
	_, _, err := c.client.Repositories.GetBranch(ctx, c.owner, c.repo, branch, 1)
	if err != nil {
		if ghErr, ok := err.(*github.ErrorResponse); ok && ghErr.Response.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type CachedClient struct {
	client *Client
	cache  map[string]cacheEntry
}

type cacheEntry struct {
	data      interface{}
	timestamp time.Time
	ttl       time.Duration
}

func NewCachedClient(owner, repo string) (*CachedClient, error) {
	client, err := NewClient(owner, repo)
	if err != nil {
		return nil, err
	}

	return &CachedClient{
		client: client,
		cache:  make(map[string]cacheEntry),
	}, nil
}

func (c *CachedClient) GetPullRequestsForBranch(ctx context.Context, branch string) ([]PullRequest, error) {
	cacheKey := "pr:" + branch

	if entry, exists := c.cache[cacheKey]; exists {
		if time.Since(entry.timestamp) < entry.ttl {
			return entry.data.([]PullRequest), nil
		}
	}

	prs, err := c.client.GetPullRequestsForBranch(ctx, branch)
	if err != nil {
		return nil, err
	}

	c.cache[cacheKey] = cacheEntry{
		data:      prs,
		timestamp: time.Now(),
		ttl:       15 * time.Minute,
	}

	return prs, nil
}

func (c *CachedClient) BranchExists(ctx context.Context, branch string) (bool, error) {
	cacheKey := "branch:" + branch

	if entry, exists := c.cache[cacheKey]; exists {
		if time.Since(entry.timestamp) < entry.ttl {
			return entry.data.(bool), nil
		}
	}

	exists, err := c.client.BranchExists(ctx, branch)
	if err != nil {
		return false, err
	}

	c.cache[cacheKey] = cacheEntry{
		data:      exists,
		timestamp: time.Now(),
		ttl:       15 * time.Minute,
	}

	return exists, nil
}
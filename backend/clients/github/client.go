package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// ErrNoToken is returned by FetchPinnedRepos when no token is configured —
// the GraphQL API (the only way to read pinned repos) requires auth even
// for public data, unlike the REST endpoints this client otherwise uses.
var ErrNoToken = errors.New("github: no token configured, cannot query pinned repos via GraphQL")

const pinnedItemsQuery = `
query($login: String!) {
  user(login: $login) {
    pinnedItems(first: 6, types: REPOSITORY) {
      nodes {
        ... on Repository {
          name
          description
          url
          stargazerCount
          forkCount
          primaryLanguage { name }
        }
      }
    }
  }
}`

const maxPages = 3

var nextLinkPattern = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)

type Client struct {
	baseURL    string
	username   string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, username, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		username:   username,
		token:      token,
		httpClient: &http.Client{},
	}
}

// FetchRepos returns the user's public repositories, following pagination
// up to maxPages (300 repos) as a safety cap.
func (c *Client) FetchRepos(ctx context.Context) ([]Repo, error) {
	url := fmt.Sprintf("%s/users/%s/repos?per_page=100&sort=updated", c.baseURL, c.username)

	var all []Repo
	for page := 0; page < maxPages && url != ""; page++ {
		repos, next, err := c.fetchPage(ctx, url)
		if err != nil {
			return nil, err
		}
		all = append(all, repos...)
		url = next
	}
	return all, nil
}

func (c *Client) fetchPage(ctx context.Context, url string) ([]Repo, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	// GitHub's REST API 403s any request without a User-Agent.
	req.Header.Set("User-Agent", "portfolio-backend")
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining == "0" {
			log.Printf("WARN: github rate limited, resets at %s", resp.Header.Get("X-RateLimit-Reset"))
		}
		return nil, "", fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, "", fmt.Errorf("decoding github response: %w", err)
	}

	return repos, parseNextLink(resp.Header.Get("Link")), nil
}

// FetchPinnedRepos returns the user's pinned repositories via the GitHub
// GraphQL API. Requires a token — see ErrNoToken.
func (c *Client) FetchPinnedRepos(ctx context.Context) ([]Repo, error) {
	if c.token == "" {
		return nil, ErrNoToken
	}

	body, err := json.Marshal(map[string]any{
		"query":     pinnedItemsQuery,
		"variables": map[string]string{"login": c.username},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "portfolio-backend")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github graphql API returned status %d", resp.StatusCode)
	}

	var parsed graphQLPinnedResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decoding github graphql response: %w", err)
	}
	if len(parsed.Errors) > 0 {
		return nil, fmt.Errorf("github graphql error: %s", parsed.Errors[0].Message)
	}

	nodes := parsed.Data.User.PinnedItems.Nodes
	repos := make([]Repo, 0, len(nodes))
	for _, n := range nodes {
		repos = append(repos, Repo{
			Name:            n.Name,
			HTMLURL:         n.URL,
			Description:     n.Description,
			StargazersCount: n.StargazerCount,
			ForksCount:      n.ForkCount,
			Language:        n.PrimaryLanguage.Name,
		})
	}
	return repos, nil
}

func parseNextLink(linkHeader string) string {
	matches := nextLinkPattern.FindStringSubmatch(linkHeader)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

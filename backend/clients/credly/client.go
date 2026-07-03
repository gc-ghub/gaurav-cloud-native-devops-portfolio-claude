package credly

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrUpstreamUnavailable wraps any failure talking to Credly's unofficial
// JSON endpoint — non-200 status, unexpected content type, or bad JSON all
// fold into this so callers have one thing to check, never a panic.
var ErrUpstreamUnavailable = errors.New("credly upstream unavailable")

type Client struct {
	baseURL    string
	userID     string
	httpClient *http.Client
}

func NewClient(baseURL, userID string) *Client {
	return &Client{
		baseURL:    baseURL,
		userID:     userID,
		httpClient: &http.Client{},
	}
}

func (c *Client) FetchBadges(ctx context.Context) ([]Badge, error) {
	url := fmt.Sprintf("%s/users/%s/badges.json", c.baseURL, c.userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "portfolio-backend")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUpstreamUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrUpstreamUnavailable, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("%w: unexpected content-type %q", ErrUpstreamUnavailable, contentType)
	}

	var parsed badgesResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("%w: decoding response: %v", ErrUpstreamUnavailable, err)
	}

	return parsed.Data, nil
}

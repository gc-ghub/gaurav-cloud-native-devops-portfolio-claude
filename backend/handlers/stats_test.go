package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"portfolio-backend/cache"
	"portfolio-backend/clients/github"
	"portfolio-backend/models"
)

func TestGetStats_AggregatesExcludingForks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fetcher := &fakeReposFetcher{repos: []github.Repo{
		{Name: "owned-1", StargazersCount: 3, ForksCount: 1},
		{Name: "owned-2", StargazersCount: 5, ForksCount: 2},
		{Name: "a-fork", Fork: true, StargazersCount: 1000, ForksCount: 1000},
	}}

	h := NewStatsHandler(fetcher, cache.NewCache[[]github.Repo](time.Minute))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/github-stats", nil)

	h.GetStats(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var stats models.GithubStats
	if err := json.Unmarshal(w.Body.Bytes(), &stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if stats.PublicRepos != 2 {
		t.Fatalf("expected 2 public repos (fork excluded), got %d", stats.PublicRepos)
	}
	if stats.TotalStars != 8 {
		t.Fatalf("expected 8 total stars, got %d", stats.TotalStars)
	}
	if stats.TotalForks != 3 {
		t.Fatalf("expected 3 total forks, got %d", stats.TotalForks)
	}
}

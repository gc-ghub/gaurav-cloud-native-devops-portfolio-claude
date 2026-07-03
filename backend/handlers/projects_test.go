package handlers

import (
	"context"
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

type fakeReposFetcher struct {
	repos []github.Repo
	err   error
}

func (f *fakeReposFetcher) FetchRepos(ctx context.Context) ([]github.Repo, error) {
	return f.repos, f.err
}

type fakePinnedFetcher struct {
	repos []github.Repo
	err   error
}

func (f *fakePinnedFetcher) FetchPinnedRepos(ctx context.Context) ([]github.Repo, error) {
	return f.repos, f.err
}

func newTestProjectsHandler(rf reposFetcher, pf pinnedFetcher, max int) *ProjectsHandler {
	return NewProjectsHandler(
		rf, pf,
		cache.NewCache[[]github.Repo](time.Minute),
		cache.NewCache[[]github.Repo](time.Minute),
		max,
	)
}

func decodeProjects(t *testing.T, w *httptest.ResponseRecorder) []models.Project {
	t.Helper()
	var body struct {
		Projects []models.Project `json:"projects"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return body.Projects
}

func TestGetProjects_AlwaysLeadsWithFeaturedProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repos := &fakeReposFetcher{repos: []github.Repo{
		{Name: featuredRepoName, StargazersCount: 7, ForksCount: 2},
	}}
	pinned := &fakePinnedFetcher{err: assertErr} // forces fallback

	h := newTestProjectsHandler(repos, pinned, 6)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	h.GetProjects(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	projects := decodeProjects(t, w)
	if len(projects) == 0 || projects[0].Name != featuredProjectTemplate.Name {
		t.Fatalf("expected the featured project first, got %+v", projects)
	}
	if projects[0].Stars != 7 || projects[0].Forks != 2 {
		t.Fatalf("expected featured project's live stars/forks to be looked up, got stars=%d forks=%d", projects[0].Stars, projects[0].Forks)
	}
}

func TestGetProjects_UsesPinnedReposWhenAvailable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repos := &fakeReposFetcher{repos: []github.Repo{
		{Name: featuredRepoName},
	}}
	pinned := &fakePinnedFetcher{repos: []github.Repo{
		{Name: "pinned-one", HTMLURL: "https://github.com/u/pinned-one", StargazersCount: 3},
	}}

	h := newTestProjectsHandler(repos, pinned, 6)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	h.GetProjects(c)

	projects := decodeProjects(t, w)
	if len(projects) != 2 || projects[1].Name != "pinned-one" {
		t.Fatalf("expected featured + pinned-one, got %+v", projects)
	}
}

func TestGetProjects_FallsBackToTopStarredWithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repos := &fakeReposFetcher{repos: []github.Repo{
		{Name: featuredRepoName},
		{Name: "a-fork", Fork: true, StargazersCount: 100},
		{Name: "archived-repo", Archived: true, StargazersCount: 50},
		{Name: "owned-high-stars", StargazersCount: 10},
		{Name: "owned-low-stars", StargazersCount: 1},
	}}
	pinned := &fakePinnedFetcher{err: github.ErrNoToken}

	h := newTestProjectsHandler(repos, pinned, 3)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	h.GetProjects(c)

	projects := decodeProjects(t, w)
	if len(projects) != 3 {
		t.Fatalf("expected 3 projects (max), got %d", len(projects))
	}
	if projects[0].Name != featuredProjectTemplate.Name {
		t.Fatalf("expected featured project first, got %q", projects[0].Name)
	}
	if projects[1].Name != "owned-high-stars" || projects[2].Name != "owned-low-stars" {
		t.Fatalf("expected fork/archived excluded and remaining sorted by stars, got %+v", projects)
	}
}

func TestGetProjects_UpstreamFailureReturns502(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repos := &fakeReposFetcher{err: assertErr}
	pinned := &fakePinnedFetcher{}
	h := newTestProjectsHandler(repos, pinned, 6)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/projects", nil)

	h.GetProjects(c)

	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", w.Code)
	}
}

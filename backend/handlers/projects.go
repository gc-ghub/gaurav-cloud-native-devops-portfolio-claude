package handlers

import (
	"context"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"portfolio-backend/cache"
	"portfolio-backend/clients/github"
	"portfolio-backend/models"
)

type reposFetcher interface {
	FetchRepos(ctx context.Context) ([]github.Repo, error)
}

type pinnedFetcher interface {
	FetchPinnedRepos(ctx context.Context) ([]github.Repo, error)
}

// featuredRepoName is manually featured (not pulled from pinned/starred
// data) — content sourced from pics/bio.txt. Stars/forks are still looked
// up live from the general repo list so the card isn't frozen at 0/0.
const featuredRepoName = "project-gc-industries-devops-superheros"

var featuredProjectTemplate = models.Project{
	Name:        "GC Industries – Cloud-Native DevOps Superheroes Platform",
	Description: "A production-style cloud-native microservices platform built on Kubernetes implementing Platform Engineering, SRE, GitOps, Security, and Observability best practices.",
	Tech:        []string{"Kubernetes", "Argo CD", "Istio", "Prometheus", "Kyverno"},
	Href:        "https://github.com/gc-ghub/" + featuredRepoName,
}

type ProjectsHandler struct {
	reposFetcher  reposFetcher
	pinnedFetcher pinnedFetcher
	reposCache    *cache.Cache[[]github.Repo]
	pinnedCache   *cache.Cache[[]github.Repo]
	maxProjects   int
}

func NewProjectsHandler(rf reposFetcher, pf pinnedFetcher, reposCache, pinnedCache *cache.Cache[[]github.Repo], max int) *ProjectsHandler {
	return &ProjectsHandler{
		reposFetcher:  rf,
		pinnedFetcher: pf,
		reposCache:    reposCache,
		pinnedCache:   pinnedCache,
		maxProjects:   max,
	}
}

func (h *ProjectsHandler) GetProjects(c *gin.Context) {
	allRepos, err := h.reposCache.GetOrFetch(func() ([]github.Repo, error) {
		return h.reposFetcher.FetchRepos(c.Request.Context())
	})
	if err != nil {
		respondError(c, http.StatusBadGateway, "upstream_unavailable", "unable to fetch projects from GitHub")
		return
	}

	featured := buildFeaturedProject(allRepos)
	others := excludeByName(h.fetchOtherRepos(c.Request.Context(), allRepos), featuredRepoName)

	projects := []models.Project{featured}
	for _, r := range others {
		if len(projects) >= h.maxProjects {
			break
		}
		projects = append(projects, toModelProject(r))
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// fetchOtherRepos returns pinned repos when available, falling back to the
// top-starred non-fork/non-archived repos (e.g. no GITHUB_TOKEN configured,
// or the GraphQL call failed) — same graceful-degradation shape used for
// Credly, so a missing/expired token never breaks this endpoint.
func (h *ProjectsHandler) fetchOtherRepos(ctx context.Context, allRepos []github.Repo) []github.Repo {
	pinned, err := h.pinnedCache.GetOrFetch(func() ([]github.Repo, error) {
		return h.pinnedFetcher.FetchPinnedRepos(ctx)
	})
	if err != nil {
		log.Printf("WARN: pinned repos unavailable (%v), falling back to top-starred repos", err)
		return filterAndSortTopStarred(allRepos)
	}
	return pinned
}

func buildFeaturedProject(allRepos []github.Repo) models.Project {
	featured := featuredProjectTemplate
	if r, ok := findRepoByName(allRepos, featuredRepoName); ok {
		featured.Stars = r.StargazersCount
		featured.Forks = r.ForksCount
	}
	return featured
}

func findRepoByName(repos []github.Repo, name string) (github.Repo, bool) {
	for _, r := range repos {
		if r.Name == name {
			return r, true
		}
	}
	return github.Repo{}, false
}

func excludeByName(repos []github.Repo, name string) []github.Repo {
	out := make([]github.Repo, 0, len(repos))
	for _, r := range repos {
		if r.Name != name {
			out = append(out, r)
		}
	}
	return out
}

func filterAndSortTopStarred(repos []github.Repo) []github.Repo {
	filtered := make([]github.Repo, 0, len(repos))
	for _, r := range repos {
		if r.Fork || r.Archived {
			continue
		}
		filtered = append(filtered, r)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].StargazersCount > filtered[j].StargazersCount
	})
	return filtered
}

func toModelProject(r github.Repo) models.Project {
	tech := r.Topics
	if len(tech) == 0 && r.Language != "" {
		tech = []string{r.Language}
	}
	return models.Project{
		Name:        r.Name,
		Description: r.Description,
		Tech:        tech,
		Href:        r.HTMLURL,
		Stars:       r.StargazersCount,
		Forks:       r.ForksCount,
	}
}

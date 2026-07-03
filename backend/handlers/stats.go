package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portfolio-backend/cache"
	"portfolio-backend/clients/github"
	"portfolio-backend/models"
)

type StatsHandler struct {
	fetcher reposFetcher
	cache   *cache.Cache[[]github.Repo]
}

func NewStatsHandler(f reposFetcher, c *cache.Cache[[]github.Repo]) *StatsHandler {
	return &StatsHandler{fetcher: f, cache: c}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	repos, err := h.cache.GetOrFetch(func() ([]github.Repo, error) {
		return h.fetcher.FetchRepos(c.Request.Context())
	})
	if err != nil {
		respondError(c, http.StatusBadGateway, "upstream_unavailable", "unable to fetch GitHub stats")
		return
	}

	stats := models.GithubStats{}
	for _, r := range repos {
		if r.Fork {
			continue
		}
		stats.PublicRepos++
		stats.TotalStars += r.StargazersCount
		stats.TotalForks += r.ForksCount
	}

	c.JSON(http.StatusOK, stats)
}

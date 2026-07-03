package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"portfolio-backend/cache"
	"portfolio-backend/clients/credly"
	"portfolio-backend/clients/github"
	"portfolio-backend/config"
	"portfolio-backend/handlers"
	"portfolio-backend/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	githubClient := github.NewClient(cfg.GitHubAPIURL, cfg.GitHubUsername, cfg.GitHubToken)
	credlyClient := credly.NewClient(cfg.CredlyAPIURL, cfg.CredlyUserID)

	// Shared by /api/projects and /api/github-stats so both endpoints draw
	// from the same GitHub call instead of doubling API usage.
	reposCache := cache.NewCache[[]github.Repo](cfg.CacheTTL)
	pinnedCache := cache.NewCache[[]github.Repo](cfg.CacheTTL)
	badgesCache := cache.NewCache[[]credly.Badge](cfg.CacheTTL)

	projectsHandler := handlers.NewProjectsHandler(githubClient, githubClient, reposCache, pinnedCache, cfg.MaxProjects)
	statsHandler := handlers.NewStatsHandler(githubClient, reposCache)
	certificationsHandler := handlers.NewCertificationsHandler(credlyClient, badgesCache)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middleware.CORS(cfg.FrontendURL))

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/api/projects", projectsHandler.GetProjects)
	router.GET("/api/certifications", certificationsHandler.GetCertifications)
	router.GET("/api/github-stats", statsHandler.GetStats)

	log.Printf("listening on :%s (env=%s)", cfg.Port, cfg.Environment)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

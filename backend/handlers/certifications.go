package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"portfolio-backend/cache"
	"portfolio-backend/clients/credly"
	"portfolio-backend/models"
)

type badgesFetcher interface {
	FetchBadges(ctx context.Context) ([]credly.Badge, error)
}

type CertificationsHandler struct {
	fetcher badgesFetcher
	cache   *cache.Cache[[]credly.Badge]
}

func NewCertificationsHandler(f badgesFetcher, c *cache.Cache[[]credly.Badge]) *CertificationsHandler {
	return &CertificationsHandler{fetcher: f, cache: c}
}

// GetCertifications degrades to an empty list + warning rather than a hard
// error: Credly's badges.json endpoint is unofficial (see backend/README.md),
// so a failure here shouldn't look like a broken page to the frontend.
func (h *CertificationsHandler) GetCertifications(c *gin.Context) {
	badges, err := h.cache.GetOrFetch(func() ([]credly.Badge, error) {
		return h.fetcher.FetchBadges(c.Request.Context())
	})
	if err != nil {
		log.Printf("WARN: credly fetch failed: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"certifications": []models.Certification{},
			"warning":        "certifications temporarily unavailable",
		})
		return
	}

	certifications := make([]models.Certification, 0, len(badges))
	for _, b := range badges {
		certifications = append(certifications, models.Certification{
			Name:     b.BadgeTemplate.Name,
			Issuer:   b.IssuerName(),
			ImageURL: b.BadgeTemplate.ImageURL,
			Url:      b.BadgeTemplate.Url,
			IssuedAt: b.IssuedAtDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{"certifications": certifications})
}

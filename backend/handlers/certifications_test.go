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
	"portfolio-backend/clients/credly"
)

type fakeBadgesFetcher struct {
	badges []credly.Badge
	err    error
}

func (f *fakeBadgesFetcher) FetchBadges(ctx context.Context) ([]credly.Badge, error) {
	return f.badges, f.err
}

func TestGetCertifications_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var badge credly.Badge
	badge.BadgeTemplate.Name = "AWS Certified Developer"
	badge.BadgeTemplate.Url = "https://www.credly.com/badge"
	badge.IssuedAtDate = "2025-01-01"

	fetcher := &fakeBadgesFetcher{badges: []credly.Badge{badge}}
	h := NewCertificationsHandler(fetcher, cache.NewCache[[]credly.Badge](time.Minute))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/certifications", nil)

	h.GetCertifications(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, hasWarning := body["warning"]; hasWarning {
		t.Fatal("did not expect a warning field on success")
	}
}

// GetCertifications must degrade gracefully — Credly's endpoint is
// undocumented, so a failure here should never look like a broken page.
func TestGetCertifications_UpstreamFailureDegradesTo200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fetcher := &fakeBadgesFetcher{err: assertErr}
	h := NewCertificationsHandler(fetcher, cache.NewCache[[]credly.Badge](time.Minute))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/certifications", nil)

	h.GetCertifications(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 even on upstream failure, got %d", w.Code)
	}

	var body struct {
		Certifications []any  `json:"certifications"`
		Warning        string `json:"warning"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(body.Certifications) != 0 {
		t.Fatalf("expected empty certifications list, got %d", len(body.Certifications))
	}
	if body.Warning == "" {
		t.Fatal("expected a warning message on upstream failure")
	}
}

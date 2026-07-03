package credly

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchBadges_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"data":[{"issued_at_date":"2025-01-01","badge_template":{"name":"Test Badge"}}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "testuser")
	badges, err := client.FetchBadges(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(badges) != 1 || badges[0].BadgeTemplate.Name != "Test Badge" {
		t.Fatalf("unexpected badges: %+v", badges)
	}
}

func TestFetchBadges_HTMLContentTypeIsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<html>not json</html>`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "testuser")
	_, err := client.FetchBadges(context.Background())
	if !errors.Is(err, ErrUpstreamUnavailable) {
		t.Fatalf("expected ErrUpstreamUnavailable, got %v", err)
	}
}

func TestFetchBadges_NonOKStatusIsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "testuser")
	_, err := client.FetchBadges(context.Background())
	if !errors.Is(err, ErrUpstreamUnavailable) {
		t.Fatalf("expected ErrUpstreamUnavailable, got %v", err)
	}
}

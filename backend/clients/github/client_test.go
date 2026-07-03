package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchRepos_FollowsPagination(t *testing.T) {
	var page1URL string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/testuser/repos", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Error("expected a User-Agent header to be set")
		}
		w.Header().Set("Link", fmt.Sprintf(`<%s/page2>; rel="next"`, page1URL))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"name":"repo-one"}]`))
	})
	mux.HandleFunc("/page2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"name":"repo-two"}]`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	page1URL = server.URL

	client := NewClient(server.URL, "testuser", "")
	repos, err := client.FetchRepos(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repos) != 2 {
		t.Fatalf("expected 2 repos across both pages, got %d", len(repos))
	}
	if repos[0].Name != "repo-one" || repos[1].Name != "repo-two" {
		t.Fatalf("unexpected repo names: %+v", repos)
	}
}

func TestFetchRepos_RateLimitedReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := NewClient(server.URL, "testuser", "")
	_, err := client.FetchRepos(context.Background())
	if err == nil {
		t.Fatal("expected an error on a rate-limited response")
	}
}

func TestFetchPinnedRepos_NoTokenReturnsErrNoToken(t *testing.T) {
	client := NewClient("https://api.github.com", "testuser", "")
	_, err := client.FetchPinnedRepos(context.Background())
	if !errors.Is(err, ErrNoToken) {
		t.Fatalf("expected ErrNoToken, got %v", err)
	}
}

func TestFetchPinnedRepos_SendsAuthAndMapsResponse(t *testing.T) {
	var gotAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")

		var reqBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if _, ok := reqBody["query"]; !ok {
			t.Fatal("expected a GraphQL query in the request body")
		}

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{
			"data": {
				"user": {
					"pinnedItems": {
						"nodes": [
							{"name":"pinned-repo","description":"a pinned repo","url":"https://github.com/testuser/pinned-repo","stargazerCount":5,"forkCount":1,"primaryLanguage":{"name":"Go"}}
						]
					}
				}
			}
		}`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "testuser", "test-token")
	repos, err := client.FetchPinnedRepos(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "Bearer test-token" {
		t.Fatalf("expected Authorization header to be set, got %q", gotAuth)
	}
	if len(repos) != 1 || repos[0].Name != "pinned-repo" || repos[0].StargazersCount != 5 {
		t.Fatalf("unexpected repos: %+v", repos)
	}
}

func TestFetchPinnedRepos_GraphQLErrorIsReturned(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":{"user":null},"errors":[{"message":"Could not resolve to a User"}]}`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "testuser", "test-token")
	_, err := client.FetchPinnedRepos(context.Background())
	if err == nil {
		t.Fatal("expected an error when the GraphQL response contains errors")
	}
}

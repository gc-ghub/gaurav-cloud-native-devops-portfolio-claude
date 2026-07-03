# Backend — DevOps Portfolio API

Go + Gin API serving the frontend's Featured Projects, Certifications, and GitHub stats sections. See the repo root `CLAUDE.md` for the full project spec.

## Setup

```bash
cp .env.example .env   # adjust GITHUB_USERNAME etc. if needed
go mod tidy
go run main.go          # listens on :8080 by default
```

## Endpoints

- `GET /healthz` — liveness check
- `GET /api/projects` — top N non-fork, non-archived repos for `GITHUB_USERNAME`, sorted by stars
- `GET /api/certifications` — Credly badges for `CREDLY_USER_ID`
- `GET /api/github-stats` — aggregate public repo/star/fork counts

No authentication, no write endpoints, no contact-form/email feature (removed by design).

## Testing

```bash
go test ./...
```
Covers the cache (TTL/stale-serve behavior), handlers (via fake fetchers), and both external clients (via `httptest.Server`, no real network calls). No integration or e2e tests — see repo root `tasks.md` for scope notes.

## Docker

```bash
docker build -t portfolio-backend .
docker run -p 8080:8080 --env-file .env portfolio-backend
```
Or build/run the whole stack via `docker compose` from the repo root.

## Known risk: Credly integration

Credly does not offer a public, documented API for anonymous badge reads — `CREDLY_API_URL`'s documented form (`api.credly.com`) does not return JSON. This backend instead calls an **undocumented** endpoint: `{CREDLY_API_URL}/users/{id}/badges.json` (i.e. the public profile URL with `.json` appended), which was verified working at implementation time. This is a website implementation detail, not a stable API contract — it can change or get rate-limited without notice, and its use for programmatic access should be weighed against Credly's Terms of Service before relying on it in production.

To limit blast radius if it breaks: `GET /api/certifications` never returns an error status for this — on any upstream failure it responds `200` with `{"certifications": [], "warning": "..."}`, and the frontend falls back to linking directly to the Credly profile.

## Rate limits

GitHub's REST API allows 60 unauthenticated requests/hour. Both `/api/projects` and `/api/github-stats` share one in-memory cache (`CACHE_TTL_MINUTES`, default 10) backed by a single GitHub call, so normal traffic won't come close to the limit. Set `GITHUB_TOKEN` to raise it to 5,000/hour if needed later.

## Pinned repos require a token

`/api/projects` always leads with one manually-featured project (`clients/github` + `handlers/projects.go` — content sourced from `pics/bio.txt`), then fills the rest from your **pinned** GitHub repos. GitHub's REST API has no endpoint for pinned repos — only the GraphQL API exposes `pinnedItems`, and GraphQL requires auth even for public data. Set `GITHUB_TOKEN` (a classic PAT with no scopes needed — this only reads public data) to enable it.

Without a token (or if the GraphQL call fails for any reason), `/api/projects` **falls back** to the same top-starred-repos behavior used before pinned-repo support existed — the endpoint never errors out just because the token is missing or expired.

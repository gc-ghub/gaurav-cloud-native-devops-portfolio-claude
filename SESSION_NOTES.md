# Session Notes — 2026-07-03

Full detail lives in `CLAUDE.md` (spec + status) and `tasks.md` (granular checklist). This is the quick-scan version.

## What was completed today

- Built the entire portfolio from scratch: Svelte 5 + Tailwind v4 frontend, Go + Gin backend, all 8 roadmap phases (setup → frontend UI → backend → integration → live GitHub/Credly data → polish/SEO/animations → backend unit tests → Dockerized + verified locally).
- Wired in **real content**, replacing every placeholder:
  - Hero bio + one "featured" project sourced from `pics/bio.txt`
  - Hero photo = real GitHub avatar
  - Featured Projects = 1 featured (live stars/forks) + real **pinned** GitHub repos, fetched via GraphQL (the only way to read pinned repos — REST has no endpoint for it). User provided a token; confirmed live pinned data returned, not a fallback.
  - Certifications = real Credly badges
- Removed entirely (not placeholders anymore): Work Experience section, Resume button, Contact Me (button/nav/footer anchor/fake email) — nav is now just Skills / Featured Projects / Achievements.
- Added real GitHub/LinkedIn/Credly brand icons (Hero + Footer).
- Added root `.gitignore` (repo isn't a git repo yet, but it's ready).

## Files changed (high level)

- **`frontend/`** — full app: `src/components/*`, `src/lib/{api,socials,reveal}.js`, `src/styles/*`, `index.html`, `Dockerfile`, `.dockerignore`, `.env.example`
- **`backend/`** — full app: `main.go`, `config/`, `models/`, `clients/{github,credly}/` (+ tests), `cache/` (+ tests), `handlers/` (+ tests), `middleware/`, `Dockerfile`, `.dockerignore`, `.env.example`, `.gitignore`
- **Root** — `docker-compose.yml`, `nginx.conf` (unwired future reverse-proxy template), `.gitignore`, `CLAUDE.md`, `tasks.md`
- **Not tracked/committed anywhere:** `backend/.env` — contains the real GitHub token, gitignored, never printed in full to a persisted file

## Current project status

- Fully working full-stack app, **verified locally only** — dev (`localhost:5173` + `:8080`) and Docker (`localhost:8081` + `:8080`). No live deployment, no domain, no cloud account.
- All backend tests pass (`go test ./...`). Frontend builds cleanly (`npm run build`). No frontend test suite (user's explicit scope choice).
- Repo is **not yet a git repository** — nothing has been committed.

## Next steps when you resume

1. **Decide whether to `git init` and make a first commit.** Everything is ready (`.gitignore` in place at root + `frontend/` + `backend/`), just hasn't been done.
2. **Real deployment** (Phase 8 production, not started): needs a VPS/cloud account + domain from you — I have no infrastructure access. Once you have one, `nginx.conf` at repo root is a ready-to-wire reverse-proxy template.
3. **Swap the placeholder SEO domain**: `frontend/index.html`'s canonical/OG URLs currently say `gaurav-chaurasia-portfolio.example.com` — update once a real domain exists.
4. Optionally revisit whether you want a lightweight "About"/work-history section back — it was removed outright per your instruction, not filled with real data (you didn't provide any).

## Commands to remember

```bash
# Frontend dev
cd frontend && npm run dev          # localhost:5173

# Backend dev
cd backend && go run main.go        # localhost:8080, needs backend/.env
cd backend && go test ./...         # run backend unit tests

# Full stack via Docker
docker compose build && docker compose up -d   # :8081 frontend, :8080 backend
docker compose down
```

## Caveats to remember

- **`.env` files: never put a comment on the same line as a value** (`KEY=value # comment`). `godotenv` doesn't strip trailing comments — it broke every GitHub call once already when `GITHUB_TOKEN=  # comment` made the comment part of the token. Comments must go on their own line above.
- **Credly integration uses an undocumented endpoint** (`credly.com/users/{id}/badges.json`, not the documented `api.credly.com`). Could break without notice — degrades gracefully (empty list + warning) if it does. See `backend/README.md`.
- **Pinned repos need `GITHUB_TOKEN`** (classic PAT, no scopes needed — public data only) in `backend/.env`. Without it, `/api/projects` silently falls back to top-starred repos instead of erroring.
- **`backend/.env` holds a real secret** — confirmed gitignored; double-check before any commit that it's not accidentally staged.

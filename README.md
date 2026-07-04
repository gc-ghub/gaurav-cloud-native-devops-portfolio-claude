# Gaurav Chaurasia — DevOps & Platform Engineering Portfolio

A full-stack portfolio site showcasing DevOps, Platform Engineering, SRE, and Cloud-Native work — built with **Svelte 5 + Tailwind CSS v4** on the frontend and **Go + Gin** on the backend. Content isn't hardcoded: projects, stars/forks, and certifications are pulled live from GitHub and Credly at request time.

- **GitHub:** https://github.com/gc-ghub
- **LinkedIn:** https://www.linkedin.com/in/gc-devops-world/
- **Credly:** https://www.credly.com/users/gaurav-chaurasia.25908149/badges

## Features

- **Live GitHub data** — Featured Projects leads with one hand-picked project (live stars/forks), then fills in from the user's real **pinned** GitHub repos via the GitHub GraphQL API, falling back to top-starred repos if no token is configured.
- **Live Credly certifications** — real badges (AWS SysOps, AWS Developer, Terraform Associate, GitHub Actions, etc.), fetched and cached server-side.
- **Three-way theme toggle** — Green (default) / Dark / Light, all driven by CSS custom properties and `localStorage`.
- **Skills & Expertise with real icons** — every skill badge (Kubernetes, Docker, Terraform, ArgoCD, ...) renders its actual brand icon, not just a text label.
- **No placeholders** — every section is backed by real content or a real API; anything that couldn't be (a work-history section, a resume download, a contact form) was removed outright rather than left as a stub.

## Tech Stack

| Layer | Stack |
|---|---|
| Frontend | Svelte 5 (runes), Vite, Tailwind CSS v4 |
| Backend | Go, Gin |
| External APIs | GitHub REST + GraphQL, Credly |
| Containerization | Docker, Docker Compose |

## Project Structure

```
.
├── frontend/       # Svelte + Vite SPA — see frontend/README.md
├── backend/        # Go + Gin API — see backend/README.md
├── docker-compose.yml
├── nginx.conf       # production reverse-proxy template (not yet wired to a live domain)
├── CLAUDE.md        # full internal spec, status, and design detail
├── tasks.md         # granular build/checklist history
└── decisions.md     # architecture/design decisions log, with rationale
```

## Running it locally

### Option 1 — Docker Compose (recommended)

```bash
docker compose build && docker compose up -d
```
Open **http://localhost:8081** (frontend, proxying API calls to the backend on `:8080` internally).

```bash
docker compose down   # when done
```

### Option 2 — Dev mode (hot reload)

```bash
# Terminal 1 — backend
cd backend
cp .env.example .env   # fill in GITHUB_USERNAME / GITHUB_TOKEN etc.
go run main.go          # localhost:8080

# Terminal 2 — frontend
cd frontend
npm install
npm run dev              # localhost:5173
```

See `backend/README.md` and `frontend/README.md` for endpoint details, environment variables, and testing commands for each subproject.

## Status

Fully working full-stack app, verified locally (dev mode and Docker Compose) with real GitHub/Credly data end-to-end. No live deployment exists yet — no domain or hosting is currently provisioned. See `CLAUDE.md` for the complete phase-by-phase status and `decisions.md` for the reasoning behind key design choices.

## License

No license file yet — all rights reserved by default. Open an issue if you'd like this changed.

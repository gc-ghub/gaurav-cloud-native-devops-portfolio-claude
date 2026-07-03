# CLAUDE.md

# Cloud-Native DevOps Portfolio Platform

**Project Name:** Gaurav Chaurasia's DevOps & Platform Engineering Portfolio

**Owner:** Gaurav Chaurasia  
**Repository:** https://github.com/gc-ghub/gaurav-cloud-native-devops-portfolio-claude  
**Deployment:** docker container / docker-compose /ngrok (local development) / Self-hosted (production)

## 📋 Project Overview

A modern, production-grade portfolio website built with **Svelte + Go** showcasing DevOps, Platform Engineering, SRE, and Cloud-Native expertise. The portfolio dynamically pulls projects from GitHub, displays certifications from Credly, and provides social integration with LinkedIn and GitHub.

**Design:** Template 3 (Dark Modern DevOps) with dark/light theme toggle  
**Live Profiles Integration:**
- GitHub: https://github.com/gc-ghub (23 public repos)
- LinkedIn: https://www.linkedin.com/in/gc-devops-world/
- Credly Certifications: https://www.credly.com/users/gaurav-chaurasia.25908149/badges

## ✅ Current Status

**All 8 phases are done**, with Phase 8 scoped to local Docker verification (see below — no live deployment exists). This is a real, working full-stack app with **real content throughout — no more placeholders**:

- **`frontend/`** — Vite + Svelte 5 + Tailwind v4 SPA styled after `pics/portfolio_template-3.png`, dark/light theme toggle, scroll-reveal animations (respects `prefers-reduced-motion`), SEO meta/OG/JSON-LD. Builds cleanly, verified via Playwright (dark/light/mobile, no console errors).
- **`backend/`** — Go + Gin API on `:8080` serving `GET /api/projects`, `GET /api/certifications`, `GET /api/github-stats`, and `GET /healthz`, with unit test coverage (`go test ./...`, all passing).
- **Real bio, real photo:** Hero bio text and one explicitly "featured" project are sourced from `pics/bio.txt` (the user's real content, not invented copy). Hero avatar is the user's real GitHub avatar (`https://avatars.githubusercontent.com/u/161034127?v=4`).
- **Real pinned + featured projects, live and verified:** `/api/projects` always leads with the featured project (from `pics/bio.txt`, with **live** stars/forks looked up from GitHub), then fills the rest from the user's real **pinned** GitHub repos via the GitHub **GraphQL** API (REST has no pinned-repos endpoint, and GraphQL requires a token even for public reads). The user provided a token and this was confirmed working live (real pinned repos returned, not the star-sorted fallback). If `GITHUB_TOKEN` is ever unset or the GraphQL call fails, `/api/projects` gracefully **falls back** to top-starred non-fork/non-archived repos — verified working both ways.
- Certifications pull real Credly badges for `gaurav-chaurasia.25908149` (AWS SysOps, AWS Developer, GitHub Actions, HashiCorp Terraform Associate).
- **Social icons** (GitHub, LinkedIn, Credly) — real brand SVGs, real links, centralized in `frontend/src/lib/socials.js`, shown in both Hero and Footer.
- **Everything fake has been removed outright, not left as a placeholder:** no `WorkExperienceSection.svelte` (deleted), no "Download My Resume" button, no "Contact Me" button/nav link/footer anchor, no fake `mailto:contact@example.com` anywhere. Nav is now just Skills / Featured Projects / Achievements.
- **Dockerized and verified locally** — `docker compose build && docker compose up -d` runs the full stack (backend on `:8080`, frontend on `:8081`), confirmed serving real data through the containers via Playwright, then torn down. No live server/domain exists — see Deployment Strategies below for what that would still take.

**Known risk:** the Credly integration uses an **undocumented** endpoint (`credly.com/users/{id}/badges.json`) — CLAUDE.md's originally-documented `api.credly.com` URL does not actually return JSON. This was a deliberate, user-approved tradeoff; see `backend/README.md` for the full risk writeup. The certifications endpoint degrades gracefully (empty list + warning, not an error) if this ever breaks.

**A note for future sessions:** `backend/.env` contains the user's real GitHub token — it's gitignored (`backend/.gitignore`) and must never be committed or printed unnecessarily. Also: never put a comment on the same line as a value in `.env`/`.env.example` (`KEY=value # comment`) — `godotenv` doesn't strip trailing comments, so they become part of the value. This broke every GitHub call once already; comments must go on their own line above.

**No longer placeholder, nothing left to fill in from the original roadmap** except the SEO canonical/OG URL (`gaurav-chaurasia-portfolio.example.com` — not a real domain, swap once one exists).

Everything currently runs locally only — dev: frontend `localhost:5173` / backend `localhost:8080`; Docker: frontend `localhost:8081` / backend `localhost:8080`.

See `tasks.md` at the repo root for the granular build checklist.

## 🎯 Core Features

### Hero Section
- Greeting: "Hi 👋, I'm Gaurav Chaurasia" (from `pics/bio.txt`)
- Professional tagline: "Platform Engineer | Kubernetes | GitOps | SRE | Cloud-Native Observability | AI-Powered Incident Response"
- Bio: real intro paragraph from `pics/bio.txt`
- Photo: real GitHub avatar (`https://avatars.githubusercontent.com/u/161034127?v=4`)
- Social icons: GitHub, LinkedIn, Credly (clickable, real brand SVGs, linking to actual profiles)
- No CTA buttons — "Download My Resume" and "Contact Me" were both removed entirely (no resume PDF, no real contact mechanism beyond the social links)

### Navigation Menu
```
Skills | Featured Projects | Achievements
```
("Work Experiences", "Resume", and "Contact Me" were removed from the original spec — see Current Status)

### Key Sections

#### 1. **Featured Projects** (Primary showcase)
`GET /api/projects` always leads with one manually-featured project (content from `pics/bio.txt`, live stars/forks), then fills the rest from the user's real GitHub **pinned** repos (via GraphQL, falls back to top-starred if no token) — up to `MAX_PROJECTS` (default 6) total. Each card shows:
- Project name and description
- GitHub repository link
- Tech stack badges (from repo topics, or primary language as a fallback)
- Stars/forks count
- Link to full GitHub profile

The originally-sketched hardcoded 5-project list (GC Industries, Stark Industries ×2, etc.) was the Phase 1-2 placeholder — GC Industries is now the one real featured project (`gc-ghub/project-gc-industries-devops-superheros`), the rest come live from pinned repos instead of being hardcoded.

#### 2. **Skills & Expertise** (Categorized display)
Core DevOps Skills:
- **Container & Orchestration:** Kubernetes, Docker, Helm, Istio, Kyverno
- **CI/CD & Automation:** GitHub Actions, Jenkins, Docker Compose
- **Infrastructure as Code:** Terraform, Ansible, Helm
- **Cloud Platforms:** AWS, Azure, GCP
- **Observability & Monitoring:** Prometheus, Grafana, Loki, Tempo, Kiali
- **Programming/Scripting:** Python, Go, Bash, FastAPI
- **Cloud-Native Tools:** Argo CD, Kafka, Gitleaks, Trivy, CodeQL

#### 3. **Certifications & Achievements**
- Pull certifications from Credly API
- Display with badges/logos
- Link to Credly profile
- Include: AWS SysOps, AWS Developer, Terraform Associate, GitHub Actions, etc.

#### 4. **Engineering Philosophy** (Values section)
- Automate everything that can be automated
- Treat infrastructure as a product
- Design for reliability, scalability, and observability
- Shift security left through DevSecOps practices
- Build platforms that empower engineering teams
- Solve operational problems with automation and AI

## 🛠️ Tech Stack

### Frontend
- **Framework:** Svelte 5, using runes (`$state`, `$props`) — ✅ in use
- **Build Tool:** Vite — ✅ in use
- **Styling:** Tailwind CSS **v4** — ✅ in use, CSS-first config (`@theme` block in `src/styles/global.css`, no `tailwind.config.js`)
- **Theme Toggle:** Dark/Light mode (CSS custom properties + `.light` class on `<html>` + `localStorage`) — ✅ implemented (`src/stores/theme.js`)
- **Icons:** Inline SVG — ✅ implemented
- **API Client:** Fetch API — ✅ implemented (`src/lib/api.js`), calls the backend for projects/certifications/stats
- **Routing:** None yet — currently a single-page SPA, all sections assembled directly in `App.svelte`. SvelteKit remains a future option if multi-page routing is needed.

### Backend
- **Language:** Go — ✅ in use
- **Framework:** Gin — ✅ in use (chose Gin over Fiber)
- **Database:** None — ✅ matches MVP plan (static/live-fetched data only, no persistence needed)
- **API Endpoints:**
  - `GET /api/projects` — ✅ implemented: 1 featured project (hardcoded content + live stats) + real pinned repos (GraphQL, token-gated) or top-starred fallback, in-memory cached
  - `GET /api/certifications` — ✅ implemented, real Credly data (via undocumented endpoint, see Current Status), in-memory cached
  - `GET /api/github-stats` — ✅ implemented, shares the projects cache (no extra API call)
  - `GET /healthz` — ✅ implemented (not in original spec, added for future deployment liveness checks)
  - ~~`POST /contact` - Contact form submission~~ — **removed**, feature dropped entirely per user decision

### External APIs
- **GitHub REST API:** ✅ implemented — unauthenticated by default, 60 req/hr limit mitigated by a shared 10-minute in-memory cache
- **GitHub GraphQL API:** ✅ implemented (`clients/github/client.go`'s `FetchPinnedRepos`) — the only way to read pinned repos; requires `GITHUB_TOKEN` (classic PAT, no scopes needed) even though the data is public. Confirmed working live with a real user-provided token.
- **Credly API:** ✅ implemented, but via an **undocumented** endpoint — see Current Status / `backend/README.md` for the risk
- ~~**Contact Form:** Email service~~ — **removed**, feature dropped entirely

## 📁 Project Structure

```
portfolio-devops/
├── frontend/                    # Svelte + Vite SPA — ✅ built (Phase 1-2, 6), Dockerized (Phase 8)
│   ├── src/
│   │   ├── components/         # One component per portfolio section
│   │   │   ├── Header.svelte           # logo, nav (Skills/Featured Projects/Achievements only), theme toggle
│   │   │   ├── HeroSection.svelte             # real bio + real GitHub avatar photo
│   │   │   ├── SkillsSection.svelte
│   │   │   ├── ProjectCard.svelte
│   │   │   ├── FeaturedProjects.svelte        # fetches live data (featured + pinned), scroll-reveal
│   │   │   ├── GithubStatsBar.svelte          # small stat pills, fetched independently
│   │   │   ├── CertificationsSection.svelte   # fetches live data, scroll-reveal
│   │   │   ├── PhilosophySection.svelte
│   │   │   ├── SocialIcon.svelte              # renders github/linkedin/credly SVGs by name
│   │   │   └── Footer.svelte
│   │   │   (no WorkExperienceSection.svelte or ContactForm.svelte — both removed entirely)
│   │   ├── lib/
│   │   │   ├── api.js                  # fetchProjects/fetchCertifications/fetchGithubStats
│   │   │   ├── socials.js              # single source of truth for GitHub/LinkedIn/Credly links
│   │   │   └── reveal.js               # IntersectionObserver scroll-reveal action, respects prefers-reduced-motion
│   │   ├── stores/
│   │   │   └── theme.js
│   │   ├── styles/
│   │   │   ├── global.css              # Tailwind v4 entry + @theme tokens
│   │   │   └── variables.css           # light-mode overrides
│   │   └── App.svelte
│   ├── public/
│   │   └── robots.txt
│   ├── Dockerfile               # multi-stage: node:22-alpine build → nginx:1.27-alpine serve
│   ├── package.json
│   ├── postcss.config.js
│   ├── svelte.config.js
│   ├── vite.config.js
│   └── README.md
│   (no pages/ directory — single-page SPA for now)
│
├── backend/                     # Go + Gin API server — ✅ built (Phase 3-5), tested (Phase 7)
│   ├── main.go                  # wiring: config → clients → caches → router → routes
│   ├── config/
│   │   └── config.go            # env loading, fails fast if required vars missing
│   ├── models/
│   │   ├── project.go
│   │   ├── certification.go
│   │   └── github_stats.go
│   ├── clients/
│   │   ├── github/{client.go,types.go,client_test.go}   # REST (repos, pagination) + GraphQL (FetchPinnedRepos)
│   │   └── credly/{client.go,types.go,client_test.go}   # real Credly badges.json calls
│   ├── cache/
│   │   └── {cache.go,cache_test.go}   # generic TTL cache, shared across handlers
│   ├── handlers/
│   │   ├── projects.go (+ _test.go)
│   │   ├── certifications.go (+ _test.go)
│   │   └── stats.go (+ _test.go)
│   ├── middleware/
│   │   └── cors.go
│   ├── Dockerfile                # multi-stage: golang:1.25-alpine build → alpine, non-root
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   └── README.md
│   (no contact.go/auth.go — contact feature removed, no auth needed for a public read-only API)
│
├── docker-compose.yml            # ✅ verified: backend (:8080) + frontend (:8081) locally
├── nginx.conf                    # documented template for a future production reverse proxy — not wired into docker-compose.yml yet (no domain to point it at)
├── CLAUDE.md                     # This file
└── tasks.md                      # granular build checklist
```

## 🚀 Build & Development Commands

### Frontend Setup
```bash
cd frontend
npm install                      # Install dependencies
npm run dev                      # Start dev server (localhost:5173)
npm run build                    # Production build
npm run preview                  # Preview production build
```
(no lint script configured yet)

### Backend Setup
```bash
cd backend
go mod download                  # Download Go dependencies
go run main.go                   # Start dev server (localhost:8080)
go build -o backend .            # Build binary
go test ./...                    # Run tests — ✅ passing (Phase 7)
```

### Docker/Container Setup — ✅ verified locally (Phase 8)
```bash
docker compose build             # Build both images
docker compose up -d             # backend on :8080, frontend on :8081
docker compose down              # Stop and remove containers
```
No production target configured — `nginx.conf` at the repo root is a documented template for a future single-domain reverse proxy, not yet wired into `docker-compose.yml`.

## 🔧 Configuration

### Environment Variables (.env)
Actual shape in use — see `backend/.env.example`:
```
GITHUB_USERNAME=gc-ghub
GITHUB_API_URL=https://api.github.com
GITHUB_TOKEN=                      # needed for pinned repos (GraphQL); falls back to top-starred repos without it

CREDLY_USER_ID=gaurav-chaurasia.25908149
CREDLY_API_URL=https://www.credly.com   # note: NOT api.credly.com — see Current Status

CACHE_TTL_MINUTES=10
MAX_PROJECTS=6

PORT=8080
FRONTEND_URL=http://localhost:5173
ENVIRONMENT=development
```
No SMTP/contact vars — that feature was removed entirely, not just deferred.

## 🎨 Design Specifications

### Theme (Dark + Light Mode)
**Dark Mode (Primary):**
- Background: #1a1a2e or #0f1419
- Primary: #6f3ce5 (Purple) or #00d4ff (Cyan)
- Text: #ffffff / #e0e0e0
- Accent: #ff006e (Magenta)

**Light Mode:**
- Background: #ffffff / #f5f5f5
- Primary: #6f3ce5 (Purple)
- Text: #1a1a1a
- Accent: #ff006e (Magenta)

### Components
- **Cards:** Glassmorphism effect (frosted glass) on dark mode
- **Buttons:** Gradient backgrounds, smooth transitions
- **Icons:** Consistent sizing, 24-32px standard
- **Typography:** Clean, modern fonts (Inter, Poppins, or system fonts)
- **Animations:** Subtle, fade-in/slide effects on scroll

## 📊 Key Integrations

### GitHub API Integration — ✅ implemented (`backend/clients/github/`)
```
GET https://api.github.com/users/gc-ghub/repos?per_page=100&sort=updated
```
Paginated (up to 300 repos as a safety cap), requires a `User-Agent` header (GitHub 403s without one), unauthenticated by default.

### Credly API Integration — ✅ implemented (`backend/clients/credly/`)
```
GET https://www.credly.com/users/gaurav-chaurasia.25908149/badges.json
```
**Not** the `api.credly.com` URL originally sketched above — that doesn't return JSON. This is an undocumented endpoint (public profile URL + `.json`), verified working but not a stable contract. See Current Status and `backend/README.md`.

### Contact Form — removed
The contact form/email feature was dropped entirely (user decision) rather than deferred. No `POST /contact`, no SMTP, no `ContactForm.svelte`, and no "Contact Me" button/nav link either — the footer's direct GitHub/LinkedIn/Credly links are the only way to reach out now.

## 🌐 Deployment Strategies

**Status:** Dockerized and verified locally (`docker compose up`) — nothing below this point has actually been provisioned. There's no cloud account, VPS, or domain available to deploy to; the options below remain the plan for whenever one exists.

### Local Development (ngrok)
```bash
# Start backend on port 8080
go run main.go

# In another terminal, expose via ngrok
ngrok http 8080

# Update frontend API URL to ngrok URL
```

### Production Options
1. **Self-Hosted (AWS EC2/VPS)**
	- Docker containers on EC2
	- nginx reverse proxy
	- SSL via Let's Encrypt

2. **Containerized Deployment**
	- Push to Docker Hub
	- Deploy via Kubernetes (fits your expertise!)
	- Or use Docker Swarm

3. **Serverless (Optional)**
	- Frontend: AWS S3 + CloudFront
	- Backend: AWS Lambda + API Gateway

## 📝 Development Workflow

1. **Phase 1 - Setup:** Project structure, build configs, theme system — ✅ Done
2. **Phase 2 - Frontend:** Build UI components based on template 3 — ✅ Done
3. **Phase 3 - Backend:** Set up Go server, API endpoints — ✅ Done
4. **Phase 4 - Integration:** Connect frontend to backend APIs — ✅ Done
5. **Phase 5 - External APIs:** GitHub & Credly data fetching — ✅ Done
6. **Phase 6 - Polish:** Animations, SEO, performance optimization — ✅ Done
7. **Phase 7 - Testing:** Unit tests, integration tests, e2e tests — ✅ Done (Go backend unit tests only, by explicit user scope choice — no frontend test suite)
8. **Phase 8 - Deployment:** Docker setup, production deployment — ✅ Dockerized and verified locally; production deployment itself is not done (no server/domain/credentials available)

See "Current Status" above for what's actually delivered and what's still placeholder.

## 🔍 Key Considerations

- **Responsive Design:** Mobile-first approach (320px to 1920px+) — ✅ verified via Playwright at mobile/desktop widths
- **Performance:** ✅ certification badge images use `loading="lazy"`/`decoding="async"`; bundle is small (~55KB JS / ~20KB gzip), no code-splitting needed for a single page
- **SEO:** ✅ meta description, Open Graph, Twitter card, JSON-LD `Person` structured data, `robots.txt` — canonical/OG URLs use a placeholder domain pending real deployment
- **Accessibility:** WCAG 2.1 AA compliance — ✅ scroll-reveal animations respect `prefers-reduced-motion` (verified via Playwright media emulation)
- **Security:** Environment variables for secrets, CORS setup, input validation — ✅ CORS configured; no user input accepted anywhere (contact form removed), so little attack surface
- **Rate Limiting:** GitHub API has 60 req/hour (unauthenticated) — ✅ mitigated via a shared 10-minute in-memory cache (`backend/cache/cache.go`)
- **Certifications:** badges refresh whenever the cache expires (10 min TTL) on next request, not a scheduled background sync

## 📞 Contact & Social Links

- **GitHub:** https://github.com/gc-ghub
- **LinkedIn:** https://www.linkedin.com/in/gc-devops-world/
- **Credly:** https://www.credly.com/users/gaurav-chaurasia.25908149/badges
- **Email:** removed — no contact form or email link exists in the UI anymore (user decision)

## 📚 Reference Files

- `pics/portfolio_template-1.png` - Colorful modern design
- `pics/portfolio_template-2.png` - Professional academic layout
- `pics/portfolio_template-3.png` - Dark modern DevOps (PRIMARY CHOICE)
- `pics/bio.txt` - **Source of truth** for the Hero bio text and the one featured project's content (name, description, repo URL)

---



# Tasks — Phase 1-8 (Full roadmap)

Tracking implementation across the plans in `.claude/plans` (DevOps Portfolio). All 8 CLAUDE.md phases are now done, with Phase 8 scoped to "Dockerized and verified locally" (see note below).

## Phase 1-2 (Setup + Frontend UI) — done
- [x] Vite + Svelte scaffold, Tailwind v4, theme toggle, all sections built with placeholder content, verified via Playwright + `npm run build`

## Phase 3-5 (Go backend, live GitHub/Credly data, frontend wiring) — done
- [x] Gin backend (`backend/`), real GitHub + Credly data, shared TTL cache, CORS, contact form removed entirely
- [x] Frontend wired to live data (`FeaturedProjects.svelte`, `CertificationsSection.svelte`, `GithubStatsBar.svelte`, `frontend/src/lib/api.js`)
- [x] `gc-ghub` confirmed as the real GitHub username; Credly integration confirmed working via undocumented `.json` endpoint (see `backend/README.md` for the risk)

## Social icons (GitHub, LinkedIn, Credly) — done
- [x] `frontend/src/lib/socials.js` (single source of truth for the real links) + `frontend/src/components/SocialIcon.svelte` (Credly's real brand SVG fetched from simple-icons and verified)
- [x] `HeroSection.svelte` and `Footer.svelte` both updated to use real icons (Credly icon was previously missing everywhere; Footer was plain text)
- [x] The fake `mailto:contact@example.com` "Email" entry was later removed entirely (see below) — only GitHub/LinkedIn/Credly remain

## Phase 6 (Polish) — done
- [x] SEO: meta description, Open Graph, Twitter card, JSON-LD `Person` structured data, `robots.txt` — canonical/og:url use a **placeholder domain** (`gaurav-chaurasia-portfolio.example.com`) since there's no real deployment URL yet; swap once one exists
- [x] Scroll-reveal animation (`frontend/src/lib/reveal.js`, IntersectionObserver-based) applied to Skills/Featured Projects/Certifications/Philosophy sections (Work Experience section was later removed entirely); verified via Playwright to respect `prefers-reduced-motion`
- [x] `loading="lazy"` + `decoding="async"` on certification badge images

## Phase 7 (Testing) — done (Go backend only, per user's explicit scope choice)
- [x] `backend/cache/cache_test.go`, `backend/handlers/{projects,stats,certifications}_test.go`, `backend/clients/{github,credly}/client_test.go`
- [x] `go test ./...` — all passing
- No frontend test suite added (user did not select Vitest/Playwright checked-in tests)

## Phase 8 (Deployment) — done, scoped to local Docker verification
- [x] `backend/Dockerfile` (multi-stage, `golang:1.25-alpine` → alpine, non-root), `frontend/Dockerfile` (multi-stage, `node:22-alpine` → `nginx:1.27-alpine`)
- [x] Root `docker-compose.yml` — verified end-to-end: `docker compose build && docker compose up -d`, confirmed real GitHub/Credly data through the containerized backend, full-page Playwright screenshot after scrolling confirms all sections + icons render correctly, then `docker compose down`
- [x] Root `nginx.conf` — documented template for a *future* single-domain production reverse proxy, **not wired into docker-compose.yml** (no real domain to point it at yet)
- **Not done, and explicitly out of scope per user decision:** actually provisioning/deploying to a live server, VPS, or cloud account — no credentials or infrastructure available to do that from here

## Deviations / notes worth remembering
- `npm ci` failed in the frontend Dockerfile because `package-lock.json` (generated on Windows) was missing Linux-specific optional deps — switched to `npm install` in the Dockerfile to fix.
- The very first `fullPage` Playwright screenshot after adding scroll-reveal showed blank sections below the fold — this was a **screenshot artifact**, not a bug: `IntersectionObserver` only fires on real scroll events, which a non-scrolling full-page capture doesn't trigger. Confirmed fine once the verification script actually scrolled through the page first.

## Real content + removed fake sections + pinned repos — done
- [x] `pics/bio.txt` (real bio + one explicitly featured project) is now the content source: Hero bio paragraph rewritten to the real text, `backend/handlers/projects.go` hardcodes the featured project (GC Industries – Cloud-Native DevOps Superheroes Platform) with **live** stars/forks looked up from the general repo list
- [x] Hero avatar replaced with the real GitHub avatar (`https://avatars.githubusercontent.com/u/161034127?v=4`, stable/keyed by user ID)
- [x] **Pinned repos implemented and verified live**: `backend/clients/github/client.go` added `FetchPinnedRepos` via the GitHub GraphQL API (the only way to read pinned repos — REST has no equivalent, and GraphQL requires a token even for public data). User provided a classic PAT with no scopes; confirmed `/api/projects` returns real pinned repos (3 pinned + 1 featured = 4 total, since the user only has 3 pinned repos, not 5). Falls back to the pre-existing top-starred-repos logic if `GITHUB_TOKEN` is unset or the GraphQL call fails — verified both paths work.
- [x] Removed entirely (not placeholders anymore): `WorkExperienceSection.svelte` (deleted), "Download My Resume" button + "Resume" nav link, "Contact Me" button + nav link + footer anchor, the fake `mailto:contact@example.com` (removed from `socials.js`, `SocialIcon.svelte`'s now-dead mail case removed too)
- [x] Nav is now just Skills / Featured Projects / Achievements — verified via Playwright
- [x] Root `.gitignore` added (repo isn't a git repo yet, but this is ready for when it is)

## Bug caught and fixed this session
- An inline comment on the same line as `GITHUB_TOKEN=` in `.env.example` (`GITHUB_TOKEN=    # needed for...`) got parsed by `godotenv` as *part of the token value* (godotenv doesn't strip trailing comments), which broke **all** GitHub REST calls, not just pinned repos, since every request suddenly sent a garbage `Authorization` header. Moved the comment to its own line above. Worth remembering for any future `.env.example` edits — never put a comment after a value on the same line, even if the value is empty.

## Remaining open assumptions (real content gaps nothing here can fill)
- Work history is gone (not filled in) — there's no work history section anymore at all, by design
- SEO canonical/OG URLs still use a placeholder domain until real deployment exists
- `backend/.env` now contains the user's real GitHub token — confirmed gitignored (`backend/.gitignore`), never committed

## Green theme + Skills icons + automated code-review hook — done
- [x] Added a third "green" theme (`frontend/src/stores/theme.js`, `frontend/src/styles/variables.css`) — toggle now cycles Green → Dark → Light → Green, **green is the default**. Colors (`#0b4640` bg, `#5cda86` accent) pixel-sampled directly from a user-provided screenshot via a Python/PIL script, deliberately excluding the screenshot's cream/tan tones per instruction — see `decisions.md`.
- [x] `Header.svelte` toggle button updated: 3-state icon cycle (leaf → moon → sun, each showing the *next* theme), dynamic `aria-label`, dynamic mobile-menu text, all driven by a new `nextTheme()` export.
- [x] Added real tech icons to every Skills & Expertise badge (`frontend/src/components/TechIcon.svelte`, `frontend/src/lib/techIcons.js`): 21 official brand SVGs fetched from simple-icons (CC0) at their real hex color; 3 generic outline icons (shield/key/mesh) for Kyverno/Gitleaks/Kiali, which have no available official logo; Docker Compose/Loki/Tempo/CodeQL reuse their parent vendor's icon (Docker/Grafana/Grafana/GitHub respectively).
- [x] Verified via `npm run build` (clean, no warnings) and a headless-browser pass (Playwright): all 3 themes screenshot correctly, default-on-load is green, and the Skills section renders all 27 badges with icons and no console errors.
- [x] Added a project-level automated code-review hook (`.claude/settings.json`, `Stop` event): runs `claude -p "/code-review medium"` in the background whenever a turn ends with uncommitted changes, guarded against recursive self-triggering via a `CLAUDE_AUTO_REVIEW` env sentinel, logs to `.claude/last-code-review.log` (already covered by the root `*.log` gitignore rule). Pipe-tested the guard logic, then live-ran it once against this session's own diff as both the hook proof and the requested review of this change — came back clean (no correctness/simplification issues at medium effort).
- [x] Documented rationale in the new `decisions.md` (theme color sourcing, icon sourcing/fallback strategy, hook design and recursion guard).

# Frontend — DevOps Portfolio

Svelte + Vite + Tailwind CSS v4 frontend for Gaurav Chaurasia's DevOps & Platform Engineering portfolio. See the repo root `CLAUDE.md` for the full project spec and roadmap.

## Status

All 8 phases complete (see repo root `CLAUDE.md`/`tasks.md`). Single-page layout styled after `pics/portfolio_template-3.png`, dark/light theme toggle, scroll-reveal animations, SEO meta/OG/JSON-LD, and live data from the backend for Featured Projects, Certifications, and GitHub stats. Bio, work history, resume PDF, and contact email remain placeholder content — no API can supply those.

## Commands

```bash
npm install       # Install dependencies
npm run dev       # Start dev server (localhost:5173) — expects the backend running on :8080
npm run build     # Production build
npm run preview   # Preview production build
```

## Structure

- `src/components/` — one component per portfolio section (Header, Hero, Skills, Work Experience, Featured Projects, Certifications, Philosophy, Footer)
- `src/lib/api.js` — fetches live data from the backend
- `src/lib/socials.js` — single source of truth for GitHub/LinkedIn/Credly/Email links, used by `SocialIcon.svelte`
- `src/lib/reveal.js` — scroll-reveal Svelte action (`use:reveal`), respects `prefers-reduced-motion`
- `src/stores/theme.js` — dark/light theme store, persisted to `localStorage`
- `src/styles/global.css` — Tailwind v4 entry point + `@theme` color tokens (dark is default)
- `src/styles/variables.css` — light-mode variable overrides (`.light` class on `<html>`)

## Docker

```bash
docker build -t portfolio-frontend --build-arg VITE_API_BASE_URL=http://localhost:8080 .
docker run -p 8081:80 portfolio-frontend
```
Or build/run the whole stack via `docker compose` from the repo root.

## Notable choices

- **Tailwind v4**, not v3 — config is CSS-first (`@theme` in `global.css`), so there's no `tailwind.config.js`.
- Theming is done via CSS custom properties, not Tailwind's `dark:` variant — components use plain utilities (`bg-bg`, `text-text`, `bg-primary`, ...) whose values swap when `.light` is toggled on `<html>`.
- No contact form — removed entirely by design; "Contact Me" scrolls to the footer's direct links instead.

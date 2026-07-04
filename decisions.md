# Decisions Log

Architecture/design decisions worth remembering the *why* behind, not just the *what* (the *what* is already visible in the code and `CLAUDE.md`/`tasks.md`). Newest first.

## Green theme: colors pixel-sampled from a reference screenshot, only the greens

**Decision:** Add a third "green" theme to the existing dark/light toggle, make it the default, and source its colors by pixel-sampling a WhatsApp promo screenshot the user provided — but only the green hues in that image, not its cream/tan accent color.

**Why:** The user explicitly said "Only take the green color from the pic. Nothing else," ruling out reproducing the screenshot's full palette (it also uses a cream/tan `#fbe9db` for headings and a CTA button) or its layout. Eyeballing a hex value from a screenshot is unreliable, so a small Python/PIL script (`Counter` over sampled pixels) was used to find the two most common green tones precisely: `#0b4640` (background, ~147k sampled pixels) and `#5cda86` (accent, sampled from the "3B+"/"2026" stat numbers). A third, slightly lighter green (`#194f4b`) was identified as the card-background tone but wasn't used directly — the theme's `--color-surface` instead uses a low-alpha tint of the accent (`rgba(92,218,134,0.08)`), matching the existing dark/light themes' pattern of a translucent glassmorphism surface rather than a flat card color.

**How to apply:** Text color stays white (unchanged from dark mode) rather than adopting the screenshot's cream tone — that's the "nothing else" part of the instruction. If the user ever asks for a similar "extract colors from an image" task again, pixel-sampling via PIL (`Counter` over a strided sample, or targeted crops for small UI elements) is the reliable approach; don't guess hex values by eye.

## Theme structure: green/light are overrides on top of a dark `@theme` base

**Decision:** Kept `global.css`'s `@theme` block (the Tailwind v4 CSS-first config) as the fallback/base — unchanged, still "dark" colors — and added `:root.green` as a third override block in `variables.css`, alongside the pre-existing `:root.light`. Did not restructure the base to make green the "true" default at the CSS level.

**Why:** Minimal diff. `theme.js` already applies the correct class synchronously on load (before first meaningful paint in this client-only SPA), so the base-theme's identity doesn't matter for the effectively-perceived default — only `getInitialTheme()`'s return value does. Restructuring which theme is the CSS base would have been a larger, purely cosmetic diff for no behavioral gain.

**How to apply:** If a 4th theme is ever added, follow the same pattern: a new `:root.<name>` block in `variables.css`, a new class toggle in `theme.js`'s subscribe callback, and an addition to `THEME_ORDER`.

## Skills & Expertise icons: real brand SVGs where they exist, generic fallbacks where they don't

**Decision:** Sourced icons from [simple-icons](https://simple-icons.org/) (CC0 1.0 licensed), fetched via the jsdelivr CDN and parsed into `frontend/src/lib/techIcons.js` as `{hex, d}` path data, rendered at each tool's real brand color. For the handful of tools with no official logo in simple-icons (confirmed via `curl` HTTP-status probing, then cross-checked against Iconify's aggregated `logos`/`devicon`/`selfhst` collections — still not found): Kyverno, Kiali, and Gitleaks. For close vendor-family siblings with no *distinct* icon (Docker Compose, Grafana Loki, Grafana Tempo, GitHub CodeQL): reused the parent vendor's existing icon (Docker, Grafana ×2, GitHub) rather than inventing a new one.

**Why:** The user asked for "their icons," implying real, recognizable brand marks — not decorative placeholders — for a portfolio a hiring manager or peer engineer will actually look at. Reusing a parent-vendor icon for e.g. Loki/Tempo is honest (they *are* Grafana Labs products) rather than fabricating a fake logo for a niche tool that doesn't have one.

**How to apply:** For the 3 fully-generic icons (shield/key/mesh), the SVGs are hand-authored outline glyphs (Feather/Lucide-style, `stroke="currentColor"`) rather than brand-colored, so they visually read as "no logo available" rather than pretending to be an official mark. If Kyverno, Kiali, or Gitleaks ever publish an official simple-icons entry, swap the `generic` reference in `techIcons.js`'s `skillIcons` map for a `slug` reference instead — no other file needs to change.

## Automated code-review hook: `Stop` event, not `PostToolUse` on every edit

**Decision:** Implemented the "review code changes automatically" request as a project-level `Stop` hook (`.claude/settings.json`, committed/team-wide) that shells out to `claude -p "/code-review medium"` in the background, gated on `git diff`/`git diff --cached` being non-empty. Did **not** implement it as a `PostToolUse` hook on `Write|Edit`.

**Why:** `PostToolUse` on `Write|Edit` would fire after *every single* file write, including iterative same-file edits mid-task — far noisier and more expensive than useful for a full review agent. The user's own phrasing for the immediate request ("once this code change is done, a sub agent should trigger to review this complete code change") describes exactly the `Stop`-event semantics: review once a turn's work is finished, not per keystroke. `Stop` also matches normal human review workflow (review the finished diff, not each intermediate edit).

**How to apply:** The hook only fires when hooks docs say `agent`/`prompt` hook types are *not* available on `Stop` (only `PreToolUse`/`PostToolUse`/`PermissionRequest` support those) — hence the `command`-type hook shells out to the `claude` CLI directly in print mode (`-p`) rather than using a native `agent`-type hook.

**Recursion guard:** The spawned `claude -p` review process runs in the same project directory and would load the same `.claude/settings.json`, including this same `Stop` hook — since a read-only review doesn't clear the git diff, an unguarded version would recurse indefinitely. Guarded via a `CLAUDE_AUTO_REVIEW=1` environment variable set on the child process invocation; the hook command checks for it first and exits immediately if present. This env var is inherited by the entire child process tree, so it holds even if the review agent itself spawns further sub-processes.

**Verification:** Pipe-tested the guard logic standalone (`echo '{}' | bash -c '...'`), then ran the real command once manually (`CLAUDE_AUTO_REVIEW=1 claude -p "/code-review medium" ...`) against this session's actual diff — this doubled as both the hook-firing proof and the user's explicit request to review the completed change. Result: clean, no findings at medium effort. Output logs to `.claude/last-code-review.log`, already covered by the existing root `*.log` gitignore rule.

**Caveat for future sessions:** if `.claude/` wasn't already being watched by the settings file watcher when the current session started, the hook may need a manual `/hooks` reload or a session restart to take effect *this session* — it will apply automatically in any fresh session regardless, since it's persisted to disk.

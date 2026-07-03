# How to Resume Work

Quick reference for picking this project back up after closing the terminal/session. See `SESSION_NOTES.md` for what was done and `CLAUDE.md`/`tasks.md` for full project detail.

## Nothing is lost when you close the terminal

- Work is committed and pushed to `https://github.com/gc-ghub/gaurav-cloud-native-devops-portfolio-claude` (branch `main`).
- Always check `git status` first — there may be small uncommitted edits waiting to be reviewed and pushed.
- No dev servers or Docker containers stay running between sessions — nothing active is ever lost by closing the terminal.

## Resuming on this machine

1. Open a terminal and `cd` into the project:
   ```bash
   cd "C:\Users\gaura\OneDrive\Documents\My-DevOps-Projects\project-claude-code\project-01"
   ```
2. Check where things stand:
   ```bash
   git status
   git log --oneline -5
   ```
3. Start whichever pieces you need, in separate terminals:
   ```bash
   cd backend && go run main.go     # localhost:8080 — needs backend/.env (already present locally)
   cd frontend && npm run dev       # localhost:5173
   ```
4. Or the full stack via Docker:
   ```bash
   docker compose build && docker compose up -d   # frontend :8081, backend :8080
   docker compose down                             # when done
   ```

## Resuming with Claude Code

Just say something like **"continue from SESSION_NOTES.md"** — Claude will read `SESSION_NOTES.md`, `tasks.md`, and `CLAUDE.md` to reload full context. No need to re-explain what's been done.

## Resuming on a different machine

1. `git clone https://github.com/gc-ghub/gaurav-cloud-native-devops-portfolio-claude.git`
2. Recreate the local env files (these are gitignored, **not** in the repo):
   ```bash
   cp backend/.env.example backend/.env      # then paste in your GitHub token
   cp frontend/.env.example frontend/.env
   ```
3. **The GitHub Personal Access Token is not stored anywhere in git** — you'll need to paste it into `backend/.env` again (or generate a new classic PAT with no scopes checked, since it only reads public data).
4. Proceed with the "Resuming on this machine" steps above.

## Known caveats worth re-reading

- Never put a comment on the same line as a value in any `.env`/`.env.example` file — `godotenv` doesn't strip trailing comments, so the comment becomes part of the value. Put comments on their own line above instead.
- The Credly integration relies on an undocumented endpoint (`credly.com/users/{id}/badges.json`) — see `backend/README.md` for the risk and its graceful-degradation behavior.
- Pinned GitHub repos require `GITHUB_TOKEN` in `backend/.env`. Without it, `/api/projects` silently falls back to top-starred repos instead of erroring — worth knowing if pinned projects seem "stale" or wrong after a fresh env setup.

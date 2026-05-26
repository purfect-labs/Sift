# sift — Operator

## What I can do
- **Build & ship** — `task build` (desktop), `task build:server` (headless), `task build:docker` (container)
- **Test & validate** — `go test ./...` for Go unit tests, `evo validate --dir .` for .evolution/ integrity
- **Branch & promote** — `git checkout -b feat/<ticket>-sift` from `main`, PR to main
- **Docs & status** — Read `.evolution/` for architecture, roadmap, bugs, tracking
- **Config & deploy** — `~/.sift/.env` for API keys, `~/.sift/jobs.db` for data
- **Explore & understand** — `go doc jobdash`, `go doc jobdash/scraper`, `go doc jobdash/db`

## How to talk to me

### Building
> build jobdash for mac
> run jobdash in dev mode
> build docker image for server mode

### Fixing
> fix the SerpAPI rate limit handling in scraper/engine.go
> the match score is too low for partial keyword matches
> add pagination to the jobs dashboard

### Shipping
> bump version and tag a release
> build for all platforms
> push to main and create a PR

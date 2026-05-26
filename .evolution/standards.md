# sift — Coding Standards

## Language Conventions
- **Import ordering**: stdlib → external → internal (`sift/db`, `sift/scraper`)
- **Error handling**: Always return errors, use `fmt.Errorf` with context. Log via `service.log()` for UI visibility.
- **Naming**: PascalCase for exported types/methods, camelCase for unexported. Service methods match Wails binding convention.
- **Comments**: Document exported functions. Internal logic self-documents via clear naming.

## State Management
- **JobService struct** holds all state: config, keywords, app reference
- **SQLite** is the single source of truth for jobs, keywords, positions
- **No global variables** except `db.DB` (package-level, initialized once)
- Frontend state via Svelte stores + Wails method calls

## Naming
- **Code**: `JobService`, `ScrapeJobs`, `computeMatchScore` (verb-noun for methods)
- **Config keys**: lowercase with underscores (`resume_path`, `serpapi_key`)
- **DB columns**: lowercase with underscores (`match_score`, `applied_date`)
- **Files**: `service.go` (main logic), `engine.go` (scraper), `jobs.go` (DB)
- **Tests**: `*_test.go` alongside source

## Git & Collaboration
- **Branch naming**: `feat/<ticket>-sift`, `fix/<ticket>-sift`
- **Commit format**: `type(scope): description` — e.g., `feat(scraper): add rate limiting`
- **PR descriptions**: Link related issues, describe testing done
- **Branch tiers**: `feat/*` → `main` (single-branch workflow)

## Config Rules
- **No hardcoded values**: All API keys, paths, defaults in `~/.sift/.env` or `db/config` table
- **Environment-aware**: `SERP_API_KEY` from env file, never in source
- **Sensible defaults**: Location defaults to "Remote", queries default to `DefaultQueries()`

## Testing
- **Unit tests**: `go test ./...` for all packages
- **Integration**: Test with real SerpAPI key in dev mode
- **No E2E yet** — manual verification via `wails3 dev`

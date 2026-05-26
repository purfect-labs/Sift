# sift — Agent Directive

## 0. Bootstrap Flow
```bash
cd ~/code/plabs/jobdash/sift
wails3 dev                    # run in dev mode
# or
task build                    # production build
```

## 1. Reading Order
1. **This directive** — operating rules
2. **`.evolution/spec.md`** — architecture reference
3. **`.evolution/readme.md`** — quick start
4. **`main.go`** — entry point (59 lines, Wails v3 bootstrap)
5. **`service.go`** — core business logic (JobService)
6. **`scraper/engine.go`** — SerpAPI + Hermes AI integration
7. **`db/jobs.go`** — SQLite schema + CRUD

## 2. Project Principles
1. **SQLite-only, zero external DB** — `modernc.org/sqlite` (pure Go, no CGO). DB at `~/.sift/jobs.db`. No Postgres, no Redis.
2. **Hermes is the AI brain** — resume parsing, keyword extraction, gap analysis, market analysis, position recommendations all go through `hermes -z <prompt>` shell calls. Never replace with direct LLM API calls.
3. **SerpAPI for job sourcing** — Google Jobs via SerpAPI. Config at `~/.sift/.env` (`SERP_API_KEY`). Remote-first filtering with keyword scoring.
4. **Wails v3 service pattern** — `JobService` implements `application.Service`. Methods exposed to frontend via Wails bindings. All state lives on the service struct.
5. **Match scoring is local** — `computeMatchScore()` uses 4-level keyword matching (exact phrase → token subset → single token → substring). No AI needed for scoring.
6. **Remote-first filtering** — positive/negative keyword lists classify jobs as remote/not_remote before Hermes analysis.

## 3. How to Make Changes
- **Before**: `git checkout -b feat/<ticket>-sift` from main
- **During**: Keep `service.go` methods as Wails-bound service methods. New AI features use `hermes -z` pattern. New scrapers follow `SearchMultiple` pattern.
- **After**: `wails3 dev` to verify. `task build` for production.

## 4. Architecture Decisions
1. **Wails v3 over Electron** — native performance, single binary, no Node runtime
2. **SQLite over Postgres** — zero-config, embedded, sufficient for single-user job tracking
3. **Hermes CLI over direct API** — leverages local LLM config, no API key management in code
4. **SerpAPI over scraping** — avoids rate limiting, CAPTCHAs, legal issues
5. **Keyword matching over embeddings** — fast, explainable, no vector DB needed
6. **Python (pymupdf) for PDF only** — sole external dependency, only used in `ReadPDFText()`

## 5. Operational Reference
- **Startup**: `service.Init()` → SQLite init → load env → load keywords → ready
- **Config**: `~/.sift/.env` for `SERP_API_KEY`, `~/.sift/jobs.db` for everything else
- **Secrets**: Never commit `.env` files. API key read from home dir only.
- **Dependencies**: Go 1.25, Wails v3 alpha, modernc.org/sqlite, Python 3 (pymupdf for PDF)

## 6. .evolution Directory
```
.evolution/
├── directive.md      ← this file
├── spec.md           ← architecture spec
├── roadmap.md        ← planned work
├── readme.md         ← project overview
├── bugs.md           ← bug log
├── tracking.md       ← implementation status
├── devqueue.md       ← dashboard submissions
├── operator.md       ← agent capabilities
├── standards.md      ← coding standards
└── features/         ← detailed feature specs
```

## 7. Roadmap & Features Workflow
- Status: `idea` → `speccing` → `spec` → `in-progress` → `completed`
- Ticket format: `### X.N-shawnM Title \`status\``
- All dates ISO 8601. NEVER leave `YYYY-MM-DD`.
- Branch: `feat/<ticket>-sift` from `main`

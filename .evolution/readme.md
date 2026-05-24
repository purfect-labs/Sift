# jobdash

AI-powered job search and application tracker. Scrape jobs via SerpAPI, extract resume keywords with Hermes AI, match and filter listings, track your pipeline from discovery to offer — all in a native desktop app.

## Quick Start
```bash
cd jobdash
wails3 dev                     # development mode
task build                     # production build
task run:server                # headless server mode
```

## Project Layout
```
jobdash/
├── main.go                    # Wails v3 entry point (59 lines)
├── service.go                 # JobService — core business logic (727 lines)
├── go.mod / go.sum            # Go 1.25, Wails v3 alpha, SQLite
├── Taskfile.yml               # Build tasks (dev, build, package, docker)
├── scraper/
│   └── engine.go              # SerpAPI client, Hermes AI integration, keyword matching
├── db/
│   └── jobs.go                # SQLite schema, CRUD, keyword storage
├── frontend/
│   ├── src/
│   │   ├── App.svelte         # Root — sidebar nav (Jobs, Insights, Settings)
│   │   ├── main.js            # Vite entry
│   │   └── lib/components/
│   │       ├── Dashboard.svelte     # Job pipeline board
│   │       ├── JobCard.svelte       # Individual job display
│   │       ├── PositionManager.svelte # Job title queries
│   │       ├── ScraperPanel.svelte  # Scrape controls + logs
│   │       ├── Insights.svelte      # Keyword stats + market analysis
│   │       └── Settings.svelte      # Config (API key, resume, location)
│   ├── package.json
│   └── vite.config.js
├── build/                     # Platform build configs (darwin, linux, windows, ios, android)
└── .evolution/                # Evo project intelligence
```

## Configuration
- `~/.jobdash/.env` — `SERP_API_KEY=your_key`
- `~/.jobdash/jobs.db` — SQLite database (auto-created)

## Building
```bash
task build          # current OS
task build:server   # headless HTTP server mode
task build:docker   # Docker image
```

## Testing
```bash
go test ./...       # Go unit tests
```

## Architecture
Three-layer Wails v3 desktop app:
1. **Svelte frontend** — 6 components, sidebar nav, real-time scrape log via Wails events
2. **Go JobService** — 17 Wails-bound methods, keyword matching engine, pipeline state machine
3. **External integrations** — SerpAPI (job sourcing), Hermes CLI (AI analysis), pymupdf (PDF reading)

## .evolution
```
.evolution/
├── directive.md      # Agent instructions
├── spec.md           # Product specification
├── roadmap.md        # Planned work
├── readme.md         # This file
├── bugs.md           # Bug log
├── tracking.md       # Implementation status
├── devqueue.md       # Dashboard submissions
├── operator.md       # Agent capabilities
├── standards.md      # Coding standards
└── features/         # Feature specs
```

# JobDash

AI-powered job search and application tracker. Scrape jobs, extract resume keywords with AI, match and filter listings, track your pipeline from discovery to offer — all in a native desktop app.

Built with [Wails v3](https://v3.wails.io/) (Go + Svelte).

## Features

- **AI Resume Parsing** — Upload your PDF resume. Hermes (local LLM) extracts 40-60 keywords automatically.
- **Smart Job Scraping** — Searches Google Jobs via SerpAPI. Multi-query support with deduplication.
- **Match Scoring** — 4-level keyword matching engine scores every job 0-100 against your resume.
- **Remote Filter** — Automatic remote/hybrid/on-site classification with positive/negative keyword scoring.
- **Gap Analysis** — Per-job skill gap analysis comparing your resume against each listing.
- **Market Analysis** — Aggregated market trends across all scraped jobs.
- **Pipeline Tracking** — Full pipeline: New → Saved → Applied → Interviewing → Offer.
- **100% Local** — SQLite database, no cloud accounts, no data leaves your machine.

## Download

Pre-built binaries for macOS, Linux, and Windows from [Releases](https://github.com/purfect-labs/jobdash/releases).

## Quick Start

### Prerequisites
- **Go 1.25+** — https://go.dev/dl/
- **Node.js 22+** — https://nodejs.org/
- **Python 3 + pymupdf** — `pip install pymupdf` (for PDF resume parsing)
- **Hermes CLI** — https://github.com/nousresearch/hermes-agent (for AI analysis)
- **SerpAPI key** — https://serpapi.com (free tier available)

**One-liner:** `./dep-install.sh` installs everything except Hermes + SerpAPI key.

```bash
# Clone
git clone https://github.com/purfect-labs/jobdash.git
cd jobdash

# Install system dependencies
./dep-install.sh

# Install frontend deps
cd frontend && npm install && cd ..

# Run in dev mode
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 dev

# Build for production
task build
```

### Build from source (no Wails CLI)

```bash
cd frontend && npm install && npm run build && cd ..
CGO_ENABLED=1 go build -o jobdash .
./jobdash
```

The frontend must be built first (`npm run build`) so it's embedded in the Go binary.

## Configuration

Paste your SerpAPI key in the Settings panel, or create `~/.jobdash/.env`:
```
SERP_API_KEY=your_serpapi_key
```

All job data is stored in `~/.jobdash/jobs.db` (SQLite, auto-created on first run).

## Architecture

```
jobdash/
├── main.go              # Wails v3 entry point
├── service.go           # JobService — core business logic (17 API methods)
├── scraper/
│   └── engine.go        # SerpAPI client, Hermes AI, keyword matching
├── db/
│   └── jobs.go          # SQLite schema, CRUD, keyword storage
├── frontend/
│   └── src/
│       └── lib/components/
│           ├── Dashboard.svelte      # Job pipeline board
│           ├── JobCard.svelte        # Individual job display
│           ├── PositionManager.svelte # Job title queries
│           ├── ScraperPanel.svelte   # Scrape controls + logs
│           ├── Insights.svelte       # Keyword stats + analysis
│           └── Settings.svelte       # Config panel
├── build/               # Cross-platform build configs
└── .evolution/          # Evo project intelligence
```

## How It Works

1. **Upload Resume** → Hermes reads PDF → extracts 40-60 keywords → stored in SQLite
2. **Recommend Positions** → Hermes suggests 5-8 job titles to search for
3. **Scrape Jobs** → SerpAPI Google Jobs search → dedup via SHA256 → compute match scores
4. **Filter Remote** → Keyword-based classification → Hermes keyword analysis on matches
5. **Review Pipeline** → Dashboard shows jobs with match scores, gap analysis, status tracking
6. **Market Analysis** → Hermes analyzes top 20 jobs for trends and skill gaps

## Tech Stack

- **Backend**: Go 1.25, Wails v3 alpha
- **Frontend**: Svelte, Vite
- **Database**: SQLite (modernc.org/sqlite — pure Go, no CGO)
- **AI**: Hermes CLI (local LLM agent)
- **Job Search**: SerpAPI (Google Jobs)
- **PDF**: Python 3 + pymupdf

## License

MIT — see [LICENSE](LICENSE)

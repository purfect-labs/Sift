# jobdash — Product Specification

## Executive Summary

Jobdash is an AI-powered job search and application tracking desktop app (Wails v3, Go + Svelte). It scrapes job listings via SerpAPI, uses Hermes (local LLM) for resume keyword extraction and gap analysis, stores everything in embedded SQLite, and provides a Svelte dashboard for managing the pipeline from discovery to offer.

## Problem

Job seekers manually search across multiple platforms, copy-paste listings, and have no automated way to match their resume against job requirements. Existing tools are either browser-based (no persistence), require cloud accounts, or lack AI-powered matching.

## Architecture

```
┌────────────────────────────────────────────┐
│  Desktop App (Wails v3 — Go + Svelte)      │
│  ┌──────────┐ ┌──────────┐ ┌───────────┐  │
│  │ Dashboard│ │ Settings │ │ Insights  │  │
│  └────┬─────┘ └────┬─────┘ └─────┬─────┘  │
│       └─────────────┼─────────────┘         │
│                     ▼                       │
│  ┌──────────────────────────────────────┐  │
│  │  JobService (Go)                     │  │
│  │  - Init, GetJobs, UpdateStatus       │  │
│  │  - ScrapeJobs, FilterRemoteJobs      │  │
│  │  - ExtractResume, RecommendPositions │  │
│  │  - GetGapAnalysis, GetMarketAnalysis │  │
│  │  - computeMatchScore, analyzeBatch   │  │
│  └────┬──────────────┬──────────────────┘  │
│       │              │                      │
│       ▼              ▼                      │
│  ┌─────────┐  ┌──────────────┐             │
│  │ SQLite  │  │ Hermes (CLI) │             │
│  │ (local) │  │ - keyword ext │             │
│  │         │  │ - gap analysis│             │
│  └─────────┘  │ - market      │             │
│               │ - positions   │             │
│               └──────┬───────┘             │
│                      │                      │
│  ┌───────────────────▼───────────────────┐ │
│  │  SerpAPI (Google Jobs)                │ │
│  │  - SearchMultiple                     │ │
│  │  - keyword filtering + dedup          │ │
│  └───────────────────────────────────────┘ │
└────────────────────────────────────────────┘
```

## Configuration

`~/.jobdash/.env`:
```
SERP_API_KEY=your_serpapi_key
```

`~/.jobdash/jobs.db` — SQLite database (auto-created on first run)

## Core Concepts

### Job Pipeline
`new` → `saved` (remote match) / `not_remote` (filtered) → `applied` → `interviewing` → `offer`

### Match Scoring
`computeMatchScore()` — 4-level keyword matching against resume keywords. Level 1: exact phrase (3pts), Level 2: all tokens present (2pts), Level 3: single token match (1pt), Level 4: substring match (1pt). Returns 0-100.

### Remote Filtering
Positive keyword list (remote, wfh, distributed, etc.) scores +2. Negative list (hybrid, on-site, relocation) scores -3. Net positive → `saved`, else → `not_remote`.

### Hermes AI Integration
Four uses of `hermes -z <prompt>`:
1. **Keyword extraction** — `ExtractKeywordsFromPDF()` reads resume PDF, Hermes outputs 40-60 keywords
2. **Position recommendations** — `RecommendPositions()` reads resume, Hermes outputs job titles
3. **Gap analysis** — `GetGapAnalysis()` compares resume keywords against job listing
4. **Market analysis** — `AnalyzeJobDemand()` aggregates job samples for trends

## Data Model

### Job (SQLite)
| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | SHA256 hash of title+company+URL |
| title | TEXT | Job title |
| company | TEXT | Company name |
| location | TEXT | Job location |
| salary | TEXT | Detected salary |
| description | TEXT | Full job description |
| skills | TEXT | Comma-separated required skills |
| url | TEXT | Application URL |
| match_score | REAL | 0-100 keyword match |
| status | TEXT | Pipeline stage |
| notes | TEXT | User notes |
| gap_analysis | TEXT | Hermes gap analysis output |
| scraped_at | TEXT | ISO 8601 timestamp |

### Config
Stored in `config` table as key-value. Keys: `location`, `resume_path`. Also `positions` table for saved job titles, `resume_keywords` for extracted keywords.

## Core Flow

1. User sets resume path + SerpAPI key in Settings
2. `ExtractResume()` → Hermes reads PDF → outputs 40-60 keywords → stored in SQLite
3. `RecommendPositions()` → Hermes suggests 5-8 job titles → stored as search queries
4. `ScrapeJobs()` → SerpAPI Google Jobs search per query → dedup via SHA256 → compute match scores → store new jobs
5. `FilterRemoteJobs()` → keyword-based remote classification → Hermes keyword analysis on remote matches
6. User reviews jobs in Dashboard, updates status, views gap analysis
7. `GetMarketAnalysis()` → Hermes analyzes top 20 jobs for trends

## API / CLI

All functionality exposed via Wails v3 service bindings to Svelte frontend:

| Method | Description |
|--------|-------------|
| GetJobs(status) | List jobs by status |
| UpdateStatus(id, status) | Move job through pipeline |
| UpdateNotes(id, notes) | Add notes to job |
| GetCounts() | Job counts by status |
| SetConfig(key, value) | Update location/resume path |
| GetConfig() | Current config |
| ScrapeJobs(query) | Run job scrape |
| FilterRemoteJobs() | Classify jobs as remote/not |
| ExtractResume(path) | Extract keywords from PDF |
| RecommendPositions() | Get job title suggestions |
| GetGapAnalysis(jobID) | Analyze skill gaps |
| GetKeywordStats() | Top 100 keywords across jobs |
| GetMarketAnalysis() | Market trend analysis |
| DeleteJob(id) | Remove a job |
| ClearAllJobs() | Wipe all jobs |
| UpdateOffer(id, salary, benefits, equity) | Record offer details |
| SavePositions(positions) | Save job title queries |

## Features

- AI-powered resume keyword extraction (Hermes + pymupdf)
- Google Jobs scraping via SerpAPI
- Multi-level keyword match scoring
- Remote job classification filter
- Skill gap analysis per job listing
- Market trend analysis across all jobs
- Pipeline management (new → saved → applied → interviewing → offer)
- Offer tracking (salary, benefits, equity)
- Keyword statistics dashboard
- Position query management

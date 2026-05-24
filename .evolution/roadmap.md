# jobdash — Roadmap

## Phase A — Foundation

### A.1-shawnM Wails v3 Desktop App Bootstrap `completed` <!-- UUID: a1b2c3d4-1111-4aaa-1111-111111111111 -->
Created: 2026-05-20
Started: 2026-05-20
Completed: 2026-05-20
Assigned: @shawnM
Wails v3 app scaffold with JobService, Svelte frontend, SQLite init, window config.

### A.2-shawnM SQLite Schema + CRUD `completed` <!-- UUID: a1b2c3d4-2222-4aaa-2222-222222222222 -->
Created: 2026-05-20
Started: 2026-05-20
Completed: 2026-05-20
Assigned: @shawnM
Jobs table, resume_keywords, config, positions tables. Full CRUD operations.

### A.3-shawnM SerpAPI Integration `completed` <!-- UUID: a1b2c3d4-3333-4aaa-3333-333333333333 -->
Created: 2026-05-20
Started: 2026-05-20
Completed: 2026-05-21
Assigned: @shawnM
Google Jobs search via SerpAPI, multi-query support, dedup via SHA256 job IDs.

### A.4-shawnM Hermes AI Integration `completed` <!-- UUID: a1b2c3d4-4444-4aaa-4444-444444444444 -->
Created: 2026-05-21
Started: 2026-05-21
Completed: 2026-05-21
Assigned: @shawnM
Resume keyword extraction, position recommendations, gap analysis, market analysis via `hermes -z`.

### A.5-shawnM Keyword Match Scoring `completed` <!-- UUID: a1b2c3d4-5555-4aaa-5555-555555555555 -->
Created: 2026-05-21
Started: 2026-05-21
Completed: 2026-05-21
Assigned: @shawnM
4-level keyword matching engine with computeMatchScore. Filler word filtering.

### A.6-shawnM Remote Job Filter `completed` <!-- UUID: a1b2c3d4-6666-4aaa-6666-666666666666 -->
Created: 2026-05-21
Started: 2026-05-21
Completed: 2026-05-22
Assigned: @shawnM
Positive/negative keyword classification, auto-filter pipeline, Hermes analysis on remote matches.

## Phase B — Frontend

### B.1-shawnM Svelte Dashboard + Pipeline `completed` <!-- UUID: b1c2d3e4-1111-4bbb-1111-111111111111 -->
Created: 2026-05-22
Started: 2026-05-22
Completed: 2026-05-22
Assigned: @shawnM
Dashboard with job pipeline view, JobCard component, status management, offer tracking.

### B.2-shawnM Scraper Panel + Settings `completed` <!-- UUID: b1c2d3e4-2222-4bbb-2222-222222222222 -->
Created: 2026-05-22
Started: 2026-05-22
Completed: 2026-05-23
Assigned: @shawnM
ScraperPanel with real-time log streaming, Settings with resume upload + config.

### B.3-shawnM Insights + Analytics `completed` <!-- UUID: b1c2d3e4-3333-4bbb-3333-333333333333 -->
Created: 2026-05-23
Started: 2026-05-23
Completed: 2026-05-23
Assigned: @shawnM
Keyword stats, market analysis view, position manager with Hermes recommendations.

## Phase C — Polish & Evo

### C.1-shawnM Evo Integration `completed` <!-- UUID: c1d2e3f4-1111-4ccc-1111-111111111111 -->
Created: 2026-05-24
Started: 2026-05-24
Completed: 2026-05-24
Assigned: @shawnM
.evolution/ bootstrapped and populated. Workspace config in plabs org.

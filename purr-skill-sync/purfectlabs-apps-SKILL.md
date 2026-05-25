---
name: purfectlabs-apps
description: Catalog of PurfectLabs apps with brand identity (name, tagline, colors, icons) extracted from each app's CSS source files. For use in product demo video close screens and portfolio overlays.
category: purfectlabs
---

# PurfectLabs Apps — Brand Catalog

Use when building portfolio close screens, app switchers, or any visual that references multiple PurfectLabs apps. Each entry includes brand colors extracted directly from each app's CSS/app.css source.

## App Roster

### Undertow
- **Tagline:** Your story. Your sound.
- **What it does:** AI music portfolio studio — memory to song
- **Source:** `plabs/undertow/undertow/frontend/src/app.css`
- **Colors:**
  - Background: `#070a08` (deep abyss)
  - Surface/card bg: `#111a11`
  - Text primary: `#d4e0d4`
  - Text muted: `#6B9A80`
  - Text body: `#a0b8a0`
  - Accent blue: `#4a9eff` (phosphor)
  - Accent gold: `#c8a44a`
  - Border: `#2a3a2a`
- **Font:** Georgia/Times New Roman serif + system-ui sans-serif
- **Icon:** ♪ (music note)

### Sift (formerly JobDash)
- **Tagline:** AI job matching that never leaves your machine.
- **What it does:** AI-powered job search orchestration — resume parsing, keyword extraction, match scoring, gap analysis, pipeline tracking. Wails v3 (Go + Svelte).
- **Source:** `plabs/jobdash/jobdash/frontend/src/App.svelte`
- **Colors:**
  - Background: `#0a0c14` (dark navy), sidebar gradient `#0f1320` → `#0a0c14`
  - Sidebar border: `#1a1d2e`
  - Text primary: `#e0e0e0`
  - Text muted: `#8888aa`
  - Logo/brand gradient: `#7c8aff` → `#b47cff` (lavender to violet)
  - Active nav text: `#8b9dff` / `#b0b8ff`
  - Active nav bg: `rgba(124,138,255,0.12)` to `rgba(180,124,255,0.08)`
  - Nav hover: `rgba(124,138,255,0.08)`
  - Nav text idle: `#6b6b8a`
  - Badge: `rgba(124,138,255,0.15)` on `#8b9dff`
  - Status new: `#555`
  - Status saved: `#4fc3f7`
  - Status applied: `#ffb74d`
  - Status interviewing: `#ab47bc`
  - Status offer: `#4caf50`
  - Status rejected: `#ef5350`
- **Font:** -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif
- **Icon:** Inline SVG logo — gradient ring with sweeping curve and center dot, no emoji
- **Icon:** 🔍 (magnifying glass — sifting through jobs)
- **macOS title bar:** 38px `.title-bar-spacer` at top of sidebar to clear traffic light buttons

### Evo
- **Tagline:** Project memory for AI coding agents.
- **What it does:** Single binary giving AI tools a persistent brain. Bootstraps `.evolution/` files. Desktop dashboard with Wails. Real-time P2P agent mesh with cross-team visibility. Native GitHub UI (PRs, branches, reviews, diffs, merges), HITL/AITL queues, boardroom, file editor, token analytics per ticket/user/session/AI model.
- **Key Svelte views discovered:**
  - `ui/src/views/NeuralMesh.svelte` — D3 force-directed galaxy graph with 2-hop focus
  - `ui/src/views/meshes/MeshOmniCommand.svelte` — D3/SVG radial org hierarchy, P2P queues, HITL
  - `ui/src/views/unified/UnifiedRoadmap.svelte` — Phase planning with token audits ($state, TokenBadge, TokenDrilldownModal)
  - `ui/src/views/ArchitectureView.svelte` — Health scoring, coupling analysis, weakness detection
  - `ui/src/views/VelocityMetrics.svelte` — Chart.js phase/status tracking
  - `ui/src/views/BugTracker.svelte` — Severity-grouped bug display with Chart.js
  - `ui/src/views/WorkspaceGenesis.svelte` — 5-step wizard: connect GitHub → discover repos → assign teams → bootstrap
  - `ui/src/views/GitHubOrgGenesis.svelte` — SSH auth, org selection, repo clone, team sync
  - `ui/src/views/EditorView.svelte` — Built-in code editor with syntax highlighting
  - `ui/src/components/BranchReview.svelte` — Full git panel: branches, diff, PR tabs, merge
  - `ui/src/components/BoardroomChannel.svelte` — Chat sessions, agent IDs, pipeline viz, diff preview
  - `ui/src/components/DiffPreview.svelte` — Git diff viewer with file tree
  - `ui/src/components/PipelineViz.svelte` — Agent pipeline visualization with status symbols
- **Colors:**
- **Colors:**
  - Background deep: `#060D08`
  - Card bg: `#0C1E10`
  - Border: `#1E3A28`
  - Text main: `#E8F0EC`
  - Text dim: `#6B9A80`
  - Input bg: `#0A1A0E`
  - Accent yellow: `#FFD600` (primary brand)
  - Active accent: `#FFD600`
  - Neon green: `#00D68F`
  - (Dark theme: bg `#050505`, card `#111`, border `#1E3A28`, cyan `#00d5ff`)
- **Font:** system-ui, sans-serif
- **Icon:** 🧠 (brain)

### Bloom
- **Tagline:** Living Architecture Ecosystem. Where architecture grows like a garden 🌱
- **What it does:** Software architecture as a living ecosystem. Seeds (DNA), Blooms (implementations), Fertilizer (requirements), Chronicle (timeline), Cambium (skill library).
- **Source:** `plabs/bloom/bloom/scripts/bloom-ui/static/css/base.css`
- **Colors:**
  - bg primary: `#1a1a1a`
  - bg secondary: `#2a2a2a`
  - bg tertiary: `#3a3a3a`
  - Text primary: `#e0e0e0`
  - Text secondary: `#b0b0b0`
  - Brand pink: `#ff69b4`
  - Brand green: `#90ee90`
  - Brand dark green: `#2d5016`
  - Brand blue: `#87ceeb`
  - Brand gold: `#ffd700`
  - Error: `#ff6b6b`
- **Font:** system-ui, sans-serif
- **Icon:** 🌸 (cherry blossom)

### Thirst
- **Tagline:** (AI video/media production — GPU pipeline)
- **What it does:** Full-stack AI media production app. Scene/composition editor with GPU rendering pipeline, character library, storyboard, TTS+dubbing, tunnel mode for remote GPU orchestration.
- **Repo structure:** `plabs/thirst/thirst/` (depth 2). Note: no `.git` at repo root — thirst is managed within the plabs workspace monorepo. **Do NOT assume depth-2 .git exists.** Verify with `find ... -name ".git" -maxdepth 4 -type d`.
- **Frontend source:** `plabs/thirst/thirst/app/desktop/frontend/src/app.css` — Wails desktop UI with Svelte
- **Colors:**
  - Background: `#0d0a0b` (deep burgundy-black)
  - Surface: `#1a0f14`
  - Surface hover: `#2a1a20`
  - Border: `#3d1f2a`
  - Text primary: `#f5ede1` (warm cream)
  - Text body: `#c4908e` (rose)
  - Text muted: `#a0807e`
  - Text dim: `#6a5a58`
  - Accent primary: `#b0585a` (rust rose)
  - Accent hover: `#c47072`
  - Danger: `#7a1e3a`
  - Danger hover: `#9e2648`
  - Badge green: `#b0585a on #7a1e3a33`
  - Badge yellow: `#e8d5b7 on #c8890e33`
  - Tier free: `#a0807e`
  - Tier individual: `#c4908e`
  - Tier creator: `#e8d5b7`
  - Tier studio: `#b0585a`
  - Dotfiles green: `#5a8a5a`
  - Progress: `#7a1e3a`
  - Terminal amber: `#d4a574`
- **Font:** -apple-system, BlinkMacSystemFont, Segoe UI, sans-serif + JetBrains Mono for terminal
- **Icon:** 💧 (droplet) or 🎬 (clapper)

### Organisync
- **Tagline:** (Infrastructure/DevOps sync — licensing daemon)
- **What it does:** Go-based cross-platform sync daemon with license-key management, dashboard UI via Wails.
- **Source:** `plabs/organisync/organisync/web/src/app.css` + inline HTML from lambda
- **Colors:**
  - Background: `#0f1117`
  - Surface: `#1a1d27`
  - Surface hover: `#22253a`
  - (licensing page uses `#0a0a0f`)
  - Accent purple: `#6c47ff`
  - Accent cyan: `#00d4ff`
  - Text primary: `#e2e8f0`
  - Text muted: `#a0aec0`
  - Code green: `#7ee787`
  - Border: `#2d2d3d`
- **Font:** system-ui, sans-serif
- **Icon:** 🔄 (sync)

### Apex
- **Status:** Empty scaffold — no live app yet
- **Colors:** Fallback `#0a0c14`
- **Icon:** 🚀 (rocket)

### Phantom Veil
- **Status:** Empty scaffold — no live app yet
- **Colors:** Fallback `#0a0c14`
- **Icon:** 👻 (phantom)

## Usage: Video Close Screen

When building a product demo video close screen, load `purfectlabs/purfectlabs-apps` skill first, then render a "Developed by PurfectLabs" brand block with horizontal app cards below:

```
Scene layout (1920×1080 canvas):

  Your life in albums.              ← tagline (gold, serif, 28px)
  Share it with the people...       ← CTA text (muted, 18px)
  ─────────────────                 ← 1px accent line

  DEVELOPED BY                     ← muted, 10px, #6B9A80
  PURFECTLABS                      ← 34px, uppercase, system-ui, #d4e0d4

  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
  │  ♪   │ │  🔍  │ │  🧠  │ │  🌸  │ │  💧  │ │  🔄  │
  │Undertow│ │ Sift │ │ Evo  │ │ Bloom│ │Thirst│ │Orgnsnc│
  │#4a9eff│ │#7c8aff││#FFD600││#ff69b4││#b0585a││#6c47ff│
  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘
```

Each app card: 2px top border in the app's accent color, dark surface bg, app icon above name. Stagger animate in with `back.out(1.5)` and ±15° rotationY. At 7s, cards + tagline fade out leaving just "DEVELOPED BY PurfectLabs" for final 3s.

Empty-status apps (Apex, Phantom Veil) rendered at lower opacity (0.4) with muted fallback border.

## Download CTA Pattern

Every video VO ends with this exact closing: `[measured] [firm] Developed by PurfectLabs. Download [AppName] now.`

The VO naturally ends ~53-57s into the video depending on length. Leave 2s of silence in the audio bed after the VO ends. This is app-agnostic — only `[AppName]` changes per video.

## Loading Convention

The `product-demo-video` skill loads `purfectlabs/purfectlabs-apps` as a related skill. Before building any close screen, call `skill_view('purfectlabs/purfectlabs-apps')` to get the brand data.

## Repo Structure Verification (CRITICAL — Do NOT Pattern-Match)

When extracting brand data from plabs apps, verify the actual repo structure BEFORE reading files:

1. **Check for .git:** `find /Users/shawn/code/plabs/<app>/ -name ".git" -maxdepth 4 -type d`
2. **Not all depth-2 dirs have .git.** Thirst (`plabs/thirst/thirst/`) has no .git — it's managed within the plabs workspace.
3. **After finding the repo root, find frontend source:** search for `app.css`, `*.svelte`, `*.css` files excluding `node_modules`
4. **Read actual CSS values** — do NOT infer or guess colors from readme emojis or vague descriptions. Each app has a real CSS/design token file with exact hex codes.
5. **If no frontend source exists** (Apex, Phantom Veil), mark as "Empty scaffold" with fallback colors. Do NOT fabricate brand colors.

Source files found in this session for reference:
- JobDash: `plabs/jobdash/jobdash/frontend/src/App.svelte` (CSS inline)
- JobDash: `plabs/jobdash/jobdash/frontend/src/lib/components/JobCard.svelte`
- Evo: `plabs/evo-app/evo/ui/src/app.css`
- Bloom: `plabs/bloom/bloom/scripts/bloom-ui/static/css/base.css`
- Thirst: `plabs/thirst/thirst/app/desktop/frontend/src/app.css`
- Organisync: `plabs/organisync/organisync/web/src/app.css`

---
name: product-demo-video
description: Create product demo videos for PurfectLabs apps — from reading the spec to delivering the MP4 with audio. Uses HyperFrames for rendering, Venice stable-audio-25 for music, Fish Audio/Venice for voiceover, and ffmpeg for audio mixing.
category: purfectlabs
tags: [video, demo, hyperframes, product, marketing, venice-music, ffmpeg]
---

# Product Demo Video — PurfectLabs Workflow

Create emotional, cinematic demo videos for PurfectLabs desktop apps using HyperFrames (HTML + GSAP → MP4). Covers the whole pipeline: narrative arc → visual identity from app CSS → HyperFrames composition → Venice music bed → ffmpeg audio mix → voiceover → delivery.

Requires the optional `hyperframes` skill for CLI commands. Load that skill separately for reference.

## When to Use

- User asks for a demo video, promotional video, or product tour
- User says "make a video for X" where X is a PurfectLabs project
- You've been given creative freedom and an app spec to work from

## The Workflow

### Phase 1: Orient — Read the Product

1. **Find the repo root** — search `~/code/plabs/<ws>/` for the product directory
2. **Check for `spec.md`** — if present, read it for product purpose, audience, pain point, and data model
3. **Fallback to `README.md`** — if no `spec.md`, the README is your primary source. Extract: what does it solve, key features, CLI reference, architecture.
4. **Read the frontend CSS** (`app.css` or equivalent) — extract design tokens:
   - Background colors (these become your scene backgrounds)
   - Text colors
   - Accent colors (phosphor blue, amber, etc.)
   - Font families
5. **Define the audience** from the product's problem statement or executive summary — is this for power users, casual creatives, agents, infrastructure teams? The video's tone depends on this.

### Phase 1b: Existing Gallery Shortcut — Skip Capture When Screenshots Already Exist

When screenshots already exist in the repo (e.g., `marketing/gallery/img/` or `marketing/gallery/`) and the user says "just pick a few" or "no image pipeline needed":

1. **Locate the gallery directory** — typically `plabs/<ws>/<name>/marketing/gallery/img/` or similar
2. **List all screenshots** with file sizes — note any large (>1MB) files as they may be rich/informative views
3. **Pick 3-5 visually compelling screenshots** based on filenames that suggest distinct, important features. Prefer variety over quantity.
4. **Do NOT run the full Playwright capture** or the one-at-a-time Y/N selection loop — the user explicitly opted out
5. **Do NOT run Venice vision evaluation** unless the user asks or there's uncertainty about what a screenshot shows
6. **Show the picks upfront** via MEDIA tags in a single message so the user can confirm
7. **If user asks to swap one**, do NOT re-send all — just swap the file and re-send that one
8. Copy confirmed screenshots to the HyperFrames assets dir.

### Phase 2: Concept — Find the Narrative Arc

Ask three questions before writing any code:

1. **What does this product DO?** (functional truth)
2. **What does this product FEEL like?** (emotional truth — this is what the video sells)
3. **Who is this for, and why do they care?** (audience hook — MUST name every audience tier)

#### Duration Guide — Match Video Length to App Complexity

| App type | Suggested length | Scene count | SS:Card ratio |
|----------|-----------------|-------------|--------------|
| Simple single-feature app (Sift, Undertow) | 30-50s | 5-7 scenes | 60:40 |
| Complex multi-feature app (Evo with mesh) | 120-180s | 14-18 scenes | 50:50 |

**Hook rules:**
- Must call out specific audience levels — C-level, project managers, devs, QA — by name in the first 8 seconds.
- Use role pills (badge-style `<span>` elements) showing each audience tier visually.

#### Duration Calibration — VO Must End 5-8s Before Video Ends

The VO script MUST finish 5-8s before the total video duration, leaving pure music + brand visuals for the close.

**How to achieve this:** Write the VO script first, estimate length (~130 words per 30s at speed 0.88), then verify by generating and measuring with ffprobe.

### Phase 3: Visual Identity — Write DESIGN.md

Use the app's actual design tokens from CSS.

### Phase 4: Compose — Write the HTML

Each scene is a `<div class="scene clip">` with `data-start`, `data-duration`, `data-track-index`. Images need `class="clip"` too. All timed elements require `clip`.

#### Phase 4b: Full-Frame Screenshots

**CRITICAL — Use `object-fit: contain` with `object-position: top`, NOT `object-fit: cover`.** Desktop app screenshots are typically 2048x1320+ — taller than the 1920x1080 canvas — so `cover` zooms in and cuts off the top chrome, sidebar, and navigation. Always add BOTH properties:

```css
.bg-image { position: absolute; top: 0; left: 0; width: 100%; height: 100%;
  object-fit: contain; object-position: top; opacity: 0; background: #0a0c14; }
```

#### Phase 4c: Full-Screen Grids for Timeline Views

Do not put timelines in a thin left column. Use a 2x2 CSS grid filling the frame.

#### Phase 4d: Motion Design — GSAP Easing Variety

Vary entrances by content type — never the same animation for adjacent scenes.

#### Phase 4e: Close Screen — "Developed by PurfectLabs" + App Portfolio + Download CTA

**App-agnostic.** Every PurfectLabs video uses the exact same close signature. Only `[AppName]` changes in the VO.

**Close scene duration must be proportional to total video length.** A 36-40s video should have ~6-8s of close. For a 180s video, ~14-18s.

**When user says close is too long, redistribute saved time to screenshot-heavy scenes** (dashboard → scraper → pipeline), in that priority order. Each gets +1-2s more hold time before its text overlay fades out.

**Timing (simple app, 6-8s close, Screenshot + Cards variant):**
1. App dashboard screenshot fades in (0.2s) with gentle slow zoom (scale:1.03, 3s)
2. At ~2.0s: **screenshot fades out** (opacity→0 over 0.6s) simultaneously as GitHub row fades in — dissolve, not static coexistence
3. "DEVELOPED BY PurfectLabs" fades in (2.8s)
4. 6 app cards stagger in (3.2-4.45s, 0.25s gap, ±15° rotationY)

### Phase 5: Validate and Render

```bash
cd <product>-video/
npx hyperframes render --video-bitrate=2M --output demo.mp4
```

### Phase 6: Generate Audio — Music Bed

Default: `stable-audio-25` ($0.24 flat up to 190s). Returns WAV, always convert to MP3.

### Phase 7: Mix Audio Into Video

**Music volume:** Start at `volume=0.15` for the music track. Shawn corrected that `volume=0.28` drowned out the VO for dark-background cinematic videos.

**Purra Heartbeat overlay — EVERY video.** Every PurfectLabs product demo MUST include the Purra heartbeat audio signature 7 seconds before the video ends, on top of all other audio. Audio file: `sift-demo/purra-heartbeat.mp3` (~5.8s). Adjust adelay: `(video_duration - 7) * 1000`ms.

```bash
ffmpeg -y -i video.mp4 -i purra-heartbeat.mp3 \
  -filter_complex "[1:a]adelay=33000|33000[hb];[0:a][hb]amix=inputs=2:duration=first:dropout_transition=0[a]" \
  -map 0:v -map "[a]" -c:v copy -c:a aac -b:a 192k output.mp4
```

**Duration calibration — NO CLIPPING:**
- VO MUST be shorter than video duration. Extend video, not the ffmpeg command.
- Never use `-shortest` — it clips the longest stream.
- Pad music bed to exceed video by at least 5s.

### Phase 8: Voiceover Integration

**"Purr's voice" convention:** When Shawn says "use purrs one," read `FISH_AUDIO_VOICE` from the profile `.env` — this IS Purr's voice. Do NOT assume Philosopher or any hardcoded ID; the `.env` is the authority.

**VO closing line:** Always use `[measured] Developed by PurfectLabs. Download [AppName] now.` — the "tonight" soft-CTA was rejected.

### Phase 9: Manifest

Create `manifest-v<N>.txt` after every render linking screenshots, music prompt, voice ID+params, composition structure, and render chain.

### Phase 10: POC Audio Mix First

Stitch music + voiceover without video to iterate levels. See `references/audio-poc-workflow.md`.

### Phase 11-12: Stitch

Render video once, then stitch with audio layers. See `references/ffmpeg-audio-mix.md`.

### Phase 13: Deliver

Push to repo on a shared branch (e.g., `dev`). Upload to Slack via inline `MEDIA:/path/to/file.mp4` in response text — NOT `send_message()`.

**Clean commit rule:** Commit ONLY source files + final MP4. Skip `work-*`, `renders/*`, stale audio files.

**Git operation: keep it SIMPLE.** When the user says "convert this branch to dev and add the file," just rename the local branch (`git branch -m oldname dev`), add the files, commit, push. Do NOT create temp branches, delete remote branches, force push, or restructure history.

### Phase 13b: Audience Variant Pipeline

When the user asks for variants targeting different demographics — fork the composition, change VO script and music bed only (same GSAP structure, same screenshots).

**CRITICAL — do NOT delegate audience HTML writing to subagents.** Shawn corrected this twice. The index.html fork must be done in the main thread with verified find-and-replace.

### Phase 13c: Rebrand/Rename — Full Product Rename

When the product has been renamed (e.g., JobDash → Sift), update all files in the demo directory. See `references/rebrand-checks.md` for the file inventory.

### Phase 13e: Deliver to Another Agent (Petal / Purr)

**1. Determine what format they need.** The user corrected this twice — "it's audio not video." Petal stitches audio overlays (brand signature, jingle, VO, heartbeat) into the final video. Do NOT assume they need the rendered MP4.

**2. Check Desktop if Slack attachment delivery fails.** The Slack gateway lacks scopes to download user-attached files. The user's workaround is to save files to the Desktop. Always check:
```bash
ls ~/Desktop/Heart*  # brand audio signatures
ls ~/Desktop/<product-name>/
```

**3. Commit audio files to the shared branch** (typically `dev`):
```bash
cp /path/to/audio.mp3 sift-demo/purra-heartbeat.mp3
git add sift-demo/purra-heartbeat.mp3
git commit -m "add Purra heartbeat audio signature"
git push origin dev
```

### Phase 14: Trace the Latest Deliverable (Post-Session Retrieval)

When someone asks for "the latest" or "the last demo video," search past sessions with `session_search(query="<product> demo video v<N> deliver")`. Do NOT guess by file timestamp — the latest file on disk may be a draft or intermediate render.

## Pitfalls

- **Close screen screenshot must fade OUT as cards/links fade IN** — do NOT let them coexist static. Dissolve transition around 2.0s into close scene.
- **Git operation: keep it SIMPLE** — just rename, add, commit, push. No branch restructuring.
- **"It's audio not video"** — when delivering to Petal, check format first. Audio files get committed; video files get Slack-delivered.
- **Slack voice messages cannot be downloaded** — gateway lacks scopes. Check `~/Desktop/` for files.
- **VO duration MUST be shorter than video duration** — extend video, not the ffmpeg command.
- **Never use `-shortest`** in the ffmpeg mix command.
- **Music volume: start at 0.15, not 0.28**.
- **Close screen VO: "Developed by PurfectLabs. Download [AppName] now."** — the "tonight" CTA was rejected.
- **Close duration redistribution** — stretch screenshot scenes, not the hook.
- **Never use `-c:v copy` for audio mixing.**
- **Fonts must be system fonts** — external requests fail in headless Chrome.
- **HyperFrames CLI at `~/.hermes/node/bin/hyperframes`.**
- **Do NOT use `$FISH_AUDIO_VOICE` from `.env` for project VO** — check manifest first.
- **Text-only models cannot do vision** — use MEDIA tags instead.

## Related Skills

- `hyperframes` (optional) — CLI commands, GSAP reference
- `purfectlabs/evo-plabs-flow` — repo structure
- `venice/venice-audio-music` — music generation
- `purfectlabs/fish-audio` — TTS, voice discovery
- `purfectlabs/purfectlabs-apps` — brand catalog
- `purfectlabs/slack-audio` — Slack delivery patterns

## References

- `references/screenshot-workflow.md` — Playwright capture, Venice evaluation, Y/N selection
- `references/interleaved-scene-pattern.md` — Complex multi-feature app scene template
- `references/venice-vision-evaluation.md` — batch screenshot scoring
- `references/undertow-composition-example.md` — GSAP patterns
- `references/master-manifest-template.md` — cross-variant master manifest template
- `references/close-screen-template.md` — universal close screen HTML/CSS/GSAP
- `references/audio-poc-workflow.md` — mix-before-render
- `references/ffmpeg-audio-mix.md` — audio muxing reference
- `references/user-screenshot-integration.md` — screenshot integration workflow
- `references/venice-music-flow.md` — music generation reference
- `references/suno-music-alternative.md` — Suno alternative

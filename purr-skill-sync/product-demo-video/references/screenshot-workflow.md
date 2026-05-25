# Screenshot Workflow — Using Real App Screenshots in Video

Capturing real app screenshots and embedding them into the HyperFrames composition. This replaces generic text cards with actual product visuals for more authentic demo videos.

## Workflow

### 0. Prerequisites — Install Playwright

The `browser_vision` tool captures static page screenshots but can't navigate SPA views (click sidebar nav, open dialogs, etc.). For full SPA navigation, use Playwright through Python:

```bash
# One-time setup
pip3 install playwright
~/.hermes/profiles/purfectlabs-tg-bot/home/Library/Python/3.9/bin/playwright install chromium
```

If Playwright path isn't on PATH, install via the full path or use `sys.path.insert()` in Python.

### 1. Run the App

Most PurfectLabs apps are Wails v2/v3 Go apps with a Svelte/Vite frontend. If `wails dev` isn't available (no Wails CLI), run the frontend dev server:

```bash
cd frontend/
npm run dev -- --host
```

Start as a background Hermes process so it stays alive during capture:

```bash
terminal(background=true, command="cd frontend/ && npm run dev -- --host", watch_patterns=["Local:"])
```

### 2. Navigate and Capture All Views

**For static pages** that render everything on load (welcomes, landing): use the Hermes `browser_navigate` + `browser_vision` combo. Even when vision fails with `unknown variant image_url`, the screenshot IS captured — check `screenshot_path` in the error response.

**For SPAs with navigation** (sidebars, clicks to reveal views): use Playwright via Python to click through every view and screenshot each one.

#### Full Playwright Screenshot Pattern

```python
import subprocess

script = '''
import sys, os
sys.path.insert(0, "/Users/shawn/.hermes/profiles/purfectlabs-tg-bot/home/Library/Python/3.9/lib/python3.9/site-packages")

from playwright.sync_api import sync_playwright

url = "http://localhost:5173"
outdir = "/Users/shawn/code/plabs/undertow/undertow/undertow-video/assets"

with sync_playwright() as pw:
    browser = pw.chromium.launch(headless=True)
    page = browser.new_page(viewport={"width": 1400, "height": 900})
    page.goto(url, wait_until="networkidle")
    page.wait_for_timeout(2000)

    # 1. Full page welcome
    page.screenshot(path=f"{outdir}/app-full-welcome.png", full_page=True)
    print("Welcome - saved")

    # 2. Click through sidebar nav — each view
    # Each button text that triggers a view navigation
    nav_buttons = [
        ("MemoryDump", "Memories"),
        ("SongExtract", "Extract"),
        ("SongBuilder", "Songs"),
        ("PlatformOutput", "Output"),
        ("PortfolioExport", "Export"),
    ]
    
    for view_name, button_text in nav_buttons:
        btn = page.get_by_text(button_text, exact=False).first
        if btn.count() > 0:
            btn.click()
            page.wait_for_timeout(1500)
            page.screenshot(path=f"{outdir}/app-view-{view_name}.png", full_page=True)
            print(f"{view_name} - saved")

    # 3. Click specific nav items (Home, My Sound, Project, LLM Settings)
    more_nav = [
        ("ErasTimeline", "Home"),
        ("StyleProfile", "My Sound"),
        ("ProjectViewer", "Project"),
        ("LLMSettings", "LLM Settings"),
    ]
    for view_name, button_text in more_nav:
        btn = page.get_by_text(button_text, exact=False).first
        if btn.count() > 0:
            btn.click()
            page.wait_for_timeout(1500)
            page.screenshot(path=f"{outdir}/app-view-{view_name}.png", full_page=True)
            print(f"{view_name} - saved")

    browser.close()
    print("\\nALL SCREENSHOTS CAPTURED")
'''

result = subprocess.run(
    ["/usr/bin/python3", "-c", script],
    capture_output=True, text=True, timeout=60
)
print(result.stdout)
# Check stderr for Playwright errors (missing elements, timeout)
```

#### Checking Which Elements Are Clickable

Before writing the navigation script, inspect the page to know what's actually rendered:

```javascript
// In browser_console:
document.querySelector('.flow-cell') ? 'found' : 'NOT found'    // Check specific elements
document.querySelectorAll('nav a, nav button, .sidebar *').length  // Count nav targets
```

The Snapshot from `browser_navigate` also lists all interactive elements with ref IDs.

#### Dealing with Svelte/React SPAs

SPAs may use Svelte `{#if $activeView === 'Welcome'}` conditionals — clicking a button changes the store, which renders a COMPLETELY different subtree. The DOM element may vanish or be replaced. Always `wait_for_timeout(1500)` after clicking to let the SPA re-render.

If a button click triggers a Go backend call that fails (no Wails runtime), handle the error silently — the screenshot still captures the view with its fallback state.

### 3. Share All Screenshots for Validation (CRITICAL — Most Error-Prone Step)

Before baking any screenshot into the video composition, share them with the user for approval. This prevents wasting a render cycle on bad screenshots.

**CRITICAL RULES — Violating any of these caused a 20-round back-and-forth this session:**

1. **Never pre-select or curate.** Do NOT pick 3-4 screenshots you think are best and present those as recommendations. The user knows the story they want to tell. Show ALL distinct screenshots.
2. **Link every screenshot by filename with inline MEDIA tags.** Do not say "here are the remaining ones" without linking. Always link.
3. **Present one at a time with explicit Y/N.** Do NOT batch-link 4+ images and ask "pick your favorites." This causes filename↔image mapping errors in Slack.
4. **Wait for user response before sending the next.** Each screenshot is a separate turn: link file → ask Y/N → wait for answer → move to next.
5. **If user rejects a numbered slot** (e.g. "4 is wrong"), do NOT guess. You will guess wrong repeatedly. Either link ALL remaining candidates, or continue the one-at-a-time flow from the remaining pool.
6. **Keep a running list of confirmed filenames** so you don't lose track mid-session. Write them down in your response as they're approved.
7. **After the user confirms their picks, clean up.** Remove stale duplicates, test screenshots, and unselected files from the assets dir. Only confirmed picks + full-page reference should remain.

Selection flow template:
```
**1.** app-full-welcome.png
MEDIA:/path/to/app-full-welcome.png
Y or N?
```
(Wait. Confirmed: Y)

```
**2.** app-view-StyleProfile.png
MEDIA:/path/to/app-view-StyleProfile.png
Y or N?
```
(Wait. Confirmed: Y)

Continue until all slots are filled. If user rejects a screen, present the next candidate from the remaining pool.

For initial orientation only: you can batch all screenshots in a single message with numbered filenames — but do NOT ask for picks in that message. Just orient.

After all confirmed picks are recorded, proceed to the composition step.

### 4. Evaluate All Screenshots via Venice Vision

Before letting the user pick or baking into the composition, evaluate every screenshot through Venice vision. Do NOT guess which ones are best — run the batch evaluation.

**Why:** The screenshots are dark-themed Wails/Svelte apps on dark backgrounds. What looks "similar" to human eyes may be very different to a vision model that actually reads the text. Let the model score them objectively.

**Model choice:** `qwen3-vl-235b-a22b` at $0.25/M input tokens. A batch of 12 screenshots costs ~$0.002. Cheaper models like `qwen3-5-9b` ($0.10/M) may not handle dark-themed app UIs as well.

**Pricing quirk:** The `/models` endpoint returns empty `pricing: {}` for qwen3-vl — but actual billing is $0.25/M input, $1.50/M output. Check usage via `GET /billing/usage` if you need confirmation.

**Capabilities quirk:** All Venice models report `capabilities.vision: false` even when they accept images. The real signals are `supportsMultipleImages: true` and `maxImages >= 10`. Ignore the `vision` field — any model with `multiImg` support will process `data:image/png;base64,...` URLs.

#### Batch evaluation — send all screenshots in one call with filename correlation

Send all screenshots in a single API call with filenames embedded in the prompt so the response can reference each one:

```python
import subprocess, json, base64

r = subprocess.run(
    ["grep", "VENICE_API_KEY", "/Users/shawn/.hermes/profiles/purfectlabs-tg-bot/.env"],
    capture_output=True, text=True
)
api_key = r.stdout.strip().split("=", 1)[1]

# Build the list: (filename, absolute_path)
images = [
    ("app-full-welcome.png", "/path/to/app-full-welcome.png"),
    ("app-view-StyleProfile.png", "/path/to/app-view-StyleProfile.png"),
    # ... all screenshots
]

# One content array with the instruction as the first part, then all images
content_parts = [{
    "type": "text",
    "text": """I'm sending you N screenshots of a dark-themed desktop app called [APP NAME].
Each image is labeled by filename in order (first = image 1, second = image 2, etc.).

For EACH screenshot, return:
- filename: the image filename (you must correlate based on image order)
- description: 1 sentence what this screen shows
- clarity: 1-5 (text readable, not cropped/blurry, dark elements visible)
- visual_impact: 1-5 (looks cinematic, impressive to show in a video)
- uniqueness: 1-5 (visually distinct from the other screens — not same layout with different text)
- story_value: 1-5 (does it communicate a step in the user journey?)

Then give me:
- recommended_order: which filenames to use in the video and in what sequence
- reasoning: 1-2 sentences why this order tells the best story"""
}]

for name, path in images:
    with open(path, 'rb') as f:
        b64 = base64.b64encode(f.read()).decode()
    content_parts.append({
        "type": "image_url",
        "image_url": {"url": f"data:image/png;base64,{b64}"}
    })

payload = json.dumps({
    "model": "qwen3-vl-235b-a22b",
    "max_completion_tokens": 2000,
    "messages": [{"role": "user", "content": content_parts}]
})

result = subprocess.run([
    "curl", "-s", "--max-time", "120",
    "https://api.venice.ai/api/v1/chat/completions",
    "-H", f"Authorization: Bearer {api_key}",
    "-H", "Content-Type: application/json",
    "-d", payload
], capture_output=True, text=True, timeout=130)

data = json.loads(result.stdout)
print(data['choices'][0]['message']['content'])

# Also print usage for cost tracking
usage = data.get('usage', {})
total_in = usage.get('prompt_tokens', 0)
cost = (total_in / 1000000) * 0.25
print(f"\\nEstimated cost: ${cost:.6f}")
```

#### Reading the results

The model returns per-screenshot scores (each 4-20) and a recommended order. Send the results + all screenshots to the user with inline MEDIA tags and let them make the final call. **Do not override the user's selection with the model's ranking** — use it as a recommendation, not a decision.

### 4A. Selection Protocol — One-At-A-Time, No Guessing

**This is the most error-prone step in the whole workflow. Follow these rules strictly.**

#### Rules:
1. **Present one screenshot at a time.** Link the file, ask "Y or N?", wait for response, then move to the next.
2. **Do NOT batch-link and say "pick your favorites"** — this causes ambiguity over which filename maps to which image, especially in Slack where MEDIA tags don't show filenames.
3. **When the user rejects a specific slot** (e.g. "4 is wrong"), do NOT guess which remaining screenshot they want. Either:
   - Link all remaining candidates and let them pick, OR
   - Continue presenting one at a time from the remaining pool.
4. **Keep a running list** of confirmed picks by filename so you don't lose track mid-session.
5. **When the user says "link them all" — do that immediately.** Do not curate, do not filter, do not explain why some are better than others. Link every distinct screenshot by filename in order. They're choosing based on visual content, not your opinion.
6. **Do not try to guess the "right" screenshots** based on Venice rankings or your own judgment. The user knows the story they want to tell. Present options, don't decide for them.
7. **After selection is done, clean up.** Remove duplicate captures, stale test screenshots, and anything in `assets/` that wasn't selected before moving to composition. Keep only the confirmed picks + the full-page reference.

#### Common failure modes:

| What went wrong this session | Fix |
|------------------------------|-----|
| Batched 4 screenshots, said "here are my picks" | Don't pre-select. Show everything, let user decide. |
| Guessed which screenshot was #4, got it wrong 6 times | Don't guess. Link candidates one at a time. |
| Said "the remaining ones" without linking | Always link when listing options. |
| Kept stale duplicate screenshots in assets dir | `rm -f` duplicates after capture cycle. Keep only clean. |

#### Evaluation criteria

The scoring criteria in the prompt map to video production concerns:

| Metric | What it catches |
|--------|----------------|
| **clarity** | Blurry captures, text cropped off-screen, dark-on-dark elements that blend into background |
| **visual_impact** | Dull/empty views, too much whitespace, purely text-based screens |
| **uniqueness** | Same layout with different text (e.g., two empty list views that look identical) |
| **story_value** | Screens that don't advance the narrative (settings screens, empty states) |

#### When to skip Venice evaluation

If the user explicitly says "just show me all of them and I'll pick" (like in this session), skip the model evaluation and go straight to sharing via MEDIA tags. The evaluation is for when you need automated triage on a large set.

#### Cost note

At qwen3-vl pricing ($0.25/M input), a batch of 12 screenshots costs under $0.002. There is no reason not to run it.

### 5. Promote Accepted Screenshots — `screenshots-final/`

After the user confirms which screenshots to use via Y/N selection:

1. **Create a `screenshots-final/` directory** inside `assets/`
2. **Copy only the confirmed screenshots** there — this is the canonical set for the video
3. **Reference `assets/screenshots-final/` in the composition HTML**, not the raw capture filenames

```bash
mkdir -p assets/screenshots-final/
cp assets/app-full-welcome.png assets/screenshots-final/
cp assets/app-view-StyleProfile.png assets/screenshots-final/
# ... etc for each confirmed screenshot
```

This way the raw captures stay in `assets/` for reference and the promoted set is clean, small, and versioned.

**Never bake raw capture filenames into the composition** — always promote first. If the user swaps a screenshot later, only the `screenshots-final/` copy changes, not the whole composition.

### 6. Alternate Screenshots with Cards

**Do NOT stack all screenshots in a block.** Alternate them with memory cards or text-based scenes for pacing: screenshot → cards → screenshot → cards → screenshot → cards.

This creates visual rhythm and prevents viewer fatigue from 4+ static screenshots in a row. Each screenshot + card pair tells one narrative beat:

| Beat | SS shows | Cards show |
|------|----------|------------|
| Identity | Welcome / hero page | Childhood memories |
| Personalize | Style Profile / genre setup | Teenage / milestone cards |
| Structure | Era Wizard / timeline | Tagged memory cards |
| Configure | LLM Settings / export | Soundtrack / people cards |

**Pattern per pair:**
- SS scene: 3.5-4s (headline fade → screenshot ~2.5s → fade out)
- Card scene: 3-3.5s (headline fade → 2 staggered cards → fade out)
- Tight gap (0-0.5s) between SS and cards so they feel like one beat

### 7. Audio Scenes Should Align with Visual Beats

Keep video scenes short (3-4s per screenshot/cards) so the VO has natural phrasing breaks. A 50s VO pairs well with a 51s video where scenes are 3-4s each — the VO pauses hit between scene transitions.

### 8. Manifest-Driven Pipeline — Track Everything

After every render, create a `manifest-v<N>.txt` in `assets/` that links ALL production artifacts. This is the source of truth — you shouldn't need to read Slack history to reproduce a render:

```txt
# Project Demo Video v<N> — Render Manifest
# Generated: YYYY-MM-DD
# Duration: Xs
# Output: product-demo-v<N>.mp4

## Screenshots
1. filename.png | one-line description of what it shows
2. filename.png | one-line description

## Music
- Source: variant-X.mp3 (Venice stable-audio-25)
- Extended to Xs via aloop + atrim
- Output: assets/audio/music-Xs.mp3
- Prompt: "The Venice generation prompt for reference"

## Voiceover
- Voice: Name (voice_id)
- Tags: voice tags from Fish API
- Temp: X.X, Speed: X.XX
- Emotion tags: [tag1] [tag2] [tag3]
- Script file: assets/audio/vo-script.txt
- Output: assets/audio/vo-<voice-name>-full.mp3

## Composition
- index.html — alternating SS/cards pattern
- N frames at 30fps
- Scene breakdown

## Render chain
1. npx hyperframes render → video-only.mp4
2. ffmpeg mix: music volume, VO adelay, amix
3. Output: product-demo-v<N>.mp4

## Prior versions
- v<N-1>: [what changed]
- v<N-2>: [what changed]
```

**Key rules:**
- Update the manifest with every render, even incremental changes
- Include the Venice music prompt — this lets you regenerate the exact same music if needed
- Include the voice ID + tags + temperature + speed — one-to-one reproducibility
- Include the paralanguage emotion tags used in the VO script
- Version numbers track across the whole project, not per session

### 9. Copy to Assets

```bash
cp <screenshot_path> <video-dir>/assets/app-<description>.png
```

Naming convention: `app-<view-description>.png` (e.g. `app-welcome-flow.png`, `app-hero.png`).

### 10. Structure Screenshots as Separate Scenes — Full-Frame Sizing

**CRITICAL: Each screenshot gets its own scene.** Do NOT stack a screenshot and memory cards in the same scene — HTML layout can't handle both without overlap at 1920x1080.

**Sizing rule: screenshots must fill most of the frame.** Users complained about "too much black dead space" with the old 1300px hero sizing. Always use 1800px for the image on a 1920px canvas.

Pattern per screenshot (full-frame variant):

```html
<div id="s2-ss" class="scene clip" data-start="3.5" data-duration="3.5" data-track-index="2"
     style="flex-direction:column;justify-content:flex-start;padding-top:10px;">
  <img class="ss-full clip" id="ss-1" data-start="3.8" data-duration="2.8" data-track-index="3"
       src="assets/screenshots-final/app-screenshot-1.png" alt="" />
  <div class="ss-label"><span>Your headline text here.</span></div>
</div>
```

CSS:
```css
.ss-full {
    width: 1800px;
    height: auto;
    max-height: 960px;
    object-fit: contain;
    border-radius: 3px;
    border: 1px solid #2a3a2a;
    box-shadow: 0 8px 40px rgba(0, 0, 0, 0.5);
    display: block;
}
.ss-label {
    position: absolute;
    bottom: 30px;
    left: 0;
    width: 1920px;
    text-align: center;
    pointer-events: none;
}
.ss-label span {
    font-family: system-ui, -apple-system, sans-serif;
    font-size: 15px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #6B9A80;
    opacity: 0.7;
}
```

**Do NOT use the old `.app-hero` sizing** (width:1300px, max-height:640px) — this leaves ~310px dead space on each side. Always use `.ss-full` with 1800px width for near-edge-to-edge coverage.

Animation — use zoom-in for smooth entrance:
```javascript
tl.from("#ss-1", { opacity: 0, scale: 1.1, duration: 1.2, ease: "power3.out" }, 0.5);
tl.from(".ss-label", { opacity: 0, y: 15, duration: 0.5, ease: "power2.out" }, 1.2);
tl.to("#s2-ss", { opacity: 0, duration: 0.5, ease: "power2.in" }, 3.0);
tl.set("#s2-ss", { opacity: 0 }, 3.5);
```

### 11. Wire GSAP Per Scene

```javascript
// Screenshot scene 1: hero
tl.from("#s2-moments .display-small", { opacity: 0, y: 20, duration: 0.8, ease: "power2.out" }, 3.7);
tl.from("#app-hero-1", { opacity: 0, scale: 0.95, duration: 1.2, ease: "power2.out" }, 4.2);
tl.to("#s2-moments", { opacity: 0, duration: 0.6, ease: "power2.in" }, 7.0);

// Screenshot scene 2: workflow
tl.from("#s2-flow .display-small", { opacity: 0, y: 20, duration: 0.6, ease: "power2.out" }, 7.7);
tl.from("#app-hero-2", { opacity: 0, scale: 0.95, duration: 1, ease: "power2.out" }, 8.2);
tl.to("#s2-flow", { opacity: 0, duration: 0.6, ease: "power2.in" }, 10.0);
```

Each screenshot scene: ~1s headline entry → ~2-3s screenshot visible → 0.6s exit.

### 12. Account for Timeline Growth

Each screenshot scene adds ~3-4 seconds to the total video duration. When adding N screenshots:

- Video grows by N × 3.5s
- Music must be regenerated or extended to match (use `aloop` + `atrim` for quick loops: `ffmpeg -i music.mp3 -filter_complex "aloop=loop=1:size=1719900,atrim=duration=45" music-extended.mp3`)
- All scene 3+ timestamps shift forward by N × 3.5s
- Update root `data-duration="<new_total>"`

## Pitfalls

- **Deduplicate nav targets** — a button like "Home" or "Output" may appear in multiple nav lists. If your script iterates both the main sidebar views AND "more_nav", you'll re-click the same view and get a duplicate screenshot. Build one flat list of unique (view_name, button_text) pairs.
- **Don't overlap screenshots with cards in same scene** — use separate scenes at separate time slots
- **Venice vision for screenshot QA** — the current deepseek-chat model cannot process images via browser_vision; use `purfectlabs/vision-ocr` skill with Venice vision models instead
- **Dev server must stay alive** during browser capture — background processes die when the parent Hermes session is interrupted
- **SPA apps without Wails backend** may not navigate when buttons are clicked — capture what's visible on initial load
- **3.7K screenshots = blank captures** — if file is tiny, page wasn't rendered yet
- **Screenshots are 1920x1080** — the composition is also 1920x1080, so images will look slightly compressed; use `max-height: 640px` to leave room for headline + footer
- **Venice vision pricing quirk** — `/models` returns empty `pricing: {}` for qwen3-vl, but actual billing is $0.25/M input tokens. Don't trust the pricing field on VL models; check `/billing/usage` for real costs.
- **Venice vision capabilities quirk** — `capabilities.vision: false` on all models, even VL models. The real signal is `supportsMultipleImages: true` + `maxImages >= 10`. See `purfectlabs/vision-ocr` reference `references/venice-vision-quirks.md` for details.
- **Venice evaluation is a recommendation, not a decision** — the model may rank purely text-based AI settings screens higher than visually rich workflow views. Always present the ranking + all screenshots to the user and let them decide. The user's judgment of narrative fit overrides any model score.

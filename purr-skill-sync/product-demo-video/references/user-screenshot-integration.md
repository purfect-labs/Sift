# User-Provided Screenshots — Integration Guide

When the user takes their own screenshots of the app and says "use these", follow this exact flow.

## Trigger Conditions

- User says "in ~/Downloads/<product>/ there are wellnamed screenshots"
- User says screenshots are on Desktop
- User says "will get screenshots" and later points you to a directory

## Step 1: Locate Files

Check in this order:
1. `~/Desktop/<product>/` (most common for recent screenshots)
2. `~/Desktop/` (loose files with product prefix)
3. `~/Downloads/<product>/`

## Step 2: List by Name

User typically names them descriptively (e.g. `sift-mainjobs-matrix.png`, `sift-settings-and-keywork-extractions.png`). Grep the output to understand what each one shows.

## Step 3: Copy to HyperFrames Assets

```bash
cp ~/Desktop/<product>/*.png <video-dir>/assets/
```

## Step 4: Map Scenes by Content

Standard 6-scene demo mapping:
- Scene 2 (Dashboard/Overview) → full app matrix / job list view
- Scene 3 (Scraper) → view with scraping console/logs visible
- Scene 4 (Pipeline) → expanded job detail or pipeline columns
- Scene 5 (Privacy/Settings) → settings / keyword extraction view

Extra screenshots not needed immediately can be held for future versions.

## Step 5: Update Composition HTML

Only change the `<img src="../assets/...">` path in each scene's HTML file. Change NOTHING else — no GSAP timings, no text, no duration.

## Step 6: CRITICAL — Fix object-fit

The old placeholder screenshots likely used `object-fit: cover`. **Desktop app screenshots require `object-fit: contain; object-position: top;`** or the top of the screenshot is cut off. Update ALL scene compositions that display screenshots.

Before:
```css
.bg-image { object-fit: cover; opacity: 0; }
```

After:
```css
.bg-image { object-fit: contain; object-position: top; opacity: 0; background: #0a0c14; }
```

Apply to scenes 2, 3, 4, and 5 (scenes 1 and 6 have no screenshots).

## Step 7: Update Docs

Update `SCRIPT.md` to reference actual filenames in the overlay descriptions. This helps future debugging.

## Step 8: Verify Before Render

Check that the screenshot files exist in `assets/` and the paths in HTML reference the exact filenames including case.

## Common Pitfalls

- **`object-fit: cover` crops the TOP** of desktop app screenshots (biggest user complaint). Desktop SS are ~1200x900+ px in 1920x1080 canvas — cover zooms in from center, cutting both top chrome and sidebar.
- **User screenshots may have different aspect ratios** than the old placeholders. `object-fit: contain` with matching `background` handles this gracefully with letterbox bars.
- **Do NOT change the composition structure** (timing, GSAP, scene durations) when integrating user screenshots — only swap image paths.

# Suno — Music Generation Alternative

Use this when Venice stable-audio-25 doesn't deliver the emotional depth the video needs. Suno produces more expressive, human-feeling music — but has no official API, so every access path is a workaround.

## Access Options (ordered by reliability)

### 1. Browser Automation (Playwright/Puppeteer)
- Drive a headless or visible Chromium browser into the Suno web UI
- Requires a Suno account with credits
- Script: navigate to suno.com, log in, paste prompt, hit Generate, poll for completion, download result
- Risk: Suno may detect automation and block, rate-limit, or CAPTCHA
- Best for: one-off generation where quality matters more than speed

### 2. Reverse-Engineered Suno API
- Open DevTools on suno.com → Network tab during a generation
- Capture the internal API calls (POST to create, GET to poll, GET to download)
- Replay via curl with the same cookies/headers/tokens
- Risk: fragile — breaks on any Suno UI update. Session tokens expire.

### 3. Third-Party Wrappers
- Community projects like `suno-api` on GitHub or `gcui-art/suno-api`
- These wrap the undocumented Suno API and expose a REST/WebSocket endpoint
- Risk: may violate Suno ToS. Wrappers can be abandoned, inject ads, or steal credits.
- Vetted: `gcui-art/suno-api` (Node.js, active), `Suno-Api` (Python, less active)

## Suno Prompt Patterns (for demo video music beds)

### Style/Genre Field
```
Instrumental cinematic underscore. Emotional, swelling strings, warm piano, gentle horn.
Builds from intimate solo to full ensemble. 70 BPM, D major.
No vocals, purely instrumental background for narrative video.
```

Key differences from Venice prompts:
- Suno understands dynamic arc descriptions better ("whisper to roar to whisper")
- Avoid artist names — they get filtered. Describe the sound
- The Style field accepts up to 1,000 chars — use them
- Always specify "Instrumental" or "No Vocals" explicitly

### Lyrics Field (for instrumental)
Put this in the lyrics field for instrumental tracks:
```
[Instrumental Intro]
[Verse - Instrumental]
[Chorus - Instrumental]
[Bridge - Instrumental]
[Outro]
```

Without any lyric text between tags, Suno generates instrumental. But it sometimes adds vocalizations anyway — use `[Instrumental]` tags and describe "purely instrumental" in the style field.

## Tradeoffs vs. Venice stable-audio-25

| Factor | Venice stable-audio-25 | Suno |
|--------|----------------------|------|
| Cost | $0.24/36s | ~10 credits/song (~$0.10-0.20) |
| Access | Clean REST API | Browser hacks or wrappers |
| Reliability | Always works | Fragile, may break |
| Expressiveness | Adequate | Lush — noticeably better for emotional content |
| Instrumental guarantee | Describe in prompt | Tags + style field, no guarantee |
| Duration control | Exact (duration_seconds) | Variable extension |

## Decision Guide

- **Use Venice stable-audio-25 first** — it's fast, reliable, $0.24
- **Try Suno when** — the music sounds flat/robotic and the video needs emotional depth the user can feel
- **Only automate Suno if** — the user explicitly asks for it and accepts the fragility

# Interleaved Scene Pattern — Complex Multi-Feature Demo Videos

For apps with 6+ distinct features and multiple audience tiers (Evo-level), a pure screenshot walkthrough is boring. The solution is **interleaving**: alternate between text cards (concept explanation) and full screenshots (visual proof).

## The Pattern

```
[Hook]         → [Text Card] → [Screenshot] → [Text Card] → [Screenshot]
→ [Text Card] → [Screenshot] → [Text Card] → [Screenshot] → [Text Card]
→ [Screenshot] → [Screenshot] → [Close]
```

Each text card makes a claim. The following screenshot proves it.

## Scene Typology

### 1. Hook (0-8s)
- Badge line: core differentiator (e.g., "Real-Time P2P Operational Mesh")
- Product name in large font
- Divider in brand accent color
- Audience tier pills: C-Level, Project, Dev, QA — each as a styled `<span>` with border
- Stagger animate: badge → name → divider → tagline → role pills

### 2. Text Cards (10-14s each)
- **Badge** at top (yellow pill, small) — single concept name
- **Accent line** — green/yellow divider
- **Headline** — 36-48px, bold claim
- **Subtitle** — 16-18px monospace, explains how
- **Optional: role pills** at bottom showing who this matters to
- Animate: badge → accent → headline → subtitle → roles

### 3. Screenshots (12-18s each)
- Full-frame image: 1800px wide, max-height 960px
- Bottom overlay label with feature name
- Zoom-in entrance: scale 1.1→1.0, `power3.out`, 1.2-1.4s
- Label fade in 0.5s after image
- Feature screenshots (rich views like Active Mesh): 16-18s
- Simple screenshots (like graphs): 12-14s

### 4. Close (8-10s)
- Tagline matching product's emotional weight (command center for ops tools)
- Audience tier reminder line
- DEVELOPED BY PurfectLabs
- App portfolio cards stagger in, fade out at ~7s
- Brand holds solo for final 3s
- 5-8s of silence after VO ends

## 180s (3-Minute) Template

| # | Type | Time | Duration | Animation timing |
|---|------|------|----------|-----------------|
| 1 | Hook | 0-8s | 8s | Enter 0.3-3.2s, exit 7.4s |
| 2 | Screenshot | 8-22s | 14s | Enter 8.3s, exit 21.3s |
| 3 | Text Card | 22-34s | 12s | Enter 22.2s, exit 33.4s |
| 4 | Screenshot | 34-52s | 18s | Enter 34.3s, exit 51.3s |
| 5 | Text Card | 52-64s | 12s | Enter 52.2s, exit 63.4s |
| 6 | Screenshot | 64-82s | 18s | Enter 64.3s, exit 81.3s |
| 7 | Text Card | 82-94s | 12s | Enter 82.2s, exit 93.4s |
| 8 | Screenshot | 94-112s | 18s | Enter 94.3s, exit 111.3s |
| 9 | Text Card | 112-124s | 12s | Enter 112.2s, exit 123.4s |
| 10 | Screenshot | 124-140s | 16s | Enter 124.3s, exit 139.3s |
| 11 | Text Card | 140-150s | 10s | Enter 140.2s, exit 149.4s |
| 12 | Screenshot | 150-164s | 14s | Enter 150.3s, exit 163.3s |
| 13 | Screenshot | 164-173s | 9s | Enter 164.3s, exit 172.4s |
| 14 | Close | 173-183s | 10s | Enter 173.1s, exit 182.0s+set |

**Total: 14 scenes, 183s (3:03).**

## GSAP Animation Rules Per Scene Type

**Hook:** `back.out(1.5+)` for product name, `power2.out` for badges/taglines, stagger role pills at 0.1s.

**Text Cards:** `power3.out` for headline (0.7s), `power3.out` for subtitle (0.6s), `back.out(1.3)` for badges/role pills with stagger.

**Screenshots:** Always `scale: 1.1, duration: 1.2-1.4, ease: power3.out`. Label fades in 0.5s later with `power2.out`.

**Close:** Same as close-screen-template.md. Brand stagger from 176s with 0.2s stagger on app cards. Cards and extra text fade out at 182s. Final `tl.set()` at 183s.

## Color Palette Per App

| App | Background | Accent 1 | Accent 2 | Text |
|-----|-----------|----------|----------|------|
| Evo | `#060D08` | `#FFD600` (yellow) | `#00D68F` (emerald) | `#E8F0EC` |
| JobDash | `#0a0c14` | `#7c8aff` (lavender) | `#4fc3f7` (cyan) | `#e0e0e0` |
| Undertow | `#070a08` | `#4a9eff` (phosphor) | `#c8a44a` (gold) | `#d4e0d4` |

## Voiceover Pacing for 180s

Target: VO ends ~174s (8-9s before close).

Speed: 0.86-0.88. Script length: ~560-600 words.

Each screen type gets roughly matching VO time:
- Hook: ~15s of VO (30 words)
- Each text card: ~8-10s of VO (20-25 words)
- Each screenshot: ~8-12s of VO (20-30 words)
- Close: ~4s of VO (last line)

Use `[pause 200]` between major sections. Use `[pause 400]` before the final "Download" line.

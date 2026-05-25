# Undertow Demo Video — Reference Composition

A 39-second cinematic product demo for Undertow (AI music portfolio studio). Dark river aesthetic, 6-scene narrative arc including a real app screenshot.

Jump to: [Critical Pacing Rule](#critical-pacing-rule) · [Real App Screenshot](#real-app-screenshot) · [GSAP Patterns](#key-gsap-patterns-used)

## Narrative Arc

| Scene | Timing | Content | Purpose |
|-------|--------|---------|---------|
| Hook | 0-3s | "Before you forget..." + "Your story. Your sound." | Emotional opener |
| App screenshot + memory cards | 3.5-15.5s | Real Undertow app screenshot ("your life, scored.") + 4 memory cards fading in at 1.8s stagger | Ground in real product + human experience |
| Era timeline | 16-24s | 4 life eras as album titles (Small World → Picking Up Speed → The Unlearning → Where The Light Gets In) | Show the product's core mechanism |
| Bridge | 24-30s | "What if there was a song for that moment?" | Emotional buy-in before payoff |
| Album reveal | 30-35s | 4 album cards with song titles (Porch Light, June Bug Summer, etc.) | The outcome |
| Closer | 35-39s | "Turn your life into an album. Share it with the people who lived it with you." + brand | CTA + sharing hook |

## Critical Pacing Rule

**Cards need 1.5-2s stagger minimum.** The original 0.8s stagger between memory cards was too fast for viewers to read the text — the user called that out during review. Always test card pacing against the voiceover timing: if a card fades out before the viewer finishes reading the first line, widen the stagger.

Correct stagger pattern:
```js
tl.from("#mem-card-1", { opacity: 0, y: 30, duration: 0.8, ease: "power2.out" }, 6.5);
tl.from("#mem-card-2", { opacity: 0, y: 30, duration: 0.8, ease: "power2.out" }, 8.3); // 1.8s stagger
tl.from("#mem-card-3", { opacity: 0, y: 30, duration: 0.8, ease: "power2.out" }, 10.1);
tl.from("#mem-card-4", { opacity: 0, y: 30, duration: 0.8, ease: "power2.out" }, 11.9);
```

## Real App Screenshot

Replace generic text-based memory cards with a real screenshot of the running app when available. This grounds the video in the actual product and adds visual interest.

**Capturing:** Run `npm run dev` in the frontend directory, navigate the browser to `localhost:5173`, screenshot the app view.

**Embedding:**
```html
<img class="app-hero clip" id="app-hero" data-start="4" data-duration="10" data-track-index="3"
     src="assets/app-screenshot-1.png"
     alt="App screenshot" />
```

**CSS:**
```css
.app-hero {
    width: 800px;
    height: 460px;
    border-radius: 4px;
    border: 1px solid #2a3a2a;
    box-shadow: 0 8px 40px rgba(0, 0, 0, 0.5);
    object-fit: cover;
}
```

**GSAP:**
```js
tl.from("#app-hero", { opacity: 0, scale: 0.95, duration: 1.2, ease: "power2.out" }, 4.2);
```

## Design Tokens (from undertow app.css)

- Background: `#070a08` (abyss) → `#0a0f0a` (deep water)
- Primary text: `#d4e0d4` (soft green-white)
- Soft text: `#a0b8a0`
- Accent: `#4a9eff` (phosphor blue)
- Warmth: `#c8a44a` (amber gold)
- Dim: `#6B9A80`
- Display font: Georgia, serif (editorial, warm)
- Body font: system-ui, sans-serif
- Card borders: `#2a3a2a` (subtle green tint)
- Card backgrounds: `rgba(17, 26, 17, 0.8)`

## Key GSAP Patterns Used

### Text entrance
```js
tl.from("#element-id", { opacity: 0, y: 30, duration: 1.2, ease: "power2.out" }, 0.2);
```

### Scale accent lines
```js
tl.from(".accent-line", { scaleX: 0, transformOrigin: "center", duration: 0.8, ease: "power2.out" }, 0.8);
```

### Scene exit + hard kill
```js
tl.to("#scene-id", { opacity: 0, duration: 0.8, ease: "power2.in" }, 20.2);
tl.set("#scene-id", { opacity: 0 }, 21.0);  // hard kill
```

### Staggered elements (slow pace — 1.8s minimum)
```js
// Cards need 1.5-2s between each entrance so viewers can read
tl.from("#mem-card-1", { opacity: 0, y: 30, duration: 0.8, ease: "power2.out" }, 6.5);
tl.from("#mem-card-2", { opacity: 0, y: 30, duration: 0.8, ease: "power2.out" }, 8.3); // 1.8s
tl.from("#mem-card-3", { opacity: 0, y: 30, duration: 0.8, ease: "power2.out" }, 10.1);
tl.from("#mem-card-4", { opacity: 0, y: 30, duration: 0.8, ease: "power2.out" }, 11.9);
```

### Image scale entrance
```js
tl.from("#app-hero", { opacity: 0, scale: 0.95, duration: 1.2, ease: "power2.out" }, 4.2);
```

## Memory Card Pattern

```html
<div class="memory-card clip" data-start="4" data-duration="8" data-track-index="3">
    "Emotional fragment text."
    <span class="tag">childhood</span>
</div>
```

### Matching CSS
```css
.memory-card {
    background: rgba(17, 26, 17, 0.8);
    border: 1px solid #2a3a2a;
    border-radius: 3px;
    padding: 24px 28px;
    width: 280px;
    font-family: Georgia, serif;
    font-size: 20px;
    font-style: italic;
    color: #d4e0d4;
    line-height: 1.5;
    text-align: center;
}
.memory-card .tag {
    font-family: system-ui, sans-serif;
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: #4a9eff;
    border: 1px solid rgba(74, 158, 255, 0.3);
}
```

## Era Timeline Pattern

A vertical timeline with colored dots. Each era has a year range, title, and color:

```html
<div class="era-line"></div>
<div class="era-marker clip" data-start="13.5" data-duration="4" data-track-index="8">
    <span class="era-year">AGES 0–12</span>
    Small World <span class="era-dot" style="background:#4A6FA5;"></span>
</div>
```

## Album Reveal Pattern

```html
<div class="album-card clip" data-start="27.5" data-duration="3" data-track-index="14">
    <div class="era-label">Album Title</div>
    <div class="song-count">5 songs</div>
    <div class="songs-list">
        Song Title 1<br/>
        Song Title 2<br/>
        Song Title 3
    </div>
</div>
```

## Brand Closer Pattern

```html
<div class="brand-word blue">U</div>
<div class="brand-word">NDERTOW</div>
```

## Source

Full composition: `purfect-labs/undertow` repo, branch `feat/undertow-demo-video`, file `undertow-video/index.html`

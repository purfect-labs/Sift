# MASTER Manifest — [Product Name] Audience Variants

## Purpose
Traces all audience-targeted variants of this demo back to the original. Evo can reconcile this manifest to understand the full video family.

---

## Variant Overview

| # | File | Audience | Music | Voice | Duration | Base Origin |
|---|------|----------|-------|-------|----------|-------------|
| 1 | `[product]-demo-v<N>.mp4` | [audience desc] | [music desc] | [voice name/ID] | [seconds]s | Original |
| 2 | `audience-[name]/[product]-audience-[name].mp4` | [audience desc] | [music desc] | Purr (`.env` voice) | [seconds]s | Forked from v<N> |
| 3 | `audience-[name]/[product]-audience-[name].mp4` | [audience desc] | [music desc] | Purr (`.env` voice) | [seconds]s | Forked from v<N> |
| 4 | `audience-[name]/[product]-audience-[name].mp4` | [audience desc] | [music desc] | Purr (`.env` voice) | [seconds]s | Forked from v<N> |

---

## Variant Detail

### Variant 1 — [Audience Name]
- **File:** `[product]-demo-v<N>.mp4`
- **Target audience:** [description]
- **Music model:** [e.g., `stable-audio-25`]
- **Music prompt:** `[exact prompt used]`
- **Music duration:** [seconds]s
- **Voice:** [name] — `[voice_id]`
- **TTS tags:** `[tag1] [tag2] [tag3]`
- **TTS params:** temp=[temp], speed=[speed]
- **VO script:**
  ```
  [full VO script text]
  ```
- **VO duration:** [seconds]s
- **Screenshots used:**
  - [filename1.png]
  - [filename2.png]
  - [filename3.png]
  - [filename4.png]
- **Composition:** [e.g., `index.html` from original HyperFrames project]
- **Duration:** [seconds]s
- **Delivery:** Slack #channel via MEDIA tag

### Variant 2 — [Audience Name]
- **File:** `audience-[name]/[product]-audience-[name].mp4`
- **Target audience:** [description]
- **Music model:** [e.g., `stable-audio-25`]
- **Music prompt:** `[exact prompt used]`
- **Music duration:** [seconds]s
- **Voice:** Purr — `[FISH_AUDIO_VOICE from .env]`
- **TTS tags:** `[tag1] [tag2] [tag3]`
- **TTS params:** temp=[temp], speed=[speed]
- **VO script:**
  ```
  [full VO script text]
  ```
- **VO duration:** [seconds]s
- **Screenshots used:** (same as base — copied from original)
- **Composition:** Forked from original `index.html` with modified memory-card copy and close CTA
- **Duration:** [seconds]s
- **Base origin:** Forked from [product]-demo-v<N>.mp4

### Variant 3 — [Audience Name]
[Same structure as Variant 2]

### Variant 4 — [Audience Name]
[Same structure as Variant 2]

---

## Render Chain

### Music generation
- [Date]: Generated via Venice `POST /audio/queue` for all 4 variants simultaneously
- Quote verified for each variant before queue

### Voiceover generation
- [Date]: Generated via Fish Audio TTS for all audience variants
- Original variant used [voice name/N], audience variants use Purr's voice

### Final mix
- [Date]: ffmpeg amix for each variant
- Command pattern:
  ```bash
  ffmpeg -y -i video.mp4 -i music-padded.mp3 -i vo.mp3 \
    -filter_complex "[1:a]volume=0.28[m];[2:a]volume=1.0,adelay=400|400[vo];[m][vo]amix=inputs=2:duration=first:dropout_transition=0[a]" \
    -map 0:v -map "[a]" -c:v copy -c:a aac -b:a 192k output.mp4
  ```

### Delivery
- [Date]: All 4 variants delivered to [platform/channel]
- Each sent via inline MEDIA tag

---

## Notes
- All audience variants share the same app screenshots and GSAP composition structure
- Differentiation is through VO script, music bed, and memory-card copy text only
- Purr's voice ID: `3de74fc2a4ec4a0a950f36b726d5ec59` (from profile `.env`)

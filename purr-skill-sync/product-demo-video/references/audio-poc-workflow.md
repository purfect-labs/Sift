# Audio POC Workflow — Mix Before You Render

Instead of running the full video stitch each time you adjust audio levels, POC the mix standalone. Saves 2-3 minutes per iteration and lets you dial in levels while the user watches.

## Why This Exists

The user explicitly asked for this pattern during the Undertow video work: "POC the audio/music combo before making the video." Stitching music + voiceover takes ~2s. Re-rendering the HyperFrames composition takes 2-3 minutes. Iterate on the mix, then commit to the full render.

## The Pattern

### Step 1 — Generate Audio Assets

- **Music bed:** Venice `stable-audio-25` (Phase 6 in main skill)
- **Voiceover:** Fish Audio TTS using the `fish-audio` skill's Voice Discovery section

### Step 2 — POC the Mix (No Video)

```bash
ffmpeg -y \
  -i assets/audio/music-bed.mp3 \
  -i assets/audio/voiceover.mp3 \
  -filter_complex \
    "[0:a]volume=0.28[m]; \
     [1:a]adelay=3000|3000,volume=1.3[v]; \
     [m][v]amix=inputs=2:duration=first:dropout_transition=3[a]" \
  -c:a aac -b:a 256k \
  -shortest \
  /tmp/audio-mix-poc.mp3
```

Deliver via MEDIA tag for user feedback.

### Step 3 — Iterate on Parameters

| Parameter | Typical range | What it does |
|-----------|---------------|--------------|
| `volume=0.25-0.35` | Music level behind VO | Higher = more emotional/dominant |
| `volume=1.1-1.4` | Voiceover clarity | Higher = punchier; avoid >1.5 (clipping) |
| `adelay=2000-5000` | VO start delay (ms) stereo pair | Lower = VO starts sooner; match to hook scene |
| `dropout_transition=2-5` | Amix transition seconds | Lower = faster ducking when VO ends |

### Step 4 — Once Dialed, Stitch Into Video

```bash
ffmpeg -y \
  -i undertow-demo-video-only.mp4 \
  -i assets/audio/music-bed.mp3 \
  -i assets/audio/voiceover.mp3 \
  -filter_complex \
    "[1:a]volume=0.28[m]; \
     [2:a]adelay=3000|3000,volume=1.3[v]; \
     [m][v]amix=inputs=2:duration=first:dropout_transition=3[a]" \
  -map 0:v -map "[a]" \
  -c:v libx264 -preset ultrafast -crf 18 \
  -c:a aac -b:a 256k \
  -shortest -movflags +faststart \
  undertow-demo-final.mp4
```

## Pitfalls

- **Same `filter_complex` rules apply** — this uses `amix`, not `afade,volume`. The `-map "[a]"` output is safe here because `amix` produces correctly formatted audio, unlike the simple single-stream filter chains that cause the 21kbps aac cliff.
- **Don't skip the POC for "just one quick render"** — the HyperFrames render is the time sink, not the audio mix. Always POC first.
- **Verify final audio bitrate:** `ffprobe -v error -show_entries stream=bit_rate -select_streams a final.mp4`. Should be ~220-256k. If ~21k, your filter_complex has the wrong chain.

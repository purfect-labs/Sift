# ffmpeg Audio Mix Patterns for Product Demo Videos

## The Two Critical Rules

1. **Never `-c:v copy`** — stream copy preserves old container metadata that confuses chat player seeking. Always re-encode video with `-c:v libx264 -preset ultrafast -crf 18`.

2. **Never `filter_complex` + `-map "[a]"` for simple audio** — see "The AAC Bitrate Cliff" below.

## The AAC Bitrate Cliff (Root Cause Analysis)

### Symptom
Video plays but audio is nearly silent or cuts off after a few seconds. The user says "audio is cutting off again short" or "it sounds like a sad church commercial."

### Cause
A `filter_complex` chain like `[1:a]afade=t=in:d=2,afade=t=out:d=3,volume=0.38[a]` with `-map 0:v -map "[a]"` causes the aac encoder to tank from 256kbps to ~21kbps. The output file is ~3MB instead of ~4MB for 36s of 1080p video. The low bitrate audio is barely audible in chat players.

### Diagnosis
```bash
ffprobe -v error -show_entries stream=bit_rate -of default=nw=1:nk=1 output.mp4
```
If the audio stream shows ~21638 bps instead of ~256000 bps, you've hit the cliff.

### Fix
Always use simple stream mapping for background music:
```bash
ffmpeg -y \
  -i video.mp4 \
  -i music.mp3 \
  -c:v libx264 -preset ultrafast -crf 18 \
  -c:a aac -b:a 256k \
  -shortest -movflags +faststart \
  output.mp4
```

### Why filter_complex Causes It
The `afade,volume` filter chain outputs audio without proper stream metadata that the aac encoder expects. The encoder falls back to a minimal bitrate by default. This does NOT happen with `amix` (multi-stream mixing) because the more complex filter graph properly initializes the encoder.

### Safe filter_complex Usage
Voiceover + music layering (the actual use case for filter_complex):
```bash
ffmpeg -i video.mp4 -i music.mp3 -i voiceover.mp3 \
  -filter_complex "[1:a]volume=0.25[m];[2:a]adelay=5000|5000,volume=1.2[v];[m][v]amix=inputs=2:duration=first[a]" \
  -map 0:v -map "[a]" \
  -c:v libx264 -preset ultrafast -crf 18 \
  -c:a aac -b:a 256k \
  -shortest -movflags +faststart \
  output.mp4
```

## Volume Adjustment (pre-processing pattern)

If you need volume control, apply it to the music file *before* mixing:

```bash
# Pre-process music with volume adjustment
ffmpeg -i music.mp3 -filter:a "volume=0.38" -c:a aac -b:a 256k music-adjusted.mp3

# Then mix cleanly (no filter_complex)
ffmpeg -i video.mp4 -i music-adjusted.mp3 \
  -c:v libx264 -preset ultrafast -crf 18 \
  -c:a aac -b:a 256k \
  -shortest -movflags +faststart \
  output.mp4
```

Recommended music volume: 0.35-0.40 (background music should support the visuals, not compete).

## Render-Pattern for Variant Iteration

When testing multiple music beds, extract video-only source once, then iterate:

```bash
# One-time: extract video-only source
ffmpeg -i original-render.mp4 -c:v libx264 -preset ultrafast -crf 18 -an video-only.mp4

# Per-variant: just mix audio
ffmpeg -i video-only.mp4 -i variant-1.mp3 \
  -c:v libx264 -preset ultrafast -crf 18 \
  -c:a aac -b:a 256k \
  -shortest -movflags +faststart \
  variant-1-output.mp4
```

This avoids re-rendering the entire composition for each music variant.

## Manifest Pattern

Store `assets/audio/manifest.txt` linking each variant file to its prompt:

```
variant-1.mp3 | Warm ambient cinematic instrumental. Gentle piano, soft strings, nostalgic and hopeful. Slowly building. 70 BPM. | 36s
variant-2.mp3 | Emotional cinematic orchestral instrumental. Rising strings, French horn swells. Builds from cello to full orchestra. 60 BPM. | 36s
```

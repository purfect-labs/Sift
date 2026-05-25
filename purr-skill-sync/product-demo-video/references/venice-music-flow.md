# Venice Music Flow — stable-audio-25

Quick reference for the cheapest music generation model. Full lifecycle reference: `venice/venice-audio-music` skill.

## Important Notes

- **stable-audio-25 does NOT support `force_instrumental`** — omit the param and describe "instrumental" in the prompt
- **Returns WAV** — always convert to MP3 before mixing
- **Poll timing** — ~15s for 36s, ~30-40s for 190s. Check `average_execution_time` from status response
- **Cost** — $0.24 flat, regardless of duration. Confirmed: 36s, 55s, and 190s all return $0.24 quote. This is the cheapest option for any length.
- **Maximum duration confirmed** — 190s works (33.5MB WAV). Unknown if longer durations work; fall back to `elevenlabs-music` or loop/concat for >190s.

## Common Prompts

- Emotional/warm: `"Warm ambient cinematic instrumental. Gentle piano melody, soft strings, nostalgic and hopeful. Slowly building from intimate to full. 70 BPM."`
- Energetic: `"Uplifting electronic pop instrumental. Driving beat, bright synths, optimistic and energetic. 120 BPM."`
- Minimal/meditative: `"Minimal ambient drone. Soft pads, gentle harmonics, meditative and spacious. 60 BPM."`
- Tech documentary (long form, 3min+): `"instrumental cinematic ambient score, deep evolving bass, atmospheric synth pads, subtle piano motifs, slow burn tech documentary feel, 90 BPM, wide stereo field, minor key"`

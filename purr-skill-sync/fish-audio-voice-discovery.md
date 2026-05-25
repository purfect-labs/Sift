# Voice Discovery — Finding Voices via GET /model

Fish Audio has 2M+ public voices. The `/model` endpoint lets you search them programmatically. No web UI needed.

## Endpoint

```
GET /model?page_size=N&page_number=1&tag=<tag1>&tag=<tag2>&tag=<tag3>&self=false
```

## Key Fields in Response

| Field | Maps to | Notes |
|-------|---------|-------|
| `_id` | TTS `reference_id` | Use this to generate speech |
| `title` | Human name | "Adrian", "Laura", "Ethan" |
| `tags` | Descriptors | Filter by gender, mood, use case, age |
| `description` | NLP description | Scan for "warm", "dramatic", "expressive" |
| `like_count` | Popularity | Higher ≠ better for your use case |
| `task_count` | Usage volume | Indicates production-tested voices |
| `languages` | Locale support | Filter for `"en"` |
| `visibility` | Must be `"public"` | Private voices need access |
| `samples[].audio` | Preview URL | Test before committing |

## Tag Filtering

Repeat `tag=` for AND logic:

```
GET /model?page_size=10&tag=female&tag=warm&tag=narration
```

Common descriminating tags:

| Category | Tags |
|----------|------|
| Gender | `male`, `female`, `neutral` |
| Age | `young`, `middle-aged`, `old` |
| Mood | `warm`, `calm`, `dramatic`, `expressive`, `serious`, `mysterious`, `cheerful` |
| Use case | `narration`, `storytelling`, `documentary`, `educational`, `entertainment`, `social-media` |
| Delivery | `professional`, `deep`, `raspy`, `smooth`, `soft`, `breathy`, `high`, `low` |
| Style | `cinematic`, `conversational`, `authoritative`, `gentle`, `character-voice` |

## Language Filtering

The `languages` field is an array. Filter in post-processing:

```python
for item in data['items']:
    if 'en' not in item.get('languages', []):
        continue
    # process
```

## Speed-Testing a Voice

Always test 2-3 candidates with one line from your actual script. Don't rely on like_count:

```python
import subprocess, json

payload = json.dumps({
    "text": "[warm, reflective] Your actual script line here.",
    "reference_id": "<VOICE_ID>",
    "temperature": 0.85,
    "prosody": {"speed": 0.88},
    "format": "mp3",
    "mp3_bitrate": 128
})

subprocess.run([
    "curl", "-s", "-w", "%{http_code}",
    "-X", "POST", "https://api.fish.audio/v1/tts",
    "-H", "Authorization: Bearer <KEY>",
    "-H", "Content-Type: application/json",
    "-H", "model: s2-pro",
    "-d", payload,
    "-o", "/tmp/voice-test.mp3"
], timeout=60)
```

## Tested Picks (English, emotional narration)

| Voice | ID | Vibe | Likes |
|-------|----|------|-------|
| Adrian | `bf322df2096a46f18c579d0baa36f41d` | Deep, measured, dramatic narrator | 1245 |
| Laura | `e3cd384158934cc9a01029cd7d278634` | Warm, confident female narrator | 658 |
| Ethan | `536d3a5e000945adb7038665781a4aca` | "A curious explainer" — professional | 1337 |

These are examples, not defaults. Each project should discover its own voice.

# Venice Vision — Batch Screenshot Evaluation for Video Production

Evaluate multiple app screenshots in one call to pick the best ones for a product demo video. Correlates scores to filenames so you know which is which.

## When to use

- You've captured 6+ screenshots of a dark-themed desktop app
- You need an objective ranking before presenting options to the user
- The screenshots look similar to human eyes (same sidebar, different content)

## When NOT to use

- User explicitly says "just show me all of them and I'll pick" — skip the evaluation, go straight to sharing
- Only 1-3 screenshots — not worth the API call

## Model choice

| Model | Cost/M input | Quality | Reason |
|-------|-------------|---------|--------|
| **qwen3-vl-235b-a22b** 🔥 | $0.25 | Best | Handles dark themes well, good at reading UI text |
| qwen3-5-9b | $0.10 | OK | Cheapest but may miss detail on dark-themed UIs |

At $0.25/M, a batch of 12 screenshots costs ~$0.002. Always use qwen3-vl.

## Pricing quirk

The `/models` endpoint returns empty `pricing: {}` for qwen3-vl — but actual billing is $0.25/M input, $1.50/M output. Check via `GET /billing/usage` if you need confirmation. The `capabilities.vision` field also returns `false` even though the model processes images — the real signal is `supportsMultipleImages: true` and `maxImages >= 10`.

## The batch call

### Build the content array

```python
import subprocess, json, base64

r = subprocess.run(
    ["grep", "VENICE_API_KEY", "/Users/shawn/.hermes/profiles/purfectlabs-tg-bot/.env"],
    capture_output=True, text=True
)
api_key = r.stdout.strip().split("=", 1)[1]

images = [
    ("app-full-welcome.png", "/path/to/app-full-welcome.png"),
    ("app-view-StyleProfile.png", "/path/to/app-view-StyleProfile.png"),
    # ... all screenshots
]

content_parts = [{
    "type": "text",
    "text": """You are evaluating screenshots of [APP NAME] for a product demo video.
Each image is labeled by filename in order (first = image 1, second = image 2, etc.).

For EACH screenshot, return:
- filename: the image filename (correlate based on image order in the input)
- description: 1 sentence what this screen shows
- clarity: 1-5 (text readable, not cropped/blurry, dark elements visible)
- visual_impact: 1-5 (looks cinematic, impressive to show)
- uniqueness: 1-5 (visually distinct from the other screens)
- story_value: 1-5 (does it communicate a step in the user journey?)

Then give me:
- recommended_order: which filenames to use in the video and in what sequence
- reasoning: 1-2 sentences why this order"""
}]

for name, path in images:
    with open(path, 'rb') as f:
        b64 = base64.b64encode(f.read()).decode()
    content_parts.append({
        "type": "image_url",
        "image_url": {"url": f"data:image/png;base64,{b64}"}
    })
```

### Send

```python
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

usage = data.get('usage', {})
total_in = usage.get('prompt_tokens', 0)
cost = (total_in / 1000000) * 0.25
print(f"\\nCost: ${cost:.6f}")
```

## Reading the results

The model returns per-screenshot blocks like:

```
filename: app-view-StyleProfile.png
description: The style profile screen lets users define...
clarity: 5
visual_impact: 4
uniqueness: 4
story_value: 5
total_score: 18
rank: 2
```

Then a final recommendation block. **This is a recommendation, not a decision.** Always present ranking + all screenshots to the user and let them make the final call.

## Pitfalls

- **Filenames must be embedded in the prompt text**, not just in code. The model correlates by image position in the input array — first image = first filename mentioned.
- **Dark themes are hard for some models.** If a model reports low clarity on all screenshots, try qwen3-vl specifically — it handles dark-on-dark better.
- **The model may rank settings screens higher than workflow screens.** Don't blindly follow the ranking. The user's narrative intent overrides the model score.
- **Always print usage/cost.** It's always under $0.01 — showing this builds confidence in using the approach.

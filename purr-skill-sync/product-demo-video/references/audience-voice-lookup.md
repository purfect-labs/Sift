# Audience Voice Lookup — Fish Audio Voices by Target Demographic

Discovered voices mapped to audience types. Each voice was found via `GET /model` tag search and tested with one line from the actual script before committing. Update this file when new voices are discovered.

## Teens — Energetic, TikTok-native, coming-of-age
| Voice | ID | Tags | Found via |
|-------|-----|------|-----------|
| Friendly Women | `b545c585f631496c914815291da4e893` | female, young, social-media, bright, energetic, enthusiastic | `tag=young&tag=female&tag=energetic` |
| Hatsune Miku (TTS) | `acc8237220d8470985ec9be6c4c480a9` | female, young, energetic, cheerful, bright, friendly, anime | `tag=young&tag=female&tag=energetic` |
| E-girl | `98655a12fa944e26b274c535e5e03842` | female, young, egirl, conversational, soft, breathy, calm | `tag=young&tag=female` |

**Best match:** Friendly Women — bright, social-media energy, not artificial/character-voice. Use temp 0.85, speed 0.92, tags [bright] [energetic] [quick] [crisp].

## Rappers — Street poets, beat-driven, hip-hop heads
| Voice | ID | Tags | Found via |
|-------|-----|------|-----------|
| horror | `ef9c79b62ef34530bf452c0e50e3c260` | male, middle-aged, narration, deep, low, serious, calm, mysterious, storytelling | `tag=male&tag=deep&tag=storytelling` |
| BOOK RECORD REGULAR | `f8dfe9c83081432386f143e2fe9767ef` | male, old, narration, deep, raspy, calm, cinematic, measured, storytelling | `tag=male&tag=deep` |
| Sleep history | `c97c6987eaf944d5b738c592c82ed0e9` | male, middle-aged, narration, deep, raspy, serious, mysterious, measured, storytelling | `tag=male&tag=deep` |

**Best match:** horror — deep, measured, mysterious, storytelling. Description: "A deep and serious middle-aged male voice with a calm, mysterious tone. Well-suited for suspenseful narration, documentaries, or cinematic storytelling." Use temp 0.80, speed 0.88, tags [steady] [measured] [mysterious] [quiet] [firm].

## 30-Somethings — Reflective, nostalgic, bittersweet
| Voice | ID | Tags | Found via |
|-------|-----|------|-----------|
| Laura | `e3cd384158934cc9a01029cd7d278634` | female, middle-aged, conversational, educational, deep, warm, calm, professional, clear | `tag=warm&tag=female&tag=narration` |
| Ringo Starr | `da8b52ae506044ce8c195a874ee572f8` | male, old, narration, British, calm, raspy, reflective, warm, storytelling | `tag=warm&tag=reflective&tag=narration` |

**Best match:** Laura — warm, calm, professional, clear. Grounded reflective narrator. Use temp 0.82, speed 0.88, tags [warm] [reflective] [quiet] [measured] [firm].

## Old / Sentimental — Original voice (memory-driven)
| Voice | ID | Tags |
|-------|-----|------|
| Philosopher | `61ca45191c804d86b7d7555fa136945a` | British, reflective, emotional (vetted, not from API search) |
| Adrian | `bf322df2096a46f18c579d0baa36f41d` | male, middle-aged, narration, deep, slow, measured, serious, dramatic, storytelling |

## Methodology
1. Search `GET /model` with tag filters matching audience demographics
2. Post-filter for `languages` containing `"en"`
3. Read `description` and `tags` to narrow candidates
4. Test-generate a short clip from the actual script before committing
5. Record chosen voice + params in each variant's `manifest.txt`

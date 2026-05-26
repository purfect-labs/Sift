# sift — Bug Log

## Active Bugs

_No active bugs._

## Resolved Bugs

### BUG-001-shawnM: App crash on Extract Keywords `closed` <!-- UUID: bug-001 -->
**Opened:** 2026-05-24
**Closed:** 2026-05-24
**Severity:** High
**Found:** macOS GUI app — clicking Extract Keywords crashed the app silently.
**Root cause:** Python3 and Hermes binaries not in sandboxed GUI app PATH. Also, empty command output caused slice bounds panic (`response[:200]` on empty string).
**Fix:** Added `exec.LookPath` for hermes and python3. Added `safeTextSlice`/`safeSlice` bounds checks. Changed `ExtractResume` to return errors in result struct instead of throwing.

### BUG-002-shawnM: Hermes returned 0 keywords `closed` <!-- UUID: bug-002 -->
**Opened:** 2026-05-24
**Closed:** 2026-05-24
**Severity:** High
**Found:** Hermes `-z` produced empty output.
**Root cause:** No default model configured in Hermes profile.
**Fix:** Ran `hermes config set model deepseek-chat`. Added Hermes setup docs to README.

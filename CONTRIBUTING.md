# Contributing to JobDash

JobDash is open source and freely given to the world. We welcome contributions but keep it simple.

## How to Contribute

1. Fork the repo
2. Create a branch: `feat/your-feature` or `fix/your-fix`
3. Make your changes
4. Run `go test ./...` and `evo validate --dir .`
5. Open a PR against `main`

## Development

```bash
wails3 dev          # hot-reload dev mode
go test ./...       # run tests
evo validate --dir . # validate .evolution/
```

## Code Style

- Go: standard conventions, `gofmt`
- Svelte: single-file components, Svelte stores for state
- Commit messages: `type(scope): description`

## Project Status

This repo is released as-is. No maintenance commitments, no SLA, no support guarantees. Use it, fork it, build on it — it's yours.

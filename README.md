# ggit

A modern, ultra-lightweight, and secure self-hosted Git service built with **Go (Gin)** and **SvelteKit**.

## Features

- **Modern Architecture** — SvelteKit frontend with a Go + Gin backend, fully decoupled.
- **SQLite Only** — Ultra-optimized SQLite with WAL mode for minimal footprint.
- **Code Review & Branch Protection** — Pull request approvals, turn-based reviewing, and inline diff comments.
- **Docs-as-Code** — Wiki generated from a `/docs` folder inside your repository (no more `.wiki.git`).
- **Kanban & Unified Inbox** — Virtual Kanban boards and multi-repo task tracking from a single inbox.
- **Secret Guardian** — Proactive detection of secrets (API keys, `.env` files) during Git pushes.

## Quick Start

```bash
# Backend
go run main.go

# Frontend (separate terminal)
cd frontend && npm install && npm run dev
```

## License

MIT — see [LICENSE](LICENSE).

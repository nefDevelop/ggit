<div align="center">
  <h1>ggit 🚀</h1>
  <p><strong>A Modern, Ultra-Lightweight, and Secure Self-Hosted Git Service</strong></p>
</div>

---

## 🔮 Vision

The **ggit** project is a complete modernization of a legacy Git service, rebuilt from the ground up for the modern era. Powered by **Go (Gin)** on the backend and **SvelteKit** on the frontend, `ggit` aims to be the most painless, fastest, and most resource-efficient self-hosted Git service available.

We embrace a **Minimalist & Lean** philosophy: dropping heavy dependencies, standardizing on **SQLite**, and shifting rendering to the client to drastically reduce the binary footprint and server RAM usage.

## ✨ Features

### 🚀 Modern Architecture
- **SvelteKit Frontend**: Completely decoupled from Go templates. Uses a modern, reactive UI with Material 3 Expressive design, glassmorphism, and deep theme support (Deep Night, Ayu, Monokai, Nord).
- **Go + Gin Backend**: Pure Go approach. Replaced heavy legacy libraries with standard library implementations (`log/slog`, native `errors`, `goccy/go-json`) for maximum performance.
- **SQLite Only**: We dropped support for MySQL/PostgreSQL to focus entirely on an ultra-optimized SQLite experience (WAL mode, high concurrency, low footprint).

### ⚡ ggit vs Legacy Gogs
If you are coming from Gogs, `ggit` is a massive leap forward. Here is what `ggit` introduces:
- **Code Review & Branch Protection**: Built-in pull request approvals, turn-based reviewing (Semáforo de turno), and inline diff commenting.
- **Docs-as-Code (No more `.wiki.git`)**: Wikis are automatically generated from a `/docs` folder right inside your main repository.
- **Kanban & Unified Inbox**: Manage issues with virtual Kanban boards and track multi-repo tasks from a single inbox.
- **Anti-OOM & Self-Healing**: Dynamic Git process limiting based on available RAM, and automatic repair of missing Git hooks.
- **Secret Guardian**: Proactive detection of secrets (API keys, `.env` files) before they are even pushed to the repository.
- **Zero Server-Side Rendering**: Markdown, syntax highlighting, and UI are fully rendered on the client to save server CPU/RAM.

### 🛡️ Security & Integrity (Top Priority)
**Security and robust file handling are our core principles.** We mandate exhaustive testing for all features:
- **Comprehensive Testing**: All operations undergo strict testing for success, failure, and malicious inputs.
- **File System & DB Integrity**: We use a **Two-Phase Commit (2PC)** pattern for atomic consistency between SQLite and the file system. Write and delete operations are strictly verified.
- **Secret Guardian**: Proactive detection of secrets (API keys, `.env` files) during Git pushes.
- **Blindaje Anti-OOM**: Dynamic resource monitoring and Git concurrency limits based on available RAM.
- **Webhooks & Network Security**: HMAC-SHA256 payload signing and strict SSRF mitigation.

### 🛠️ Developer Productivity
- **Code Review & Branch Protection**: Require PRs, inline commenting, and approval workflows.
- **Docs-as-Code**: Legacy wikis have been replaced by a dynamic "Wiki-Folder" system (reading from `/docs` or `.docs` directly from your repository).
- **Unified Inbox & Kanban**: Manage issues with virtual Kanban boards and track all your tasks across repositories in one place.
- **Integrated Gists & Technotes**: Built-in support for code snippets and ADRs (Architecture Decision Records).

## 🛠️ Tech Stack

### Frontend
- **Framework**: [Svelte 5](https://svelte.dev/) + [SvelteKit](https://kit.svelte.dev/).
- **Tooling**: [Vite 8](https://vitejs.dev/) (Bundler), [TypeScript 6](https://www.typescriptlang.org/).
- **UI & Markdown**: [Lucide Svelte](https://lucide.dev/) (Icons), [Marked](https://marked.js.org/) + [DOMPurify](https://github.com/cure53/DOMPurify) (Client-side markdown rendering), [Highlight.js](https://highlightjs.org/) (Syntax highlighting).
- **Validation**: [Zod](https://zod.dev/) for robust client-side schemas.
- **Testing**: [Vitest](https://vitest.dev/) (Unit) + [Playwright](https://playwright.dev/) (E2E).

### Backend
- **Core**: [Go 1.25+](https://go.dev/).
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) for high-performance routing.
- **Database**: [GORM](https://gorm.io/) with a Pure Go [SQLite](https://github.com/glebarez/sqlite) driver.
- **Performance & Observability**: [OpenTelemetry](https://opentelemetry.io/) for tracing, `slog` for structured logging, `goccy/go-json` for ultra-fast JSON processing.

## 💾 Hardware Requirements

`ggit` is designed to run on a potato. 
- **Minimum**: A Raspberry Pi or a $5 VPS (512MB RAM, 1 CPU core) is more than enough for small teams.
- **Optimization**: Thanks to our aggressive removal of background workers, OOM protections, and L2 cache trimming, `ggit` maintains a near-zero idle RAM footprint.

## 📦 Installation

*(Installation instructions are currently being updated for the new SvelteKit + Gin architecture. Stay tuned!)*

For development, you can run the backend and frontend separately:
1. **Backend**: `go run main.go`
2. **Frontend**: `cd frontend && npm install && npm run dev`

## 🧪 Testing Guidelines

Contributors **must** follow our strict testing philosophy:
- **Exhaustive Scenarios**: Tests must cover "happy paths" (success), edge cases (failures), and malicious attempts.
- **Security First**: Any feature involving authentication, authorization, or cryptography must have rigorous security tests.
- **File Operations**: Features involving reading, writing, or deleting files on the host system are critical. They must have robust tests simulating permission errors, path traversals, and disk full scenarios.

## ⚖️ License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for the full license text.

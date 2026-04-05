# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

`admin-lab` is **not a production business system**. It is an external-network experimentation and asset incubation platform serving three purposes:

1. **需求仿真 (Simulation)** — Full-stack prototypes of real business requirements before migrating to internal systems
2. **组件示例 (Component Demo)** — Isolated validation of unfamiliar components, libraries, or interactions
3. **复用组件 (Reusable)** — Polished, migratable assets: generic table, search bars, file upload, charts, etc.

Before adding any feature, verify it serves real backend R&D needs and has experimental, reuse, or migration value. Read `docs/project-handbook.md` for full governance rules.

## Tech Stack

**Frontend:** Vue 3 + Vite 6 + Pinia + Element Plus + UnoCSS + ECharts  
**Backend:** Go 1.22+ + Gin + GORM + MySQL 8.0 + Casbin + Redis (optional) + Zap

## Commands

### Frontend (`cd web`)

```bash
npm install          # Install dependencies
npm run dev          # Dev server at http://localhost:8080 (proxies API to :8888)
npm run build        # Production build → dist/
npm run preview      # Preview production build
npx eslint .         # Lint
```

### Backend (`cd server`)

```bash
go run .             # Start dev server at :8888
go build             # Compile binary
go test ./...        # Run tests
```

### Full-stack containerized build (root)

```bash
make build-local          # Build web + server
make build-web-local      # Build frontend only
make build-server-local   # Build backend only
```

## Architecture

### Backend layering

Strict one-way dependency: `router → api → service → model`

Lab modules are organized by zone:
- `server/api/v1/lab/simulation/` — simulation handlers
- `server/api/v1/lab/componentdemo/` — component demo handlers
- `server/api/v1/lab/reusable/` — reusable module handlers
- Same pattern for `router/`, `service/`, `model/lab/`

**Do not** add lab code to `example/` packages.

### Frontend structure

- `src/view/lab/simulation/` — simulation pages
- `src/view/lab/component-demo/` — component demo pages
- `src/view/lab/reusable/` — reusable module pages
- `src/api/` — all HTTP request modules
- `src/components/lab/` — shared lab components (table-pro, crud-form-dialog, list-query-bar, etc.)
- `src/pinia/` — global state stores (app, user, router, dictionary, params)

### Routing

Routing mode is set in `src/core/config.js`. Currently static menu mode is enabled (`openStaticMenu: true`). Mock login is also configurable here for standalone frontend testing.

### Backend startup

`initialize/` runs on startup: DB connection, Casbin policy load, auto-migration (`AutoMigrate`), and seed data injection. Configuration is loaded from `config.yaml` (override with `config.debug.yaml`).

## Feature Classification Flow

New capabilities should evolve through zones in order:
1. First attempt → `组件示例` (componentdemo)
2. Grows into a full page/flow → `需求仿真` (simulation)
3. Has cross-project reuse value → `复用组件` (reusable)

Each zone maps to its own menu path under `实验室 / <zone> / <feature>`.

## Key Configuration

- `web/.env.development` — frontend dev env vars (API base URL, ports)
- `server/config.yaml` — MySQL, Redis, JWT, file storage, rate limiting
- MySQL database name: `admin_lab`, default port `3306`
- Backend port: `8888`, frontend dev port: `8080`

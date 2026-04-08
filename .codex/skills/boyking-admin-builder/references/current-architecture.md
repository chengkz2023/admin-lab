# Current Architecture

Use this reference when working on the current BoyKing Admin scaffold.

## Positioning

BoyKing Admin is a slimmed full-stack admin scaffold derived from gin-vue-admin.
It is no longer treated as the original upstream project.

## Scaffold source repository

Use this repository as the current scaffold source when the user asks to initialize a blank project from the scaffold:

`https://github.com/chengkz2023/My-GVA.git`

## Blank-project bootstrap guidance

### Safe cases

You can usually proceed directly when:

- the target directory is empty, or
- the target directory only contains `.git` and no real project files yet

### Caution cases

Pause and confirm with the user when:

- the directory already contains source files
- there is existing git history the user may want to preserve
- pulling the scaffold would mix with unrelated files

### Bootstrap intent

After the scaffold is pulled, the project should already contain the current BoyKing Admin baseline rather than the original gin-vue-admin defaults.

## Current technical baseline

### Backend

- Go + Gin
- GORM with MySQL as the main relational database
- Redis as optional cache/session support
- Local file storage only
- Casbin-based permission model

### Frontend

- Vue 3
- Vite
- Pinia
- Element Plus

## Current scope that should be preserved

- Super-admin centered management flow
- User, role, menu, API, dictionary, parameter, and operation log management
- Backend startup seeding of default system data
- Minimal configuration and minimal runtime branding

## Removed or intentionally trimmed capabilities

Do not casually restore these:

- plugin architecture
- AI or MCP features
- auto-code generation
- swagger docs exposure
- multi-OSS support
- Mongo and extra database initializers
- frontend first-install initialization entry
- unrelated demo, marketing, and example pages

## Configuration rules

### Commit-safe config

Use `server/config.yaml` for safe defaults only.

Expected style:

- localhost or placeholder database host
- placeholder passwords
- no real production credentials
- project naming aligned with BoyKing Admin

### Local debug config

Use `server/config.debug.yaml` for local developer overrides.
Do not commit it.

## Backend implementation checklist

When adding or changing backend features, check whether you need to update:

- `server/model/...`
- `server/service/...`
- `server/api/v1/...`
- `server/router/...`
- matching `enter.go` files
- `server/source/system/...` for seeded APIs, menus, or roles
- startup seeding logic if defaults changed

## Frontend implementation checklist

When adding or changing frontend features, check whether you need to update:

- `web/src/api/...`
- `web/src/view/...`
- `web/src/router/...`
- `web/src/pinia/...`
- menu or authority-linked identifiers that must match backend data

## Menu and permission guidance

Prefer keeping new admin pages under the current management structure.
If a page is permission-controlled:

- define the backend API and permission point first
- keep menu names, route names, and button identifiers stable
- update seeded menu/API data when the feature should exist in a fresh install

## Validation commands

Use these by default after substantial changes:

- backend: `go build -buildvcs=false ./core ./initialize ./service ./service/system ./source/system`
- frontend: `npm run build`

## When to push back

Pause and confirm with the user before proceeding if the request would:

- reintroduce removed platform features
- broaden the scaffold into a general plugin marketplace again
- add non-MySQL primary database assumptions
- add production secrets to tracked config files
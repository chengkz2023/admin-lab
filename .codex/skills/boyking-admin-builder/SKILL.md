---
name: boyking-admin-builder
description: Extend, customize, or bootstrap the BoyKing Admin full-stack scaffold that was slimmed from gin-vue-admin. Use when working on BoyKing Admin projects to initialize a blank project from the BoyKing Admin remote repository, add business modules, adjust system management features, update menus, permissions, seed data, configuration, or keep new code aligned with the current scaffold boundaries and conventions.
---

# BoyKing Admin Builder

Use this skill when building on the current BoyKing Admin scaffold rather than the original upstream gin-vue-admin feature set.

## Quick start

1. Read `references/current-architecture.md` before making substantial changes.
2. If the user is starting from a blank project, check the workspace state first and use the bootstrap workflow below.
3. Inspect the existing implementation in the target area before adding files.
4. Preserve the current slimmed scope unless the user explicitly asks to restore a removed capability.
5. Implement changes end-to-end: backend, frontend, seed data, and validation.

## Bootstrap a blank project

Use this flow when the user says to initialize a new project from the BoyKing Admin scaffold.

### Detect the workspace state

- If the current directory is empty, clone the scaffold directly into it.
- If the current directory only contains `.git` and no actual project files, connect it to the BoyKing Admin remote and pull the scaffold into place.
- If the directory already contains user files, stop and confirm before doing anything that could overwrite or mix contents.

### Bootstrap source

Use this repository as the scaffold source:

`https://github.com/chengkz2023/My-GVA.git`

### Preferred behavior

- In a truly empty directory, prefer cloning the scaffold into the current directory.
- In a git-initialized but otherwise blank directory, prefer wiring the remote and checking out `main`.
- After bootstrapping, verify the expected top-level structure exists and then continue as a normal BoyKing Admin project.

### Post-bootstrap checks

After pulling the scaffold, inspect and mention at least these items:

- `server/`
- `web/`
- `server/config.yaml`
- ignored local debug config expectations
- whether any project-specific rename or remote replacement is still needed

## Working rules

- Treat the current project as a slimmed full-stack scaffold, not a plugin platform.
- Reuse existing patterns in `server/api`, `server/service`, `server/router`, `server/model`, `web/src/api`, and `web/src/view` before inventing new ones.
- Keep configuration minimal. Use `server/config.yaml` for safe defaults and `server/config.debug.yaml` for local debugging only.
- Keep secrets, local IPs, and developer-only files out of committed config.
- Prefer preserving the existing super-admin permission and menu model unless the user asks for a broader redesign.

## Delivery workflow

### 1. Build context

- Identify whether the request is bootstrap only, backend only, frontend only, or full-stack.
- Inspect neighboring files first and follow their naming and layering patterns.
- Check whether the change needs menu data, API permissions, Casbin rules, or seeded defaults.

### 2. Backend changes

When adding or changing a backend module, keep the dependency flow:

`router -> api -> service -> model`

Typical touch points:

- `server/model/<domain>/`
- `server/model/<domain>/request/`
- `server/service/<domain>/`
- `server/api/v1/<domain>/`
- `server/router/<domain>/`
- corresponding `enter.go` aggregators

Rules:

- Keep HTTP concerns in API handlers.
- Keep database and business logic in services.
- Return structured errors upward instead of formatting HTTP responses in services.
- Follow existing request and response shapes already used in the scaffold.

### 3. Frontend changes

Typical touch points:

- `web/src/api/`
- `web/src/view/`
- `web/src/router/`
- `web/src/pinia/`

Rules:

- Use the existing request wrapper in `web/src/utils/request.js`.
- Follow current Element Plus and Vue 3 patterns already present in the repo.
- Keep pages aligned with the current admin structure rather than bringing back removed demo pages.
- If a page needs menu access or button permissions, keep names and identifiers consistent with backend seed data.

### 4. Seed data and permissions

If a change affects startup defaults, update the relevant files in `server/source/system/` and any startup seeding logic that loads them.

Typical cases:

- new menu or button permissions
- new protected APIs
- default dictionaries or parameters
- default role bindings

Do not assume the old first-install frontend flow exists. The current scaffold relies on backend startup seeding.

### 5. Guardrails

Do not reintroduce these removed capabilities unless the user explicitly asks:

- plugin system
- AI or MCP integration
- auto-code generation
- swagger exposure
- multi-database support beyond the current minimal setup
- object storage providers beyond local storage
- legacy demo or marketing pages
- frontend database initialization entry

### 6. Validation

Prefer validating with the same commands already used in this scaffold:

- backend: `go build -buildvcs=false ./core ./initialize ./service ./service/system ./source/system`
- frontend: `npm run build`

If a change touches other packages, expand validation as needed, but avoid claiming success without running something relevant.

## Response style for future use

When using this skill on real tasks:

- state assumptions briefly
- make concrete file edits instead of only giving advice when implementation is expected
- call out any place where a request would conflict with the current slimmed architecture
- mention validation status clearly

## References

- Read `references/current-architecture.md` for the current BoyKing Admin boundaries, retained modules, bootstrap notes, and implementation checklist.
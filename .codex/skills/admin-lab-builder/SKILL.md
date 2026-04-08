---
name: admin-lab-builder
description: Use this skill whenever you are coding, reviewing, or planning anything in the admin-lab project — including adding new pages, backend APIs, reusable modules, seed data, menus, Casbin rules, or migration notes. Always activate for tasks that mention admin-lab, lab modules, 需求仿真, 组件示例, 复用组件, BoyKing Admin lab extensions, or any front/back-end development inside this repository, even if the user just says "add a page" or "build a feature" without naming the skill explicitly.
---

# Admin Lab Builder

Apply this skill for any code, review, or planning task inside the `admin-lab` repository.

## Why this project exists

`admin-lab` is an **external experiment workspace** that serves three purposes:

1. **Simulate** real company requirements before porting them to the internal project.
2. **Validate** unfamiliar components, libraries, or interaction patterns at low risk.
3. **Extract** reusable frontend/backend assets ready for intranet migration.

It is not a production system. Keep that constraint in mind when scoping every task.

---

## Step 1 — Classify the work first

Every feature belongs to exactly one area:

| Area | Chinese name | Purpose |
|---|---|---|
| Demand simulation | 需求仿真 | Realistic end-to-end business page or flow prototypes |
| Component demo | 组件示例 | Minimal, standalone technical validation of a component |
| Reusable module | 复用组件 | Cross-project, configurable, migration-ready capabilities |

**Evolution path:** a capability usually moves `组件示例 → 需求仿真 → 复用组件` as it matures. Never force a feature into a later stage before it is ready.

**Admission test** — before starting, confirm the feature passes at least most of these:
- Clear background (why are we building this?)
- Assigned to one area
- Minimum runnable outcome defined
- Expected to be migrated to the intranet later
- Someone could understand it from docs alone

---

## Step 2 — Implementation patterns

### Backend: `router → api → service → model`

Maintain strict separation:

- **router** — only registers routes, nothing else
- **api** — receives params, calls service, returns response; never formats responses inside service
- **service** — all business logic and data processing
- **model** — struct definitions only

**Package paths for new work:**

```
server/api/v1/lab/simulation/
server/api/v1/lab/componentdemo/
server/api/v1/lab/reusable/

server/router/lab/simulation/
server/router/lab/componentdemo/
server/router/lab/reusable/

server/service/lab/simulation/
server/service/lab/componentdemo/
server/service/lab/reusable/

server/model/lab/
```

Never dump lab code into vague packages like `example`.

### Frontend: Vue 3 + Vite + Pinia + Element Plus

**Page location pattern:**

```
web/src/view/lab/simulation/<feature>.vue
web/src/view/lab/component-demo/<feature>.vue
web/src/view/lab/reusable/<feature>.vue
```

**API file location:** `web/src/api/<feature>.js`

Use the existing request wrapper. Match Element Plus patterns already in the project. Keep page state, API calls, and rendering logic clearly separated.

---

## Step 3 — Seed data (always check these four files)

When adding a new feature, update seed data so fresh installs and existing databases both get the feature automatically.

### 1. Menu — `server/source/system/menu.go`

Follow the existing naming convention exactly:

```go
// Parent level: lab, labSimulation, labComponentDemo, labReusable
// Child level: lab<Area><FeaturePascalCase>
// Example for a new reusable page called "batch-toolbar":
{
    MenuLevel: 2,
    Hidden:    false,
    ParentId:  menuNameMap["labReusable"],
    Path:      "batch-toolbar",
    Name:      "labReusableBatchToolbar",
    Component: "view/lab/reusable/batch-toolbar.vue",
    Sort:      <next available number>,
    Meta:      Meta{Title: "批量操作工具条", Icon: "<element-plus-icon>"},
},
```

Add the Name string (`"labReusableBatchToolbar"`) to the role-menu binding slice below.

### 2. Role-menu binding — `server/source/system/authorities_menus.go`

Append the new menu Name to the default admin role's slice so it is visible after a fresh install.

### 3. API seed — `server/source/system/api_<feature>.go`

Copy an existing file (e.g. `api_table_pro.go`) and adapt. Key things to set:
- `initOrder<Feature>` constant = previous feature's order + 1
- `InitializerName()` returns a unique string like `"sys_apis_batch_toolbar"`
- `entities` slice listing all route paths and HTTP methods

### 4. Casbin seed — `server/source/system/casbin_<feature>.go`

Copy `casbin_table_pro.go` and adapt. Rules use `V0: "888"` (super-admin role) by default:

```go
{Ptype: "p", V0: "888", V1: "/batchToolbar/list", V2: "POST"},
```

---

## Step 4 — New feature checklist

Use this as a completion gate before considering the task done:

- [ ] Area classified (`simulation` / `componentdemo` / `reusable`)
- [ ] Backend: router, api, service, model files created in correct packages
- [ ] Frontend: page vue file + api js file in correct paths
- [ ] Menu entry added to `menu.go`
- [ ] Menu name added to role binding in `authorities_menus.go`
- [ ] API seed file created (`api_<feature>.go`)
- [ ] Casbin seed file created (`casbin_<feature>.go`)
- [ ] For reusable modules: usage notes, configurable props, migration notes documented
- [ ] Build passes: `go build -buildvcs=false ./core ./initialize ./service ./service/system ./source/system ./api/v1/... ./router/...`
- [ ] Frontend build passes: `npm run build` (run in `web/`)

---

## Step 5 — Migration notes (reusable modules)

If the module belongs to `复用组件`, capture at minimum:

- What problem it solves
- How to integrate it (props / events / API contract)
- Configurable options
- Files to copy when migrating: vue component, api.js, backend service+model, seed data files
- Known limitations or intranet-specific adjustments needed

---

## Guardrails

Push back clearly when a request would:

- Turn `admin-lab` into a formal production system
- Add showcase pages with no experiment or reuse value
- Introduce a conflicting architecture without clear justification
- Commit real credentials, production data, or unsafe config values
- Create one-off throwaway code with no demo, simulation, or reuse value

---

## Configuration files

- `server/config.yaml` — only safe, committable default values
- `server/config.debug.yaml` — local-only overrides, not committed
- Never commit real DB, Redis, tokens, or secrets

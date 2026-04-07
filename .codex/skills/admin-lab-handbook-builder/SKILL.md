---
name: admin-lab-handbook-builder
description: 基于 admin-lab 协作手册实施、重构、评审实验室功能。Use when the task needs handbook-governed delivery, including module classification (`simulation` / `component-demo` / `reusable`), BoyKing Admin style package layout, menu/API/Casbin/seed-data alignment, and migration-oriented documentation for `admin-lab`.
---

# Admin Lab Handbook Builder

按以下流程执行任务，避免“只写页面/只写接口”的半成品交付。

## 1) 先判断是否应该进入仓库

先用两问过滤需求：

1. 是否服务团队未来真实后台研发？
2. 是否具备实验价值、复用价值或迁移价值？

若两问都是否，先明确建议“本地试验即可，不建议入仓”。

## 2) 先归类再动手

每个功能只能归属一个主分区：

- `simulation`：真实业务流程预演（场景完整闭环）
- `component-demo`：新组件/新交互验证（最小独立示例）
- `reusable`：跨项目可搬运能力（可配置、可裁剪）

若需求描述混合多个目标，先拆成主目标 + 次目标，再只为主目标选分区。

## 3) 代码落地遵循固定骨架

后端统一遵循 `router -> api -> service -> model`，并按实验室分区建包：

- `server/router/lab/{simulation|componentdemo|reusable}`
- `server/api/v1/lab/{simulation|componentdemo|reusable}`
- `server/service/lab/{simulation|componentdemo|reusable}`
- `server/model/lab/...`

前端按菜单分区组织：

- `web/src/view/lab/simulation/...`
- `web/src/view/lab/component-demo/...`
- `web/src/view/lab/reusable/...`
- `web/src/api/...`

禁止把实验室新代码继续堆进语义模糊目录（如 `example`）。

## 4) 不遗漏“菜单 + 权限 + 种子”

新增实验室功能时，默认同步检查：

- 菜单初始化：`server/source/system/menu.go`
- 角色菜单绑定：`server/source/system/authorities_menus.go`
- API/Casbin 种子：`server/source/system/` 下对应文件

目标是同时满足：

- Fresh install 能看到基础入口
- 老库升级后能补齐缺失菜单/权限/API 规则

## 5) 交付必须面向迁移

输出结果时，始终补这几项：

- 功能背景与目标
- 分区归属与原因
- 前后端改动清单
- 菜单/API/权限/种子变更点
- 迁移到内网需同步的文件、依赖、配置
- 已知限制与后续演进建议

若是 `reusable`，额外补：

- 使用场景
- 接入方式
- 可配置项

## 6) 执行和验证约束

默认验证命令：

- 后端：`go build -buildvcs=false ./core ./initialize ./service ./service/system ./source/system ./api/v1/... ./router/...`
- 前端：`npm run build`

若功能涉及交互，再明确建议补充：

- 页面手工验证
- 菜单显示验证
- 权限验证
- 导入导出等真实文件验证（如适用）

## 7) 响应输出模板

完成任务后，按以下结构给出结果：

1. `归属分区`：`simulation` / `component-demo` / `reusable`
2. `目标闭环`：本次覆盖了哪些页面、接口、数据流
3. `改动摘要`：前端、后端、菜单权限种子各改了什么
4. `验证记录`：执行了哪些命令，结果如何
5. `迁移提示`：搬运内网时要同步什么

## 8) 参考手册

优先读取：

- `docs/project-handbook.md`
- `README.md`
- 仓库根 `AGENTS.md`

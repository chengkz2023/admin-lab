# 字典消费端 + 业务操作日志 — 设计文档

**日期：** 2026-04-05  
**目标：** 以提升内网开发效率为核心，沉淀两个高频复用模块到 admin-lab 复用组件区，供内网项目直接迁入。

---

## 背景

当前 admin-lab 已收录：Table Pro、List Query Bar、CRUD Form Dialog、Security Chart Panel、Excel I/O、Reliable Upload Kit。

新增两个方向：

1. **字典消费端组件**：系统字典管理后台已有，痛点在每次新项目都要手写 select、tag、翻译逻辑，需要统一消费层。
2. **业务操作日志**：系统级请求日志已有，缺少"谁对哪条记录做了什么"的业务语义日志及其前端展示组件。

---

## 模块一：字典消费端

### 目标

提供三种开箱即用的字典消费方式，后台配置字典后，业务页面直接引用，无需手写映射。

### 方案选择

采用**登录后全量预加载 + Pinia 全局缓存**方案。

- 企业后台字典类型通常 < 50 个，一次加载代价可忽略
- 组件读取同步，无异步状态管理负担
- 迁入成本最低：复制 store + 两个组件即可

### 文件清单

| 文件 | 职责 |
|---|---|
| `web/src/api/dict.js` | 调用系统字典列表接口 |
| `web/src/pinia/modules/dict.js` | 全局字典 store，提供 `getItems` / `getLabel` |
| `web/src/components/lab/dict-select.vue` | 字典下拉选择器组件 |
| `web/src/components/lab/dict-tag.vue` | 字典值 Tag 展示组件 |
| `web/src/view/lab/reusable/dict-usage.vue` | demo 页，展示三种消费用法 |

后端**不新增代码**，复用系统字典已有接口。

### 数据流

```
登录成功
  → permission.js 调用 dictStore.loadAll()
  → 请求系统字典接口
  → 存入 dictMap: { [dictType]: [{label, value, cssClass}] }

业务页面
  DictSelect  → store.getItems(dictType)    → 渲染 el-select options
  DictTag     → store.getLabel(dictType, v) → 渲染带颜色的 el-tag
  纯文字翻译  → store.getLabel(dictType, v) → 模板插值
```

### 组件 API

**DictSelect**
```vue
<dict-select dict-type="order_status" v-model="form.status" />
```
Props: `dictType` (string, required), 透传所有 el-select 原生 props。

**DictTag**
```vue
<dict-tag dict-type="order_status" :value="row.status" />
```
Props: `dictType` (string, required), `value` (required)。颜色来自字典值的 `cssClass` 字段（`success / warning / danger / info`），由后台配置，前端直接映射到 `el-tag type`。

**纯翻译**
```js
import { useDictStore } from '@/pinia/modules/dict'
const dictStore = useDictStore()
dictStore.getLabel('order_status', row.status) // → "待审批"
```

### 错误处理

- 字典加载失败：打印 warn，store 返回空数组，组件渲染空 select/显示原始值，不阻断页面。
- 找不到 dictType：`getItems` 返回 `[]`，`getLabel` 返回原始 value 字符串。

---

## 模块二：业务操作日志

### 目标

提供轻量的业务语义日志记录工具和通用时间线展示组件，可嵌入任意业务详情页。

### 方案选择

采用 **Service 层手动埋点 + 通用日志表**方案。

- 日志有完整业务语义，可读性强
- `LogRecorder` 异步写入，不侵入业务事务
- 一张通用表，前端组件只需 `module + entityId` 即可查询

### 文件清单

**后端**

| 文件 | 职责 |
|---|---|
| `server/model/lab/reusable/biz_log.go` | BizLog 表结构定义 |
| `server/utils/bizlog/recorder.go` | LogRecorder 埋点工具 |
| `server/service/lab/reusable/biz_log.go` | 按 module + entityID 查询日志 |
| `server/api/v1/lab/reusable/biz_log.go` | HTTP handler |
| `server/router/lab/reusable/biz_log.go` | 路由注册 |
| `server/source/system/api_biz_log.go` | API seed 数据 |
| `server/source/system/casbin_biz_log.go` | Casbin 权限 seed |

**前端**

| 文件 | 职责 |
|---|---|
| `web/src/api/bizLog.js` | 请求日志列表接口 |
| `web/src/components/lab/biz-log-timeline.vue` | 时间线展示组件 |
| `web/src/view/lab/reusable/biz-log.vue` | demo 页 |

### 表结构

```sql
CREATE TABLE biz_log (
  id            BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
  module        VARCHAR(64)  NOT NULL COMMENT '业务模块，如 order / contract',
  entity_id     VARCHAR(64)  NOT NULL COMMENT '业务记录 ID',
  action        VARCHAR(64)  NOT NULL COMMENT '操作动作，如 approve / reject / edit',
  operator_id   BIGINT       NOT NULL COMMENT '操作人 ID',
  operator_name VARCHAR(64)  NOT NULL COMMENT '操作人姓名（冗余）',
  remark        VARCHAR(512) NOT NULL DEFAULT '' COMMENT '业务语义描述',
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_module_entity (module, entity_id),
  INDEX idx_created_at (created_at)
);
```

### 埋点方式

迁入内网后，service 层关键操作结束时调用一行：

```go
bizlog.Record(ctx, bizlog.Entry{
    Module:       "order",
    EntityID:     strconv.FormatUint(order.ID, 10),
    Action:       "approve",
    OperatorID:   claims.BaseClaims.ID,
    OperatorName: claims.Username,
    Remark:       fmt.Sprintf("订单 #%d 审批通过", order.ID),
})
```

`LogRecorder.Record` 内部：
- 异步写入（goroutine），不阻塞主流程
- 写失败仅打 warn 日志，不影响业务事务
- admin-lab 中使用全局 DB 实例，迁入时替换为项目自身 DB

### 查询接口

```
GET /bizLog/list?module=order&entityId=1001&page=1&pageSize=20
```

返回按 `created_at DESC` 排列的分页结果。

### 前端组件用法

```vue
<!-- 嵌入详情页，组件自己负责请求和渲染，父页面零感知 -->
<biz-log-timeline module="order" :entity-id="String(orderId)" />
```

时间线每条记录展示：操作人姓名、action badge（使用 el-tag）、remark 文字、相对时间（如"3 小时前"）。

### 错误处理

- 接口失败：组件内展示空状态提示，不抛出到父页面。
- `remark` 为空时：仅显示操作人 + action，不报错。

---

## 迁移指引

### 字典消费端迁入步骤

1. 确认内网项目系统字典接口路径，修改 `api/dict.js` 的 URL。
2. 复制 `pinia/modules/dict.js`，在 `permission.js` 登录成功后调用 `dictStore.loadAll()`。
3. 复制 `dict-select.vue` 和 `dict-tag.vue`，全局注册或按需引入。

### 操作日志迁入步骤

1. 执行建表 SQL（或通过 GORM AutoMigrate）。
2. 复制 `server/utils/bizlog/recorder.go`，修改 DB 实例来源。
3. 在需要记录的 service 方法中埋点。
4. 注册路由和 Casbin 权限。
5. 复制前端 `biz-log-timeline.vue` 和 `api/bizLog.js`，在详情页引用组件。

---

## 不在本次范围内

- 字典管理后台（CRUD）：已有，不重复建设。
- 操作日志的删除/清理接口：admin-lab 中只做查询展示，清理策略由各内网项目自行决定。
- 日志的字段级 diff 记录：当前 remark 为人工描述字符串，不做自动 diff。

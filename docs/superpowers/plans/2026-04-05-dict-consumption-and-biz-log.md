# 字典消费端 + 业务操作日志 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 admin-lab 复用组件区新增字典消费端三件套（DictSelect / DictTag / getLabel）和业务操作日志模块（后端通用表 + LogRecorder + 前端时间线组件），供内网项目直接迁入。

**Architecture:** 字典消费端直接复用已有的 `useDictionaryStore`（按需懒加载 + 缓存），只新增两个无状态展示组件和一个 demo 页。操作日志遵循 router → api → service → model 的现有分层，后端新增 `biz_log` 表和异步写入工具 `utils/bizlog/recorder.go`，前端封装 `BizLogTimeline` 组件，按 `module + entityId` 拉取数据并渲染时间线。

**Tech Stack:** Vue 3 + Element Plus + useDictionaryStore（前端）；Go + Gin + GORM + global.GVA_DB（后端）

---

## 文件清单

### 新增文件

| 文件 | 职责 |
|---|---|
| `web/src/components/lab/dict-select.vue` | 字典下拉选择器，封装 el-select |
| `web/src/components/lab/dict-tag.vue` | 字典值 Tag，根据 extend 字段着色 |
| `web/src/view/lab/reusable/dict-usage.vue` | 字典消费 demo 页 |
| `server/model/lab/reusable/biz_log.go` | BizLog 表结构 |
| `server/utils/bizlog/recorder.go` | 异步写入工具 |
| `server/service/lab/reusable/biz_log.go` | 分页查询 Service |
| `server/api/v1/lab/reusable/biz_log.go` | HTTP handler |
| `server/router/lab/reusable/biz_log.go` | 路由注册 |
| `server/source/system/api_biz_log.go` | API seed |
| `server/source/system/casbin_biz_log.go` | Casbin seed |
| `web/src/api/bizLog.js` | 前端 API 调用 |
| `web/src/components/lab/biz-log-timeline.vue` | 时间线展示组件 |
| `web/src/view/lab/reusable/biz-log.vue` | 操作日志 demo 页 |

### 修改文件

| 文件 | 变更 |
|---|---|
| `web/src/router/staticRoutes.js` | 新增两条菜单路由 |
| `server/service/lab/reusable/enter.go` | 注册 BizLogService |
| `server/api/v1/lab/reusable/enter.go` | 注册 BizLogApi + bizLogService 变量 |
| `server/router/lab/reusable/enter.go` | 注册 BizLogRouter + bizLogApi 变量 |
| `server/initialize/gorm_biz.go` | 加入 BizLog AutoMigrate |

---

## Task 1: DictSelect 组件

**Files:**
- Create: `web/src/components/lab/dict-select.vue`

- [ ] **Step 1: 创建组件文件**

```vue
<!-- web/src/components/lab/dict-select.vue -->
<template>
  <el-select v-bind="$attrs" :loading="loading" :disabled="loading || $attrs.disabled">
    <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    />
  </el-select>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { useDictionaryStore } from '@/pinia/modules/dictionary'

  defineOptions({ name: 'DictSelect', inheritAttrs: false })

  const props = defineProps({
    dictType: { type: String, required: true }
  })

  const dictStore = useDictionaryStore()
  const options = ref([])
  const loading = ref(false)

  onMounted(async () => {
    loading.value = true
    try {
      const items = await dictStore.getDictionary(props.dictType, 1)
      options.value = items || []
    } finally {
      loading.value = false
    }
  })
</script>
```

- [ ] **Step 2: 验证组件语法无错误**

```bash
cd web && npx vue-tsc --noEmit 2>&1 | grep dict-select || echo "no errors"
```

- [ ] **Step 3: Commit**

```bash
git add web/src/components/lab/dict-select.vue
git commit -m "feat(reusable): add DictSelect component backed by useDictionaryStore"
```

---

## Task 2: DictTag 组件

**Files:**
- Create: `web/src/components/lab/dict-tag.vue`

- [ ] **Step 1: 创建组件文件**

`extend` 字段若设为 `success / warning / danger / info` 则映射到 el-tag type；否则降级为默认样式。

```vue
<!-- web/src/components/lab/dict-tag.vue -->
<template>
  <el-tag :type="tagType">{{ label }}</el-tag>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { useDictionaryStore } from '@/pinia/modules/dictionary'

  defineOptions({ name: 'DictTag' })

  const props = defineProps({
    dictType: { type: String, required: true },
    value: { required: true }
  })

  const VALID_TYPES = ['success', 'warning', 'danger', 'info', '']
  const dictStore = useDictionaryStore()
  const label = ref(String(props.value))
  const tagType = ref('')

  onMounted(async () => {
    const items = await dictStore.getDictionary(props.dictType, 1)
    const matched = (items || []).find((i) => String(i.value) === String(props.value))
    if (matched) {
      label.value = matched.label
      tagType.value = VALID_TYPES.includes(matched.extend) ? matched.extend : ''
    }
  })
</script>
```

- [ ] **Step 2: Commit**

```bash
git add web/src/components/lab/dict-tag.vue
git commit -m "feat(reusable): add DictTag component with extend-based color mapping"
```

---

## Task 3: 字典消费 demo 页 + 路由

**Files:**
- Create: `web/src/view/lab/reusable/dict-usage.vue`
- Modify: `web/src/router/staticRoutes.js`

- [ ] **Step 1: 创建 demo 页**

```vue
<!-- web/src/view/lab/reusable/dict-usage.vue -->
<template>
  <div class="lab-page">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / 字典消费端</p>
        <h2>字典消费组件</h2>
        <p class="subtitle">
          基于 useDictionaryStore，提供 DictSelect、DictTag 和 getLabel
          三种消费方式，字典数据按需加载并自动缓存，迁入只需复制两个组件。
        </p>
      </div>
      <div class="hero-meta">
        <el-tag type="success">前端组件</el-tag>
        <el-tag type="primary">可迁移</el-tag>
        <el-tag>零后端依赖</el-tag>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="12">
        <el-card shadow="hover">
          <template #header><span class="card-title">DictSelect — 下拉选择</span></template>
          <p class="hint">用法：<code>&lt;dict-select dict-type="sex" v-model="val" /&gt;</code></p>
          <div class="demo-row">
            <dict-select dict-type="sex" v-model="selectVal" style="width: 200px" placeholder="请选择" />
            <span class="result">当前值：{{ selectVal ?? '—' }}</span>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header><span class="card-title">DictTag — Tag 展示</span></template>
          <p class="hint">用法：<code>&lt;dict-tag dict-type="sex" :value="1" /&gt;</code></p>
          <div class="demo-row" style="gap: 8px; flex-wrap: wrap">
            <dict-tag dict-type="sex" :value="1" />
            <dict-tag dict-type="sex" :value="2" />
            <dict-tag dict-type="sex" :value="999" />
          </div>
          <p class="hint" style="margin-top: 8px">
            颜色来自字典值 <code>extend</code> 字段，在字典管理后台配置为
            <code>success / warning / danger / info</code> 即可生效。
          </p>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="12">
        <el-card shadow="hover">
          <template #header><span class="card-title">getLabel — 纯文字翻译</span></template>
          <p class="hint">直接调用 store 方法，适合在 js/模板插值中翻译枚举值。</p>
          <pre class="code-block">import { useDictionaryStore } from '@/pinia/modules/dictionary'

const dictStore = useDictionaryStore()
const items = await dictStore.getDictionary('sex', 1)
const label = items.find(i => String(i.value) === String(val))?.label ?? val</pre>
        </el-card>

        <el-card shadow="hover">
          <template #header><span class="card-title">迁移说明</span></template>
          <div class="step-list">
            <div class="step-item"><span>1</span><div>确认内网项目已有 <code>useDictionaryStore</code>（gin-vue-admin 标准内置）。</div></div>
            <div class="step-item"><span>2</span><div>复制 <code>components/lab/dict-select.vue</code> 和 <code>dict-tag.vue</code>，全局注册或按需引入。</div></div>
            <div class="step-item"><span>3</span><div>在字典管理后台为每个字典值的 <code>extend</code> 字段填入颜色标记（可选）。</div></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import DictSelect from '@/components/lab/dict-select.vue'
  import DictTag from '@/components/lab/dict-tag.vue'

  defineOptions({ name: 'LabReusableDictUsage' })

  const selectVal = ref(null)
</script>

<style scoped>
  .lab-page { display: flex; flex-direction: column; gap: 16px; }
  .hero {
    display: flex; justify-content: space-between; gap: 16px;
    padding: 24px; border-radius: 16px;
    border: 1px solid #dbeafe;
    background: linear-gradient(135deg, #f0f9ff 0%, #f8fafc 100%);
  }
  .eyebrow { margin: 0 0 8px; color: #0369a1; font-size: 13px; font-weight: 700; letter-spacing: 0.08em; }
  .hero h2 { margin: 0 0 8px; color: #0f172a; font-size: 28px; }
  .subtitle { margin: 0; color: #475569; line-height: 1.7; max-width: 600px; }
  .hero-meta { display: flex; flex-wrap: wrap; align-content: flex-start; gap: 8px; }
  .card-title { font-weight: 600; }
  .hint { margin: 0 0 12px; color: #475569; font-size: 13px; }
  .demo-row { display: flex; align-items: center; gap: 16px; }
  .result { color: #0f172a; font-size: 14px; }
  .code-block {
    margin: 0; padding: 12px 16px; border-radius: 8px;
    background: #f1f5f9; font-family: Consolas, 'Courier New', monospace;
    font-size: 13px; line-height: 1.65; white-space: pre-wrap; word-break: break-word;
  }
  .step-list { display: flex; flex-direction: column; gap: 10px; }
  .step-item { display: grid; grid-template-columns: 28px 1fr; gap: 12px; align-items: start; }
  .step-item span {
    display: inline-flex; align-items: center; justify-content: center;
    width: 28px; height: 28px; border-radius: 999px;
    background: #0ea5e9; color: #fff; font-size: 13px; font-weight: 700;
  }
</style>
```

- [ ] **Step 2: 在 staticRoutes.js 中添加路由**

在 `reliable-upload` 路由条目之后插入：

```js
{ path: 'dict-usage', name: 'labReusableDictUsage', meta: { title: '字典消费组件', icon: 'collection-tag' }, component: 'view/lab/reusable/dict-usage.vue' },
```

即 `web/src/router/staticRoutes.js` 中 `labReusable` 的 `children` 数组末尾加一条。

- [ ] **Step 3: 启动前端 dev server 验证菜单和页面可正常渲染**

```bash
cd web && npm run dev
```

打开浏览器，确认侧边栏出现"字典消费组件"菜单，DictSelect 能展示 options，DictTag 能展示文字。

- [ ] **Step 4: Commit**

```bash
git add web/src/view/lab/reusable/dict-usage.vue web/src/router/staticRoutes.js
git commit -m "feat(reusable): add dict consumption demo page and static route"
```

---

## Task 4: BizLog 模型 + LogRecorder

**Files:**
- Create: `server/model/lab/reusable/biz_log.go`
- Create: `server/utils/bizlog/recorder.go`
- Modify: `server/initialize/gorm_biz.go`

- [ ] **Step 1: 创建 BizLog 表模型**

```go
// server/model/lab/reusable/biz_log.go
package reusable

import "time"

// BizLog 业务操作日志，按 module + entity_id 查询。
// 迁入内网时：执行 AutoMigrate 或手动建表；recorder 替换 DB 实例即可。
type BizLog struct {
	ID           uint      `json:"id"           gorm:"primaryKey;autoIncrement"`
	Module       string    `json:"module"       gorm:"size:64;not null;index:idx_module_entity"`
	EntityID     string    `json:"entityId"     gorm:"column:entity_id;size:64;not null;index:idx_module_entity"`
	Action       string    `json:"action"       gorm:"size:64;not null"`
	OperatorID   uint      `json:"operatorId"   gorm:"not null"`
	OperatorName string    `json:"operatorName" gorm:"size:64;not null"`
	Remark       string    `json:"remark"       gorm:"size:512"`
	CreatedAt    time.Time `json:"createdAt"    gorm:"index:idx_created_at"`
}

func (BizLog) TableName() string { return "biz_log" }
```

- [ ] **Step 2: 将 BizLog 加入 AutoMigrate**

修改 `server/initialize/gorm_biz.go`：

```go
package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	reusableModel "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable"
)

func bizModel() error {
	db := global.GVA_DB
	return db.AutoMigrate(&reusableModel.BizLog{})
}
```

- [ ] **Step 3: 创建 LogRecorder 工具**

```go
// server/utils/bizlog/recorder.go
package bizlog

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	reusableModel "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable"
	"go.uber.org/zap"
)

// Entry 是一条业务日志记录。
type Entry struct {
	Module       string
	EntityID     string
	Action       string
	OperatorID   uint
	OperatorName string
	Remark       string
}

// Record 异步写入业务日志。写失败只打 warn，不影响调用方业务流程。
func Record(ctx context.Context, e Entry) {
	go func() {
		row := reusableModel.BizLog{
			Module:       e.Module,
			EntityID:     e.EntityID,
			Action:       e.Action,
			OperatorID:   e.OperatorID,
			OperatorName: e.OperatorName,
			Remark:       e.Remark,
		}
		if err := global.GVA_DB.WithContext(ctx).Create(&row).Error; err != nil {
			global.GVA_LOG.Warn("bizlog: write failed", zap.String("module", e.Module), zap.String("entityId", e.EntityID), zap.Error(err))
		}
	}()
}
```

- [ ] **Step 4: 验证编译通过**

```bash
cd server && go build ./... 2>&1
```

期望：无输出（无错误）。

- [ ] **Step 5: Commit**

```bash
git add server/model/lab/reusable/biz_log.go server/utils/bizlog/recorder.go server/initialize/gorm_biz.go
git commit -m "feat(reusable): add BizLog model, AutoMigrate, and async LogRecorder"
```

---

## Task 5: BizLog Service + API + Router

**Files:**
- Create: `server/service/lab/reusable/biz_log.go`
- Modify: `server/service/lab/reusable/enter.go`
- Create: `server/api/v1/lab/reusable/biz_log.go`
- Modify: `server/api/v1/lab/reusable/enter.go`
- Create: `server/router/lab/reusable/biz_log.go`
- Modify: `server/router/lab/reusable/enter.go`

- [ ] **Step 1: 创建 BizLog Service**

```go
// server/service/lab/reusable/biz_log.go
package reusable

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	reusableModel "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable"
)

type BizLogService struct{}

type BizLogListResult struct {
	List  []reusableModel.BizLog `json:"list"`
	Total int64                  `json:"total"`
}

// List 按 module + entityID 倒序分页查询。
func (s *BizLogService) List(module, entityID string, page, pageSize int) BizLogListResult {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var rows []reusableModel.BizLog
	var total int64

	db := global.GVA_DB.Model(&reusableModel.BizLog{}).
		Where("module = ? AND entity_id = ?", module, entityID)

	db.Count(&total)
	db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows)

	return BizLogListResult{List: rows, Total: total}
}
```

- [ ] **Step 2: 注册到 ServiceGroup**

修改 `server/service/lab/reusable/enter.go`：

```go
package reusable

type ServiceGroup struct {
	BizLogService
	ExcelIOService
	ReliableUploadService
	SecurityDashboardService
	TableProService
}
```

- [ ] **Step 3: 创建 BizLog API Handler**

```go
// server/api/v1/lab/reusable/biz_log.go
package reusable

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type BizLogApi struct{}

func (b *BizLogApi) List(c *gin.Context) {
	module := c.Query("module")
	entityID := c.Query("entityId")
	if module == "" || entityID == "" {
		response.FailWithMessage("module 和 entityId 不能为空", c)
		return
	}

	page := 1
	pageSize := 20
	if v := c.Query("page"); v != "" {
		if n, err := parsePositiveInt(v); err == nil {
			page = n
		}
	}
	if v := c.Query("pageSize"); v != "" {
		if n, err := parsePositiveInt(v); err == nil {
			pageSize = n
		}
	}

	result := bizLogService.List(module, entityID, page, pageSize)
	response.OkWithDetailed(result, "获取业务日志成功", c)
}

func parsePositiveInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("invalid")
	}
	return n, nil
}
```

在文件顶部 import 中加 `"fmt"`。完整 import 块：

```go
import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)
```

- [ ] **Step 4: 注册到 API ApiGroup + 添加 bizLogService 变量**

修改 `server/api/v1/lab/reusable/enter.go`：

```go
package reusable

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	BizLogApi
	ExcelIOApi
	ReliableUploadApi
	SecurityDashboardApi
	TableProApi
}

var (
	bizLogService        = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.BizLogService
	excelIOService       = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ExcelIOService
	reliableUploadService    = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.ReliableUploadService
	securityDashboardService = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.SecurityDashboardService
	tableProService      = service.ServiceGroupApp.LabServiceGroup.ReusableServiceGroup.TableProService
)
```

- [ ] **Step 5: 创建 BizLog Router**

```go
// server/router/lab/reusable/biz_log.go
package reusable

import "github.com/gin-gonic/gin"

type BizLogRouter struct{}

func (r *BizLogRouter) InitBizLogRouter(Router *gin.RouterGroup) {
	bizLogRouter := Router.Group("bizLog")
	{
		bizLogRouter.GET("list", bizLogApi.List)
	}
}
```

- [ ] **Step 6: 注册到 RouterGroup + 添加 bizLogApi 变量**

修改 `server/router/lab/reusable/enter.go`：

```go
package reusable

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	BizLogRouter
	ExcelIORouter
	ReliableUploadRouter
	SecurityDashboardRouter
	TableProRouter
}

var (
	bizLogApi            = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.BizLogApi
	excelIOApi           = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.ExcelIOApi
	reliableUploadApi    = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.ReliableUploadApi
	securityDashboardApi = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.SecurityDashboardApi
	tableProApi          = api.ApiGroupApp.LabApiGroup.ReusableApiGroup.TableProApi
)
```

- [ ] **Step 7: 在路由初始化中注册 BizLog 路由**

搜索 `server/initialize/router_biz.go` 中 `InitReliableUploadRouter` 的调用位置，在其后加一行：

```go
reusableRouter.InitBizLogRouter(privateGroup)
```

- [ ] **Step 8: 验证编译通过**

```bash
cd server && go build ./... 2>&1
```

期望：无输出。

- [ ] **Step 9: Commit**

```bash
git add server/service/lab/reusable/biz_log.go server/service/lab/reusable/enter.go \
        server/api/v1/lab/reusable/biz_log.go server/api/v1/lab/reusable/enter.go \
        server/router/lab/reusable/biz_log.go server/router/lab/reusable/enter.go
git commit -m "feat(reusable): add BizLog service, API handler, and router"
```

---

## Task 6: BizLog Seed 数据

**Files:**
- Create: `server/source/system/api_biz_log.go`
- Create: `server/source/system/casbin_biz_log.go`

- [ ] **Step 1: 创建 API seed**

```go
// server/source/system/api_biz_log.go
package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApiBizLog = initOrderApiReliableUpload + 1

type initApiBizLog struct{}

func init() {
	system.RegisterInit(initOrderApiBizLog, &initApiBizLog{})
}

func (i *initApiBizLog) InitializerName() string { return "sys_apis_biz_log" }

func (i *initApiBizLog) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApiBizLog) TableCreated(ctx context.Context) bool { return false }

func (i *initApiBizLog) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entity := sysModel.SysApi{ApiGroup: "BizLog", Method: "GET", Path: "/bizLog/list", Description: "查询业务操作日志"}
	var existing sysModel.SysApi
	err := db.Where("path = ? AND method = ?", entity.Path, entity.Method).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx, db.Create(&entity).Error
	}
	return ctx, err
}

func (i *initApiBizLog) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return !errors.Is(db.Where("path = ? AND method = ?", "/bizLog/list", "GET").First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound)
}
```

- [ ] **Step 2: 创建 Casbin seed**

```go
// server/source/system/casbin_biz_log.go
package system

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderCasbinBizLog = initOrderCasbinReliableUpload + 1

type initCasbinBizLog struct{}

func init() {
	system.RegisterInit(initOrderCasbinBizLog, &initCasbinBizLog{})
}

func (i *initCasbinBizLog) InitializerName() string { return "casbin_biz_log" }

func (i *initCasbinBizLog) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initCasbinBizLog) TableCreated(ctx context.Context) bool { return false }

func (i *initCasbinBizLog) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/bizLog/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/bizLog/list", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/bizLog/list", V2: "GET"},
	}
	for _, entity := range entities {
		var existing adapter.CasbinRule
		err := db.Where(adapter.CasbinRule{Ptype: entity.Ptype, V0: entity.V0, V1: entity.V1, V2: entity.V2}).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if createErr := db.Create(&entity).Error; createErr != nil {
				return ctx, createErr
			}
			continue
		}
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func (i *initCasbinBizLog) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return !errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/bizLog/list", V2: "GET"}).First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound)
}
```

- [ ] **Step 3: 验证编译通过**

```bash
cd server && go build ./... 2>&1
```

期望：无输出。

- [ ] **Step 4: Commit**

```bash
git add server/source/system/api_biz_log.go server/source/system/casbin_biz_log.go
git commit -m "feat(reusable): add BizLog API and Casbin seed data"
```

---

## Task 7: 前端组件 + Demo 页

**Files:**
- Create: `web/src/api/bizLog.js`
- Create: `web/src/components/lab/biz-log-timeline.vue`
- Create: `web/src/view/lab/reusable/biz-log.vue`
- Modify: `web/src/router/staticRoutes.js`

- [ ] **Step 1: 创建 API 调用文件**

```js
// web/src/api/bizLog.js
import service from '@/utils/request'

export const getBizLogList = (params) =>
  service({ url: '/bizLog/list', method: 'get', params })
```

- [ ] **Step 2: 创建 BizLogTimeline 组件**

```vue
<!-- web/src/components/lab/biz-log-timeline.vue -->
<template>
  <div v-loading="loading">
    <el-empty v-if="!loading && !list.length" description="暂无操作记录" :image-size="60" />
    <el-timeline v-else>
      <el-timeline-item
        v-for="item in list"
        :key="item.id"
        :timestamp="formatTime(item.createdAt)"
        placement="top"
        type="primary"
      >
        <div class="log-item">
          <span class="operator">{{ item.operatorName }}</span>
          <el-tag size="small" style="margin: 0 6px">{{ item.action }}</el-tag>
          <span class="remark">{{ item.remark }}</span>
        </div>
      </el-timeline-item>
    </el-timeline>
    <div v-if="total > pageSize" class="pagination">
      <el-pagination
        small
        layout="prev, pager, next"
        :total="total"
        :page-size="pageSize"
        v-model:current-page="page"
        @current-change="load"
      />
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { getBizLogList } from '@/api/bizLog'

  defineOptions({ name: 'BizLogTimeline' })

  const props = defineProps({
    module: { type: String, required: true },
    entityId: { type: String, required: true }
  })

  const loading = ref(false)
  const list = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = 20

  const formatTime = (ts) => {
    if (!ts) return ''
    return new Date(ts).toLocaleString('zh-CN', { hour12: false })
  }

  const load = async () => {
    loading.value = true
    try {
      const res = await getBizLogList({ module: props.module, entityId: props.entityId, page: page.value, pageSize })
      if (res.code === 0) {
        list.value = res.data.list || []
        total.value = res.data.total || 0
      }
    } finally {
      loading.value = false
    }
  }

  onMounted(load)
</script>

<style scoped>
  .log-item { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
  .operator { font-weight: 600; color: #0f172a; }
  .remark { color: #475569; }
  .pagination { margin-top: 12px; display: flex; justify-content: center; }
</style>
```

- [ ] **Step 3: 创建 Demo 页**

```vue
<!-- web/src/view/lab/reusable/biz-log.vue -->
<template>
  <div class="lab-page">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / 业务操作日志</p>
        <h2>业务操作日志</h2>
        <p class="subtitle">
          通用日志表 + 异步 LogRecorder + BizLogTimeline 前端时间线组件。
          Service 层埋点一行代码，详情页嵌入组件一个标签，迁入只需复制
          utils/bizlog、model 和前端组件。
        </p>
      </div>
      <div class="hero-meta">
        <el-tag type="success">前后端</el-tag>
        <el-tag type="primary">可迁移</el-tag>
        <el-tag>通用日志表</el-tag>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="14">
        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>日志时间线 Demo</span>
              <div style="display: flex; gap: 8px; align-items: center">
                <el-input v-model="demoModule" placeholder="module" style="width: 120px" size="small" />
                <el-input v-model="demoEntityId" placeholder="entityId" style="width: 100px" size="small" />
                <el-button size="small" @click="writeTestLog" :loading="writing">写入测试日志</el-button>
                <el-button size="small" type="primary" @click="refreshTimeline">刷新</el-button>
              </div>
            </div>
          </template>
          <biz-log-timeline :key="timelineKey" :module="demoModule" :entity-id="demoEntityId" />
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="10">
        <el-card shadow="hover">
          <template #header><span class="panel-title">后端埋点示例</span></template>
          <pre class="code-block">import "github.com/your-org/server/utils/bizlog"

// 在 service 层操作完成后调用
bizlog.Record(ctx, bizlog.Entry{
    Module:       "order",
    EntityID:     strconv.FormatUint(order.ID, 10),
    Action:       "approve",
    OperatorID:   utils.GetUserID(c),
    OperatorName: utils.GetUserName(c),
    Remark:       fmt.Sprintf("订单 #%d 审批通过", order.ID),
})</pre>
        </el-card>

        <el-card shadow="hover">
          <template #header><span class="panel-title">前端组件用法</span></template>
          <pre class="code-block">&lt;biz-log-timeline
  module="order"
  :entity-id="String(orderId)"
/&gt;</pre>
          <p class="hint">组件自行请求数据，父页面零感知。支持分页（默认每页 20 条）。</p>
        </el-card>

        <el-card shadow="hover">
          <template #header><span class="panel-title">迁移说明</span></template>
          <div class="step-list">
            <div class="step-item"><span>1</span><div>执行建表或开启 AutoMigrate，确认 <code>biz_log</code> 表存在。</div></div>
            <div class="step-item"><span>2</span><div>复制 <code>server/utils/bizlog/recorder.go</code>，修改 DB 实例来源为项目自身。</div></div>
            <div class="step-item"><span>3</span><div>在 service 层关键操作后调用 <code>bizlog.Record()</code>。</div></div>
            <div class="step-item"><span>4</span><div>注册路由和 Casbin 权限，复制前端组件，在详情页引用。</div></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import BizLogTimeline from '@/components/lab/biz-log-timeline.vue'
  import service from '@/utils/request'

  defineOptions({ name: 'LabReusableBizLog' })

  const demoModule = ref('demo')
  const demoEntityId = ref('1')
  const timelineKey = ref(0)
  const writing = ref(false)

  const refreshTimeline = () => { timelineKey.value++ }

  // 直接调用 demo 接口写入测试日志（仅 admin-lab 内使用）
  const writeTestLog = async () => {
    writing.value = true
    try {
      const res = await service({
        url: '/bizLog/writeDemo',
        method: 'post',
        data: { module: demoModule.value, entityId: demoEntityId.value }
      })
      if (res.code === 0) {
        ElMessage.success('测试日志已写入')
        refreshTimeline()
      }
    } finally {
      writing.value = false
    }
  }
</script>

<style scoped>
  .lab-page { display: flex; flex-direction: column; gap: 16px; }
  .hero {
    display: flex; justify-content: space-between; gap: 16px;
    padding: 24px; border-radius: 16px;
    border: 1px solid #dbeafe;
    background: linear-gradient(135deg, #f0f9ff 0%, #f8fafc 100%);
  }
  .eyebrow { margin: 0 0 8px; color: #0369a1; font-size: 13px; font-weight: 700; letter-spacing: 0.08em; }
  .hero h2 { margin: 0 0 8px; color: #0f172a; font-size: 28px; }
  .subtitle { margin: 0; color: #475569; line-height: 1.7; max-width: 600px; }
  .hero-meta { display: flex; flex-wrap: wrap; align-content: flex-start; gap: 8px; }
  .panel-header { display: flex; justify-content: space-between; align-items: center; font-weight: 600; }
  .panel-title { font-weight: 600; }
  .code-block {
    margin: 0; padding: 12px 16px; border-radius: 8px;
    background: #f1f5f9; font-family: Consolas, 'Courier New', monospace;
    font-size: 12px; line-height: 1.65; white-space: pre-wrap; word-break: break-word;
  }
  .hint { margin: 8px 0 0; color: #475569; font-size: 13px; }
  .step-list { display: flex; flex-direction: column; gap: 10px; }
  .step-item { display: grid; grid-template-columns: 28px 1fr; gap: 12px; align-items: start; }
  .step-item span {
    display: inline-flex; align-items: center; justify-content: center;
    width: 28px; height: 28px; border-radius: 999px;
    background: #0ea5e9; color: #fff; font-size: 13px; font-weight: 700;
  }
</style>
```

> **注意：** demo 页中的"写入测试日志"按钮调用了 `/bizLog/writeDemo` 接口。需要在 Task 5 完成后，在 `biz_log.go` handler 中额外添加一个 `WriteDemo` 方法（POST，仅 admin-lab 内用于演示）并在 router 中注册 `bizLogRouter.POST("writeDemo", bizLogApi.WriteDemo)`。`WriteDemo` 实现：直接调用 `bizlog.Record` 写入一条 action="demo_action"、remark="这是一条测试日志" 的记录。Casbin seed 中同样加上对应 POST 权限。

- [ ] **Step 4: 在 staticRoutes.js 中添加路由**

在 `dict-usage` 条目之后插入：

```js
{ path: 'biz-log', name: 'labReusableBizLog', meta: { title: '业务操作日志', icon: 'document' }, component: 'view/lab/reusable/biz-log.vue' },
```

- [ ] **Step 5: 启动 dev server 验证**

```bash
cd web && npm run dev
```

确认：菜单出现"业务操作日志"，BizLogTimeline 可正常显示加载态和空状态；后端 `go run .` 后点击"写入测试日志"能成功写入并刷新时间线。

- [ ] **Step 6: Commit**

```bash
git add web/src/api/bizLog.js web/src/components/lab/biz-log-timeline.vue \
        web/src/view/lab/reusable/biz-log.vue web/src/router/staticRoutes.js
git commit -m "feat(reusable): add BizLogTimeline component and biz-log demo page"
```

---

## 自检结果

| 检查项 | 状态 |
|---|---|
| 字典消费端 spec 要求：DictSelect / DictTag / getLabel | ✅ Task 1-3 |
| 操作日志 spec 要求：通用日志表 | ✅ Task 4（BizLog model） |
| 操作日志 spec 要求：LogRecorder 异步写入 | ✅ Task 4（recorder.go） |
| 操作日志 spec 要求：前端时间线组件 | ✅ Task 7（BizLogTimeline） |
| 操作日志 spec 要求：按 module + entityId 查询 | ✅ Task 5（service List 方法） |
| 路由注册 | ✅ Task 3 / Task 7（staticRoutes.js） |
| 后端 seed（API + Casbin） | ✅ Task 6 |
| WriteDemo 接口（demo 页测试按钮依赖） | ⚠️ Task 7 注意事项中已标出，需在 Task 5 中补充 |
| 无 TBD / 无占位符 | ✅ |
| 类型一致性（BizLog 字段名在 service/model/前端一致） | ✅ createdAt / operatorName / entityId 统一 |

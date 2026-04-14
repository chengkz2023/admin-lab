<template>
  <div class="page-wrap">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / ElTablePro</p>
        <h2>面向复杂数据场景的 Element Plus 表格增强封装</h2>
        <p class="subtitle">
          这套组件提供多类型单元格渲染、多级表头、详情抽屉、分页联动和
          `useElTablePro` 数据装配 Hook，适合日志、监控、配置巡检等重数据列表页。
        </p>
      </div>
      <div class="hero-tags">
        <el-tag type="primary">复杂列渲染</el-tag>
        <el-tag type="success">Hook 驱动</el-tag>
        <el-tag>可迁移</el-tag>
      </div>
    </div>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>组件能力总览</span></div>
      </template>
      <el-table :data="featureRows" border>
        <el-table-column prop="feature" label="能力" min-width="200" />
        <el-table-column prop="desc" label="说明" min-width="380" />
        <el-table-column prop="config" label="关键点" min-width="260" />
      </el-table>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>接入要点</span></div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="组件路径">`web/src/components/ElTablePro/ElTablePro.vue`</el-descriptions-item>
        <el-descriptions-item label="渲染器目录">`web/src/components/ElTablePro/renderers/*`</el-descriptions-item>
        <el-descriptions-item label="Hook 路径">`web/src/hooks/useElTablePro.js`</el-descriptions-item>
        <el-descriptions-item label="列配置入口">通过 `columns` 描述列类型、层级、格式化和自定义渲染。</el-descriptions-item>
        <el-descriptions-item label="主要插槽">`toolbar / expand / actions / 自定义列名插槽`。</el-descriptions-item>
        <el-descriptions-item label="分页绑定">`v-model:current-page + v-model:page-size + @page-change`。</el-descriptions-item>
        <el-descriptions-item label="详情能力">`show-detail` 开启内置详情抽屉，适合表格行快速查看。</el-descriptions-item>
        <el-descriptions-item label="迁移建议">迁入内网时同步复制组件目录、Hook 导出、页面示例及依赖图标/Element Plus 配置。</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>Hook 使用方式</span></div>
      </template>
      <pre class="code-block"><code>{{ hookCode }}</code></pre>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>页面演示</span></div>
      </template>
      <component-example />
    </el-card>
  </div>
</template>

<script setup>
  import ComponentExample from '@/components/ElTablePro/Example.vue'

  defineOptions({
    name: 'LabReusableElTablePro'
  })

  const featureRows = [
    { feature: '多类型渲染器', desc: '内置状态、进度、IP、协议、端口、JSON、URL、开关、聚合单元格等常见后台字段展示。', config: 'renderers/*.vue + type 字段' },
    { feature: '多级表头', desc: '支持通过 children 递归定义分组列，适合监控指标和组合字段场景。', config: 'columns[].children' },
    { feature: '插槽扩展', desc: '支持 toolbar、expand、actions 及列名同名插槽，便于按业务补充操作区与定制单元格。', config: '#toolbar / #expand / #actions / #字段名' },
    { feature: '详情抽屉', desc: '支持内置详情查看按钮，快速浏览整行数据。', config: 'showDetail + actionsConfig' },
    { feature: '分页联动 Hook', desc: '通过 useElTablePro 统一维护分页、loading、请求参数拼装和结果转换。', config: 'useElTablePro(service, options)' }
  ]

  const hookCode = `import { onMounted, ref } from "vue";
import ElTablePro from "@/components/ElTablePro";
import { useElTablePro } from "@/hooks/useElTablePro";
import { getUserList as userListApi } from "@/api/user";

// 1. 查询表单统一用 ref 包裹对象。
// 这样在 setup 中传递、重置、整体替换时更直接。
const searchForm = ref({
  keyword: "",
  dateRange: [],
});

// 2. 初始化表格 Hook。
// Hook 内部负责维护 data、total、loading、currentPage、pageSize。
const { tableState, loadTableData, search } = useElTablePro(userListApi, {
  pageSize: 20,

  // 3. params 每次请求前都会执行一次，
  // 适合在这里把页面表单转换成后端需要的参数结构。
  params: () => ({
    key: searchForm.value.keyword,
    startTime: searchForm.value.dateRange?.[0] || "",
    endTime: searchForm.value.dateRange?.[1] || "",
  }),

  // 4. transform 用来处理接口返回列表，
  // 把展示层需要的派生字段提前算好，避免模板里塞过多判断。
  transform: (list) =>
    list.map((item) => ({
      ...item,
      _tagType: item.status === 1 ? "success" : "info",
      _statusText: item.status === 1 ? "启用" : "禁用",
    })),
});

// 5. 首次进入页面时主动加载一次数据。
onMounted(() => {
  loadTableData();
});

// 6. 查询按钮通常直接调用 search。
// search 会先把页码重置为 1，再触发 loadTableData。
const handleSearch = () => {
  search();
};

// 7. 重置时建议整体替换 ref.value，
// 再调用 search，保证筛选条件和分页都回到初始状态。
const handleReset = () => {
  searchForm.value = {
    keyword: "",
    dateRange: [],
  };
  search();
};`
</script>

<style scoped>
  .page-wrap {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .hero {
    display: flex;
    justify-content: space-between;
    gap: 16px;
    padding: 24px;
    border-radius: 16px;
    border: 1px solid #dbeafe;
    background: linear-gradient(135deg, #eff6ff 0%, #f8fafc 100%);
  }

  .eyebrow {
    margin: 0 0 8px;
    color: #1d4ed8;
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.08em;
  }

  .hero h2 {
    margin: 0 0 8px;
    font-size: 24px;
    color: #0f172a;
  }

  .subtitle {
    margin: 0;
    max-width: 760px;
    color: #475569;
    line-height: 1.75;
  }

  .hero-tags {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    gap: 8px;
  }

  .panel-title {
    font-weight: 600;
  }

  .code-block {
    margin: 0;
    padding: 14px;
    border-radius: 10px;
    border: 1px solid #e2e8f0;
    background: #0f172a;
    color: #e2e8f0;
    line-height: 1.6;
    white-space: pre-wrap;
    word-break: break-word;
    font-size: 13px;
  }

  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
    }
  }
</style>

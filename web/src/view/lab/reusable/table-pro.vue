<template>
  <div class="page-wrap">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / Table Pro</p>
        <h2>高复用 el-table 二次封装</h2>
        <p class="subtitle">
          Table Pro 专注表格本体能力，不内置查询栏和导出；推荐与独立查询组件、独立导出逻辑组合使用。
        </p>
      </div>
      <div class="hero-tags">
        <el-tag type="primary">纯表格容器</el-tag>
        <el-tag type="success">高可扩展</el-tag>
        <el-tag>内网可迁移</el-tag>
      </div>
    </div>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>外部查询栏演示</span></div>
      </template>
      <lab-list-query-bar
        v-model="queryModel"
        :items="queryItems"
        @search="handleSearch"
        @reset="handleReset"
      />
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>能力总览</span></div>
      </template>
      <el-table :data="featureRows" border>
        <el-table-column prop="feature" label="能力" min-width="200" />
        <el-table-column prop="desc" label="说明" min-width="380" />
        <el-table-column prop="config" label="关键配置" min-width="280" />
      </el-table>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>Props 配置</span></div>
      </template>
      <el-table :data="propsRows" border>
        <el-table-column prop="name" label="属性" min-width="180" />
        <el-table-column prop="type" label="类型" width="200" />
        <el-table-column prop="required" label="必填" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.required ? 'danger' : 'info'">
              {{ row.required ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="defaultValue" label="默认值" min-width="160" />
        <el-table-column prop="desc" label="说明" min-width="340" />
      </el-table>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>Events / Slots</span></div>
      </template>
      <el-tabs>
        <el-tab-pane label="事件 Events">
          <el-table :data="eventRows" border>
            <el-table-column prop="name" label="事件名" min-width="220" />
            <el-table-column prop="payload" label="参数" min-width="260" />
            <el-table-column prop="desc" label="说明" min-width="320" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="插槽 Slots">
          <el-table :data="slotRows" border>
            <el-table-column prop="name" label="插槽名" min-width="220" />
            <el-table-column prop="scope" label="作用域参数" min-width="260" />
            <el-table-column prop="desc" label="说明" min-width="320" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>columns 字段配置</span></div>
      </template>
      <el-table :data="columnRows" border>
        <el-table-column prop="name" label="字段" min-width="180" />
        <el-table-column prop="required" label="必填" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.required ? 'danger' : 'info'">
              {{ row.required ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="desc" label="说明" min-width="360" />
        <el-table-column prop="example" label="示例" min-width="260" />
      </el-table>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>接入示例</span></div>
      </template>
      <pre class="code-block"><code>{{ usageCode }}</code></pre>
    </el-card>

    <lab-table-pro
      ref="tableRef"
      title="订单列表（演示）"
      :columns="columns"
      :query-params="queryModel"
      :fetcher="fetchTableData"
      :reserve-selection="true"
      :row-class-name="rowClassName"
      persist-key="reusable-table-pro-demo"
      @selection-change="handleSelectionChange"
    >
      <template #toolbar-right="{ selectedRows }">
        <el-button
          type="warning"
          :disabled="!selectedRows.length"
          @click="mockBatchDone(selectedRows)"
        >
          批量完成
        </el-button>
      </template>

      <template #error="{ reload }">
        <el-result
          icon="error"
          title="加载异常"
          sub-title="可以重试或调整查询条件后再试"
        >
          <template #extra>
            <el-button type="primary" @click="reload({ resetPage: false })">重试</el-button>
          </template>
        </el-result>
      </template>

      <template #empty>
        <el-empty description="暂无匹配数据" />
      </template>

      <template #cell-status="{ value }">
        <el-tag size="small" :type="statusTag(value)">{{ value }}</el-tag>
      </template>

      <template #cell-amount="{ value }">
        <span>¥{{ Number(value || 0).toFixed(2) }}</span>
      </template>
    </lab-table-pro>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import LabTablePro from '@/components/lab/table-pro.vue'
  import LabListQueryBar from '@/components/lab/list-query-bar.vue'
  import { getTableProPage } from '@/api/tablePro'

  defineOptions({
    name: 'LabReusableTablePro'
  })

  const tableRef = ref()

  const queryModel = ref({
    keyword: '',
    status: '',
    owner: '',
    createdDateRange: []
  })

  const queryItems = ref([
    { prop: 'keyword', label: '关键字', type: 'input', placeholder: '订单号 / 客户 / 来源', span: 8 },
    {
      prop: 'status',
      label: '状态',
      type: 'select',
      options: [
        { label: '待处理', value: 'pending' },
        { label: '处理中', value: 'processing' },
        { label: '已完成', value: 'done' }
      ],
      span: 5
    },
    {
      prop: 'owner',
      label: '负责人',
      type: 'select',
      options: [
        { label: 'Alice', value: 'Alice' },
        { label: 'Bob', value: 'Bob' },
        { label: 'Carol', value: 'Carol' },
        { label: 'David', value: 'David' }
      ],
      span: 5
    },
    { prop: 'createdDateRange', label: '创建日期', type: 'dateRange', span: 8 }
  ])

  const columns = ref([
    { type: 'selection', configurable: false, width: 48 },
    { key: 'orderNo', label: '订单号', minWidth: 180, sortable: 'custom' },
    { key: 'customerName', label: '客户', minWidth: 140, sortable: 'custom' },
    { key: 'status', label: '状态', minWidth: 110, sortable: 'custom' },
    { key: 'priority', label: '优先级', minWidth: 100, sortable: 'custom' },
    { key: 'owner', label: '负责人', minWidth: 120, sortable: 'custom' },
    { key: 'source', label: '来源渠道', minWidth: 160, sortable: 'custom', showOverflowTooltip: true },
    { key: 'amount', label: '金额', minWidth: 120, sortable: 'custom' },
    { key: 'createdAt', label: '创建时间', minWidth: 120, sortable: 'custom' }
  ])

  const featureRows = [
    { feature: '外部 loading 覆盖', desc: '支持父级统一控制加载态，也保留组件内部请求 loading。', config: 'loading' },
    { feature: '空态与错态扩展', desc: '支持空态与错态槽位，避免每个业务页重复写兜底 UI。', config: '#empty / #error' },
    { feature: '分页字段可映射', desc: '可配置请求分页字段和响应列表/总数字段，适配不同后端协议。', config: 'pageField / pageSizeField / listField / totalField' },
    { feature: '排序映射钩子', desc: '支持自定义排序参数映射，如 asc/desc -> ASC/DESC。', config: 'sortMapper' },
    { feature: '跨页勾选保留', desc: '支持跨页勾选状态保留，适合批量处理场景。', config: 'reserveSelection + rowKey' },
    { feature: '行样式透传', desc: '支持按行动态样式，便于高亮风险行、置灰禁用行。', config: 'rowClassName' },
    { feature: '列设置持久化', desc: '支持列显示隐藏、排序与 localStorage 持久化。', config: 'enableColumnSetting / persistKey' }
  ]

  const propsRows = [
    { name: 'title', type: 'String', required: false, defaultValue: '结果列表', desc: '表格标题。' },
    { name: 'columns', type: 'Array', required: true, defaultValue: '[]', desc: '表格列配置数组。' },
    { name: 'queryParams', type: 'Object', required: false, defaultValue: '{}', desc: '外部查询参数对象。' },
    { name: 'fetcher', type: 'Function', required: true, defaultValue: '-', desc: '数据请求函数，返回列表和总数。' },
    { name: 'loading', type: 'Boolean', required: false, defaultValue: 'undefined', desc: '外部覆盖 loading 状态。' },
    { name: 'rowKey', type: 'String', required: false, defaultValue: 'id', desc: '行主键字段。' },
    { name: 'rowClassName', type: 'String | Function', required: false, defaultValue: "''", desc: '行 className 透传。' },
    { name: 'reserveSelection', type: 'Boolean', required: false, defaultValue: 'false', desc: 'selection 列是否跨页保留勾选。' },
    { name: 'immediate', type: 'Boolean', required: false, defaultValue: 'true', desc: '挂载后是否自动请求。' },
    { name: 'showPagination', type: 'Boolean', required: false, defaultValue: 'true', desc: '是否显示分页器。' },
    { name: 'pageSize / pageSizes', type: 'Number / Array', required: false, defaultValue: '10 / [10,20,50]', desc: '分页大小配置。' },
    { name: 'pageField / pageSizeField', type: 'String / String', required: false, defaultValue: 'page / pageSize', desc: '请求分页字段映射。' },
    { name: 'listField / totalField', type: 'String / String', required: false, defaultValue: 'list / total', desc: '响应字段映射，支持路径格式。' },
    { name: 'sortMapper', type: 'Function', required: false, defaultValue: 'null', desc: '排序参数映射钩子。' },
    { name: 'border / stripe', type: 'Boolean / Boolean', required: false, defaultValue: 'true / false', desc: '表格外观控制。' },
    { name: 'tableProps', type: 'Object', required: false, defaultValue: '{}', desc: '透传 el-table 其他属性。' },
    { name: 'enableColumnSetting / persistKey', type: 'Boolean / String', required: false, defaultValue: "true / ''", desc: '列设置与持久化配置。' },
    { name: 'showTotalTag', type: 'Boolean', required: false, defaultValue: 'true', desc: '是否显示总数标签。' }
  ]

  const eventRows = [
    { name: 'selection-change', payload: 'selectedRows', desc: '勾选行变化时触发。' },
    { name: 'loaded', payload: '{ rows, total, payload }', desc: '列表加载成功后触发。' },
    { name: 'error', payload: 'error', desc: '请求异常时触发。' }
  ]

  const slotRows = [
    { name: 'title', scope: '-', desc: '自定义标题。' },
    { name: 'toolbar-left / toolbar-right', scope: '{ selectedRows, query, reload }', desc: '工具栏扩展。' },
    { name: 'table-before / table-after', scope: '{ rows, selectedRows }', desc: '表格上下扩展。' },
    { name: 'empty', scope: '-', desc: '空数据态内容。' },
    { name: 'error', scope: '{ error, reload }', desc: '异常态内容与重试入口。' },
    { name: 'cell-字段名 或 columns[].slot', scope: '{ row, column, index, value }', desc: '单元格渲染扩展。' }
  ]

  const columnRows = [
    { name: 'key', required: true, desc: '列唯一 key，建议映射后端字段。', example: "key: 'orderNo'" },
    { name: 'prop', required: false, desc: '显示字段，缺省等于 key。', example: "prop: 'orderNo'" },
    { name: 'label', required: false, desc: '列标题。', example: "label: '订单号'" },
    { name: 'type', required: false, desc: '特殊列：selection / index。', example: "type: 'selection'" },
    { name: 'width / minWidth', required: false, desc: '列宽配置。', example: 'minWidth: 160' },
    { name: 'fixed / align', required: false, desc: '固定列及对齐方式。', example: "fixed: 'right'" },
    { name: 'sortable', required: false, desc: '建议 custom 走服务端排序。', example: "sortable: 'custom'" },
    { name: 'showOverflowTooltip', required: false, desc: '超长文本提示。', example: 'showOverflowTooltip: true' },
    { name: 'formatter', required: false, desc: '默认格式化函数。', example: '(row, col, value) => value' },
    { name: 'slot', required: false, desc: '自定义单元格槽位名。', example: "slot: 'cell-status'" },
    { name: 'configurable', required: false, desc: '是否参与列设置。', example: 'configurable: false' },
    { name: 'selectable', required: false, desc: 'selection 列可选控制。', example: '(row) => row.status !== "done"' }
  ]

  const statusTag = (status) => {
    if (status === 'done') {
      return 'success'
    }
    if (status === 'processing') {
      return 'warning'
    }
    return 'info'
  }

  const rowClassName = ({ row }) => {
    if (row.priority === 'P0') {
      return 'row-priority-high'
    }
    return ''
  }

  const fetchTableData = async (payload) => {
    const dateRange = payload?.filters?.createdDateRange || []
    const result = await getTableProPage({
      page: payload.page,
      pageSize: payload.pageSize,
      keyword: payload.filters?.keyword || '',
      status: payload.filters?.status || '',
      owner: payload.filters?.owner || '',
      startDate: dateRange[0] || '',
      endDate: dateRange[1] || '',
      sortBy: payload.sortBy || '',
      sortOrder: payload.sortOrder || ''
    })
    return {
      list: result?.data?.list || [],
      total: result?.data?.total || 0
    }
  }

  const handleSearch = () => {
    tableRef.value?.reload?.({ resetPage: true })
  }

  const handleReset = () => {
    tableRef.value?.reload?.({ resetPage: true })
  }

  const handleSelectionChange = (rows) => {
    void rows
  }

  const mockBatchDone = (rows) => {
    ElMessage.success(`已提交 ${rows.length} 条记录进行批量处理（演示）。`)
  }

  const usageCode = `<lab-list-query-bar
  v-model="queryModel"
  :items="queryItems"
  @search="() => tableRef?.reload({ resetPage: true })"
  @reset="() => tableRef?.reload({ resetPage: true })"
/>

<lab-table-pro
  ref="tableRef"
  :columns="columns"
  :query-params="queryModel"
  :fetcher="fetchTableData"
  :loading="outerLoading"
  :reserve-selection="true"
  :row-class-name="rowClassName"
  page-field="current"
  page-size-field="size"
  list-field="data.list"
  total-field="data.total"
  :sort-mapper="({ prop, order }) => ({ sortBy: prop, sortOrder: order === 'ascending' ? 'ASC' : order === 'descending' ? 'DESC' : '' })"
>
  <template #error="{ reload }">
    <el-button @click="reload()">重试</el-button>
  </template>
  <template #empty>
    <el-empty description="暂无数据" />
  </template>
</lab-table-pro>`
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

  :deep(.row-priority-high) {
    --el-table-tr-bg-color: #fff7ed;
  }

  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
    }
  }
</style>

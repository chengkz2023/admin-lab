<template>
  <div class="page-wrap">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / 列表查询栏</p>
        <h2>大数据列表页通用查询栏</h2>
        <p class="subtitle">
          通过配置驱动统一列表页搜索栏，支持输入、下拉、日期范围、级联和多选能力，减少业务页面重复开发。
        </p>
      </div>
      <div class="hero-tags">
        <el-tag type="primary">可复用</el-tag>
        <el-tag type="success">配置驱动</el-tag>
        <el-tag>低耦合</el-tag>
      </div>
    </div>

    <lab-list-query-bar
      v-model="queryForm"
      :items="queryItems"
      :loading="searching"
      :show-export="true"
      @search="handleSearch"
      @reset="handleReset"
      @export="handleExport"
    >
      <template #actions-left>
        <el-alert
          type="info"
          :closable="false"
          title="演示页中的查询是前端过滤，真实项目建议在 search 事件中对接后端分页接口。"
        />
      </template>
    </lab-list-query-bar>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>结果列表</span>
          <span class="total-text">共 {{ filteredTotal }} 条</span>
        </div>
      </template>
      <el-table :data="pagedRows" border>
        <el-table-column prop="orderNo" label="订单号" min-width="180" />
        <el-table-column prop="customerName" label="客户名称" min-width="160" />
        <el-table-column prop="statusText" label="状态" width="120" />
        <el-table-column prop="owner" label="负责人" width="120" />
        <el-table-column prop="source" label="来源系统" min-width="180" />
        <el-table-column prop="regionText" label="区域" min-width="200" />
        <el-table-column prop="createdAt" label="创建日期" width="140" />
      </el-table>

      <div class="pager-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="filteredTotal"
        />
      </div>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>items 配置能力一览</span>
        </div>
      </template>
      <el-table :data="itemFeatureRows" border>
        <el-table-column prop="name" label="配置项" width="180" />
        <el-table-column prop="required" label="必填" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.required ? 'danger' : 'info'">
              {{ row.required ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="desc" label="说明" min-width="320" />
        <el-table-column prop="example" label="示例" min-width="260" />
      </el-table>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>items 配置示例</span>
        </div>
      </template>
      <el-tabs>
        <el-tab-pane label="输入框">
          <pre class="code-block"><code>{{ itemCodeSamples.input }}</code></pre>
        </el-tab-pane>
        <el-tab-pane label="下拉单选">
          <pre class="code-block"><code>{{ itemCodeSamples.selectSingle }}</code></pre>
        </el-tab-pane>
        <el-tab-pane label="下拉多选">
          <pre class="code-block"><code>{{ itemCodeSamples.selectMultiple }}</code></pre>
        </el-tab-pane>
        <el-tab-pane label="级联多选">
          <pre class="code-block"><code>{{ itemCodeSamples.cascaderMultiple }}</code></pre>
        </el-tab-pane>
        <el-tab-pane label="日期范围">
          <pre class="code-block"><code>{{ itemCodeSamples.dateRange }}</code></pre>
        </el-tab-pane>
        <el-tab-pane label="完整 items">
          <pre class="code-block"><code>{{ itemCodeSamples.fullItems }}</code></pre>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>接入要点</span>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="组件路径">`web/src/components/lab/list-query-bar.vue`</el-descriptions-item>
        <el-descriptions-item label="核心输入">`v-model + items`，字段渲染由 `items` 配置驱动。</el-descriptions-item>
        <el-descriptions-item label="字段类型">支持 `input / select / dateRange / cascader`。</el-descriptions-item>
        <el-descriptions-item label="多选能力">`select` 与 `cascader` 均支持 `multiple: true`。</el-descriptions-item>
        <el-descriptions-item label="级联数据">`cascader` 通过 `options + cascaderProps` 支持自定义层级结构。</el-descriptions-item>
        <el-descriptions-item label="核心事件">`search / reset / export`，页面只需处理数据请求与状态。</el-descriptions-item>
        <el-descriptions-item label="插槽扩展">
          通过 `field-字段名` 覆盖字段渲染，通过 `actions-left/actions-right` 扩展操作区。
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
  import { computed, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import LabListQueryBar from '@/components/lab/list-query-bar.vue'

  defineOptions({
    name: 'LabReusableListQueryBar'
  })

  const searching = ref(false)
  const queryForm = ref({
    keyword: '',
    status: '',
    owner: '',
    source: [],
    regionPath: [],
    createdDateRange: []
  })

  const regionOptions = [
    {
      label: '华东',
      value: 'east',
      children: [
        {
          label: '浙江',
          value: 'zhejiang',
          children: [
            { label: '杭州', value: 'hangzhou' },
            { label: '宁波', value: 'ningbo' }
          ]
        },
        {
          label: '江苏',
          value: 'jiangsu',
          children: [
            { label: '南京', value: 'nanjing' },
            { label: '苏州', value: 'suzhou' }
          ]
        }
      ]
    },
    {
      label: '华南',
      value: 'south',
      children: [
        {
          label: '广东',
          value: 'guangdong',
          children: [
            { label: '广州', value: 'guangzhou' },
            { label: '深圳', value: 'shenzhen' }
          ]
        },
        {
          label: '福建',
          value: 'fujian',
          children: [
            { label: '福州', value: 'fuzhou' },
            { label: '厦门', value: 'xiamen' }
          ]
        }
      ]
    }
  ]

  const queryItems = [
    { prop: 'keyword', label: '关键字', type: 'input', placeholder: '订单号 / 客户名称 / 手机号', span: 8 },
    {
      prop: 'status',
      label: '状态',
      type: 'select',
      span: 6,
      options: [
        { label: '待处理', value: 'pending' },
        { label: '处理中', value: 'processing' },
        { label: '已完成', value: 'done' }
      ]
    },
    {
      prop: 'owner',
      label: '负责人',
      type: 'select',
      span: 6,
      options: [
        { label: '张楠', value: '张楠' },
        { label: '李聪', value: '李聪' },
        { label: '王敏', value: '王敏' }
      ]
    },
    {
      prop: 'source',
      label: '来源系统',
      type: 'select',
      multiple: true,
      span: 6,
      options: [
        { label: '运营平台', value: '运营平台' },
        { label: '渠道系统', value: '渠道系统' },
        { label: '大数据平台', value: '大数据平台' }
      ]
    },
    {
      prop: 'regionPath',
      label: '区域',
      type: 'cascader',
      multiple: true,
      span: 8,
      options: regionOptions,
      placeholder: '支持多选级联区域',
      cascaderProps: { value: 'value', label: 'label', children: 'children' }
    },
    { prop: 'createdDateRange', label: '创建日期', type: 'dateRange', span: 8 }
  ]

  const itemFeatureRows = [
    { name: 'prop', required: true, desc: '字段唯一标识，用于绑定 modelValue。', example: 'keyword' },
    { name: 'label', required: true, desc: '表单项标签文本。', example: '关键字' },
    { name: 'type', required: false, desc: '字段类型：input/select/dateRange/cascader。', example: 'select' },
    { name: 'placeholder', required: false, desc: '输入提示文案。', example: '请选择状态' },
    { name: 'defaultValue', required: false, desc: '重置时回到默认值。', example: "defaultValue: []" },
    { name: 'span', required: false, desc: '栅格宽度，等价 lg。', example: 'span: 8' },
    { name: 'sm/md/lg', required: false, desc: '分端栅格宽度，覆盖 span。', example: 'sm: 12, md: 8, lg: 6' },
    { name: 'options', required: false, desc: 'select/cascader 的选项数据。', example: '[{ label, value }]' },
    { name: 'multiple', required: false, desc: 'select/cascader 多选。', example: 'multiple: true' },
    { name: 'pickerType', required: false, desc: 'dateRange 使用的日期类型。', example: "'daterange'" },
    { name: 'startPlaceholder', required: false, desc: '日期范围开始占位文案。', example: '开始日期' },
    { name: 'endPlaceholder', required: false, desc: '日期范围结束占位文案。', example: '结束日期' },
    { name: 'cascaderProps', required: false, desc: '级联字段映射和行为配置。', example: "{ value, label, children }" },
    { name: 'checkStrictly', required: false, desc: '级联父子不关联选择。', example: 'checkStrictly: true' },
    { name: 'emitPath', required: false, desc: '级联回传完整路径还是叶子值。', example: 'emitPath: true' }
  ]

  const itemCodeSamples = {
    input: `{
  prop: 'keyword',
  label: '关键字',
  type: 'input',
  placeholder: '订单号 / 客户名称 / 手机号',
  span: 8
}`,
    selectSingle: `{
  prop: 'status',
  label: '状态',
  type: 'select',
  options: [
    { label: '待处理', value: 'pending' },
    { label: '处理中', value: 'processing' },
    { label: '已完成', value: 'done' }
  ],
  span: 6
}`,
    selectMultiple: `{
  prop: 'source',
  label: '来源系统',
  type: 'select',
  multiple: true,
  options: [
    { label: '运营平台', value: '运营平台' },
    { label: '渠道系统', value: '渠道系统' },
    { label: '大数据平台', value: '大数据平台' }
  ],
  defaultValue: [],
  span: 6
}`,
    cascaderMultiple: `{
  prop: 'regionPath',
  label: '区域',
  type: 'cascader',
  multiple: true,
  options: regionOptions,
  emitPath: true,
  checkStrictly: false,
  cascaderProps: {
    value: 'value',
    label: 'label',
    children: 'children'
  },
  defaultValue: [],
  span: 8
}`,
    dateRange: `{
  prop: 'createdDateRange',
  label: '创建日期',
  type: 'dateRange',
  pickerType: 'daterange',
  startPlaceholder: '开始日期',
  endPlaceholder: '结束日期',
  defaultValue: [],
  span: 8
}`,
    fullItems: `const queryItems = [
  { prop: 'keyword', label: '关键字', type: 'input', placeholder: '订单号 / 客户名称 / 手机号', span: 8 },
  { prop: 'status', label: '状态', type: 'select', options: statusOptions, span: 6 },
  { prop: 'owner', label: '负责人', type: 'select', options: ownerOptions, span: 6 },
  { prop: 'source', label: '来源系统', type: 'select', multiple: true, options: sourceOptions, defaultValue: [], span: 6 },
  { prop: 'regionPath', label: '区域', type: 'cascader', multiple: true, options: regionOptions, emitPath: true, cascaderProps: { value: 'value', label: 'label', children: 'children' }, defaultValue: [], span: 8 },
  { prop: 'createdDateRange', label: '创建日期', type: 'dateRange', startPlaceholder: '开始日期', endPlaceholder: '结束日期', defaultValue: [], span: 8 }
]`
  }

  const statusTexts = {
    pending: '待处理',
    processing: '处理中',
    done: '已完成'
  }

  const sourcePool = ['运营平台', '渠道系统', '大数据平台']
  const ownerPool = ['张楠', '李聪', '王敏']
  const statusPool = ['pending', 'processing', 'done']
  const regionPool = [
    ['east', 'zhejiang', 'hangzhou'],
    ['east', 'zhejiang', 'ningbo'],
    ['east', 'jiangsu', 'nanjing'],
    ['east', 'jiangsu', 'suzhou'],
    ['south', 'guangdong', 'guangzhou'],
    ['south', 'guangdong', 'shenzhen'],
    ['south', 'fujian', 'fuzhou'],
    ['south', 'fujian', 'xiamen']
  ]
  const regionLabelMap = {
    east: '华东',
    zhejiang: '浙江',
    hangzhou: '杭州',
    ningbo: '宁波',
    jiangsu: '江苏',
    nanjing: '南京',
    suzhou: '苏州',
    south: '华南',
    guangdong: '广东',
    guangzhou: '广州',
    shenzhen: '深圳',
    fujian: '福建',
    fuzhou: '福州',
    xiamen: '厦门'
  }

  const fullRows = Array.from({ length: 180 }).map((_, index) => {
    const createdAt = new Date(2026, 0, (index % 28) + 1).toISOString().slice(0, 10)
    const status = statusPool[index % statusPool.length]
    const regionPath = regionPool[index % regionPool.length]
    return {
      id: index + 1,
      orderNo: `SO2026${String(index + 1).padStart(6, '0')}`,
      customerName: `客户-${String(index + 1).padStart(4, '0')}`,
      status,
      statusText: statusTexts[status],
      owner: ownerPool[index % ownerPool.length],
      source: sourcePool[index % sourcePool.length],
      regionPath,
      regionText: regionPath.map((node) => regionLabelMap[node] || node).join(' / '),
      createdAt
    }
  })

  const pagination = ref({
    page: 1,
    pageSize: 10
  })

  const filteredRows = computed(() => {
    const { keyword, status, owner, source, regionPath, createdDateRange } = queryForm.value
    return fullRows.filter((row) => {
      const matchKeyword = !keyword || [row.orderNo, row.customerName].some((text) => text.includes(keyword))
      const matchStatus = !status || row.status === status
      const matchOwner = !owner || row.owner === owner
      const matchSource = !source?.length || source.includes(row.source)
      const matchRegion = !regionPath?.length || regionPath.some((path) => path.join('/') === row.regionPath.join('/'))
      const matchDateRange = !createdDateRange?.length || (row.createdAt >= createdDateRange[0] && row.createdAt <= createdDateRange[1])
      return matchKeyword && matchStatus && matchOwner && matchSource && matchRegion && matchDateRange
    })
  })

  const filteredTotal = computed(() => filteredRows.value.length)
  const pagedRows = computed(() => {
    const start = (pagination.value.page - 1) * pagination.value.pageSize
    const end = start + pagination.value.pageSize
    return filteredRows.value.slice(start, end)
  })

  const simulateSearchDelay = async () => {
    searching.value = true
    await new Promise((resolve) => {
      setTimeout(resolve, 280)
    })
    searching.value = false
  }

  const resetPager = () => {
    pagination.value.page = 1
  }

  const handleSearch = async () => {
    resetPager()
    await simulateSearchDelay()
    ElMessage.success('查询完成，演示数据已按条件过滤。')
  }

  const handleReset = async () => {
    resetPager()
    await simulateSearchDelay()
    ElMessage.success('已重置筛选条件。')
  }

  const handleExport = () => {
    ElMessage.info('导出事件已触发，可在业务页中对接后端导出接口。')
  }
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
    max-width: 780px;
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
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    font-weight: 600;
  }

  .total-text {
    color: #64748b;
    font-size: 13px;
    font-weight: 500;
  }

  .pager-wrap {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
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

    .pager-wrap {
      justify-content: flex-start;
    }
  }
</style>

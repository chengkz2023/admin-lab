<template>
  <div class="page-wrap">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / 新增编辑弹窗</p>
        <h2>通用新增编辑表单弹窗</h2>
        <p class="subtitle">
          一套配置同时覆盖新增与编辑，自动回显、统一校验和提交流程，减少管理模块中重复开发表单弹窗。
        </p>
      </div>
      <div class="hero-tags">
        <el-tag type="primary">新增编辑复用</el-tag>
        <el-tag type="success">配置驱动</el-tag>
        <el-tag>可扩展</el-tag>
      </div>
    </div>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>演示列表</span>
          <el-button type="primary" @click="openAddDialog">新增数据</el-button>
        </div>
      </template>
      <el-table :data="rows" border>
        <el-table-column prop="name" label="名称" min-width="180" />
        <el-table-column prop="statusText" label="状态" width="120" />
        <el-table-column prop="owner" label="负责人" width="120" />
        <el-table-column prop="regionText" label="区域" min-width="180" />
        <el-table-column prop="effectiveDate" label="生效日期" width="130" />
        <el-table-column prop="remark" label="备注" min-width="220" show-overflow-tooltip />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEditDialog(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>组件接入要点</span>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="组件路径">`web/src/components/lab/crud-form-dialog.vue`</el-descriptions-item>
        <el-descriptions-item label="核心输入">`mode + formData + items + rules`。</el-descriptions-item>
        <el-descriptions-item label="新增/编辑">`mode=add` 用于新增，`mode=edit` 自动回显 `formData`。</el-descriptions-item>
        <el-descriptions-item label="核心事件">`submit` 返回通过校验后的数据，`cancel` 关闭弹窗。</el-descriptions-item>
        <el-descriptions-item label="字段类型">支持 `input / textarea / select / switch / date / dateRange / cascader`。</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>组件级配置项</span>
        </div>
      </template>
      <el-table :data="dialogFeatureRows" border>
        <el-table-column prop="name" label="配置项" width="180" />
        <el-table-column prop="required" label="必填" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.required ? 'danger' : 'info'">
              {{ row.required ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="desc" label="说明" min-width="320" />
        <el-table-column prop="example" label="示例" min-width="280" />
      </el-table>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>items 字段配置项</span>
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
        <el-table-column prop="example" label="示例" min-width="280" />
      </el-table>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">
          <span>配置示例</span>
        </div>
      </template>
      <el-tabs>
        <el-tab-pane label="弹窗组件">
          <pre class="code-block"><code>{{ codeSamples.dialog }}</code></pre>
        </el-tab-pane>
        <el-tab-pane label="items 完整示例">
          <pre class="code-block"><code>{{ codeSamples.items }}</code></pre>
        </el-tab-pane>
        <el-tab-pane label="rules 示例">
          <pre class="code-block"><code>{{ codeSamples.rules }}</code></pre>
        </el-tab-pane>
        <el-tab-pane label="新增/编辑调用">
          <pre class="code-block"><code>{{ codeSamples.openActions }}</code></pre>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <lab-crud-form-dialog
      v-model="dialogVisible"
      v-model:form-data="formData"
      :mode="dialogMode"
      :items="formItems"
      :rules="formRules"
      :loading="saving"
      add-title="新增管理数据"
      edit-title="编辑管理数据"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import LabCrudFormDialog from '@/components/lab/crud-form-dialog.vue'

  defineOptions({
    name: 'LabReusableCrudFormDialog'
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
        }
      ]
    }
  ]

  const statusOptions = [
    { label: '启用', value: 'enabled' },
    { label: '停用', value: 'disabled' }
  ]

  const ownerOptions = [
    { label: '张楠', value: '张楠' },
    { label: '李聪', value: '李聪' },
    { label: '王敏', value: '王敏' }
  ]

  const formItems = [
    { prop: 'name', label: '名称', type: 'input', placeholder: '请输入名称', span: 12 },
    { prop: 'status', label: '状态', type: 'select', options: statusOptions, span: 12 },
    { prop: 'owner', label: '负责人', type: 'select', options: ownerOptions, span: 12 },
    { prop: 'effectiveDate', label: '生效日期', type: 'date', span: 12 },
    {
      prop: 'regionPath',
      label: '区域',
      type: 'cascader',
      options: regionOptions,
      placeholder: '请选择区域',
      span: 12,
      cascaderProps: { value: 'value', label: 'label', children: 'children' }
    },
    { prop: 'notifyUsers', label: '通知对象', type: 'select', multiple: true, options: ownerOptions, defaultValue: [], span: 12 },
    { prop: 'enabled', label: '是否启用', type: 'switch', defaultValue: true, span: 12 },
    { prop: 'remark', label: '备注', type: 'textarea', rows: 3, span: 24 }
  ]

  const formRules = {
    name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
    status: [{ required: true, message: '请选择状态', trigger: 'change' }],
    owner: [{ required: true, message: '请选择负责人', trigger: 'change' }],
    effectiveDate: [{ required: true, message: '请选择生效日期', trigger: 'change' }],
    regionPath: [{ required: true, message: '请选择区域', trigger: 'change' }]
  }

  const dialogFeatureRows = [
    { name: 'modelValue', required: true, desc: '控制弹窗显示与隐藏。', example: 'v-model="dialogVisible"' },
    { name: 'mode', required: true, desc: '弹窗模式：add/edit。', example: 'mode="add"' },
    { name: 'formData', required: true, desc: '表单数据，edit 模式用于回显。', example: 'v-model:form-data="formData"' },
    { name: 'items', required: true, desc: '字段配置数组，决定表单结构。', example: ':items="formItems"' },
    { name: 'rules', required: false, desc: 'Element Plus 表单校验规则。', example: ':rules="formRules"' },
    { name: 'loading', required: false, desc: '提交按钮 loading 状态。', example: ':loading="saving"' },
    { name: 'addTitle/editTitle', required: false, desc: '新增/编辑标题文案。', example: 'add-title="新增数据"' },
    { name: 'confirmText/cancelText', required: false, desc: '底部按钮文案覆盖。', example: 'confirm-text="提交"' },
    { name: 'width', required: false, desc: '弹窗宽度。', example: 'width="760px"' },
    { name: 'labelWidth', required: false, desc: '表单标签宽度。', example: 'label-width="96px"' },
    { name: 'destroyOnClose', required: false, desc: '关闭销毁内容，默认 true。', example: ':destroy-on-close="true"' },
    { name: '@submit', required: true, desc: '提交事件，返回校验后的 payload。', example: '@submit="handleSubmit"' },
    { name: '@cancel / @closed', required: false, desc: '取消与关闭回调。', example: '@cancel="handleCancel"' }
  ]

  const itemFeatureRows = [
    { name: 'prop', required: true, desc: '字段唯一标识，绑定 formData。', example: 'name' },
    { name: 'label', required: true, desc: '表单标签。', example: '名称' },
    { name: 'type', required: false, desc: '字段类型：input/textarea/select/switch/date/dateRange/cascader。', example: 'select' },
    { name: 'placeholder', required: false, desc: '输入提示文案。', example: '请输入名称' },
    { name: 'defaultValue', required: false, desc: '新增和重置时默认值。', example: 'defaultValue: []' },
    { name: 'span/sm/md/lg', required: false, desc: '栅格布局。', example: 'span: 12' },
    { name: 'options', required: false, desc: 'select/cascader 的选项数据。', example: '[{ label, value }]' },
    { name: 'multiple', required: false, desc: 'select/cascader 多选。', example: 'multiple: true' },
    { name: 'rows', required: false, desc: 'textarea 行数。', example: 'rows: 3' },
    { name: 'activeText/inactiveText', required: false, desc: 'switch 文案。', example: 'activeText: 启用' },
    { name: 'startPlaceholder/endPlaceholder', required: false, desc: 'dateRange 文案。', example: '开始日期 / 结束日期' },
    { name: 'cascaderProps', required: false, desc: '级联字段映射与行为。', example: "{ value, label, children }" },
    { name: 'checkStrictly', required: false, desc: '级联父子节点不关联。', example: 'checkStrictly: true' },
    { name: 'emitPath', required: false, desc: '级联回传完整路径或叶子值。', example: 'emitPath: true' }
  ]

  const codeSamples = {
    dialog: `<lab-crud-form-dialog
  v-model="dialogVisible"
  v-model:form-data="formData"
  :mode="dialogMode"
  :items="formItems"
  :rules="formRules"
  :loading="saving"
  add-title="新增管理数据"
  edit-title="编辑管理数据"
  @submit="handleSubmit"
/>`,
    items: `const formItems = [
  { prop: 'name', label: '名称', type: 'input', placeholder: '请输入名称', span: 12 },
  { prop: 'status', label: '状态', type: 'select', options: statusOptions, span: 12 },
  { prop: 'owner', label: '负责人', type: 'select', options: ownerOptions, span: 12 },
  { prop: 'effectiveDate', label: '生效日期', type: 'date', span: 12 },
  { prop: 'regionPath', label: '区域', type: 'cascader', options: regionOptions, cascaderProps: { value: 'value', label: 'label', children: 'children' }, span: 12 },
  { prop: 'notifyUsers', label: '通知对象', type: 'select', multiple: true, options: ownerOptions, defaultValue: [], span: 12 },
  { prop: 'enabled', label: '是否启用', type: 'switch', defaultValue: true, span: 12 },
  { prop: 'remark', label: '备注', type: 'textarea', rows: 3, span: 24 }
]`,
    rules: `const formRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  owner: [{ required: true, message: '请选择负责人', trigger: 'change' }],
  effectiveDate: [{ required: true, message: '请选择生效日期', trigger: 'change' }],
  regionPath: [{ required: true, message: '请选择区域', trigger: 'change' }]
}`,
    openActions: `const openAddDialog = () => {
  dialogMode.value = 'add'
  formData.value = { name: '', status: 'enabled', owner: '' }
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  dialogMode.value = 'edit'
  formData.value = { ...row }
  dialogVisible.value = true
}`
  }

  const rows = ref([
    {
      id: 1,
      name: '客户分层策略',
      status: 'enabled',
      statusText: '启用',
      owner: '张楠',
      regionPath: ['east', 'zhejiang', 'hangzhou'],
      regionText: '华东 / 浙江 / 杭州',
      effectiveDate: '2026-03-01',
      notifyUsers: ['张楠', '李聪'],
      enabled: true,
      remark: '用于重点客户分层运营。'
    },
    {
      id: 2,
      name: '库存预警策略',
      status: 'disabled',
      statusText: '停用',
      owner: '王敏',
      regionPath: ['south', 'guangdong', 'shenzhen'],
      regionText: '华南 / 广东 / 深圳',
      effectiveDate: '2026-03-05',
      notifyUsers: ['王敏'],
      enabled: false,
      remark: '当前暂未启用。'
    }
  ])

  const dialogVisible = ref(false)
  const dialogMode = ref('add')
  const saving = ref(false)
  const editingId = ref(0)
  const formData = ref({})

  const toRegionText = (regionPath = []) => {
    const labels = {
      east: '华东',
      zhejiang: '浙江',
      hangzhou: '杭州',
      ningbo: '宁波',
      south: '华南',
      guangdong: '广东',
      guangzhou: '广州',
      shenzhen: '深圳'
    }
    return regionPath.map((item) => labels[item] || item).join(' / ')
  }

  const openAddDialog = () => {
    dialogMode.value = 'add'
    editingId.value = 0
    formData.value = {
      name: '',
      status: 'enabled',
      owner: '',
      effectiveDate: '',
      regionPath: [],
      notifyUsers: [],
      enabled: true,
      remark: ''
    }
    dialogVisible.value = true
  }

  const openEditDialog = (row) => {
    dialogMode.value = 'edit'
    editingId.value = row.id
    formData.value = {
      name: row.name,
      status: row.status,
      owner: row.owner,
      effectiveDate: row.effectiveDate,
      regionPath: row.regionPath || [],
      notifyUsers: row.notifyUsers || [],
      enabled: row.enabled,
      remark: row.remark
    }
    dialogVisible.value = true
  }

  const handleSubmit = async (payload) => {
    saving.value = true
    await new Promise((resolve) => {
      setTimeout(resolve, 280)
    })

    const merged = {
      ...payload,
      statusText: payload.status === 'enabled' ? '启用' : '停用',
      regionText: toRegionText(payload.regionPath)
    }

    if (dialogMode.value === 'add') {
      rows.value.unshift({
        id: Date.now(),
        ...merged
      })
      ElMessage.success('新增成功')
    } else {
      rows.value = rows.value.map((item) => (item.id === editingId.value ? { ...item, ...merged } : item))
      ElMessage.success('编辑成功')
    }

    dialogVisible.value = false
    saving.value = false
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
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    font-weight: 600;
  }

  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
    }
  }
</style>

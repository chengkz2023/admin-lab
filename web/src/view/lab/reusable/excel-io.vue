<template>
  <div class="lab-page">
    <div class="lab-hero reusable-hero">
      <div>
        <p class="lab-eyebrow">复用组件 / Excel 实验面板</p>
        <h2>Excel 导入导出实验台</h2>
        <p class="lab-subtitle">
          这个模块专门模拟公司里常见的 Excel 导入导出需求。现在模板下载、模板选择、导入解析和结果展示都已经按模板联动。
        </p>
      </div>
      <div class="hero-actions">
        <el-button type="primary" @click="handleDownloadDefaultTemplate" :loading="downloadingTemplate">
          下载默认模板
        </el-button>
        <el-button @click="openTemplateSelector" :loading="templateOptionsLoading">
          选择模板下载
        </el-button>
        <el-button @click="handleExportSample" :loading="downloadingSample">
          导出默认示例
        </el-button>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="8">
        <el-card class="h-full" shadow="hover">
          <template #header>
            <div class="panel-title">适用场景</div>
          </template>
          <el-timeline>
            <el-timeline-item v-for="item in scenarios" :key="item" type="primary">
              {{ item }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="16">
        <el-card class="h-full" shadow="hover">
          <template #header>
            <div class="panel-header">
              <span class="panel-title">导入实验</span>
              <el-select
                v-model="selectedImportTemplateKey"
                class="template-select"
                placeholder="选择导入模板"
                :loading="templateOptionsLoading"
              >
                <el-option
                  v-for="template in templateOptions"
                  :key="template.key"
                  :label="template.name"
                  :value="template.key"
                />
              </el-select>
            </div>
          </template>

          <div v-if="currentImportTemplate" class="current-template-card">
            <div class="current-template-top">
              <div>
                <div class="current-template-title">
                  {{ currentImportTemplate.name }}
                  <el-tag v-if="currentImportTemplate.isDefault" size="small" type="primary">默认</el-tag>
                </div>
                <div class="current-template-scene">{{ currentImportTemplate.scene }}</div>
              </div>
              <el-button link type="primary" @click="handleDownloadTemplate(currentImportTemplate)">
                下载当前模板
              </el-button>
            </div>
            <div class="current-template-desc">{{ currentImportTemplate.description }}</div>
            <div class="template-columns">
              <el-tag v-for="column in currentImportTemplate.columns" :key="column" size="small" effect="plain">
                {{ column }}
              </el-tag>
            </div>
          </div>

          <div class="upload-panel">
            <el-upload
              ref="uploadRef"
              :auto-upload="false"
              :show-file-list="false"
              accept=".xlsx"
              @change="handleImportChange"
            >
              <el-button type="success" :loading="importing">
                上传 Excel 并解析
              </el-button>
            </el-upload>
            <span class="upload-tip">上传时会按当前选中的模板进行字段映射和校验。</span>
          </div>

          <el-alert
            class="mb-4"
            type="info"
            :closable="false"
            :title="importNotice"
          />

          <div v-if="summary" class="summary-grid">
            <div class="summary-card">
              <span class="summary-label">总行数</span>
              <strong>{{ summary.totalRows }}</strong>
            </div>
            <div class="summary-card success">
              <span class="summary-label">成功行</span>
              <strong>{{ summary.successRows }}</strong>
            </div>
            <div class="summary-card danger">
              <span class="summary-label">失败行</span>
              <strong>{{ summary.failedRows }}</strong>
            </div>
          </div>

          <el-empty v-if="!rows.length" description="上传 Excel 后，这里会按当前模板显示解析结果和校验信息。" />

          <el-table v-else :data="rows" class="mt-4" border>
            <el-table-column prop="rowNumber" label="行号" width="80" fixed="left" />
            <el-table-column
              v-for="column in resultColumns"
              :key="column"
              :label="column"
              min-width="140"
            >
              <template #default="{ row }">
                {{ row.values?.[column] || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="校验结果" min-width="280" fixed="right">
              <template #default="{ row }">
                <div v-if="row.errorFields?.length" class="error-tags">
                  <el-tag v-for="item in row.errorFields" :key="item" type="danger" size="small">
                    {{ item }}
                  </el-tag>
                </div>
                <el-tag v-else type="success" size="small">通过</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">可选模板</div>
      </template>
      <el-empty v-if="!templateOptions.length && !templateOptionsLoading" description="暂无模板数据" />
      <div v-else class="template-grid">
        <div v-for="template in templateOptions" :key="template.key" class="template-card">
          <div class="template-header">
            <div>
              <div class="template-title">
                {{ template.name }}
                <el-tag v-if="template.isDefault" size="small" type="primary">默认</el-tag>
              </div>
              <div class="template-scene">{{ template.scene }}</div>
            </div>
            <div class="template-actions">
              <el-button link type="primary" @click="selectedImportTemplateKey = template.key">
                设为当前
              </el-button>
              <el-button link type="primary" @click="handleDownloadTemplate(template)">
                下载
              </el-button>
            </div>
          </div>
          <div class="template-description">{{ template.description }}</div>
          <div class="template-columns">
            <el-tag v-for="column in template.columns" :key="column" size="small" effect="plain">
              {{ column }}
            </el-tag>
          </div>
        </div>
      </div>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title">搬运到内网时建议替换的部分</div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="模板策略">可以按不同业务维护不同模板，例如客户、订单、商品、库存、财务单据。</el-descriptions-item>
        <el-descriptions-item label="字段定义">把示例字段替换成你业务里的真实列头和字段映射。</el-descriptions-item>
        <el-descriptions-item label="校验规则">把模板内的演示校验改成真实业务规则。</el-descriptions-item>
        <el-descriptions-item label="导入落库">现在只做解析预演，进内网后可接入数据库写入、重复校验和事务处理。</el-descriptions-item>
        <el-descriptions-item label="导出逻辑">把示例导出替换成真实查询条件、权限过滤和业务数据导出。</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-dialog v-model="templateSelectorVisible" title="选择模板下载" width="720px">
      <div class="selector-list">
        <label v-for="template in templateOptions" :key="template.key" class="selector-item">
          <input v-model="selectedTemplateKey" type="radio" :value="template.key">
          <div class="selector-content">
            <div class="selector-title">
              {{ template.name }}
              <el-tag v-if="template.isDefault" size="small" type="primary">默认</el-tag>
            </div>
            <div class="selector-desc">{{ template.description }}</div>
            <div class="selector-scene">{{ template.scene }}</div>
          </div>
        </label>
      </div>
      <template #footer>
        <el-button @click="templateSelectorVisible = false">取消</el-button>
        <el-button type="primary" :loading="downloadingTemplate" @click="handleConfirmTemplateDownload">
          下载选中模板
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import {
    downloadExcelTemplate,
    exportExcelSample,
    getExcelTemplateOptions,
    importExcelFile
  } from '@/api/excelIO'

  defineOptions({
    name: 'LabReusableExcelIO'
  })

  const scenarios = [
    '批量导入客户、用户、商品、字典、组织结构等基础数据。',
    '导出列表查询结果，交给运营、业务或实施同学二次处理。',
    '在内网正式接库前，先把 Excel 模板、字段映射和校验规则跑通。',
    '沉淀一套可复用的 Excel 导入导出处理骨架，减少重复开发。'
  ]

  const downloadingTemplate = ref(false)
  const downloadingSample = ref(false)
  const importing = ref(false)
  const templateOptionsLoading = ref(false)
  const templateOptions = ref([])
  const templateSelectorVisible = ref(false)
  const selectedTemplateKey = ref('')
  const selectedImportTemplateKey = ref('')
  const rows = ref([])
  const summary = ref(null)
  const resultColumns = ref([])
  const uploadRef = ref(null)

  const currentImportTemplate = computed(() => {
    return templateOptions.value.find((item) => item.key === selectedImportTemplateKey.value) || null
  })

  const importNotice = computed(() => {
    if (!currentImportTemplate.value) {
      return '请选择一个导入模板后再上传 Excel。'
    }
    return `当前导入模板：${currentImportTemplate.value.name}。字段为：${currentImportTemplate.value.columns.join('、')}。`
  })

  const saveBlob = (response, fallbackName) => {
    const blob = new Blob([response.data], { type: response.headers['content-type'] || 'application/octet-stream' })
    const disposition = response.headers['content-disposition'] || ''
    const match = disposition.match(/filename=([^;]+)/i)
    const fileName = match ? decodeURIComponent(match[1].replace(/"/g, '')) : fallbackName
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  }

  const loadTemplateOptions = async () => {
    templateOptionsLoading.value = true
    try {
      const result = await getExcelTemplateOptions()
      templateOptions.value = result.data || []
      const defaultTemplate = templateOptions.value.find((item) => item.isDefault)
      selectedTemplateKey.value = defaultTemplate?.key || templateOptions.value[0]?.key || ''
      selectedImportTemplateKey.value = selectedImportTemplateKey.value || defaultTemplate?.key || templateOptions.value[0]?.key || ''
    } finally {
      templateOptionsLoading.value = false
    }
  }

  const downloadTemplateByKey = async (templateKey, fallbackName) => {
    downloadingTemplate.value = true
    try {
      const response = await downloadExcelTemplate(templateKey)
      saveBlob(response, fallbackName)
      ElMessage.success('模板下载成功')
    } finally {
      downloadingTemplate.value = false
    }
  }

  const handleDownloadDefaultTemplate = async () => {
    const defaultTemplate = templateOptions.value.find((item) => item.isDefault)
    await downloadTemplateByKey(defaultTemplate?.key, defaultTemplate?.fileName || 'admin-lab-excel-template.xlsx')
  }

  const openTemplateSelector = async () => {
    if (!templateOptions.value.length) {
      await loadTemplateOptions()
    }
    templateSelectorVisible.value = true
  }

  const handleDownloadTemplate = async (template) => {
    await downloadTemplateByKey(template.key, template.fileName)
  }

  const handleConfirmTemplateDownload = async () => {
    const template = templateOptions.value.find((item) => item.key === selectedTemplateKey.value)
    if (!template) {
      ElMessage.warning('请先选择一个模板')
      return
    }
    await downloadTemplateByKey(template.key, template.fileName)
    templateSelectorVisible.value = false
  }

  const handleExportSample = async () => {
    downloadingSample.value = true
    try {
      const response = await exportExcelSample()
      saveBlob(response, 'admin-lab-excel-sample.xlsx')
      ElMessage.success('示例导出成功')
    } finally {
      downloadingSample.value = false
    }
  }

  const handleImportChange = async (file) => {
    if (!file?.raw) {
      return
    }
    if (!selectedImportTemplateKey.value) {
      ElMessage.warning('请先选择导入模板')
      uploadRef.value?.clearFiles()
      return
    }

    importing.value = true
    try {
      const formData = new FormData()
      formData.append('file', file.raw)
      formData.append('templateKey', selectedImportTemplateKey.value)
      const result = await importExcelFile(formData)
      rows.value = result.data.rows || []
      resultColumns.value = result.data.columns || []
      summary.value = {
        totalRows: result.data.totalRows || 0,
        successRows: result.data.successRows || 0,
        failedRows: result.data.failedRows || 0
      }
      ElMessage.success(`Excel 解析完成，当前模板：${result.data.templateName || currentImportTemplate.value?.name || ''}`)
    } finally {
      importing.value = false
      uploadRef.value?.clearFiles()
    }
  }

  onMounted(() => {
    loadTemplateOptions()
  })
</script>

<style scoped>
  .lab-page {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .lab-hero {
    display: flex;
    justify-content: space-between;
    gap: 16px;
    padding: 24px;
    border-radius: 16px;
    border: 1px solid #d5efe1;
    background: linear-gradient(135deg, #eefbf4 0%, #f3f8ff 100%);
  }

  .lab-eyebrow {
    margin: 0 0 8px;
    color: #059669;
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.08em;
  }

  .lab-hero h2 {
    margin: 0 0 8px;
    font-size: 24px;
    line-height: 1.4;
    color: #1f2937;
  }

  .lab-subtitle {
    margin: 0;
    max-width: 760px;
    color: #4b5563;
    line-height: 1.75;
  }

  .hero-actions {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    gap: 12px;
  }

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
  }

  .panel-title {
    font-weight: 600;
    color: #111827;
  }

  .template-select {
    width: 240px;
  }

  .current-template-card {
    margin-bottom: 16px;
    padding: 16px;
    border-radius: 14px;
    border: 1px solid #d9f1e3;
    background: #f7fcf9;
  }

  .current-template-top {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 10px;
  }

  .current-template-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
    font-weight: 600;
    color: #111827;
  }

  .current-template-scene {
    color: #059669;
    font-size: 13px;
  }

  .current-template-desc {
    margin-bottom: 12px;
    color: #4b5563;
    line-height: 1.7;
    font-size: 14px;
  }

  .upload-panel {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
  }

  .upload-tip {
    color: #6b7280;
    font-size: 13px;
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
    margin-bottom: 16px;
  }

  .summary-card {
    padding: 16px;
    border-radius: 12px;
    background: #f8fafc;
    border: 1px solid #e5e7eb;
  }

  .summary-card.success {
    background: #ecfdf5;
    border-color: #a7f3d0;
  }

  .summary-card.danger {
    background: #fef2f2;
    border-color: #fecaca;
  }

  .summary-label {
    display: block;
    margin-bottom: 6px;
    color: #6b7280;
    font-size: 13px;
  }

  .summary-card strong {
    font-size: 24px;
    color: #111827;
  }

  .error-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .template-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 16px;
  }

  .template-card {
    padding: 16px;
    border: 1px solid #e5e7eb;
    border-radius: 14px;
    background: #f8fafc;
  }

  .template-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 10px;
  }

  .template-actions {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
  }

  .template-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
    font-weight: 600;
    color: #111827;
  }

  .template-scene {
    color: #059669;
    font-size: 13px;
  }

  .template-description {
    margin-bottom: 12px;
    color: #4b5563;
    line-height: 1.7;
    font-size: 14px;
  }

  .template-columns {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .selector-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .selector-item {
    display: flex;
    gap: 12px;
    padding: 14px;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    cursor: pointer;
  }

  .selector-content {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .selector-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    color: #111827;
  }

  .selector-desc,
  .selector-scene {
    color: #6b7280;
    line-height: 1.7;
    font-size: 14px;
  }

  @media (max-width: 768px) {
    .lab-hero,
    .panel-header,
    .current-template-top {
      flex-direction: column;
    }

    .template-select {
      width: 100%;
    }

    .summary-grid {
      grid-template-columns: 1fr;
    }

    .template-actions {
      align-items: flex-start;
    }
  }
</style>

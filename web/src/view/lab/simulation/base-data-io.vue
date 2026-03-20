<template>
  <div class="sim-page">
    <div class="hero">
      <div>
        <p class="eyebrow">需求仿真 / 基础数据导入导出</p>
        <h2>基础数据综合导入导出仿真</h2>
        <p class="subtitle">
          支持机房、服务用户、其他用户三类模板。当前版本已实现导入校验与结果展示，但仍不落库。
        </p>
      </div>
      <el-tag type="warning" effect="dark">Simulation</el-tag>
    </div>

    <el-card shadow="never">
      <template #header>
        <div class="title">模板类型</div>
      </template>
      <el-radio-group v-model="selectedTemplateKey">
        <el-radio-button
          v-for="template in templateOptions"
          :key="template.key"
          :label="template.key"
        >
          {{ template.name }}
        </el-radio-button>
      </el-radio-group>
      <div v-if="currentTemplate" class="template-info">
        <el-alert
          :title="currentTemplate.description"
          :description="`适用场景：${currentTemplate.scene}`"
          type="info"
          :closable="false"
          show-icon
        />
        <div class="field-tags">
          <el-tag v-for="field in currentTemplate.columns" :key="field" size="small" effect="plain">
            {{ field }}
          </el-tag>
        </div>
      </div>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="title">模板下载 / 数据导出</div>
      </template>
      <div class="action-row">
        <el-button type="primary" :loading="downloading" @click="onDownloadTemplate">
          下载模板
        </el-button>
        <el-button type="success" :loading="exporting" @click="onExportData">
          导出数据
        </el-button>
      </div>
      <p class="hint">模板会标记必填字段（*），导出内容为仿真数据。</p>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="title">导入校验（不落库）</div>
      </template>
      <div class="action-row">
        <el-upload
          ref="uploadRef"
          :auto-upload="false"
          :show-file-list="false"
          accept=".xlsx"
          @change="onImportChange"
        >
          <el-button type="warning" :loading="importing">上传并导入</el-button>
        </el-upload>
      </div>

      <el-alert
        v-if="importMessage"
        :title="importMessage"
        type="info"
        :closable="false"
        show-icon
        class="mb-3"
      />

      <div v-if="summary" class="summary-grid">
        <div class="summary-card">
          <span class="summary-label">总行数</span>
          <strong>{{ summary.totalRows }}</strong>
        </div>
        <div class="summary-card success">
          <span class="summary-label">通过行</span>
          <strong>{{ summary.successRows }}</strong>
        </div>
        <div class="summary-card danger">
          <span class="summary-label">失败行</span>
          <strong>{{ summary.failedRows }}</strong>
        </div>
      </div>

      <el-empty v-if="!rows.length" description="上传 Excel 后展示逐行校验结果" />

      <el-table v-else :data="rows" border>
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
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import {
    downloadBaseDataTemplate,
    exportBaseData,
    getBaseDataTemplateOptions,
    importBaseData
  } from '@/api/baseDataSimulation'

  defineOptions({ name: 'LabSimulationBaseDataIO' })

  const templateOptions = ref([])
  const selectedTemplateKey = ref('')
  const downloading = ref(false)
  const exporting = ref(false)
  const importing = ref(false)
  const uploadRef = ref(null)

  const rows = ref([])
  const resultColumns = ref([])
  const summary = ref(null)
  const importMessage = ref('')

  const currentTemplate = computed(() => {
    return templateOptions.value.find((item) => item.key === selectedTemplateKey.value) || null
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

  const loadTemplates = async () => {
    const res = await getBaseDataTemplateOptions()
    templateOptions.value = res.data || []
    const defaultItem = templateOptions.value.find((item) => item.isDefault)
    selectedTemplateKey.value = defaultItem?.key || templateOptions.value[0]?.key || ''
  }

  const onDownloadTemplate = async () => {
    if (!selectedTemplateKey.value) return
    downloading.value = true
    try {
      const res = await downloadBaseDataTemplate(selectedTemplateKey.value)
      saveBlob(res, 'base-data-template.xlsx')
      ElMessage.success('模板下载成功')
    } finally {
      downloading.value = false
    }
  }

  const onExportData = async () => {
    if (!selectedTemplateKey.value) return
    exporting.value = true
    try {
      const res = await exportBaseData(selectedTemplateKey.value)
      saveBlob(res, 'base-data-export.xlsx')
      ElMessage.success('导出成功')
    } finally {
      exporting.value = false
    }
  }

  const onImportChange = async (file) => {
    if (!file?.raw || !selectedTemplateKey.value) return
    importing.value = true
    try {
      const formData = new FormData()
      formData.append('file', file.raw)
      formData.append('templateKey', selectedTemplateKey.value)
      const res = await importBaseData(formData)
      const data = res.data || {}
      rows.value = data.rows || []
      resultColumns.value = data.columns || []
      summary.value = {
        totalRows: data.totalRows || 0,
        successRows: data.successRows || 0,
        failedRows: data.failedRows || 0
      }
      importMessage.value = res.msg || '导入校验完成'
      ElMessage.success('导入校验完成')
    } finally {
      importing.value = false
      uploadRef.value?.clearFiles()
    }
  }

  onMounted(async () => {
    await loadTemplates()
  })
</script>

<style scoped>
  .sim-page {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .hero {
    display: flex;
    justify-content: space-between;
    gap: 16px;
    padding: 20px 24px;
    border: 1px solid #dbeafe;
    border-radius: 14px;
    background: linear-gradient(135deg, #f5faff 0%, #f8fafc 100%);
  }

  .eyebrow {
    margin: 0 0 8px;
    color: #2563eb;
    font-size: 13px;
    font-weight: 700;
  }

  .hero h2 {
    margin: 0 0 8px;
    color: #0f172a;
    font-size: 24px;
  }

  .subtitle {
    margin: 0;
    color: #475569;
    line-height: 1.7;
  }

  .title {
    font-weight: 600;
    color: #111827;
  }

  .template-info {
    margin-top: 12px;
  }

  .field-tags {
    margin-top: 12px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .action-row {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 12px;
  }

  .hint {
    margin: 0;
    color: #6b7280;
    font-size: 13px;
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
    margin: 12px 0;
  }

  .summary-card {
    padding: 14px;
    border-radius: 12px;
    border: 1px solid #e5e7eb;
    background: #f8fafc;
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
    font-size: 22px;
    color: #111827;
  }

  .error-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .mb-3 {
    margin-bottom: 12px;
  }

  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
    }

    .summary-grid {
      grid-template-columns: 1fr;
    }
  }
</style>


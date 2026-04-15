<template>
  <div class="lab-page" v-loading="loadingProfile || running">
    <div class="hero">
      <div>
        <p class="eyebrow">{{ profile.classification || '复用组件' }}</p>
        <h2>{{ profile.title || '目录文件处理流水线' }}</h2>
        <p class="subtitle">
          {{
            profile.summary ||
            '配置输入输出目录后，可直接执行一次目录文件处理流水线。'
          }}
        </p>
      </div>
      <div class="hero-tags">
        <el-tag type="success">单次执行</el-tag>
        <el-tag type="primary">可复用</el-tag>
        <el-tag>后端能力</el-tag>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="11">
        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">能力亮点</div>
          </template>
          <div class="bullet-list">
            <div v-for="item in profile.highlights || []" :key="item" class="bullet-item">
              {{ item }}
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">快速上手</div>
          </template>
          <div class="step-list">
            <div v-for="(item, index) in profile.quickSteps || []" :key="item" class="step-item">
              <span>{{ index + 1 }}</span>
              <div>{{ item }}</div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">可用处理器</div>
          </template>
          <div class="processor-list">
            <div v-for="item in profile.processors || []" :key="item.key" class="processor-item">
              <div class="processor-head">
                <strong>{{ item.name }}</strong>
                <el-tag size="small">{{ item.key }}</el-tag>
              </div>
              <div class="muted">{{ item.description }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="13">
        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">执行一次</div>
          </template>
          <el-form :model="form" label-position="top">
            <el-row :gutter="12">
              <el-col :xs="24" :md="12">
                <el-form-item label="输入目录">
                  <el-input v-model="form.inputDir" placeholder="D:/runtime/in" clearable />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="输出目录">
                  <el-input v-model="form.outputDir" placeholder="D:/runtime/out" clearable />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="12">
              <el-col :xs="24" :md="12">
                <el-form-item label="错误目录（可选）">
                  <el-input v-model="form.errorDir" placeholder="默认：output/_errors" clearable />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="归档目录（可选）">
                  <el-input v-model="form.archiveDir" placeholder="默认：input/_archive" clearable />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="12">
              <el-col :xs="24" :md="8">
                <el-form-item label="文件匹配模式">
                  <el-input v-model="form.filePattern" placeholder="*.txt" clearable />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="8">
                <el-form-item label="单次最大文件数">
                  <el-input-number v-model="form.maxFiles" :min="1" :max="1000" style="width: 100%" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="8">
                <el-form-item label="稳定等待（ms）">
                  <el-input-number v-model="form.stableWaitMs" :min="100" :max="10000" style="width: 100%" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="12">
              <el-col :xs="24" :md="12">
                <el-form-item label="处理器">
                  <el-select v-model="form.processor" style="width: 100%">
                    <el-option
                      v-for="item in profile.processors || []"
                      :key="item.key"
                      :label="item.name"
                      :value="item.key"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :md="12">
                <el-form-item label="输出后缀">
                  <el-input v-model="form.outputSuffix" placeholder="_processed" clearable />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item>
              <el-button type="primary" :loading="running" @click="runOnce">执行一次</el-button>
              <el-button :disabled="running" @click="resetToDefault">恢复默认</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card shadow="hover" v-if="result">
          <template #header>
            <div class="panel-header">最近一次执行结果</div>
          </template>
          <el-row :gutter="12" class="metrics">
            <el-col :xs="12" :md="6"><div class="metric">扫描数：{{ result.scanned }}</div></el-col>
            <el-col :xs="12" :md="6"><div class="metric">成功数：{{ result.processed }}</div></el-col>
            <el-col :xs="12" :md="6"><div class="metric">跳过数：{{ result.skipped }}</div></el-col>
            <el-col :xs="12" :md="6"><div class="metric">失败数：{{ result.failed }}</div></el-col>
          </el-row>

          <div class="result-section" v-if="result.outputFiles?.length">
            <h4>输出文件</h4>
            <div class="path-list">
              <code v-for="item in result.outputFiles" :key="item">{{ item }}</code>
            </div>
          </div>

          <div class="result-section" v-if="result.archivedFiles?.length">
            <h4>归档文件</h4>
            <div class="path-list">
              <code v-for="item in result.archivedFiles" :key="item">{{ item }}</code>
            </div>
          </div>

          <div class="result-section" v-if="result.failureItems?.length">
            <h4>失败明细</h4>
            <el-table :data="result.failureItems" size="small">
              <el-table-column prop="file" label="文件" min-width="220" show-overflow-tooltip />
              <el-table-column prop="stage" label="阶段" width="130" />
              <el-table-column prop="reason" label="原因" min-width="220" show-overflow-tooltip />
              <el-table-column prop="errorFile" label="错误目录文件" min-width="220" show-overflow-tooltip />
            </el-table>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { onMounted, reactive, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getDirFilePipelineProfile, runDirFilePipelineOnce } from '@/api/dirFilePipeline'

  defineOptions({
    name: 'LabReusableDirFilePipeline'
  })

  const loadingProfile = ref(false)
  const running = ref(false)
  const result = ref(null)
  const profile = ref({
    processors: [],
    highlights: [],
    quickSteps: [],
    defaultConfig: {}
  })

  const form = reactive({
    inputDir: '',
    outputDir: '',
    errorDir: '',
    archiveDir: '',
    filePattern: '*.txt',
    maxFiles: 20,
    stableWaitMs: 1000,
    processor: 'copy',
    outputSuffix: '_processed'
  })

  const applyDefaultConfig = (config = {}) => {
    form.inputDir = config.inputDir || ''
    form.outputDir = config.outputDir || ''
    form.errorDir = config.errorDir || ''
    form.archiveDir = config.archiveDir || ''
    form.filePattern = config.filePattern || '*.txt'
    form.maxFiles = config.maxFiles || 20
    form.stableWaitMs = config.stableWaitMs || 1000
    form.processor = config.processor || 'copy'
    form.outputSuffix = config.outputSuffix || '_processed'
  }

  const loadProfile = async () => {
    loadingProfile.value = true
    try {
      const res = await getDirFilePipelineProfile()
      profile.value = res.data || {}
      applyDefaultConfig(profile.value.defaultConfig)
    } catch (error) {
      ElMessage.error(error?.message || '加载介绍信息失败')
    } finally {
      loadingProfile.value = false
    }
  }

  const resetToDefault = () => {
    applyDefaultConfig(profile.value.defaultConfig || {})
  }

  const runOnce = async () => {
    running.value = true
    try {
      const payload = {
        inputDir: form.inputDir,
        outputDir: form.outputDir,
        errorDir: form.errorDir,
        archiveDir: form.archiveDir,
        filePattern: form.filePattern,
        maxFiles: form.maxFiles,
        stableWaitMs: form.stableWaitMs,
        processor: form.processor,
        outputSuffix: form.outputSuffix
      }
      const res = await runDirFilePipelineOnce(payload)
      result.value = res.data || null
      ElMessage.success('执行完成')
    } catch (error) {
      ElMessage.error(error?.message || '执行失败')
    } finally {
      running.value = false
    }
  }

  onMounted(() => {
    loadProfile()
  })
</script>

<style scoped>
  .lab-page {
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
    color: #2563eb;
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.08em;
  }

  .hero h2 {
    margin: 0 0 8px;
    color: #0f172a;
    font-size: 28px;
  }

  .subtitle {
    margin: 0;
    color: #475569;
    line-height: 1.7;
    max-width: 760px;
  }

  .hero-tags {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    gap: 8px;
  }

  .panel-header {
    font-weight: 600;
  }

  .bullet-list,
  .step-list,
  .processor-list,
  .path-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .bullet-item,
  .processor-item {
    padding: 12px 14px;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
    background: #f8fafc;
  }

  .step-item {
    display: grid;
    grid-template-columns: 28px 1fr;
    gap: 10px;
    align-items: start;
    padding: 10px 12px;
    border-radius: 10px;
    border: 1px dashed #bfdbfe;
    background: #eff6ff;
  }

  .step-item > span {
    width: 24px;
    height: 24px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 999px;
    background: #2563eb;
    color: #fff;
    font-size: 12px;
    font-weight: 600;
  }

  .processor-head {
    display: flex;
    justify-content: space-between;
    gap: 8px;
    margin-bottom: 6px;
  }

  .muted {
    color: #64748b;
    line-height: 1.7;
  }

  .metrics {
    margin-bottom: 12px;
  }

  .metric {
    padding: 10px;
    border-radius: 10px;
    text-align: center;
    background: #f1f5f9;
    border: 1px solid #e2e8f0;
    font-weight: 600;
  }

  .result-section {
    margin-top: 12px;
  }

  .result-section h4 {
    margin: 0 0 8px;
    color: #334155;
  }

  .path-list code {
    display: block;
    padding: 8px 10px;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
    background: #f8fafc;
    font-size: 12px;
    word-break: break-all;
  }

  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
    }
  }
</style>

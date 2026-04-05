<template>
  <div class="lab-page" v-loading="loading">
    <div class="hero">
      <div>
        <p class="eyebrow">{{ profile.classification }} / 后端可靠上报框架</p>
        <h2>{{ profile.title || 'Reliable Upload Framework' }}</h2>
        <p class="subtitle">{{ profile.summary }}</p>
      </div>
      <div class="hero-meta">
        <el-tag type="success">独立开源</el-tag>
        <el-tag type="primary">可迁移</el-tag>
        <el-tag>后端组件</el-tag>
        <el-button
          v-if="profile.githubUrl"
          type="primary"
          plain
          size="small"
          :icon="Link"
          @click="openGithub"
        >GitHub 仓库</el-button>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="15">
        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>核心能力</span>
              <div class="path-list">
                <code>{{ profile.packagePath }}</code>
              </div>
            </div>
          </template>
          <div class="bullet-list">
            <div v-for="item in profile.capabilityPoints || []" :key="item" class="bullet-item">
              {{ item }}
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>任务模型</span>
            </div>
          </template>
          <div class="task-grid">
            <div v-for="task in profile.taskModels || []" :key="task.key" class="task-card">
              <div class="task-top">
                <strong>{{ task.name }}</strong>
                <el-tag size="small">{{ task.triggerMode }}</el-tag>
              </div>
              <p class="muted">{{ task.scene }}</p>
              <div class="mini-list">
                <div v-for="point in task.highlights || []" :key="point">{{ point }}</div>
              </div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>核心契约</span>
            </div>
          </template>
          <div class="interface-list">
            <div v-for="item in profile.interfaces || []" :key="item.name" class="interface-card">
              <div class="interface-name">{{ item.name }}</div>
              <pre>{{ item.signature }}</pre>
              <div class="mini-list">
                <div v-for="note in item.notes || []" :key="note">{{ note }}</div>
              </div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>示例代码</span>
            </div>
          </template>
          <div class="snippet-list">
            <div v-for="item in profile.codeSnippets || []" :key="item.title" class="snippet-card">
              <div class="snippet-title">{{ item.title }}</div>
              <pre>{{ item.code }}</pre>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="9">
        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>引擎入口</span>
            </div>
          </template>
          <el-timeline>
            <el-timeline-item v-for="item in profile.engineEntries || []" :key="item.name" type="primary">
              <div class="timeline-title">{{ item.name }}</div>
              <div class="muted">{{ item.description }}</div>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>集成步骤</span>
            </div>
          </template>
          <div class="number-list">
            <div v-for="(item, index) in profile.integrationSteps || []" :key="item" class="number-item">
              <span>{{ index + 1 }}</span>
              <div>{{ item }}</div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>迁移注意点</span>
            </div>
          </template>
          <div class="bullet-list">
            <div v-for="item in profile.migrationNotes || []" :key="item" class="bullet-item">
              {{ item }}
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>开源文件</span>
              <el-button
                v-if="profile.githubUrl"
                link
                type="primary"
                size="small"
                :icon="Link"
                @click="openGithub"
              >在 GitHub 查看</el-button>
            </div>
          </template>
          <div class="file-list">
            <div v-for="item in profile.includedFiles || []" :key="item.path" class="file-item">
              <code>{{ item.path }}</code>
              <div class="muted">{{ item.role }}</div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>示例配置</span>
            </div>
          </template>
          <div class="config-list">
            <div v-for="item in profile.exampleConfigs || []" :key="item.taskCode" class="config-card">
              <div class="task-top">
                <strong>{{ item.taskCode }}</strong>
                <el-tag size="small" type="success">{{ item.taskType }}</el-tag>
              </div>
              <div class="muted">batchSize: {{ item.batchSize }} / maxRetry: {{ item.maxRetry }}</div>
              <div class="muted">prefix: {{ item.filePrefix }}</div>
              <div class="muted">remote: {{ item.sftpSubdir }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { onMounted, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Link } from '@element-plus/icons-vue'
  import { getReliableUploadProfile } from '@/api/reliableUpload'

  defineOptions({
    name: 'LabReusableReliableUpload'
  })

  const loading = ref(false)
  const profile = ref({})

  const loadProfile = async () => {
    loading.value = true
    try {
      const result = await getReliableUploadProfile()
      profile.value = result.data || {}
    } catch (error) {
      ElMessage.error(error?.message || '加载可靠上报框架资料失败')
    } finally {
      loading.value = false
    }
  }

  const openGithub = () => {
    window.open(profile.value.githubUrl, '_blank', 'noopener,noreferrer')
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
    background: linear-gradient(135deg, #f0f9ff 0%, #f8fafc 100%);
  }

  .eyebrow {
    margin: 0 0 8px;
    color: #0369a1;
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.08em;
  }

  .hero h2 {
    margin: 0 0 8px;
    color: #0f172a;
    font-size: 28px;
  }

  .subtitle,
  .muted {
    color: #475569;
    line-height: 1.7;
  }

  .subtitle {
    margin: 0;
    max-width: 760px;
  }

  .hero-meta {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    gap: 8px;
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    align-items: center;
    font-weight: 600;
  }

  .path-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
    align-items: flex-end;
    font-size: 12px;
  }

  .bullet-list,
  .mini-list,
  .file-list,
  .config-list,
  .snippet-list,
  .interface-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .bullet-item,
  .task-card,
  .interface-card,
  .file-item,
  .config-card,
  .snippet-card {
    padding: 14px 16px;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
    background: #f8fafc;
  }

  .task-grid {
    display: grid;
    gap: 12px;
  }

  .task-top {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    align-items: center;
    margin-bottom: 8px;
  }

  .number-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .number-item {
    display: grid;
    grid-template-columns: 28px 1fr;
    gap: 12px;
    align-items: start;
  }

  .number-item span {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 999px;
    background: #0ea5e9;
    color: #fff;
    font-size: 13px;
    font-weight: 700;
  }

  .timeline-title,
  .interface-name,
  .snippet-title {
    margin-bottom: 6px;
    font-weight: 600;
    color: #0f172a;
  }

  pre {
    margin: 0;
    white-space: pre-wrap;
    word-break: break-word;
    line-height: 1.65;
    color: #0f172a;
    font-family: Consolas, 'Courier New', monospace;
  }

  @media (max-width: 768px) {
    .hero,
    .panel-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .path-list {
      align-items: flex-start;
    }
  }
</style>

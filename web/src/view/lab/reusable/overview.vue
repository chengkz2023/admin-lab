<template>
  <div class="lab-page">
    <div class="lab-hero reusable-hero">
      <div>
        <p class="lab-eyebrow">复用组件</p>
        <h2>把公司高频业务需求沉淀成可搬运、可扩展的通用资产</h2>
        <p class="lab-subtitle">
          这里存放已经验证过、值得反复复用的模块。它们既可以直接搬进内网项目，也可以作为业务开发时的参考底座。
        </p>
      </div>
      <div class="lab-badges">
        <el-tag type="success">可沉淀</el-tag>
        <el-tag type="primary">可迁移</el-tag>
        <el-tag>可扩展</el-tag>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :md="8" v-for="card in cards" :key="card.title">
        <el-card class="lab-card" shadow="hover">
          <template #header>
            <div class="lab-card-header">
              <span>{{ card.title }}</span>
              <el-tag size="small" :type="card.tagType">{{ card.tag }}</el-tag>
            </div>
          </template>
          <p class="lab-card-text">{{ card.description }}</p>
          <p class="lab-card-tip">{{ card.tip }}</p>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never">
      <template #header>
        <div class="lab-card-header">
          <span>建议优先沉淀的复用资产</span>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col :xs="24" :md="12" :lg="8" v-for="asset in assets" :key="asset.title">
          <div class="asset-item">
            <div class="asset-title">{{ asset.title }}</div>
            <div class="asset-desc">{{ asset.description }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
  defineOptions({
    name: 'LabReusableOverview'
  })

  const cards = [
    {
      title: '通用处理骨架',
      tag: '基础层',
      tagType: 'primary',
      description: '把导入导出、附件处理、字典映射、状态转换等基础能力先抽成稳定骨架。',
      tip: '后续接业务时主要替换字段、规则和服务调用。'
    },
    {
      title: '页面级复用模块',
      tag: '页面层',
      tagType: 'success',
      description: '把搜索栏、表格工具条、批量操作栏、图表卡片区等页面级片段统一封装。',
      tip: '适合多业务线共用的后台页面结构。'
    },
    {
      title: '迁移说明模板',
      tag: '方法论',
      tagType: 'warning',
      description: '每个模块都记录依赖、菜单、接口、权限和环境差异，搬到内网时不再重复摸索。',
      tip: '让复用资产不止能用，还能稳定迁移。'
    }
  ]

  const assets = [
    { title: 'Excel 导入导出', description: '模板下载、批量导入、解析校验、错误反馈和数据导出。' },
    { title: '通用查询工具栏', description: '收拢筛选、重置、展开、导出、快捷搜索等常见交互。' },
    { title: '可配置表格能力', description: '列显隐、列顺序、固定列、批量操作和空状态处理。' },
    { title: '附件上传模块', description: '上传限制、回显、预览、重传和权限控制。' },
    { title: '图表容器组件', description: '统一卡片结构、筛选栏、空态、加载态和导出能力。' },
    { title: '业务状态组件', description: '统一状态标签、进度展示和不同业务的枚举映射。' }
  ]
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

  .lab-badges {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    gap: 8px;
  }

  .lab-card {
    height: 100%;
  }

  .lab-card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    font-weight: 600;
  }

  .lab-card-text {
    margin: 0 0 12px;
    color: #374151;
    line-height: 1.75;
  }

  .lab-card-tip {
    margin: 0;
    color: #6b7280;
    font-size: 13px;
  }

  .asset-item {
    height: 100%;
    margin-bottom: 16px;
    padding: 16px;
    border-radius: 12px;
    background: #f8fafc;
    border: 1px solid #e5e7eb;
  }

  .asset-title {
    margin-bottom: 8px;
    font-weight: 600;
    color: #111827;
  }

  .asset-desc {
    color: #6b7280;
    line-height: 1.7;
    font-size: 14px;
  }

  @media (max-width: 768px) {
    .lab-hero {
      flex-direction: column;
    }
  }
</style>

<template>
  <div class="lab-page">
    <div class="lab-hero">
      <div>
        <p class="lab-eyebrow">需求仿真</p>
        <h2>先在外网快速复刻真实需求，再搬进内网细化落地</h2>
        <p class="lab-subtitle">
          这里适合放和公司真实任务高度接近的页面样例，例如列表页、详情页、复杂表单、导入导出、图表分析和审批流页面。
        </p>
      </div>
      <div class="lab-badges">
        <el-tag type="primary">真实任务预演</el-tag>
        <el-tag>快速验证</el-tag>
        <el-tag type="success">便于搬运</el-tag>
      </div>
    </div>

    <el-row :gutter="16" class="mb-4">
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
          <span>推荐进入标准</span>
        </div>
      </template>
      <el-timeline>
        <el-timeline-item v-for="item in checklist" :key="item" type="primary">
          {{ item }}
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup>
  defineOptions({
    name: 'LabSimulationOverview'
  })

  const cards = [
    {
      title: '业务页面样例',
      tag: '页面级',
      tagType: 'primary',
      description: '先按公司里的真实需求做一版完整页面，把交互、字段和接口流转都跑通。',
      tip: '适合列表页、详情页、编辑页、审批页。'
    },
    {
      title: '复杂交互验证',
      tag: '高频',
      tagType: 'success',
      description: '把复杂筛选、批量操作、步骤流转、附件处理等高频业务交互先验证一遍。',
      tip: '减少内网环境里边试边改的成本。'
    },
    {
      title: '搬运前清单',
      tag: '迁移',
      tagType: 'warning',
      description: '把接口依赖、字段映射、权限点和内网环境差异记录下来，后续搬运更稳。',
      tip: '建议每个仿真页面都补迁移说明。'
    }
  ]

  const checklist = [
    '页面结构至少能完整走通核心业务流程。',
    '接口字段和状态流转要尽量贴近真实项目。',
    '记录哪些代码是通用能力，哪些是业务定制。',
    '搬运进内网前补上依赖、菜单、权限和配置差异说明。'
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
    background: linear-gradient(135deg, #f4f8ff 0%, #eef7f2 100%);
    border: 1px solid #dbe7ff;
  }

  .lab-eyebrow {
    margin: 0 0 8px;
    color: #4e80ee;
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
    max-width: 720px;
    color: #4b5563;
    line-height: 1.7;
  }

  .lab-badges {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    justify-content: flex-end;
    gap: 8px;
    min-width: 180px;
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

  @media (max-width: 768px) {
    .lab-hero {
      flex-direction: column;
    }

    .lab-badges {
      justify-content: flex-start;
      min-width: auto;
    }
  }
</style>

<template>
  <div class="lab-page">
    <div class="lab-hero reusable-hero">
      <div>
        <p class="lab-eyebrow">复用组件</p>
        <h2>把公司里反复碰到的需求沉淀成可复制、可裁剪的资产</h2>
        <p class="lab-subtitle">
          这里不只是展示组件，还应该附带接入方式、依赖说明、接口约定和迁移提示，方便你直接搬进内网复用。
        </p>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :md="12" v-for="group in groups" :key="group.title">
        <el-card class="lab-card" shadow="hover">
          <template #header>
            <div class="lab-card-header">
              <span>{{ group.title }}</span>
            </div>
          </template>
          <div class="chip-list">
            <el-tag v-for="item in group.items" :key="item" class="chip">{{ item }}</el-tag>
          </div>
          <p class="group-desc">{{ group.description }}</p>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never">
      <template #header>
        <div class="lab-card-header">
          <span>建议每个复用资产附带的信息</span>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="组件职责">这个组件解决什么问题，适用于哪些场景。</el-descriptions-item>
        <el-descriptions-item label="依赖清单">依赖哪些前端库、后端接口、字典或权限点。</el-descriptions-item>
        <el-descriptions-item label="接入方式">最小接入步骤、必要 props、返回结构和示例代码。</el-descriptions-item>
        <el-descriptions-item label="迁移提示">搬到内网时需要替换哪些 API、配置和业务字段。</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
  defineOptions({
    name: 'LabReusable'
  })

  const groups = [
    {
      title: '通用业务能力',
      items: ['导入导出', '附件上传', '字典选择器', '状态标签', '批量操作栏'],
      description: '优先沉淀最常复用、跨业务都能用上的能力组件。'
    },
    {
      title: '页面装配能力',
      items: ['通用搜索栏', '表格列配置', '图表容器', '详情布局', '弹窗表单'],
      description: '把页面层高频重复结构抽出来，后续开发会明显提速。'
    }
  ]
</script>

<style scoped>
  .lab-page {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .lab-hero {
    padding: 24px;
    border-radius: 16px;
    border: 1px solid #d5efe1;
    background: linear-gradient(135deg, #f3fbf6 0%, #f2f8ff 100%);
  }

  .reusable-hero {
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
    color: #4b5563;
    line-height: 1.7;
  }

  .lab-card {
    height: 100%;
  }

  .lab-card-header {
    font-weight: 600;
  }

  .chip-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
  }

  .chip {
    margin-right: 0;
  }

  .group-desc {
    margin: 0;
    color: #6b7280;
    line-height: 1.7;
  }
</style>

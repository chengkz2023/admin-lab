<template>
  <div class="lab-page">
    <div class="lab-hero demo-hero">
      <div>
        <p class="lab-eyebrow">组件示例</p>
        <h2>集中试新组件、新交互和你暂时还不熟的能力点</h2>
        <p class="lab-subtitle">
          这里更适合做最小可运行示例，重点是快速理解组件的能力边界、接入成本和后续在业务里的使用方式。
        </p>
      </div>
    </div>

    <el-card shadow="never">
      <template #header>
        <div class="lab-card-header">
          <span>建议优先建设的示例类型</span>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col :xs="24" :md="12" :lg="8" v-for="example in examples" :key="example.title">
          <div class="example-item">
            <div class="example-title">{{ example.title }}</div>
            <div class="example-desc">{{ example.description }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <el-alert
      type="info"
      :closable="false"
      title="落地建议：每个示例都尽量附一段“适用场景 / 关键配置 / 踩坑点 / 是否推荐进业务”的说明。"
    />
  </div>
</template>

<script setup>
  defineOptions({
    name: 'LabComponentDemo'
  })

  const examples = [
    { title: '高级表格', description: '列配置、固定列、汇总、虚拟滚动、批量操作。' },
    { title: '富文本与附件', description: '编辑器、预览、上传限制、回显和权限控制。' },
    { title: '图表与可视化', description: '折线图、柱状图、联动筛选、空数据态和导出。' },
    { title: '拖拽与排序', description: '拖拽卡片、拖拽列表、可视化布局编排。' },
    { title: '复杂表单', description: '联动字段、动态区块、步骤表单和校验封装。' },
    { title: '文档能力', description: 'Excel、PDF、Word 预览、导入、导出和模板下载。' }
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
    border: 1px solid #eadfff;
    background: linear-gradient(135deg, #fbf7ff 0%, #f4fbff 100%);
  }

  .demo-hero {
    background: linear-gradient(135deg, #f9f4ff 0%, #eef9ff 100%);
  }

  .lab-eyebrow {
    margin: 0 0 8px;
    color: #8b5cf6;
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

  .lab-card-header {
    font-weight: 600;
  }

  .example-item {
    height: 100%;
    margin-bottom: 16px;
    padding: 16px;
    border-radius: 12px;
    background: #f8fafc;
    border: 1px solid #e5e7eb;
  }

  .example-title {
    margin-bottom: 8px;
    font-weight: 600;
    color: #111827;
  }

  .example-desc {
    color: #6b7280;
    line-height: 1.7;
    font-size: 14px;
  }
</style>

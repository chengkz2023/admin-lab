<template>
  <div class="page-wrap">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / 通用图表组件</p>
        <h2>面向后台高频统计分析场景的趋势图与分布图封装</h2>
        <p class="subtitle">
          收录 `TrendChart` 与 `PieChart` 两个通用图表组件，覆盖多序列趋势、双轴、时间对齐、饼图占比、柱图分布、格式化与点击事件透出等常见后台可视化能力。
        </p>
      </div>
      <div class="hero-tags">
        <el-tag type="primary">ECharts</el-tag>
        <el-tag type="success">可迁移</el-tag>
        <el-tag>通用图表底座</el-tag>
      </div>
    </div>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>组件能力总览</span></div>
      </template>
      <el-table :data="featureRows" border>
        <el-table-column prop="component" label="组件" min-width="160" />
        <el-table-column prop="feature" label="能力" min-width="220" />
        <el-table-column prop="desc" label="说明" min-width="360" />
      </el-table>
    </el-card>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="14">
        <el-card shadow="never">
          <template #header>
            <div class="panel-title"><span>TrendChart 演示</span></div>
          </template>
          <trend-chart
            :data="trendData"
            :series="trendSeries"
            x-axis-key="time"
            y-axis-name="访问量"
            y-axis-right-name="转化率"
            legend="top"
            theme="techBlue"
            :height="320"
            @registry-event="handleTrendRegistry"
          />
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="10">
        <el-card shadow="never">
          <template #header>
            <div class="panel-title"><span>PieChart 演示</span></div>
          </template>
          <pie-chart
            :data="distributionData"
            chart-type="pie"
            name-key="name"
            value-key="value"
            legend="bottom"
            :show-pie-percent="true"
            :height="320"
            @chart-click="handlePieClick"
          />
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>接入要点</span></div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="组件目录">`web/src/components/charts`</el-descriptions-item>
        <el-descriptions-item label="TrendChart 适用场景">趋势分析、双轴指标、时序监控、周期统计。</el-descriptions-item>
        <el-descriptions-item label="PieChart 适用场景">分布占比、枚举统计、渠道结构、风险等级可视化。</el-descriptions-item>
        <el-descriptions-item label="配色来源">`generateColorStyles.js` 统一生成渐变色与图表风格。</el-descriptions-item>
        <el-descriptions-item label="事件接入">TrendChart 通过 `@registry-event` 暴露实例；PieChart 通过 `@chart-click` 暴露点击事件。</el-descriptions-item>
        <el-descriptions-item label="迁移建议">迁入内网时同步复制组件目录、配色工具文件、示例页，并统一核对 ECharts 版本与主题风格。</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="panel-title"><span>TrendChart 最小接入示例</span></div>
          </template>
          <pre class="code-block"><code>{{ trendUsageCode }}</code></pre>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="panel-title"><span>PieChart 最小接入示例</span></div>
          </template>
          <pre class="code-block"><code>{{ pieUsageCode }}</code></pre>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { ElMessage } from 'element-plus'
  import TrendChart from '@/components/charts/TrendChart.vue'
  import PieChart from '@/components/charts/PieChart.vue'

  defineOptions({
    name: 'LabReusableCharts'
  })

  const featureRows = [
    { component: 'TrendChart', feature: '多序列趋势图', desc: '支持折线 / 柱状混合配置，适合趋势分析与指标对比。' },
    { component: 'TrendChart', feature: '双 Y 轴', desc: '支持左右双轴配置，便于同时展示数量与比例。' },
    { component: 'TrendChart', feature: '时间对齐', desc: '支持开始时间、结束时间、粒度与补齐策略。' },
    { component: 'TrendChart', feature: '实例透出', desc: '通过 registry-event 向外暴露 ECharts 实例，便于扩展交互。' },
    { component: 'PieChart', feature: '饼图 / 柱图切换', desc: '同一组件支持 pie / bar 两种图形形态。' },
    { component: 'PieChart', feature: '名称映射与截断', desc: '适合后端字段名国际化或超长名称显示控制。' },
    { component: 'PieChart', feature: '点击事件透出', desc: '对外暴露 chart-click，方便实现筛选联动与下钻。' },
    { component: 'PieChart', feature: '方向与标签控制', desc: '柱状图支持横向 / 纵向，以及柱顶值显示开关。' }
  ]

  const trendData = [
    { time: '04-10', visits: 120, rate: 18 },
    { time: '04-11', visits: 168, rate: 22 },
    { time: '04-12', visits: 156, rate: 19 },
    { time: '04-13', visits: 210, rate: 25 },
    { time: '04-14', visits: 246, rate: 29 }
  ]

  const trendSeries = {
    visits: { label: '访问量', type: 'bar', yAxis: 'left' },
    rate: { label: '转化率', type: 'line', yAxis: 'right', formatter: (value) => `${value}%` }
  }

  const distributionData = [
    { name: '高危', value: 38 },
    { name: '中危', value: 46 },
    { name: '低危', value: 16 }
  ]

  const handleTrendRegistry = () => {
    // 保留接入点，示例页不额外挂复杂事件
  }

  const handlePieClick = (params) => {
    ElMessage.info(`点击了 ${params?.name || '未知项'}：${params?.value ?? '-'} `)
  }

  const trendUsageCode = `<TrendChart
  :data="trendData"
  :series="{
    orderCount: { label: '订单量', type: 'bar', yAxis: 'left' },
    successRate: { label: '成功率', type: 'line', yAxis: 'right', formatter: (v) => v + '%' }
  }"
  x-axis-key="date"
  y-axis-name="数量"
  y-axis-right-name="比例"
  :height="320"
/>`

  const pieUsageCode = `<PieChart
  :data="distributionData"
  chart-type="pie"
  name-key="name"
  value-key="value"
  legend="bottom"
  :show-pie-percent="true"
  :height="320"
  @chart-click="handleClick"
/>`
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
    font-weight: 600;
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
  }
</style>

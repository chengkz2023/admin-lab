<template>
  <div class="page-wrap">
    <div class="hero">
      <div>
        <p class="eyebrow">组件示例 / 通用图表示例</p>
        <h2>聚焦组件 API、配置变化与交互边界的纯能力演示</h2>
        <p class="subtitle">
          这一页不强调业务接入，而是集中展示 `TrendChart` 和 `PieChart` 的常见玩法，方便快速理解 props、事件、图形切换和边界表现。
        </p>
      </div>
      <div class="hero-tags">
        <el-tag type="primary">Playground</el-tag>
        <el-tag type="success">最小示例</el-tag>
        <el-tag>组件能力验证</el-tag>
      </div>
    </div>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>在线控制台</span></div>
      </template>
      <el-row :gutter="16">
        <el-col :xs="24" :md="6">
          <div class="control-item">
            <div class="control-label">PieChart 形态</div>
            <el-radio-group v-model="pieChartType">
              <el-radio-button label="饼图" value="pie" />
              <el-radio-button label="柱图" value="bar" />
            </el-radio-group>
          </div>
        </el-col>
        <el-col :xs="24" :md="6">
          <div class="control-item">
            <div class="control-label">柱图方向</div>
            <el-radio-group v-model="barDirection">
              <el-radio-button label="纵向" value="vertical" />
              <el-radio-button label="横向" value="horizontal" />
            </el-radio-group>
          </div>
        </el-col>
        <el-col :xs="24" :md="6">
          <div class="control-item">
            <div class="control-label">图例位置</div>
            <el-select v-model="legendPosition" style="width: 100%">
              <el-option label="顶部" value="top" />
              <el-option label="底部" value="bottom" />
              <el-option label="右侧" value="right" />
              <el-option label="左侧" value="left" />
            </el-select>
          </div>
        </el-col>
        <el-col :xs="24" :md="6">
          <div class="control-item switch-group">
            <el-switch v-model="showPiePercent" active-text="显示饼图百分比" />
            <el-switch v-model="showBarValue" active-text="显示柱顶数值" />
          </div>
        </el-col>
      </el-row>
    </el-card>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="panel-title"><span>示例 A：TrendChart 基础折线</span></div>
          </template>
          <trend-chart
            :data="basicTrendData"
            :series="basicTrendSeries"
            x-axis-key="time"
            legend="top"
            theme="default"
            :height="300"
          />
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="panel-title"><span>示例 B：TrendChart 柱线双轴</span></div>
          </template>
          <trend-chart
            :data="mixTrendData"
            :series="mixTrendSeries"
            x-axis-key="time"
            y-axis-name="数量"
            y-axis-right-name="比例"
            legend="top"
            theme="techBlue"
            :height="300"
            @registry-event="handleRegistryEvent"
          />
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="panel-title"><span>示例 C：PieChart 交互切换</span></div>
          </template>
          <pie-chart
            :data="distributionData"
            :chart-type="pieChartType"
            name-key="name"
            value-key="value"
            :legend="legendPosition"
            :show-pie-percent="showPiePercent"
            :show-bar-value="showBarValue"
            :bar-direction="barDirection"
            :height="320"
            @chart-click="handleChartClick"
          />
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="panel-title"><span>示例 D：名称映射 / 超长文本 / 格式化</span></div>
          </template>
          <pie-chart
            :data="mappedDistributionData"
            chart-type="bar"
            name-key="key"
            value-key="count"
            legend="none"
            :name-map="nameMap"
            :name-max-length="8"
            :formatter="formatCount"
            :show-bar-value="true"
            bar-direction="horizontal"
            :show-y-axis="true"
            :height="320"
          />
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never">
      <template #header>
        <div class="panel-title"><span>能力清单</span></div>
      </template>
      <el-table :data="featureRows" border>
        <el-table-column prop="component" label="组件" min-width="160" />
        <el-table-column prop="feature" label="能力点" min-width="220" />
        <el-table-column prop="example" label="当前示例" min-width="220" />
        <el-table-column prop="desc" label="说明" min-width="360" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import TrendChart from '@/components/charts/TrendChart.vue'
  import PieChart from '@/components/charts/PieChart.vue'

  defineOptions({
    name: 'LabComponentDemoCharts'
  })

  const pieChartType = ref('pie')
  const barDirection = ref('vertical')
  const legendPosition = ref('bottom')
  const showPiePercent = ref(true)
  const showBarValue = ref(true)

  const basicTrendData = [
    { time: 'Mon', cpu: 35 },
    { time: 'Tue', cpu: 48 },
    { time: 'Wed', cpu: 42 },
    { time: 'Thu', cpu: 56 },
    { time: 'Fri', cpu: 51 }
  ]

  const basicTrendSeries = {
    cpu: { label: 'CPU 使用率', type: 'line', yAxis: 'left', formatter: (value) => `${value}%` }
  }

  const mixTrendData = [
    { time: '04-10', events: 120, rate: 18 },
    { time: '04-11', events: 168, rate: 23 },
    { time: '04-12', events: 146, rate: 20 },
    { time: '04-13', events: 210, rate: 29 },
    { time: '04-14', events: 188, rate: 26 }
  ]

  const mixTrendSeries = {
    events: { label: '事件量', type: 'bar', yAxis: 'left' },
    rate: { label: '处置率', type: 'line', yAxis: 'right', formatter: (value) => `${value}%` }
  }

  const distributionData = [
    { name: '高危', value: 38 },
    { name: '中危', value: 46 },
    { name: '低危', value: 16 }
  ]

  const mappedDistributionData = [
    { key: 'veryLongThreatCategoryNameA', count: 12600 },
    { key: 'veryLongThreatCategoryNameB', count: 9800 },
    { key: 'veryLongThreatCategoryNameC', count: 7600 }
  ]

  const nameMap = {
    veryLongThreatCategoryNameA: '恶意扫描流量',
    veryLongThreatCategoryNameB: '暴力破解尝试',
    veryLongThreatCategoryNameC: '异常登录行为'
  }

  const formatCount = (value) => `${Number(value || 0).toLocaleString()} 次`

  const handleChartClick = (params) => {
    ElMessage.info(`当前点击：${params?.name || '未知'} / ${params?.value ?? '-'}`)
  }

  const handleRegistryEvent = () => {
    // 保留实例注册扩展点
  }

  const featureRows = [
    { component: 'TrendChart', feature: '基础折线图', example: '示例 A', desc: '验证最小接入的数据、series、x-axis-key 配置。' },
    { component: 'TrendChart', feature: '柱线双轴', example: '示例 B', desc: '验证左右双轴、不同图形类型组合与 formatter。' },
    { component: 'TrendChart', feature: '实例事件透出', example: '示例 B', desc: '通过 registry-event 暴露图表实例，便于外部扩展点击与缩放。' },
    { component: 'PieChart', feature: '饼图 / 柱图切换', example: '示例 C', desc: '验证 chart-type、legend、百分比与柱顶值联动。' },
    { component: 'PieChart', feature: '横向 / 纵向柱图', example: '示例 C', desc: '验证 bar-direction 对展示布局的影响。' },
    { component: 'PieChart', feature: '名称映射与格式化', example: '示例 D', desc: '验证 name-map、name-max-length 和 formatter 的组合效果。' }
  ]
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
    border: 1px solid #eadfff;
    background: linear-gradient(135deg, #f9f4ff 0%, #eef9ff 100%);
  }

  .eyebrow {
    margin: 0 0 8px;
    color: #8b5cf6;
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.08em;
  }

  .hero h2 {
    margin: 0 0 8px;
    font-size: 24px;
    color: #1f2937;
  }

  .subtitle {
    margin: 0;
    max-width: 760px;
    color: #4b5563;
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

  .control-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .control-label {
    font-size: 13px;
    color: #6b7280;
  }

  .switch-group {
    padding-top: 22px;
    gap: 12px;
  }

  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
    }
  }
</style>

<template>
  <div class="page-wrap">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / 网安可视化面板</p>
        <h2>企业级图表组件全特性展示</h2>
        <p class="subtitle">
          覆盖网安场景高频能力：多柱/多线、维度切换、双轴混合、点击下钻、图表与表格切换。示例数据已切换为后端接口返回。
        </p>
      </div>
      <div class="hero-tags">
        <el-tag type="primary">Apache ECharts</el-tag>
        <el-tag type="success">可复用</el-tag>
        <el-tag type="warning">企业场景</el-tag>
      </div>
    </div>

    <lab-list-query-bar
      v-model="queryForm"
      :items="queryItems"
      :loading="loading"
      :show-export="true"
      @search="loadData"
      @reset="handleReset"
      @export="handleExport"
    />

    <el-row :gutter="14">
      <el-col v-for="item in summaryCards" :key="item.key" :xs="24" :sm="12" :lg="6">
        <el-card class="kpi-card" shadow="never">
          <div class="kpi-title">{{ item.title }}</div>
          <div class="kpi-value">{{ item.value }}</div>
          <div class="kpi-sub">{{ item.sub }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="feature-card" shadow="never">
      <template #header>
        <div class="feature-title">交互控制台</div>
      </template>
      <el-row :gutter="12">
        <el-col :xs="24" :md="8">
          <div class="control-item">
            <div class="control-label">图表类型</div>
            <el-radio-group v-model="playgroundType">
              <el-radio-button value="bar">柱状图</el-radio-button>
              <el-radio-button value="line">折线图</el-radio-button>
              <el-radio-button value="pie">饼图</el-radio-button>
            </el-radio-group>
          </div>
        </el-col>
        <el-col :xs="24" :md="8">
          <div class="control-item">
            <div class="control-label">图表高度</div>
            <el-select v-model="playgroundHeight" style="width: 100%">
              <el-option v-for="item in heightOptions" :key="item" :label="item" :value="item" />
            </el-select>
          </div>
        </el-col>
        <el-col :xs="24" :md="8">
          <div class="control-item">
            <div class="control-label">加载状态</div>
            <el-switch v-model="playgroundLoading" />
          </div>
        </el-col>
      </el-row>
    </el-card>

    <lab-security-chart-panel
      title="能力示例 A：基础能力"
      subtitle="演示 type、height、loading、slot:extra"
      :type="playgroundType"
      :height="playgroundHeight"
      :loading="playgroundLoading || loading"
      :data="playgroundData"
    >
      <template #extra>
        <el-button type="primary" link @click="loadData">刷新数据</el-button>
      </template>
    </lab-security-chart-panel>

    <el-row :gutter="14">
      <el-col :xs="24" :lg="12">
        <lab-security-chart-panel
          title="能力示例 B：维度切换（受控）"
          subtitle="通过 v-model:dimension-value + @dimension-change 获取当前维度"
          type="line"
          :loading="loading"
          :dimension-options="lineDimensionOptions"
          v-model:dimension-value="trendDimension"
          :dimension-data-map="lineDataMap"
          @dimension-change="handleDimensionChange"
        >
          <template #extra>
            <el-tag type="info">当前维度：{{ currentTrendDimensionLabel }}</el-tag>
          </template>
        </lab-security-chart-panel>
      </el-col>
      <el-col :xs="24" :lg="12">
        <lab-security-chart-panel
          title="能力示例 C：点击下钻"
          subtitle="点击图元进入事件明细（chart-click）"
          type="bar"
          :loading="loading"
          :dimension-options="commonDimensionOptions"
          :dimension-data-map="barDimensionDataMap"
          @chart-click="openDrilldown"
        />
      </el-col>
    </el-row>

    <el-row :gutter="14">
      <el-col :xs="24" :lg="12">
        <lab-security-chart-panel
          title="能力示例 D：饼图维度切换"
          subtitle="同一组件在饼图下支持多维度映射"
          type="pie"
          :loading="loading"
          :dimension-options="commonDimensionOptions"
          :dimension-data-map="pieDimensionDataMap"
        />
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card class="integration-card" shadow="never">
          <template #header>
            <div class="panel-title">组件配置清单</div>
          </template>
          <el-table :data="featureTableData" size="small" border>
            <el-table-column prop="name" label="配置项 / 事件" min-width="170" />
            <el-table-column prop="desc" label="说明" min-width="240" />
            <el-table-column prop="example" label="示例值" min-width="220" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="feature-card" shadow="never">
      <template #header>
        <div class="feature-title">能力示例 E：双轴混合图 + 图表/表格切换</div>
      </template>
      <div class="mix-tools">
        <el-radio-group v-model="mixViewMode" size="small">
          <el-radio-button value="chart">图表视图</el-radio-button>
          <el-radio-button value="table">表格视图</el-radio-button>
        </el-radio-group>
      </div>
      <lab-security-chart-panel
        v-if="mixViewMode === 'chart'"
        title="攻防态势混合分析"
        subtitle="柱状展示事件量，折线展示处置率（双轴）"
        type="mix"
        height="340px"
        :loading="loading"
        :data="mixData"
      />
      <el-table v-else :data="mixTableData" border size="small">
        <el-table-column prop="time" label="时间" min-width="120" />
        <el-table-column prop="eventCount" label="事件量" min-width="110" />
        <el-table-column prop="blockedCount" label="拦截量" min-width="110" />
        <el-table-column prop="handleRate" label="处置率(%)" min-width="120" />
      </el-table>
    </el-card>

    <el-drawer v-model="drilldownVisible" size="52%" title="下钻事件明细">
      <div class="drilldown-title">{{ drilldownContext }}</div>
      <el-table :data="drilldownTableData" border size="small">
        <el-table-column prop="eventId" label="事件编号" min-width="150" />
        <el-table-column prop="attackType" label="攻击类型" min-width="120" />
        <el-table-column prop="sourceIp" label="来源IP" min-width="130" />
        <el-table-column prop="status" label="状态" min-width="100" />
        <el-table-column prop="level" label="等级" min-width="100" />
        <el-table-column prop="occurAt" label="发生时间" min-width="170" />
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup>
  import { computed, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import LabListQueryBar from '@/components/lab/list-query-bar.vue'
  import LabSecurityChartPanel from '@/components/lab/security-chart-panel.vue'
  import { getPanelBinding, normalizeDashboardPayload } from '@/components/lab/security-chart-adapter'
  import { getSecurityDashboardDrilldown, getSecurityDashboardPanel } from '@/api/securityDashboard'

  defineOptions({
    name: 'LabReusableSecurityEcharts'
  })

  const loading = ref(false)
  const playgroundLoading = ref(false)
  const playgroundType = ref('bar')
  const playgroundHeight = ref('300px')
  const heightOptions = ['240px', '300px', '360px', '420px']

  const mixViewMode = ref('chart')
  const drilldownVisible = ref(false)
  const drilldownContext = ref('')
  const drilldownTableData = ref([])

  const queryForm = ref({
    keyword: '',
    threatLevel: [],
    dataSource: '',
    timeRange: []
  })

  const queryItems = [
    { prop: 'keyword', label: '关键字', type: 'input', placeholder: 'IP / 域名 / 事件编号', span: 8 },
    {
      prop: 'threatLevel',
      label: '风险等级',
      type: 'select',
      multiple: true,
      defaultValue: [],
      span: 6,
      options: [
        { label: '高危', value: 'high' },
        { label: '中危', value: 'medium' },
        { label: '低危', value: 'low' }
      ]
    },
    {
      prop: 'dataSource',
      label: '数据源',
      type: 'select',
      span: 5,
      options: [
        { label: 'SOC 平台', value: 'soc' },
        { label: '流量探针', value: 'probe' },
        { label: '态势平台', value: 'situation' }
      ]
    },
    { prop: 'timeRange', label: '时间范围', type: 'dateRange', span: 7 }
  ]

  const summaryCards = ref([])
  const baseBarData = ref({ categories: [], series: [] })
  const baseLineData = ref({ categories: [], series: [] })
  const basePieData = ref({ data: [] })
  const mixData = ref({ categories: [], barSeries: [], lineSeries: [] })
  const mixTableData = ref([])

  const trendDimension = ref('hour')
  const trendDimensionOptions = [
    { label: '小时', value: 'hour' },
    { label: '天', value: 'day' },
    { label: '周', value: 'week' }
  ]
  const dimensionOptionsSimple = [
    { label: '近24小时', value: '24h' },
    { label: '近7天', value: '7d' }
  ]
  const lineDataMap = ref({})
  const barDimensionDataMap = ref({})
  const pieDimensionDataMap = ref({})
  const lineDimensionOptions = ref([...trendDimensionOptions])
  const commonDimensionOptions = ref([...dimensionOptionsSimple])

  const featureTableData = [
    { name: 'title / subtitle', desc: '图表主副标题', example: '攻防态势混合分析' },
    { name: 'type', desc: '图表类型', example: 'bar / line / pie / mix' },
    { name: 'data.series', desc: '多柱/多线数据', example: '[{ name, values, color }]' },
    { name: 'height / loading', desc: '高度与加载状态', example: '300px / true' },
    { name: 'dimension-options', desc: '维度选项列表', example: '[{ label, value }]' },
    { name: 'v-model:dimension-value', desc: '受控维度切换', example: 'v-model:dimension-value="state"' },
    { name: 'dimension-data-map', desc: '各维度数据映射', example: '{ hour: {...}, day: {...} }' },
    { name: '@dimension-change', desc: '维度切换回调', example: '(value) => {}' },
    { name: '@chart-click', desc: '图表点击事件（下钻）', example: '(payload) => {}' },
    { name: 'slot:extra', desc: '头部右侧操作区', example: '刷新按钮 / 状态标签' },
    { name: 'mix 数据结构', desc: '双轴混合图', example: '{ categories, barSeries, lineSeries }' }
  ]

  const currentTrendDimensionLabel = computed(() => trendDimensionOptions.find((item) => item.value === trendDimension.value)?.label || '未知')

  const playgroundData = computed(() => {
    if (playgroundType.value === 'line') return baseLineData.value
    if (playgroundType.value === 'pie') return basePieData.value
    return baseBarData.value
  })

  const formatRangeValue = (value) => {
    if (!value) return ''
    return typeof value === 'string' ? value : `${value}`
  }

  const buildPanelQueryParams = () => {
    const [timeFrom, timeTo] = queryForm.value.timeRange || []
    return {
      keyword: queryForm.value.keyword || '',
      threatLevels: Array.isArray(queryForm.value.threatLevel) ? queryForm.value.threatLevel.join(',') : '',
      dataSource: queryForm.value.dataSource || '',
      timeFrom: formatRangeValue(timeFrom),
      timeTo: formatRangeValue(timeTo),
      trendDimension: trendDimension.value
    }
  }

  const applyDashboardPayload = (payload) => {
    const normalized = normalizeDashboardPayload(payload)
    const barBinding = getPanelBinding(normalized, 'bar')
    const lineBinding = getPanelBinding(normalized, 'line')
    const pieBinding = getPanelBinding(normalized, 'pie')
    const mixBinding = getPanelBinding(normalized, 'mix')

    summaryCards.value = normalized.summary
    baseBarData.value = barBinding.data
    baseLineData.value = lineBinding.data
    basePieData.value = pieBinding.data
    lineDataMap.value = lineBinding.dimensionDataMap
    barDimensionDataMap.value = barBinding.dimensionDataMap
    pieDimensionDataMap.value = pieBinding.dimensionDataMap
    lineDimensionOptions.value = lineBinding.dimensionOptions.length ? lineBinding.dimensionOptions : [...trendDimensionOptions]
    commonDimensionOptions.value = barBinding.dimensionOptions.length ? barBinding.dimensionOptions : [...dimensionOptionsSimple]
    if (lineBinding.dimensionValue !== null && lineBinding.dimensionValue !== undefined) {
      trendDimension.value = lineBinding.dimensionValue
    }
    mixData.value = mixBinding.data
    mixTableData.value = normalized.meta?.mixTableData ?? []
  }

  const loadData = async () => {
    loading.value = true
    try {
      const result = await getSecurityDashboardPanel(buildPanelQueryParams())
      applyDashboardPayload(result?.data ?? {})
    } catch (error) {
      ElMessage.error(error?.response?.data?.msg || error?.message || '加载可视化面板失败')
    } finally {
      loading.value = false
    }
  }

  const openDrilldown = async ({ params, dimensionValue }) => {
    const category = params?.name || '未知类别'
    const series = params?.seriesName || '未知序列'
    drilldownContext.value = `维度：${dimensionValue || '默认'} / 类别：${category} / 序列：${series}`
    try {
      const [timeFrom, timeTo] = queryForm.value.timeRange || []
      const result = await getSecurityDashboardDrilldown({
        category,
        series,
        dimension: dimensionValue || '',
        keyword: queryForm.value.keyword || '',
        threatLevels: Array.isArray(queryForm.value.threatLevel) ? queryForm.value.threatLevel.join(',') : '',
        dataSource: queryForm.value.dataSource || '',
        timeFrom: formatRangeValue(timeFrom),
        timeTo: formatRangeValue(timeTo)
      })
      drilldownTableData.value = Array.isArray(result?.data) ? result.data : []
      drilldownVisible.value = true
    } catch (error) {
      ElMessage.error(error?.response?.data?.msg || error?.message || '加载下钻明细失败')
    }
  }

  const handleDimensionChange = (value) => {
    ElMessage.info(`已切换趋势维度：${value}`)
  }

  const handleReset = async () => {
    await loadData()
    ElMessage.success('已重置筛选并刷新图表。')
  }

  const handleExport = () => {
    ElMessage.info('导出事件已触发，可对接后端报表导出接口。')
  }

  loadData()
</script>

<style scoped>
  .page-wrap {
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .hero {
    display: flex;
    justify-content: space-between;
    gap: 16px;
    padding: 22px;
    border-radius: 16px;
    border: 1px solid #c9dbff;
    background: linear-gradient(140deg, #eef4ff 0%, #f8fbff 46%, #f1f7ff 100%);
  }

  .eyebrow {
    margin: 0 0 8px;
    color: #1f4ed8;
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
    max-width: 820px;
    color: #475569;
    line-height: 1.7;
  }

  .hero-tags {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    gap: 8px;
  }

  .kpi-card {
    border: 1px solid #dbe7ff;
    border-radius: 12px;
    background: linear-gradient(180deg, #fff 0%, #f8fbff 100%);
  }

  .kpi-title {
    color: #64748b;
    font-size: 13px;
  }

  .kpi-value {
    margin: 8px 0 6px;
    color: #0f172a;
    font-weight: 700;
    font-size: 28px;
    line-height: 1.2;
  }

  .kpi-sub {
    color: #0ea5a0;
    font-size: 12px;
  }

  .feature-card {
    border: 1px solid #dbe7ff;
    border-radius: 12px;
    background: linear-gradient(180deg, #fff 0%, #f9fcff 100%);
  }

  .feature-title {
    font-size: 15px;
    font-weight: 600;
    color: #0f172a;
  }

  .control-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 6px;
  }

  .control-label {
    font-size: 13px;
    color: #475569;
    font-weight: 600;
  }

  .integration-card {
    border: 1px solid #dbe7ff;
    border-radius: 14px;
    background: linear-gradient(180deg, #ffffff 0%, #f8fbff 100%);
  }

  .panel-title {
    font-size: 16px;
    font-weight: 600;
    color: #0f172a;
  }

  .mix-tools {
    margin-bottom: 10px;
  }

  .drilldown-title {
    margin-bottom: 10px;
    color: #334155;
    font-size: 13px;
  }

  @media (max-width: 768px) {
    .hero {
      flex-direction: column;
    }
  }
</style>

const securityPalette = ['#1f6feb', '#00a3bf', '#36cfc9', '#73d13d', '#faad14', '#ff7a45', '#9254de']
// 语义配色映射：后端返回 metricKey，前端统一决定视觉样式。
const semanticColorMap = {
  handled_count: '#1f6feb',
  pending_count: '#faad14',
  external_attack_count: '#00a3bf',
  internal_anomaly_count: '#1f6feb',
  event_count: '#1f6feb',
  blocked_count: '#36cfc9',
  handle_rate: '#faad14'
}
const semanticAreaMetricKeys = new Set(['external_attack_count'])

const axisStyle = {
  axisLine: { lineStyle: { color: '#d9e2f3' } },
  axisTick: { show: false },
  axisLabel: { color: '#64748b', fontSize: 12 },
  splitLine: { lineStyle: { color: '#edf2f7' } }
}

const hexToRgba = (hex, alpha) => {
  const value = hex.replace('#', '')
  const short = value.length === 3
  const full = short
    ? value
        .split('')
        .map((item) => item + item)
        .join('')
    : value
  const r = parseInt(full.slice(0, 2), 16)
  const g = parseInt(full.slice(2, 4), 16)
  const b = parseInt(full.slice(4, 6), 16)
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

const resolveSeriesColor = (item, index) => {
  if (item?.color) return item.color
  if (item?.metricKey && semanticColorMap[item.metricKey]) return semanticColorMap[item.metricKey]
  return securityPalette[index % securityPalette.length]
}

const resolveLineArea = (item) => {
  if (typeof item?.area === 'boolean') return item.area
  if (item?.metricKey) return semanticAreaMetricKeys.has(item.metricKey)
  return false
}

// 企业级柱图通用样式，保证各业务图表视觉风格一致。
const buildBarSeriesStyle = (seriesColor) => ({
  color: {
    type: 'linear',
    x: 0,
    y: 0,
    x2: 0,
    y2: 1,
    colorStops: [
      { offset: 0, color: seriesColor },
      { offset: 1, color: hexToRgba(seriesColor, 0.72) }
    ]
  },
  borderRadius: [8, 8, 2, 2],
  shadowBlur: 10,
  shadowColor: hexToRgba(seriesColor, 0.25),
  shadowOffsetY: 4
})

export const buildBarOption = ({
  categories = [],
  seriesName = '数量',
  values = [],
  color = securityPalette[0],
  series = []
} = {}) => {
  // 向后兼容：同时支持旧版单序列（values）与新版多序列（series[]）。
  const hasMultiSeries = Array.isArray(series) && series.length > 0
  const normalizedSeries = hasMultiSeries
    ? series.map((item, index) => ({
        name: item?.name || `系列${index + 1}`,
        type: 'bar',
        stack: item?.stack,
        barMaxWidth: item?.barMaxWidth ?? 22,
        barMinHeight: 2,
        data: item?.values ?? [],
        itemStyle: buildBarSeriesStyle(resolveSeriesColor(item, index)),
        showBackground: true,
        backgroundStyle: { color: '#f2f6ff' },
        emphasis: { focus: 'series' }
      }))
    : [
        {
          name: seriesName,
          type: 'bar',
          barMaxWidth: 22,
          barMinHeight: 2,
          data: values,
          itemStyle: buildBarSeriesStyle(color),
          showBackground: true,
          backgroundStyle: { color: '#f2f6ff' },
          emphasis: { focus: 'series' }
        }
      ]

  return {
    color: securityPalette,
    tooltip: { trigger: 'axis' },
    legend:
      normalizedSeries.length > 1
        ? { top: 2, icon: 'roundRect', itemWidth: 22, itemHeight: 10, textStyle: { color: '#475569', fontSize: 12 } }
        : undefined,
    grid: { left: 36, right: 18, top: normalizedSeries.length > 1 ? 50 : 26, bottom: 34 },
    xAxis: { type: 'category', data: categories, ...axisStyle },
    yAxis: { type: 'value', ...axisStyle },
    series: normalizedSeries
  }
}

export const buildLineOption = ({
  categories = [],
  seriesName = '趋势',
  values = [],
  color = securityPalette[1],
  series = []
} = {}) => {
  // 向后兼容：同时支持单折线与多折线配置。
  const hasMultiSeries = Array.isArray(series) && series.length > 0
  const normalizedSeries = hasMultiSeries
    ? series.map((item, index) => ({
        name: item?.name || `系列${index + 1}`,
        type: 'line',
        smooth: item?.smooth ?? true,
        symbol: item?.symbol || 'circle',
        symbolSize: item?.symbolSize ?? 7,
        data: item?.values ?? [],
        lineStyle: { color: resolveSeriesColor(item, index) },
        itemStyle: { color: resolveSeriesColor(item, index) },
        areaStyle: resolveLineArea(item) ? { opacity: item?.opacity ?? 0.12 } : undefined
      }))
    : [
        {
          name: seriesName,
          type: 'line',
          smooth: true,
          symbol: 'circle',
          symbolSize: 7,
          data: values,
          lineStyle: { color },
          itemStyle: { color },
          areaStyle: { opacity: 0.12 }
        }
      ]

  return {
    color: securityPalette,
    tooltip: { trigger: 'axis' },
    legend: normalizedSeries.length > 1 ? { top: 0, textStyle: { color: '#475569', fontSize: 12 } } : undefined,
    grid: { left: 36, right: 18, top: normalizedSeries.length > 1 ? 44 : 26, bottom: 34 },
    xAxis: { type: 'category', data: categories, ...axisStyle },
    yAxis: { type: 'value', ...axisStyle },
    series: normalizedSeries
  }
}

export const buildPieOption = ({
  seriesName = '占比',
  data = []
} = {}) => ({
  color: securityPalette,
  tooltip: { trigger: 'item' },
  legend: { bottom: 0, textStyle: { color: '#475569', fontSize: 12 } },
  series: [
    {
      name: seriesName,
      type: 'pie',
      radius: ['45%', '72%'],
      center: ['50%', '45%'],
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { color: '#334155', formatter: '{b}: {d}%' },
      data
    }
  ]
})

export const buildMixOption = ({
  categories = [],
  leftName = '事件数',
  rightName = '处置率',
  barSeries = [],
  lineSeries = []
} = {}) => {
  // 混合图面向企业常见 KPI：左轴展示数量，右轴展示比率。
  const normalizedBars = (Array.isArray(barSeries) ? barSeries : []).map((item, index) => ({
    name: item?.name || `柱系列${index + 1}`,
    type: 'bar',
    yAxisIndex: 0,
    barMaxWidth: item?.barMaxWidth ?? 20,
    data: item?.values ?? [],
    itemStyle: buildBarSeriesStyle(resolveSeriesColor(item, index))
  }))
  const normalizedLines = (Array.isArray(lineSeries) ? lineSeries : []).map((item, index) => ({
    name: item?.name || `线系列${index + 1}`,
    type: 'line',
    yAxisIndex: 1,
    smooth: item?.smooth ?? true,
    symbol: item?.symbol || 'circle',
    symbolSize: item?.symbolSize ?? 7,
    data: item?.values ?? [],
    lineStyle: { color: resolveSeriesColor(item, index + 2), width: 2 },
    itemStyle: { color: resolveSeriesColor(item, index + 2) }
  }))
  return {
    color: securityPalette,
    tooltip: { trigger: 'axis' },
    legend: { top: 0, textStyle: { color: '#475569', fontSize: 12 } },
    grid: { left: 40, right: 40, top: 44, bottom: 34 },
    xAxis: { type: 'category', data: categories, ...axisStyle },
    yAxis: [
      { type: 'value', name: leftName, ...axisStyle },
      { type: 'value', name: rightName, ...axisStyle, axisLabel: { color: '#64748b', formatter: '{value}%' } }
    ],
    series: [...normalizedBars, ...normalizedLines]
  }
}

export const securityChartPalette = securityPalette

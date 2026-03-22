const toArray = (value) => (Array.isArray(value) ? value : [])

const normalizeDimensionOptions = (value) =>
  toArray(value)
    .map((item) => ({
      label: item?.label ?? String(item?.value ?? ''),
      value: item?.value
    }))
    .filter((item) => item.value !== undefined && item.value !== null)

const normalizeChartNode = (node = {}) => {
  const data = node?.data ?? {}
  return {
    data: typeof data === 'object' && data !== null ? data : {},
    dimensionOptions: normalizeDimensionOptions(node?.dimensionOptions),
    dimensionDataMap:
      typeof node?.dimensionDataMap === 'object' && node?.dimensionDataMap !== null
        ? node.dimensionDataMap
        : {},
    dimensionValue: node?.dimensionValue ?? null
  }
}

export const normalizeDashboardPayload = (payload = {}) => {
  // 将后端返回归一化为前端可预期的数据契约。
  // 即使上游字段缺失一部分，页面层代码也能保持稳定。
  const summary = toArray(payload?.summary).map((item, index) => ({
    key: item?.key ?? `summary_${index}`,
    title: item?.title ?? '-',
    value: item?.value ?? '-',
    sub: item?.sub ?? ''
  }))

  const charts = {
    bar: normalizeChartNode(payload?.charts?.bar),
    line: normalizeChartNode(payload?.charts?.line),
    pie: normalizeChartNode(payload?.charts?.pie),
    mix: normalizeChartNode(payload?.charts?.mix)
  }

  return {
    summary,
    charts,
    meta: typeof payload?.meta === 'object' && payload.meta !== null ? payload.meta : {}
  }
}

export const getPanelBinding = (normalizedPayload, chartType) => {
  // 便捷绑定助手：可直接对接 <lab-security-chart-panel>。
  const chart = normalizedPayload?.charts?.[chartType] ?? {}
  return {
    data: chart?.data ?? {},
    dimensionOptions: chart?.dimensionOptions ?? [],
    dimensionDataMap: chart?.dimensionDataMap ?? {},
    dimensionValue: chart?.dimensionValue ?? null
  }
}

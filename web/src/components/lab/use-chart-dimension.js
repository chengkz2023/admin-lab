import { computed, ref, watch } from 'vue'

export const useChartDimension = (props, emit) => {
  // 非受控模式下的兜底状态；受控模式始终以 props.dimensionValue 为准。
  const innerDimension = ref(null)
  const hasControlledDimension = computed(() => props.dimensionValue !== null && props.dimensionValue !== undefined)

  const normalizeOptionValues = (options) =>
    (Array.isArray(options) ? options : [])
      .map((item) => item?.value)
      .filter((item) => item !== undefined && item !== null)

  watch(
    () => props.dimensionOptions,
    (options) => {
      if (hasControlledDimension.value) {
        return
      }
      // 选项列表动态变化（如接口刷新）时，保证内部维度值始终有效。
      const values = normalizeOptionValues(options)
      if (!values.includes(innerDimension.value)) {
        innerDimension.value = values[0] ?? null
      }
    },
    { immediate: true, deep: true }
  )

  const currentDimension = computed(() => (hasControlledDimension.value ? props.dimensionValue : innerDimension.value))

  const panelDimension = computed({
    get: () => currentDimension.value,
    set: (value) => {
      if (!hasControlledDimension.value) {
        innerDimension.value = value
      }
      // 无论受控/非受控都触发事件，便于父组件统一感知维度变化。
      emit('update:dimensionValue', value)
      emit('dimension-change', value)
    }
  })

  const currentData = computed(() => {
    // 开启维度模式时优先取维度映射数据，否则回退到普通 data。
    if (props.dimensionOptions.length && currentDimension.value !== null && currentDimension.value !== undefined) {
      return props.dimensionDataMap?.[currentDimension.value] ?? props.data
    }
    return props.data
  })

  return {
    panelDimension,
    currentDimension,
    currentData
  }
}

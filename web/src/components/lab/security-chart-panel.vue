<template>
  <el-card class="security-chart-panel" shadow="never" v-loading="loading">
    <template #header>
      <div class="panel-header">
        <div>
          <div class="panel-title">{{ title }}</div>
          <div v-if="subtitle" class="panel-subtitle">{{ subtitle }}</div>
        </div>
        <div class="panel-actions">
          <el-radio-group
            v-if="dimensionOptions.length"
            v-model="panelDimension"
            size="small"
          >
            <el-radio-button
              v-for="item in dimensionOptions"
              :key="item.value"
              :value="item.value"
            >
              {{ item.label }}
            </el-radio-button>
          </el-radio-group>
          <slot name="extra" />
        </div>
      </div>
    </template>
    <v-chart class="chart" :option="option" autoresize :style="{ height }" @click="handleChartClick" />
  </el-card>
</template>

<script setup>
  import { computed } from 'vue'
  import VChart from 'vue-echarts'
  import { use } from 'echarts/core'
  import { CanvasRenderer } from 'echarts/renderers'
  import { BarChart, LineChart, PieChart } from 'echarts/charts'
  import {
    GridComponent,
    TooltipComponent,
    LegendComponent,
    TitleComponent
  } from 'echarts/components'
  import { buildBarOption, buildLineOption, buildPieOption, buildMixOption } from '@/components/lab/security-chart-options'
  import { useChartDimension } from '@/components/lab/use-chart-dimension'

  // 仅注册当前用到的渲染器/图表/组件，避免运行时报 “Renderer undefined”。
  use([CanvasRenderer, BarChart, LineChart, PieChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent])

  defineOptions({
    name: 'LabSecurityChartPanel'
  })

  const emit = defineEmits(['update:dimensionValue', 'dimension-change', 'chart-click'])

  const props = defineProps({
    title: {
      type: String,
      default: ''
    },
    subtitle: {
      type: String,
      default: ''
    },
    type: {
      type: String,
      default: 'line'
    },
    data: {
      type: Object,
      default: () => ({})
    },
    loading: {
      type: Boolean,
      default: false
    },
    dimensionOptions: {
      type: Array,
      default: () => []
    },
    dimensionValue: {
      type: [String, Number],
      default: null
    },
    dimensionDataMap: {
      type: Object,
      default: () => ({})
    },
    height: {
      type: String,
      default: '300px'
    }
  })

  const { panelDimension, currentDimension, currentData } = useChartDimension(props, emit)

  const option = computed(() => {
    // 将图表类型分发集中在这里，父组件只需要传入 `type` 即可。
    if (props.type === 'bar') {
      return buildBarOption(currentData.value)
    }
    if (props.type === 'pie') {
      return buildPieOption(currentData.value)
    }
    if (props.type === 'mix') {
      return buildMixOption(currentData.value)
    }
    return buildLineOption(currentData.value)
  })

  const handleChartClick = (params) => {
    // 对外暴露统一结构的点击事件数据，方便做下钻联动。
    emit('chart-click', {
      params,
      dimensionValue: currentDimension.value,
      data: currentData.value
    })
  }
</script>

<style scoped>
  .security-chart-panel {
    border: 1px solid #dbe7ff;
    border-radius: 14px;
    background: linear-gradient(180deg, #ffffff 0%, #f8fbff 100%);
  }

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
  }

  .panel-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .panel-title {
    font-size: 16px;
    font-weight: 600;
    color: #0f172a;
  }

  .panel-subtitle {
    margin-top: 4px;
    font-size: 12px;
    color: #64748b;
  }

  .chart {
    width: 100%;
  }
</style>

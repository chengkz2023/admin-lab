<template>
  <UnitCell
    :value="value"
    :column="enhancedColumn"
  />
</template>

<script setup>
import { computed } from 'vue';
import UnitCell from './UnitCell.vue';

/**
 * @file FlowRateCell.vue
 * @description 流量速率单元格渲染器，用于格式化和显示流量速率数据（如 B/s, KB/s, MB/s）。
 *              它基于 UnitCell 实现，并预设了流量相关的单位和换算基数。
 *
 * @usage
 * {
 *   prop: 'downloadSpeed',
 *   label: '下载速度',
 *   type: 'Bps',
 *
 *   // [可选]
 *   decimals: 1, // 精度，默认为 1
 * }
 */


const props = defineProps({
  value: { type: Number, default: null },
  column: { type: Object, default: () => ({}) },
});

const throughputPresetConfig = {
  unitBase: 1000,
  units: ['B/s', 'KB/s', 'MB/s', 'GB/s'],
};

const enhancedColumn = computed(() => {
  return {
    ...props.column,
    config: {
      ...throughputPresetConfig,
      decimals: props.column?.decimals || 1,
    },
  };
});
</script>

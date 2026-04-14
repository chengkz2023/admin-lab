<template>
  <UnitCell :value="value" :column="enhancedColumn" />
</template>

<script setup>
import { computed } from "vue";
import UnitCell from "./UnitCell.vue";

/**
 * @file PktRateCell.vue
 * @description 包速率单元格渲染器，用于格式化和显示网络包速率数据（如 pps, kpps, Mpps）。
 *              它基于 UnitCell 实现，并预设了包速率相关的单位和换算基数。
 *
 * @usage
 * {
 *   prop: 'packetRate',
 *   label: '包速率',
 *   type: 'pps',
 *
 *   // [可选]
 *   decimals: 0, // 精度，默认为 0 (不显示小数)
 * }
 */

const props = defineProps({
  value: { type: Number, default: null },
  column: { type: Object, default: () => ({}) },
});

const packetPresetConfig = {
  unitBase: 1000,
  units: ["pps", "kpps", "Mpps", "Gpps"],
};

const enhancedColumn = computed(() => {
  return {
    ...props.column,
    config: {
      ...packetPresetConfig,
      decimals: props.column.decimals || 0,
    },
  };
});
</script>

<template>
  <UnitCell :value="value" :column="enhancedColumn" />
</template>

<script setup>
import { computed } from "vue";
import UnitCell from "./UnitCell.vue";

/**
 * @file BitrateCell.vue
 * @description 比特率单元格渲染器，用于格式化和显示比特率数据（如 bps, Kbps, Mbps）。
 *              它基于 UnitCell 实现，并预设了比特率相关的单位和换算基数。
 *
 * @usage
 * {
 *   prop: 'networkBitrate',
 *   label: '网络比特率',
 *   type: 'bps',
 *
 *   // [可选]
 *   decimals: 1, // 精度，默认为 1
 * }
 */

const props = defineProps({
  value: { type: Number, default: null },
  column: { type: Object, default: () => ({}) },
});

// 定义 rate 类型的预设配置
const ratePresetConfig = {
  unitBase: 1000,
  units: ["bps", "Kbps", "Mbps", "Gbps", "Tbps"],
};

// 将预设配置与用户传入的 config 进行合并
const enhancedColumn = computed(() => {
  return {
    ...props.column,
    config: {
      ...ratePresetConfig,
      decimals: props.column?.decimals,
    },
  };
});
</script>

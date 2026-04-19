<template>
  <UnitCell :value="value" :column="enhancedColumn" />
</template>

<script setup>
import { computed } from "vue";
import UnitCell from "./UnitCell.vue";

/**
 * @file SizeCell.vue
 * @description 文件大小单元格渲染器，用于格式化和显示文件大小数据（如 Bytes, KB, MB, GB）。
 *              它基于 UnitCell 实现，并预设了文件大小相关的单位和换算基数。
 *
 * @usage
 * {
 *   prop: 'fileSize',
 *   label: '文件大小',
 *   type: 'size|byte',
 *
 *   // [可选]
 *   decimals: 1, // 精度，默认为 1
 * }
 */

const props = defineProps({
  value: { type: Number, default: null },
  column: { type: Object, default: () => ({}) },
});

// 定义 size 类型的预设配置
const sizePresetConfig = {
  unitBase: 1000,
  units: ["Bytes", "KB", "MB", "GB", "TB", "PB"],
};

// 将预设配置与用户传入的 config 进行合并
const enhancedColumn = computed(() => {
  return {
    ...props.column,
    config: {
      ...sizePresetConfig,
      decimals: props.column.decimals,
    },
  };
});
</script>

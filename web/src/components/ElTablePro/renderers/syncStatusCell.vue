<template>
  <el-tag :type="tagType">
    {{ statusLabel }}
  </el-tag>
</template>

<script setup>
import { computed } from "vue";

/**
 * @file syncStatusCell.vue
 * @description 同步状态单元格渲染器，用于展示同步状态（待同步、已同步、同步失败、待删除等）
 *
 * @usage
 * {
 *   prop: 'syncStatus',
 *   label: '同步状态',
 *   type: 'sync_status',
 *   width: 120,
 *
 *   // [可选] 提供一个 config 对象来自定义状态映射
 *   config: {
 *     // 示例: 覆盖默认的状态映射
 *     0: { label: '待同步', type: 'info' },
 *     1: { label: '已同步', type: 'success' },
 *     2: { label: '同步失败', type: 'danger' },
 *     3: { label: '待删除', type: 'warning' },
 *     // 也可以添加其他状态码
 *     4: { label: '同步中', type: 'primary' },
 *   }
 * }
 *
 * // 示例数据:
 * // { syncStatus: 0 } -> 显示灰色"待同步"标签
 * // { syncStatus: 1 } -> 显示绿色"已同步"标签
 * // { syncStatus: 2 } -> 显示红色"同步失败"标签
 * // { syncStatus: 3 } -> 显示橙色"待删除"标签
 */

const props = defineProps({
  value: { type: [String, Number], default: null },
  column: { type: Object, default: () => ({}) },
});

// =============================================================================
//  1. 默认状态配置 (Default Config)
// =============================================================================
const defaultConfig = {
  0: { label: "待同步", type: "info" },
  1: { label: "已同步", type: "success" },
  2: { label: "同步失败", type: "danger" },
  3: { label: "待删除", type: "warning" },
  4: { label: "过期已解绑", type: "warning" },
  5: { label: "已停用", type: "info" },
  6: { label: "未生效", type: "info" },
};

// =============================================================================
//  2. 核心计算属性
// =============================================================================

const effectiveConfig = computed(() => ({
  ...defaultConfig,
  ...(props.column.config || {}),
}));

const statusInfo = computed(() => {
  const key = props.value;
  const config = effectiveConfig.value[key];

  if (config) {
    return {
      label: config.label,
      type: config.type,
    };
  }

  // 如果没有找到配置，返回默认值
  return {
    label: `未知(${key})`,
    type: "info",
  };
});

const statusLabel = computed(() => statusInfo.value.label);
const tagType = computed(() => statusInfo.value.type);
</script>

<template>
  <div
    v-if="resolvedStatus"
    class="flex items-center space-x-2"
  >
    <!-- 外环 -->
    <div
      class="flex h-3 w-3 items-center justify-center rounded-full"
      :class="stylePreset.ringClass"
    >
      <!-- Ping 动画 -->
      <span
        v-if="resolvedStatus.isPing"
        class="animate-ping absolute inline-flex h-3 w-3 rounded-full"
        :class="stylePreset.dotClass"
      />

      <!-- 内点 -->
      <span
        class="h-1.5 w-1.5 rounded-full"
        :class="stylePreset.dotClass"
      />
    </div>

    <!-- 文本 -->
    <span
      class="text-sm font-medium"
      :class="stylePreset.textClass"
    >
      {{ resolvedStatus.label }}
    </span>
  </div>

  <span v-else>{{ value }}</span>
</template>

<script setup>
import { computed } from 'vue';

/**
 * @file StatusCell.vue
 * @description 状态单元格渲染器，根据状态值显示带有不同颜色和动画的指示点及文本。
 *              支持预设状态（如 'up', 'down', 'pending'）和自定义状态映射。
 *
 * @usage
 * {
 *   prop: 'serviceStatus',
 *   label: '服务状态',
 *   type: 'status',
 *
 *   // [可选] 提供一个 config 对象来自定义状态映射和样式
 *   config: {
 *     // 示例: 覆盖默认的 'up' 状态样式或添加新的状态
 *     'up': { label: '运行中', type: 'success', isPing: true },
 *     'maintenance': { label: '维护中', type: 'info', isPing: false },
 *     'critical': { label: '严重', type: 'danger', isPing: true },
 *   }
 * }
 *
 * // 示例数据:
 * // { serviceStatus: 'up' }      -> 显示绿色“在线”状态，带 Ping 动画
 * // { serviceStatus: 'down' }    -> 显示红色“离线”状态
 * // { serviceStatus: 'pending' } -> 显示橙色“处理中”状态，带 Ping 动画
 */

const props = defineProps({
  value: { type: [String, Number, Boolean], default: null },
  column: { type: Object, default: () => ({}) },
});

// =============================================================================
//  1. 预设样式集 (Style Presets)
// =============================================================================
const stylePresets = {
  success: {
    ringClass: 'bg-green-100 dark:bg-green-500/30',
    dotClass:  'bg-green-500 dark:bg-green-400',
    textClass: 'text-green-700 dark:text-slate-200',
  },
  danger: {
    ringClass: 'bg-red-100 dark:bg-red-500/30',
    dotClass:  'bg-red-500 dark:bg-red-400',
    textClass: 'text-red-700 dark:text-slate-200',
  },
  warning: {
    // 统一使用 orange-500 色系，视觉更和谐
    ringClass: 'bg-orange-100 dark:bg-orange-500/30',
    dotClass:  'bg-orange-500 dark:bg-orange-400',
    textClass: 'text-orange-700 dark:text-slate-200',
  },
  primary: {
    ringClass: 'bg-blue-100 dark:bg-blue-500/30',
    dotClass:  'bg-blue-500 dark:bg-blue-400',
    textClass: 'text-blue-700 dark:text-slate-200',
  },
  info: {
    ringClass: 'bg-slate-100 dark:bg-slate-500/30',
    dotClass:  'bg-slate-500 dark:bg-slate-400',
    textClass: 'text-slate-700 dark:text-slate-400',
  },
};

// =============================================================================
//  2. 默认状态配置 (Default Config)
// =============================================================================
const defaultConfig = {
  'Completed': { label: 'Completed', type: 'success' },
  'up':        { label: '在线', type: 'success', isPing: true },
  'active':    { label: '活跃', type: 'success', isPing: true },
  'connected': { label: '已连接', type: 'success', isPing: true },
  'ok':        { label: '正常', type: 'success' },
  'pending':   { label: '处理中', type: 'warning', isPing: true },
  'warning':   { label: '告警', type: 'warning' },
  'Error':     { label: 'Error', type: 'danger' },
  'down':      { label: '离线', type: 'danger' },
  'error':     { label: '错误', type: 'danger' },
  'inactive':  { label: '不活跃', type: 'info' },
  'disabled':  { label: '已禁用', type: 'info' },
};

// =============================================================================
//  3. 核心计算属性
// =============================================================================

const effectiveConfig = computed(() => ({
  ...defaultConfig,
  ...(props.column.config || {}),
}));

const resolvedStatus = computed(() => {
  const key = String(props.value);
  const statusConfig = effectiveConfig.value[key] || effectiveConfig.value[key.toLowerCase()];

  if (!statusConfig) {
    return { label: key, type: 'info', isPing: false };
  }
  return {
    label: statusConfig.label || key,
    type: statusConfig.type,
    isPing: statusConfig.isPing || false,
  };
});

const stylePreset = computed(() => {
  const type = resolvedStatus.value.type || 'info';
  return stylePresets[type] || stylePresets.info;
});
</script>

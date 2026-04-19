<template>
  <div
      v-if="resolvedData.label"
      class="tag-badge inline-flex h-6 items-center rounded-md px-2
           text-xs text-center font-medium transition-colors duration-200"
      :class="[resolvedStyle.bgClass, resolvedStyle.textClass, resolvedStyle.hoverBgClass]"
  >
    {{ resolvedData.label }}
  </div>
  <span v-else>-</span>
</template>

<script setup>
import { computed } from 'vue';

/**
 * @file optionsCell.vue
 * @description options 标签渲染器。
 *              接收一个 options 数组，根据当前 value 查找对应的 label 并显示为标签。
 *
 * @usage
 * {
 *   prop: 'status',
 *   label: '状态',
 *   type: 'options',
 *   options: [
 *     { label: '启用', value: 1, type: 'success' },
 *     { label: '禁用', value: 0, type: 'danger' },
 *     { label: '待审核', value: 2 }
 *   ]
 * }
 */

const props = defineProps({
  value: { type: [String, Number, Boolean], default: null },
  column: { type: Object, default: () => ({}) },
});

const tagStyles = {
  primary: { bgClass: 'bg-blue-100 dark:bg-blue-800', textClass: 'text-blue-800 dark:text-blue-100', hoverBgClass: 'hover:bg-blue-200 dark:hover:bg-blue-700' },
  success: { bgClass: 'bg-green-100 dark:bg-green-800', textClass: 'text-green-800 dark:text-green-100', hoverBgClass: 'hover:bg-green-200 dark:hover:bg-green-700' },
  warning: { bgClass: 'bg-orange-100 dark:bg-orange-800', textClass: 'text-orange-800 dark:text-orange-100', hoverBgClass: 'hover:bg-orange-200 dark:hover:bg-orange-700' },
  danger:  { bgClass: 'bg-red-100 dark:bg-red-800', textClass: 'text-red-800 dark:text-red-100', hoverBgClass: 'hover:bg-red-200 dark:hover:bg-red-700' },
  info:    { bgClass: 'bg-slate-100 dark:bg-slate-700', textClass: 'text-slate-800 dark:text-slate-200', hoverBgClass: 'hover:bg-slate-200 dark:hover:bg-slate-600' },
};

const resolvedData = computed(() => {
  const options = props.column.options || [];
  const found = options.find(item => item.value == props.value);

  if (found) {
    return {
      label: found.label,
      type: found.type || 'info'
    };
  }

  return {
    label: String(props.value ?? ''),
    type: 'info'
  };
});

const resolvedStyle = computed(() => {
  return tagStyles[resolvedData.value.type] || tagStyles.info;
});
</script>

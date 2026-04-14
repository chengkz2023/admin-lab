<template>
  <div
    v-if="resolvedTag.label"
    class="tag-badge inline-flex h-6 items-center rounded-md px-2
           text-xs font-medium transition-colors duration-200"
    :class="[resolvedStyle.bgClass, resolvedStyle.textClass, resolvedStyle.hoverBgClass]"
  >
    {{ resolvedTag.label }}
  </div>
</template>

<script setup>
import { computed } from 'vue';

/**
 * @file TagCell.vue
 * @description 标签单元格渲染器，根据标签值显示带有不同颜色和文本的徽章。
 *              适用于展示分类、状态、类型等信息。
 *
 * @usage
  {
    prop: 'category',
    label: '分类',
    type: 'tag',

    // [可选] 提供一个 config 对象来自定义标签映射和样式
    config: {
      'frontend': { label: '前端', type: 'primary' },
      'backend': { label: '后端', type: 'success' },
      'database': { label: '数据库', type: 'info' },
      'urgent': { label: '紧急', type: 'danger' },
    }
  }
 *
 * // 示例数据:
 * // { category: 'frontend' } -> 显示蓝色“前端”标签
 * // { category: 'urgent' }   -> 显示红色“紧急”标签
 */

const props = defineProps({
  value: { type: [String, Number, Boolean], default: null },
  column: { type: Object, default: () => ({}) },
});

// --- 内部样式映射表 ---
// 将 'type' 映射到具体的 Tailwind CSS 类名。这是组件的核心。
const tagStyles = {
  primary: { bgClass: 'bg-blue-100 dark:bg-blue-800', textClass: 'text-blue-800 dark:text-blue-100', hoverBgClass: 'hover:bg-blue-200 dark:hover:bg-blue-700' },
  success: { bgClass: 'bg-green-100 dark:bg-green-800', textClass: 'text-green-800 dark:text-green-100', hoverBgClass: 'hover:bg-green-200 dark:hover:bg-green-700' },
  warning: { bgClass: 'bg-orange-100 dark:bg-orange-800', textClass: 'text-orange-800 dark:text-orange-100', hoverBgClass: 'hover:bg-orange-200 dark:hover:bg-orange-700' },
  danger:  { bgClass: 'bg-red-100 dark:bg-red-800', textClass: 'text-red-800 dark:text-red-100', hoverBgClass: 'hover:bg-red-200 dark:hover:bg-red-700' },
  info:    { bgClass: 'bg-slate-100 dark:bg-slate-700', textClass: 'text-slate-800 dark:text-slate-200', hoverBgClass: 'hover:bg-slate-200 dark:hover:bg-slate-600' },
};

// --- 计算属性 1: 解析标签的 label 和 type ---
// 这部分逻辑与您提供的原版完全一致，负责从 config 中提取信息。
const resolvedTag = computed(() => {
  const config = props.column.config;
  const key = String(props.value); // 确保 value 是字符串以便作为 key

  if (config && config[key]) {
    return {
      label: config[key].label || key,
      type: config[key].type || 'primary', // 默认为 'primary'
    };
  }

  // 如果没有配置或找不到匹配项，直接显示 value 值
  return {
    label: key,
    type: 'info', // 对于未配置的值，使用中性的 'info' 样式
  };
});


// --- 计算属性 2: 根据解析出的 type 应用样式 ---
// 这是新的核心逻辑，它将 type 转换为最终的 CSS 类。
const resolvedStyle = computed(() => {
  // 从 tagStyles 映射表中查找对应 type 的样式
  // 如果传入的 type 无效 (如 'custom-color')，则优雅地降级到 'info' 样式
  return tagStyles[resolvedTag.value.type] || tagStyles.info;
});
</script>

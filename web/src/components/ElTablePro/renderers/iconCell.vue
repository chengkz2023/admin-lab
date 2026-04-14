<template>
  <div class="flex items-center space-x-2">
    <el-icon
      :size="16"
      :class="iconInfo.class"
    >
      <component :is="iconInfo.component" />
    </el-icon>
    <span class="dark:text-slate-300">{{ iconInfo.text }}</span>
  </div>
</template>
<script setup>
import { computed } from 'vue';
import { Cpu, Coin, DataLine, MessageBox, Document } from '@element-plus/icons-vue';

/**
 * @file IconCell.vue
 * @description 图标单元格渲染器，根据传入的 `value` 显示预设的图标和文本。
 *              适用于需要将特定字符串值映射为带有图标的视觉表示的场景。
 *
 * @usage
 * {
 *   prop: 'serviceType',
 *   label: '服务类型',
 *   type: 'icon',
 *   // value 对应 map 中的 key
 * }
 *
 * // 示例数据:
 * // { serviceType: 'database' } -> 显示数据库图标和“数据库”文本
 * // { serviceType: 'app' }      -> 显示应用服务图标和“应用服务”文本
 */

const props = defineProps({
  value: { type: String, default: null },
  column: { type: Object, default: () => ({}) },
});
const iconInfo = computed(() => {
  const map = {
    'database': { component: Coin, text: '数据库', class: 'text-blue-600 dark:text-blue-400' },
    'redis': { component: MessageBox, text: '缓存', class: 'text-red-500 dark:text-red-400' },
    'app': { component: Cpu, text: '应用服务', class: 'text-green-600 dark:text-green-400' },
    'log': { component: Document, text: '日志', class: 'text-gray-500 dark:text-gray-400' },
  };
  return map[props.value] || { component: DataLine, text: props.value, class: '' };
});
</script>

<template>
  <div
    class="flex items-center space-x-1"
    :class="mainContainerClass"
  >
    <div :class="barContainerClass">
      <span
        v-for="i in totalBlocks"
        :key="i"
        class="h-1.5 w-3 rounded-sm"
        :class="getBlockClass(i)"
      />
    </div>
    <span
      class="font-mono text-xs font-medium text-right"
      :class="styles.textColor"
    >
      {{ formattedValue }}
    </span>
  </div>
</template>

<script setup>
import { computed } from 'vue';

/**
 * @file SignalCell.vue
 * @description 一个通用的“信号强度”渲染器，支持水平/垂直排列和反向指标。
 *
 * @usage
 * {
 *   prop: 'qualityScore',
 *   label: '链路质量',
 *   type: 'signal',
 *
 *   // [可选] 提供一个 config 对象来自定义行为
 *   config: {
 *     valueType: 'decimal',   // 定义传入值的类型: 'percent' (0-100) 或 'decimal' (0-1)，默认为 'percent'
 *     decimals: 1,      // 显示的百分比小数位数, 默认为 1
 *     inverse: true,        // 是否为反向指标, 默认为 false (值越小信号越强)
 *     direction: 'vertical', // 信号条方向: 'horizontal' (默认) 或 'vertical'
 *     thresholds: {         // 自定义颜色阈值 (基于 0-100 的信号强度)
 *       warning: 50,    // 信号强度低于 50% 为警告色
 *       danger: 25,     // 信号强度低于 25% 为危险色
 *     }
 *   }
 * }
 */

const props = defineProps({
  value: { type: Number, default: null },
  column: { type: Object, default: () => ({}) },
});

// --- 默认配置 ---
const defaultThresholds = {
  warning: 50,
  danger: 25,
};

// --- 从列配置中读取或使用默认值 ---
const config = computed(() => props.column.config || {});
const valueType = computed(() => config.value.valueType || 'decimal');
const decimals = computed(() => config.value.decimals ?? 1);
const thresholds = computed(() => config.value.thresholds || defaultThresholds);
const isInverse = computed(() => config.value.inverse === true);
const direction = computed(() => config.value.direction || 'horizontal');

// --- 核心计算属性 ---
const totalBlocks = computed(() => {
  return direction.value === 'vertical' ? 3 : 4;
});

/**
 * 【核心优化】创建内部百分比值 (internalPercentage)
 * 无论传入的是 0-1 还是 0-100，都将其统一转换为 0-100 的格式。
 */
const internalPercentage = computed(() => {
  if (props.value === null || !isFinite(props.value)) {
    return 0;
  }
  if (valueType.value === 'decimal') {
    return Math.max(0, Math.min(100, props.value * 100));
  }
  return Math.max(0, Math.min(100, props.value));
});

/**
 * 信号强度现在基于 internalPercentage 计算
 */
const signalStrength = computed(() => {
  return isInverse.value ? 100 - internalPercentage.value : internalPercentage.value;
});

const styles = computed(() => {
  const strength = signalStrength.value;
  const th = thresholds.value;

  let textColor = 'text-green-600 dark:text-green-400';
  if (strength < th.danger) textColor = 'text-red-600 dark:text-red-400';
  else if (strength < th.warning) textColor = 'text-orange-600 dark:text-orange-400';

  let blockColor = 'bg-green-500 dark:bg-green-400';
  if (strength < th.danger) blockColor = 'bg-red-500 dark:bg-red-400'; // 修正：危险时用红色
  else if (strength < th.warning) blockColor = 'bg-orange-500 dark:bg-orange-400'; // 修正：警告时用橙色

  return { textColor, blockColor };
});

const formattedValue = computed(() => {
  if (props.value === null || !isFinite(props.value)) {
    return '—';
  }
  return `${internalPercentage.value.toFixed(decimals.value)}%`;
});

const mainContainerClass = computed(() => {
  if (direction.value === 'vertical') {
    return 'space-x-1.5';
  }
  return 'space-x-2';
});

const barContainerClass = computed(() => {
  if (direction.value === 'vertical') {
    return 'flex flex-col-reverse space-y-0.5';
  }
  // 【优化】将信号条容器与右侧文本分开，让它占据剩余空间
  return 'flex flex-1 space-x-0.5';
});

// --- 方法 ---
function getBlockClass(index) {
  // 当信号为0时，不应该有任何点亮的块
  if (signalStrength.value <= 0) {
    return 'bg-gray-200 dark:bg-gray-700';
  }

  const activeBlocks = Math.ceil((signalStrength.value / 100) * totalBlocks.value);

  if (index <= activeBlocks) {
    return styles.value.blockColor;
  }

  return 'bg-gray-200 dark:bg-gray-700';
}
</script>

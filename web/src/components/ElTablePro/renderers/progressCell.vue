<template>
  <div
    class="flex items-center w-full"
    :class="containerClass"
  >
    <el-progress
      :class="progressClass"
      :percentage="internalValue"
      :show-text="showText"
      :text-inside="textInside"
      :stroke-width="strokeWidth"
      :status="progressStatus"
      :type="progressType"
      :width="circleWidth"
      :format="formatProgressText"
      :color="circleProgressColor"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { ElProgress } from 'element-plus';

/**
 * @file ProgressCell.vue
 * @description 使用 el-progress 组件来可视化百分比数据，支持线性和圆形进度条。
 *
 * 注意：默认显示文字，线性进度条文字显示在进度条内部，圆形进度条文字显示在圆形中心。
 *
 * @usage
 * {
 *   prop: 'usage',
 *   label: '使用率',
 *   type: 'progress',
 *
 *   // [可选] 提供一个 config 对象来自定义行为
 *   config: {
 *     // 定义传入值的类型: 'percent' (0-100) 或 'decimal' (0-1)
 *     valueType: 'decimal', // 默认为 'percent'
 *
 *     // 进度条类型: 'line' 或 'circle'
 *     type: 'circle', // 默认为 'line'
 *
 *     // 圆形进度条的尺寸
 *     circleWidth: 80, // 默认为 60，仅对 type='circle' 有效
 *
 *     decimals: 1, // 显示的小数位数, 默认为 0
 *     strokeWidth: 8, // 进度条粗细, 默认为 6
 *     thresholds: {
 *       warning: 80, // 80% 及以上为警告色
 *       danger: 95, // 95% 及以上为危险色
 *     }
 *   }
 * }
 */

const props = defineProps({
  value: {
    type: Number,
    default: null
  },
  column: {
    type: Object,
    default: () => ({})
  },
});

// --- 默认配置 ---
const defaultThresholds = {
  warning: 70,
  danger: 90,
};

// --- 从列配置中读取或使用默认值 ---
const config = computed(() => props.column.config || {});
const valueType = computed(() => config.value.valueType || 'percent');
const progressType = computed(() => config.value.type || 'line');
const circleWidth = computed(() => config.value.circleWidth || 80);
const strokeWidth = computed(() => config.value.strokeWidth || progressType.value === 'circle'? 12 :  16);

const decimals = computed(() => config.value.decimals ?? 0);
const thresholds = computed(() => config.value.thresholds || defaultThresholds);

// --- 核心计算属性 ---
/**
 * 【核心优化】创建内部值 (internalValue)
 * 无论传入的是 0-1 还是 0-100，都将其统一转换为 0-100 的格式。
 * 组件的所有后续逻辑都基于此值，实现了与输入格式的解耦。
 */
const internalValue = computed(() => {
  if (props.value === null || !isFinite(props.value)) {
    return 0;
  }

  if (valueType.value === 'decimal') {
    // 如果是小数格式 (0-1)，则乘以 100
    return Math.max(0, Math.min(100, props.value * 100));
  }

  // 如果是百分比格式 (0-100)，则直接使用，但要确保在 0-100 范围内
  return Math.max(0, Math.min(100, props.value));
});

/**
 * 动态样式现在基于 internalValue 进行判断。
 */
const dynamicStyle = computed(() => {
  const val = internalValue.value;
  const th = thresholds.value;

  if (val >= th.danger) {
    return {
      progressStatus: 'exception',
      textColor: 'text-red-600 dark:text-red-400'
    };
  }

  if (val >= th.warning) {
    return {
      progressStatus: 'warning',
      textColor: 'text-orange-600 dark:text-orange-400'
    };
  }

  // 默认使用 success 状态，即使值为0
  return {
    progressStatus: 'success',
    textColor: 'text-green-600 dark:text-green-400'
  };
});

/**
 * 进度条状态
 * 对于圆形进度条，如果想显示百分比文本而不是状态图标，则不设置 status
 * 对于线性进度条，可以设置 status 来显示不同的颜色
 */
const progressStatus = computed(() => {
  if (progressType.value === 'circle') {
    // 圆形进度条不设置 status，这样就会显示百分比文本而不是图标
    return undefined;
  }
  // 线性进度条可以使用 status 来显示不同颜色
  return dynamicStyle.value.progressStatus;
});

/**
 * 圆形进度条的颜色
 * 由于圆形进度条不使用 status 属性，我们通过 color 属性来控制颜色
 */
const circleProgressColor = computed(() => {
  if (progressType.value !== 'circle') {
    return undefined;
  }

  const val = internalValue.value;
  const th = thresholds.value;

  if (val >= th.danger) {
    return '#f56565'; // 红色
  }

  if (val >= th.warning) {
    return '#ed8936'; // 橙色
  }

  return '#48bb78'; // 绿色
});

/**
 * 格式化进度条内部显示文字的函数
 * 确保内部文字显示为百分比格式，覆盖默认的状态图标
 */
const formatProgressText = (percentage) => {
  if (!isFinite(percentage)) {
    return '—';
  }
  // 对于圆形进度条，始终返回百分比文本，不显示状态图标
  return `${percentage.toFixed(decimals.value)}%`;
};

/**
 * 格式化输出也基于 internalValue，确保显示正确。
 */
const formattedValue = computed(() => {
  // 当原始值为 null 时，显示占位符，更友好
  if (props.value === null || !isFinite(props.value)) {
    return '—';
  }

  return `${internalValue.value.toFixed(decimals.value)}%`;
});

/**
 * 是否显示文字
 * 默认显示文字，无论是线性还是圆形进度条
 */
const showText = computed(() => {
  return true;
});

/**
 * 是否在进度条内部显示文字
 * 线性进度条：在进度条内部显示文字
 * 圆形进度条：文字显示在圆形中心
 */
const textInside = computed(() => {
  return progressType.value === 'line';
});

/**
 * 容器样式类
 */
const containerClass = computed(() => {
  return 'justify-center';
});

/**
 * 进度条样式类
 */
const progressClass = computed(() => {
  if (progressType.value === 'circle') {
    return 'shrink-0'; // 圆形进度条固定尺寸，不伸缩
  }
  return 'flex-1'; // 线性进度条占满剩余空间
});
</script>

<style scoped>
/* 进度条背景色保持柔和 */
.el-progress :deep(.el-progress-bar__outer) {
  background-color: #f3f4f6; /* Tailwind gray-100 */
}

.dark .el-progress :deep(.el-progress-bar__outer) {
  background-color: #4b5563; /* Tailwind gray-700 */
}

/* 圆形进度条文字样式优化 */
.el-progress--circle :deep(.el-progress__text) {
  font-size: 12px;
  font-weight: 500;
  font-family: ui-monospace, 'SF Mono', monospace;
}

/* 线性进度条内部文字样式优化 */
.el-progress--line :deep(.el-progress-bar__innerText) {
  font-size: 11px;
  font-weight: 500;
  font-family: ui-monospace, 'SF Mono', monospace;
}
</style>

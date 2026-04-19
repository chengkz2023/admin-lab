<template>
  <!--
    使用 el-tooltip 作为根元素，包裹所有内容。
    :disabled="!isTruncated" 是关键，确保只在需要时显示 Tooltip。
  -->
  <el-tooltip
    :content="String(formattedValue)"
    placement="top"
    :disabled="!isTruncated"
  >
    <!--
      文本容器：
      - ref="textRef" 用于在 script 中获取 DOM 元素。
      - truncate 类应用了溢出截断样式。
      - w-full 确保 span 能撑满单元格，以便正确计算截断。
    -->
    <span
      ref="textRef"
      class="text-sm text-slate-700 dark:text-slate-300 truncate block"
    >
      {{ formattedValue }}
    </span>
  </el-tooltip>
</template>

<script setup>
import { ref, computed, onMounted, onUpdated } from 'vue';
import { ElTooltip } from 'element-plus';

/**
 * @file DefaultCell.vue
 * @description 默认的单元格渲染器，内置溢出截断和 Tooltip 功能，
 *              并支持通过 `format` 提供自定义格式化函数。
 *
 * @usage
 * {
 *   prop: 'longDescription',
 *   label: '详细描述',
 *   width: 200, // 必须设置宽度以触发截断
 * }
 * {
 *   prop: 'orderId',
 *   label: '订单ID',
 *   format: (value) => `ORD-${value}`
 * }
 */

const props = defineProps({
  value: { type: [String, Number, Boolean, Object, Array], default: null },
  row: { type: Object, default: () => ({}) },
  column: { type: Object, default: () => ({}) },
});

// --- Tooltip 状态检测 ---
const textRef = ref(null);
const isTruncated = ref(false);

const checkTruncation = () => {
  // 使用 requestAnimationFrame 确保在 DOM 更新后执行
  requestAnimationFrame(() => {
    const el = textRef.value;
    if (el) {
      // 比较元素的滚动宽度和可见宽度
      isTruncated.value = el.scrollWidth > el.offsetWidth;
    }
  });
};

// 在组件挂载和更新后都检查一次
onMounted(checkTruncation);
onUpdated(checkTruncation);


// --- 核心计算属性 ---
const formattedValue = computed(() => {
  const formatter = props.column?.formatter;

  if (typeof formatter === 'function') {
    return formatter(
        props.row,
        props.column,
        props.row?.[props.column?.prop],
        props.row?._index,
    );
  }

  if (props.value === null || props.value === undefined) {
    return '—';
  }

  return String(props.value);
});
</script>

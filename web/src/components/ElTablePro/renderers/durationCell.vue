<template>
  <div
    v-if="formatted"
    class="flex items-baseline space-x-1"
  >
    <!-- 数值容器 -->
    <span
      class="inline-block text-right text-sm font-semibold font-mono text-slate-700 dark:text-slate-300"
      :class="valueWidthClass"
    >
      {{ formatted.value }}
    </span>
    <!-- 单位容器 -->
    <span
      class="inline-block text-left text-xs font-mono text-slate-500 dark:text-slate-400"
      :class="unitWidthClass"
    >
      {{ formatted.unit }}
    </span>
  </div>
  <span
    v-else
    class="text-gray-400 dark:text-gray-500"
  >—</span>
</template>

<script setup>
import { computed } from 'vue';

/**
 * @file DurationCell.vue
 * @description 持续时间单元格渲染器，用于格式化和显示时间长度数据（如 ms, s, min, h）。
 *              它能根据数值大小自动选择合适的单位。
 *
 * @usage
 * {
 *   prop: 'responseTime',
 *   label: '响应时间',
 *   type: 'duration|delay',
 *
 *   // [可选] 提供一个 config 对象来自定义格式化行为
 *   config: {
 *     unit: 'ms',    // 输入值的单位，默认为 'ms' (毫秒)。可选 'ns', 's', 'm', 'h'
 *     decimals: 2,        // 显示的小数位数，默认为 2
 *     valueWidth: 'w-14', // 数值部分的 Tailwind CSS 宽度类，默认为 'w-14'
 *     unitWidth: 'w-10',  // 单位部分的 Tailwind CSS 宽度类，默认为 'w-10'
 *   }
 * }
 */

const props = defineProps({
  value: { type: Number, default: null },
  column: { type: Object, default: () => ({}) },
});

// --- 从列配置中读取或使用默认值 ---
const valueWidthClass = computed(() => props.column?.config?.valueWidth || 'w-14');
const unitWidthClass = computed(() => props.column?.config?.unitWidth || 'w-10');

// --- 核心格式化函数 ---
// 首先将输入值转换为毫秒，然后根据单位进行格式化
function formatDuration(rawValue, unitIn = 'ms', dec = 2) {
  if (typeof rawValue !== 'number' || !isFinite(rawValue)) return null;
  if (rawValue === 0) return { value: '0', unit: unitIn };

  // 根据输入单位将值转化为毫秒
  let valueInMs = rawValue;
  switch (unitIn) {
    case 'ns': valueInMs = rawValue / 1000000; break; // 纳秒 -> 毫秒
    case 's': valueInMs = rawValue * 1000; break; // 秒 -> 毫秒
    case 'm': valueInMs = rawValue * 60000; break; // 分钟 ->
    case 'h': valueInMs = rawValue * 3600000; break; // 小时 -> 毫秒
  }

  const ONE_MICROSECOND = 0.001;  // 微秒 （1微秒 = 0.001毫秒）
  const ONE_MILLISECOND = 1;  // 毫秒 （基准单位）
  const ONE_SECOND = 1000;  // 秒 （1秒 = 1000毫秒）
  const ONE_MINUTE = 60 * ONE_SECOND; // 分钟 （1分钟 = 60秒）
  const ONE_HOUR = 60 * ONE_MINUTE; // 小时 （1小时 = 60分钟）

  if (valueInMs < ONE_MICROSECOND) {
    return {value: (valueInMs * 1000000).toFixed(0), unit: 'ns'}; // 纳秒
  }
  if (valueInMs < ONE_MILLISECOND) {
    return {value: (valueInMs * 1000).toFixed(dec), unit: 'μs'}; // 微秒
  }
  if (valueInMs < ONE_SECOND) {
    return { value: valueInMs.toFixed(0), unit: 'ms' }; // 毫秒
  }
  if (valueInMs < ONE_MINUTE) {
    return { value: (valueInMs / ONE_SECOND).toFixed(dec), unit: 's' }; // 秒
  }
  if (valueInMs < ONE_HOUR) {
    return { value: (valueInMs / ONE_MINUTE).toFixed(dec), unit: 'min' }; // 分钟
  }
  return { value: (valueInMs / ONE_HOUR).toFixed(dec), unit: 'h' }; // 小时
}

const formatted = computed(() => {
  // 从 column.config 中读取配置，并提供默认值
  const config = props.column?.config || {}

  return formatDuration(
      props.value,
      config?.unit || 'ms',
      config?.decimals ?? 2
  );
});
</script>

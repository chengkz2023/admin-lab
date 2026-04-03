<template>
  <div v-if="formatted" class="flex items-baseline space-x-1">
    <span
      class="text-sm font-semibold font-mono text-slate-700 dark:text-slate-300"
    >
      {{ formatted.value }}
    </span>
    <span class="text-xs font-mono text-slate-500 dark:text-slate-400">
      {{ formatted.unit }}
    </span>
  </div>
  <span v-else class="text-gray-400 dark:text-gray-500">—</span>
</template>

<script setup>
import { computed } from "vue";

/**
 * @file UnitCell.vue
 * @description 一个通用的、用于格式化带单位的数值的渲染器。
 *              它能根据数值大小和配置的换算基数、单位数组自动选择合适的单位。
 *
 * @usage
 * {
 *   prop: 'traffic',
 *   label: '流量',
 *   type: 'unit', // 通用单位格式化类型
 *
 *   // [可选] 提供一个 config 对象来自定义格式化行为
 *   config: {
 *     decimals: 2,          // 小数位数，默认为 2
 *     unitBase: 1000,       // 换算基数，默认为 1000 (例如，1000 用于 SI 单位，1000 用于二进制单位)
 *     units: ['B', 'KB', 'MB', 'GB', 'TB'], // 单位数组，按从小到大顺序排列
 *     suffix: '/s',         // 单位后缀，例如流量的 '/s'
 *   }
 * }
 *
 * // 示例数据:
 * // { traffic: 12345678 } -> 根据配置显示为 12.35 MB 或 11.77 MiB
 */

const props = defineProps({
  value: { type: Number, default: null },
  column: { type: Object, default: () => ({}) },
});

/**
 * 核心格式化函数，将数值和单位分离返回
 * @param {number | null | undefined} rawValue - 原始数值
 * @param {object} options - 配置项
 * @returns {{value: string, unit: string} | null}
 */
function formatValueAndUnit(rawValue, options) {
  // 1. 输入验证
  if (typeof rawValue !== "number" || !isFinite(rawValue)) {
    return null;
  }

  // 2. 解构并设置默认配置
  const {
    units = ["B", "KB", "MB", "GB", "TB"],
    base = 1000,
    decimals = 2,
    suffix = "",
  } = options;

  const dm = Math.max(0, decimals); // 确保小数位不为负数

  // 3. 处理 0 的情况
  if (rawValue === 0) {
    return { value: (0).toFixed(dm), unit: `${units[0]}${suffix}` };
  }

  // 4. 计算单位索引
  // 防止对小于1的数取对数导致索引为负
  const index =
    rawValue < 1
      ? 0
      : Math.min(
          units.length - 1, // 防止索引越界
          Math.floor(Math.log(rawValue) / Math.log(base)),
        );

  // 5. 计算最终值和单位
  const finalValue = (rawValue / Math.pow(base, index)).toFixed(dm);
  const finalUnit = `${units[index] || ""}${suffix}`;

  return { value: finalValue, unit: finalUnit };
}

const formatted = computed(() => {
  // 从 column.config 中读取配置，并提供默认值
  const config = props.column?.config || {};

  return formatValueAndUnit(props.value, {
    // 为兼容旧版配置名，同时检查新版名称
    units: config.units || ["Bytes", "KB", "MB", "GB", "TB"],
    base: config.base || config.unitBase || 1000,
    decimals: config.decimals ?? 2, // 默认小数位改为2
    suffix: config.suffix || "",
  });
});
</script>

<template>
  <!--
    经典 '#' 风格 ID 展示:
    - 整体是一个简洁、统一的文本元素，视觉上干净利落。
    - 通过 '#' 符号和不同的颜色，巧妙地区分了“标识符”和“ID值”。
  -->
  <div class="id-cell inline-flex items-center">
    <span class="mr-0.5 text-sm font-medium text-slate-400 dark:text-slate-500">#</span>
    <span class="text-sm font-mono font-medium text-slate-700 dark:text-slate-300">
      {{ value }}
    </span>
  </div>
</template>

<script setup>

/**
 * @file IdCell.vue
 * @description ID 单元格渲染器，以统一的“# ID值”格式显示标识符。
 *              适用于展示各种系统 ID、订单号等。
 *
 * @usage
 * {
 *   prop: 'userId',
 *   label: '用户ID',
 *   type: 'id',
 * }
 *
 * // 示例数据:
 * // { userId: 12345 } -> 显示 # 12345
 * // { userId: 'ABC-987' } -> 显示 # ABC-987
 */

defineProps({
  value: { type: [String, Number], default: null },
  // 保留 column prop，以备不时之需。
  column: { type: Object, default: () => ({}) },
});
</script>

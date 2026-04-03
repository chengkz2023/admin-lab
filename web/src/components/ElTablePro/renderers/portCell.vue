<template>
  <div
    class="relative px-1.5 py-0.5 group inline-flex items-center space-x-5 rounded-md
           bg-slate-100 dark:bg-slate-700
           hover:bg-slate-200 dark:hover:bg-slate-600 transition-colors duration-200"
  >
    <svg
      class="absolute top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500 dark:text-slate-400"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 22 22"
      fill="none"
      aria-hidden="true"
    >
      <path
        stroke="currentColor"
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="m8 8-4 4 4 4m8 0 4-4-4-4m-2-3-4 14"
      />
    </svg>
    <span class="text-xs font-mono font-medium text-slate-800 dark:text-slate-200">
      {{ value }}
    </span>
  </div>
</template>

<script setup>

/**
 * @file PortCell.vue
 * @description 端口号单元格渲染器，以带有图标的徽章形式显示端口号。
 *              适用于展示网络端口、服务监听端口等。
 *
 * @usage
 * {
 *   prop: 'servicePort',
 *   label: '服务端口',
 *   type: 'port',
 * }
 *
 * // 示例数据:
 * // { servicePort: 80 }   -> 显示 </> 80
 * // { servicePort: 443 }  -> 显示 </> 443
 */

defineProps({
  value: { type: [ String, Number ], default: null },
  column: { type: Object, default: () => ({}) },
});
</script>

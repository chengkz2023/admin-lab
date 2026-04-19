<template>
  <div
      class="relative px-1.5 py-0.5 group inline-flex items-center space-x-5 rounded-md
           bg-slate-100 dark:bg-slate-700
           hover:bg-slate-200 dark:hover:bg-slate-600 transition-colors duration-200"
  >
    <span class="text-xs font-mono font-medium text-slate-800 dark:text-slate-200">
      {{ formattedValue }}
    </span>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUpdated } from 'vue';
import { ElTooltip } from 'element-plus';

const props = defineProps({
  value: { type: [String, Number, Boolean, Object, Array], default: null },
  row: { type: Object, default: () => ({}) },
  column: { type: Object, default: () => ({}) },
});


// --- 核心计算属性 ---
const formattedValue = computed(() => {
  const formatter = props.column?.formatter;

  if (typeof formatter === 'function') {
    return formatter(props.value)
  }

  if (props.value === null || props.value === undefined) {
    return '—';
  }

  return String(props.value);
});
</script>

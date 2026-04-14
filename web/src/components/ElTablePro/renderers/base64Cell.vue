<template>
  <DefaultCell
      :value="value"
      :row="row"
      :column="enhancedColumn"
  />
</template>

<script setup>
import { computed } from 'vue';
import TypeCell from "@/components/ElTablePro/renderers/typeCell.vue";
import DefaultCell from "@/components/ElTablePro/renderers/defaultCell.vue";


const props = defineProps({
  value: { type: [ String, Number, Boolean, Object, Array ], default: null },
  row: { type: Object, default: () => ({}) },
  column: { type: Object, default: () => ({}) },
});

const formatter = (row, col, cell) => {
  if (!cell) return ""

  try {
    return decodeURIComponent(escape(atob(cell)))
  } catch (e) {
    console.warn('Base64 decode failed', e)
    return cell
  }
}

// 将预设配置与用户传入的 config 进行合并
const enhancedColumn = computed(() => {
  return {
    ...props.column,
    formatter: formatter
  }
});
</script>

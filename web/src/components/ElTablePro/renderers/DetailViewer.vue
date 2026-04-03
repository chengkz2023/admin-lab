<!-- DetailViewer.vue -->
<template>
  <el-drawer
      v-model="visible"
      destroy-on-close
      append-to-body
      resizable
      size="45%"
      :with-header="false"
      :show-close="false"
  >
      <div class="p-2 flex justify-between">
        <el-text size="large">详情</el-text>
        <el-button link circle type="danger" @click="visible=false">
          <el-icon :size="22" class="el-icon--left"><CircleCloseFilled/></el-icon>
        </el-button>
      </div>

    <el-descriptions
        :column="1"
        class=""
        size="small"
        border
    >
      <el-descriptions-item
          v-for="(col, index) in flattenColumns"
          :key="col.prop || col.label || index"
          :label="col.label"
      >
        <!-- 复用 renderers 目录下的组件 -->
        <component
            :is="resolveCellComponent(col.type)"
            :value="row[col.prop]"
            :row="row"
            :column="col"
            class="flex items-center"
        />
      </el-descriptions-item>
    </el-descriptions>

    <template #footer>
      <div class="flex justify-end">
        <el-button @click="visible = false">关闭</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup>
import { computed } from 'vue';
import { resolveCellComponent, isBuiltInType } from '../renderer.js';
import { CircleCloseFilled } from "@element-plus/icons-vue";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  row: { type: Object, default: () => ({}) },
  columns: { type: Array, default: () => [] },
});

const emit = defineEmits(['update:modelValue']);

// 双向绑定
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
});

// 扁平化 Columns 逻辑 (保持不变)
const flattenColumns = computed(() => {
  const result = [];
  const traverse = (cols) => {
    cols.forEach((col) => {
      // 过滤非数据列
      if (
        isBuiltInType(col.type) ||
        col.type === 'action' ||
        col.type === 'switch'
      ) return;

      if (col.children && col.children.length > 0) {
        traverse(col.children);
      } else {
        result.push(col);
      }
    });
  };
  traverse(props.columns);
  return result;
});
</script>

<style scoped>
:deep(.el-drawer__header) {
  @apply p-0 mb-0
}
</style>

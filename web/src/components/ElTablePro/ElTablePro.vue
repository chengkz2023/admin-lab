<template>
  <div
    class="el-table-pro-wrapper dark:bg-slate-800"
    :class="{ 'is-empty': !Array.isArray(data) || data.length === 0 }"
  >
    <!-- 顶部工具栏插槽 -->
    <div v-if="$slots.toolbar" class="px-4 pb-2 pt-4">
      <slot name="toolbar" />
    </div>

    <el-table
      class="rounded-md"
      :data="data"
      v-bind="$attrs"
      header-cell-class-name="el-table-pro-header-cell--custom"
      cell-class-name="el-table-pro-cell--custom"
      size="small"
      :height="height"
    >
      <!--
        启动递归渲染：
        - 循环顶层列配置。
        - 将 ProTable 的 $slots 作为 prop 传递给 ColumnRenderer。
      -->
      <ColumnRenderer
        v-for="column in columns"
        :key="column.label"
        :column="column"
        :table-slots="$slots"
      />

      <!-- 2. 检查是否存在 'actions' 插槽，如果存在，则渲染一个固定的操作列 -->
      <el-table-column
        v-if="$slots.actions || showDetail"
        v-bind="actionsSlotConfig"
      >
        <template #default="{ row, $index }">
          <el-button
            v-if="showDetail"
            link
            title="详情"
            type="success"
            @click="openDetail(row)"
          >
            <el-icon :size="12" class="el-icon--left"><Document /></el-icon>
          </el-button>

          <!-- 将 'actions' 插槽的内容渲染到这里 -->
          <!-- 并将 row 和 $index 作为作用域参数传出 -->
          <slot name="actions" :row="row" :index="$index" />
        </template>
      </el-table-column>

      <!-- 空状态 -->
      <template #empty>
        <span class="text-xl">暂无数据</span>
      </template>
    </el-table>

    <!-- 分页器 -->
    <div v-if="showPagination" class="flex justify-end px-2">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        :page-sizes="pageSizes"
        :layout="layout"
        size="small"
        @update:current-page="handlePageChange('currentPage', $event)"
        @update:page-size="handlePageChange('pageSize', $event)"
      />
    </div>

    <detail-viewer
      v-model="detailVisible"
      :row="detailRow"
      :columns="columns"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import ColumnRenderer from "./renderers/ColumnRenderer.vue";
import DetailViewer from "@/components/ElTablePro/renderers/DetailViewer.vue";
import { Document } from "@element-plus/icons-vue";

// --- Props & Emits ---
const props = defineProps({
  data: { type: Array, required: true },
  columns: { type: Array, required: true },
  showPagination: { type: Boolean, default: true },
  height: { type: [Number, String], default: "800" },
  total: { type: Number, default: 0 },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 20 },
  pageSizes: { type: Array, default: () => [20, 25, 50, 100] },
  layout: { type: String, default: "total, sizes, prev, pager, next, jumper" },
  actionsConfig: {
    type: Object,
    default: () => ({
      width: 140,
    }),
  },
  showDetail: { type: Boolean, default: false },
});
const emit = defineEmits([
  "update:currentPage",
  "update:pageSize",
  "page-change",
]);

const handlePageChange = (type, value) => {
  if (type === "currentPage") emit("update:currentPage", value);
  if (type === "pageSize") emit("update:pageSize", value);

  emit("page-change");
};

const actionsSlotConfig = computed(() => {
  const defaults = {
    label: "操作",
    width: 120,
    align: "center",
    fixed: "right",
  };

  return { ...defaults, ...props.actionsConfig };
});

const detailVisible = ref(false);
const detailRow = ref({});
const openDetail = (row) => {
  detailRow.value = { ...row };
  detailVisible.value = true;
};

onMounted(() => {
  emit("update:currentPage", props.currentPage);
  emit("update:pageSize", props.pageSize);
});
</script>

<style scoped>
/*
 * ==========================================================================
 *  ElTablePro 样式作用域
 * ==========================================================================
 *  使用 `scoped` 属性确保此组件的样式不会泄露到全局，也不会被全局样式轻易污染。
 *  使用 `:deep()` 伪类选择器来“穿透”作用域，以修改子组件（Element Plus）的内部样式。
 *  所有针对子组件的样式都以根组件的 .el-table-pro-wrapper 类作为起点，增加了样式的特异性，
 *  确保能覆盖 Element Plus 的默认样式。
 * ==========================================================================
 */

/*
 * --------------------------------------------------------------------------
 *  1. 全局结构重置 (Structural Resets)
 * --------------------------------------------------------------------------
 */

/* 移除 Element Plus 在表格外部和分组时可能添加的伪元素边框线 */
/* 这让我们的自定义边框系统成为唯一的视觉标准，使外观更统一、干净 */
.el-table-pro-wrapper :deep(.el-table--border::after),
.el-table-pro-wrapper :deep(.el-table--group::after),
.el-table-pro-wrapper :deep(.el-table::before) {
  content: none;
}

.el-table-pro-wrapper :deep(.el-table__inner-wrapper) {
  @apply border rounded-md border-slate-200 dark:border-slate-700 !important;
  box-sizing: border-box;
}

.el-table-pro-wrapper :deep(.el-table__empty-block) {
  @apply dark:bg-slate-800 overflow-x-hidden !important;
}

.el-table-pro-wrapper.is-empty :deep(.el-scrollbar__bar.is-horizontal) {
  display: none !important;
}

/*
 * --------------------------------------------------------------------------
 *  2. 表头样式 (Header Styling)
 * --------------------------------------------------------------------------
 */

/*
  【核心】自定义表头单元格的样式
  通过 `header-cell-class-name` prop 应用
*/
.el-table-pro-wrapper :deep(.el-table-pro-header-cell--custom) {
  /*
   * 背景与文字:
   * - bg-slate-50: 使用一个非常柔和、不刺眼的浅灰色作为背景。
   * - text-slate-500: 文字颜色稍暗，与背景形成舒适的对比度。
   * - font-medium: 适中的字重，既有强调效果又不过于粗重。
   */
  @apply bg-gray-300 text-slate-800 font-bold dark:bg-slate-700 dark:text-slate-300 !important;

  /*
   * 边距与字体:
   * - py-3: 增加垂直内边距，使表头呼吸感更强，更显专业。
   * - px-4: 设定水平内边距 (1rem / 16px)。
   * - text-sm: 设定字体大小为 14px，这是UI表格中常见的、易读的尺寸。
   */
  @apply text-sm !important;

  /*
   * 边框:
   * - border-b border-slate-300: 只保留清晰的下边框作为与内容区的分隔。
   */
  @apply border-b border-slate-300 dark:border-slate-600 !important;
}

/*
 * --------------------------------------------------------------------------
 *  3. 内容单元格样式 (Cell Styling)
 * --------------------------------------------------------------------------
 */

/*
  【核心】自定义内容单元格的样式
  通过 `cell-class-name` prop 应用
*/
.el-table-pro-wrapper :deep(.el-table-pro-cell--custom) {
  /*
   * 文字样式:
   * - text-slate-600: 内容文字颜色比表头更深，确保数据清晰易读。
   * - text-sm: 与表头保持一致的字体大小(14px)，确保视觉统一性。
   */
  @apply text-slate-600 text-xs dark:text-slate-400 !important;

  /*
   * 边框:
   * - border-b border-slate-100: 使用比表头分隔线更浅的颜色作为行间分隔线，
   *   既能区分行，又不会造成过多的视觉干扰。
   */
  @apply border-b border-slate-100 dark:border-slate-700 !important;
}

/*
 * --------------------------------------------------------------------------
 *  4. 交互效果 (Interaction Effects)
 * --------------------------------------------------------------------------
 */

/*
  【核心】行悬停效果 (Row Hover Effect)
  当鼠标悬停在某一行上时，该行所有单元格的背景色改变。
*/
.el-table-pro-wrapper :deep(.el-table__body tr:hover > td) {
  /*
   * bg-slate-50: 使用与表头背景相同的柔和颜色作为高亮色，
   * 既能提供清晰的视觉反馈，又保持了整体设计的和谐统一。
   */
  @apply bg-slate-50 dark:bg-slate-700 !important;
}

/*
 * --------------------------------------------------------------------------
 *  分页器样式
 * --------------------------------------------------------------------------
 */

/* 分页器整体容器样式 */
.el-table-pro-wrapper :deep(.el-pagination) {
  @apply bg-transparent border-0 shadow-none px-0 py-0;
}

/* 分页按钮基础样式 */
.el-table-pro-wrapper :deep(.el-pagination .el-pager li) {
  @apply mx-1 px-3 py-1 text-base transition-all duration-150 font-medium text-slate-600 dark:text-slate-300;
}

/* 当前页按钮高亮 */
.el-table-pro-wrapper :deep(.el-pagination .el-pager li.is-active) {
  @apply bg-blue-500 text-white dark:bg-blue-400 dark:text-slate-900 font-bold shadow-md;
}

/* 悬停效果 */
.el-table-pro-wrapper :deep(.el-pagination .el-pager li:not(.is-active):hover) {
  @apply bg-blue-100 text-blue-700 dark:bg-slate-700 dark:text-blue-300 shadow;
}

/* 其他按钮 */
.el-table-pro-wrapper :deep(.el-pagination button) {
  @apply mx-1 px-2 py-1 text-slate-500 dark:text-slate-300 transition-all duration-150;
}

.el-table-pro-wrapper :deep(.el-pagination button:hover:not(:disabled)) {
  @apply bg-blue-100 text-blue-700 dark:bg-slate-700 dark:text-blue-300 shadow;
}

/* Jumper、Sizes 下拉等输入框样式 */

.el-table-pro-wrapper :deep(.el-pagination .el-input__inner) {
  @apply dark:bg-slate-800 text-slate-700 dark:text-slate-200;
  height: 24px;
}

.el-table-pro-wrapper :deep(.el-pagination .el-select) {
  @apply pr-4;
}

.el-table-pro-wrapper :deep(.el-pagination .el-select__wrapper) {
  @apply px-2.5 py-0.5;
}
</style>

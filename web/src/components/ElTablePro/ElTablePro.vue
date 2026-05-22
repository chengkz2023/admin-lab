<template>
  <div
    class="el-table-pro-wrapper"
    :class="{ 'is-empty': !Array.isArray(data) || data.length === 0 }"
  >
    <div v-if="$slots.toolbar" class="el-table-pro-toolbar">
      <slot name="toolbar" />
    </div>

    <el-table
      class="el-table-pro"
      :data="data"
      v-bind="$attrs"
      header-cell-class-name="el-table-pro-header-cell--custom"
      cell-class-name="el-table-pro-cell--custom"
      size="small"
      :height="height"
    >
      <ColumnRenderer
        v-for="column in columns"
        :key="column.prop || column.label || column.type"
        :column="column"
        :table-slots="$slots"
      />

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
          <slot name="actions" :row="row" :index="$index" />
        </template>
      </el-table-column>

      <template #empty>
        <span class="el-table-pro-empty">暂无数据</span>
      </template>
    </el-table>

    <div v-if="showPagination" class="el-table-pro-pagination">
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
import { computed, onMounted, ref } from 'vue'
import { Document } from '@element-plus/icons-vue'
import ColumnRenderer from './renderers/ColumnRenderer.vue'
import DetailViewer from '@/components/ElTablePro/renderers/DetailViewer.vue'

const props = defineProps({
  data: { type: Array, required: true },
  columns: { type: Array, required: true },
  showPagination: { type: Boolean, default: true },
  height: { type: [Number, String], default: '800' },
  total: { type: Number, default: 0 },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 20 },
  pageSizes: { type: Array, default: () => [20, 25, 50, 100] },
  layout: { type: String, default: 'total, sizes, prev, pager, next, jumper' },
  actionsConfig: {
    type: Object,
    default: () => ({
      width: 140
    })
  },
  showDetail: { type: Boolean, default: false }
})

const emit = defineEmits([
  'update:currentPage',
  'update:pageSize',
  'page-change'
])

const handlePageChange = (type, value) => {
  if (type === 'currentPage') emit('update:currentPage', value)
  if (type === 'pageSize') emit('update:pageSize', value)
  emit('page-change')
}

const actionsSlotConfig = computed(() => ({
  label: '操作',
  width: 120,
  align: 'center',
  fixed: 'right',
  ...props.actionsConfig
}))

const detailVisible = ref(false)
const detailRow = ref({})

const openDetail = (row) => {
  detailRow.value = { ...row }
  detailVisible.value = true
}

onMounted(() => {
  emit('update:currentPage', props.currentPage)
  emit('update:pageSize', props.pageSize)
})
</script>

<style scoped>
.el-table-pro-wrapper {
  background: var(--el-bg-color);
  color: var(--el-text-color-primary);
}

.el-table-pro-toolbar {
  padding: 16px 16px 8px;
}

.el-table-pro-wrapper :deep(.el-table--border::after),
.el-table-pro-wrapper :deep(.el-table--group::after),
.el-table-pro-wrapper :deep(.el-table::before) {
  content: none;
}

.el-table-pro-wrapper :deep(.el-table__inner-wrapper) {
  overflow: hidden;
  box-sizing: border-box;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: var(--el-border-radius-base);
}

.el-table-pro-wrapper :deep(.el-table),
.el-table-pro-wrapper :deep(.el-table__empty-block),
.el-table-pro-wrapper :deep(.el-table__body td.el-table__cell),
.el-table-pro-wrapper :deep(.el-table__fixed-right),
.el-table-pro-wrapper :deep(.el-table__fixed-left) {
  background: var(--el-bg-color) !important;
}

.el-table-pro-wrapper :deep(.el-table) {
  --el-table-header-bg-color: var(--el-fill-color-light);
  --el-table-header-text-color: var(--el-text-color-primary);
  --el-table-border-color: var(--el-border-color-lighter);
  --el-table-row-hover-bg-color: var(--el-fill-color-lighter);
}

.el-table-pro-wrapper.is-empty :deep(.el-scrollbar__bar.is-horizontal) {
  display: none !important;
}

.el-table-pro-wrapper :deep(.el-table-pro-header-cell--custom) {
  background: var(--el-fill-color-light) !important;
  color: var(--el-text-color-primary) !important;
  font-size: 13px !important;
  font-weight: 600 !important;
  border-bottom: 1px solid var(--el-border-color) !important;
}

.el-table-pro-wrapper :deep(.el-table__header th.el-table__cell),
.el-table-pro-wrapper :deep(.el-table__fixed-header-wrapper th.el-table__cell),
.el-table-pro-wrapper :deep(.el-table__fixed-right th.el-table__cell),
.el-table-pro-wrapper :deep(.el-table__fixed-left th.el-table__cell),
.el-table-pro-wrapper :deep(.el-table__fixed-right-patch) {
  background: var(--el-fill-color-light) !important;
  color: var(--el-text-color-primary) !important;
}

.el-table-pro-wrapper :deep(.el-table-pro-header-cell--custom .cell) {
  display: flex;
  min-height: 40px;
  align-items: center;
  line-height: 1.2;
}

.el-table-pro-wrapper :deep(.el-table-pro-header-cell--custom.is-center .cell) {
  justify-content: center;
  text-align: center;
}

.el-table-pro-wrapper :deep(.el-table-pro-header-cell--custom.is-right .cell) {
  justify-content: flex-end;
  text-align: right;
}

.el-table-pro-wrapper :deep(.el-table-pro-header-cell--custom.is-left .cell) {
  justify-content: flex-start;
  text-align: left;
}

.el-table-pro-wrapper :deep(.el-table-pro-header-cell--custom .cell .el-checkbox) {
  margin: 0 auto;
}

.el-table-pro-wrapper :deep(.el-table-pro-cell--custom) {
  color: var(--el-text-color-regular) !important;
  font-size: 12px !important;
  border-bottom: 1px solid var(--el-border-color-lighter) !important;
}

.el-table-pro-wrapper :deep(.el-table__body tr:hover > td.el-table__cell) {
  background: var(--el-fill-color-lighter) !important;
}

.el-table-pro-wrapper :deep(.el-table__fixed-right::before),
.el-table-pro-wrapper :deep(.el-table__fixed-left::before) {
  background: var(--el-border-color-lighter);
}

.el-table-pro-wrapper :deep(.el-table__fixed),
.el-table-pro-wrapper :deep(.el-table__fixed-right) {
  box-shadow: var(--el-box-shadow-lighter);
}

.el-table-pro-empty {
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.el-table-pro-pagination {
  display: flex;
  justify-content: flex-end;
  padding: 0 8px;
}

.el-table-pro-wrapper :deep(.el-pagination) {
  margin-top: 14px;
  padding: 0;
  background: transparent;
  border: 0;
  box-shadow: none;
}

.el-table-pro-wrapper :deep(.el-pagination .el-pager li),
.el-table-pro-wrapper :deep(.el-pagination button) {
  margin: 0 4px;
  color: var(--el-text-color-regular);
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color);
  border-radius: var(--el-border-radius-small);
  transition: color 0.15s, background-color 0.15s, border-color 0.15s;
}

.el-table-pro-wrapper :deep(.el-pagination .el-pager li.is-active) {
  color: #fff !important;
  background: var(--el-color-primary) !important;
  border-color: var(--el-color-primary) !important;
}

.el-table-pro-wrapper :deep(.el-pagination .el-pager li:not(.is-active):hover),
.el-table-pro-wrapper :deep(.el-pagination button:hover:not(:disabled)) {
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
  border-color: var(--el-color-primary-light-5);
}

.el-table-pro-wrapper :deep(.el-pagination button:disabled) {
  color: var(--el-text-color-disabled);
  background: var(--el-disabled-bg-color);
  border-color: var(--el-border-color-light);
}

.el-table-pro-wrapper :deep(.el-pagination .el-input__inner) {
  height: 24px;
  color: var(--el-text-color-regular);
  background: var(--el-bg-color);
}

.el-table-pro-wrapper :deep(.el-pagination .el-select) {
  padding-right: 16px;
}

.el-table-pro-wrapper :deep(.el-pagination .el-select__wrapper) {
  padding: 2px 10px;
  background: var(--el-bg-color);
}
</style>

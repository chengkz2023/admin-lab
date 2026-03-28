<template>
  <div class="lab-table-pro">
    <el-card shadow="never">
      <template #header>
        <div class="panel-header">
          <div class="panel-title-wrap">
            <slot name="title">
              <span>{{ title }}</span>
            </slot>
            <el-tag v-if="showTotalTag">总数 {{ total }}</el-tag>
          </div>
          <div class="panel-actions">
            <slot
              name="toolbar-left"
              :selected-rows="selectedRows"
              :query="queryParams"
              :reload="reload"
            />
            <el-button
              v-if="enableColumnSetting"
              type="primary"
              plain
              @click="columnSettingVisible = true"
            >
              列设置
            </el-button>
            <slot
              name="toolbar-right"
              :selected-rows="selectedRows"
              :query="queryParams"
              :reload="reload"
            />
          </div>
        </div>
      </template>

      <slot name="table-before" :rows="rows" :selected-rows="selectedRows" />

      <div v-if="errorInfo" class="table-error-wrap">
        <slot name="error" :error="errorInfo" :reload="reload">
          <el-alert title="表格数据加载失败，请稍后重试" type="error" show-icon :closable="false" />
        </slot>
      </div>

      <el-table
        v-loading="tableLoading"
        :data="rows"
        :row-key="rowKey"
        :row-class-name="rowClassName"
        :border="border"
        :stripe="stripe"
        v-bind="tableProps"
        @sort-change="handleSortChange"
        @selection-change="handleSelectionChange"
      >
        <template #empty>
          <slot name="empty">暂无数据</slot>
        </template>

        <template v-for="column in visibleColumns" :key="column.key">
          <el-table-column
            v-if="column.type === 'selection'"
            type="selection"
            :width="column.width || 48"
            :fixed="column.fixed"
            :align="column.align"
            :selectable="column.selectable"
            :reserve-selection="reserveSelection"
          />

          <el-table-column
            v-else-if="column.type === 'index'"
            type="index"
            :label="column.label"
            :width="column.width || 60"
            :fixed="column.fixed"
            :align="column.align"
          />

          <el-table-column
            v-else
            :prop="column.prop || column.key"
            :label="column.label"
            :min-width="column.minWidth || 120"
            :width="column.width"
            :fixed="column.fixed"
            :align="column.align"
            :sortable="column.sortable || false"
            :show-overflow-tooltip="column.showOverflowTooltip"
          >
            <template #default="scope">
              <slot
                :name="column.slot || `cell-${column.prop || column.key}`"
                :row="scope.row"
                :column="column"
                :index="scope.$index"
                :value="scope.row[column.prop || column.key]"
              >
                {{
                  formatCellValue(
                    scope.row,
                    column,
                    scope.row[column.prop || column.key]
                  )
                }}
              </slot>
            </template>
          </el-table-column>
        </template>
      </el-table>

      <slot name="table-after" :rows="rows" :selected-rows="selectedRows" />

      <div v-if="showPagination" class="pager">
        <el-pagination
          v-model:current-page="pager.page"
          v-model:page-size="pager.pageSize"
          :page-sizes="pageSizes"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @current-change="fetchPage"
          @size-change="handlePageSizeChange"
        />
      </div>
    </el-card>

    <el-drawer
      v-if="enableColumnSetting"
      v-model="columnSettingVisible"
      title="列设置"
      size="420px"
    >
      <div class="column-tip">支持控制列显示与排序。</div>
      <div class="column-list">
        <div
          v-for="(column, index) in configurableColumns"
          :key="column.key"
          class="column-item"
        >
          <el-checkbox
            :model-value="visibleColumnKeys.includes(column.key)"
            @change="(val) => toggleColumn(column.key, val)"
          >
            {{ column.label }}
          </el-checkbox>
          <div class="column-order">
            <el-button text size="small" @click="moveColumnUp(index)">上移</el-button>
            <el-button text size="small" @click="moveColumnDown(index)">下移</el-button>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref, watch } from 'vue'
  import { ElMessage } from 'element-plus'

  defineOptions({
    name: 'LabTablePro'
  })

  const props = defineProps({
    title: {
      type: String,
      default: '结果列表'
    },
    columns: {
      type: Array,
      default: () => []
    },
    queryParams: {
      type: Object,
      default: () => ({})
    },
    fetcher: {
      type: Function,
      required: true
    },
    loading: {
      type: Boolean,
      default: undefined
    },
    rowKey: {
      type: String,
      default: 'id'
    },
    rowClassName: {
      type: [Function, String],
      default: ''
    },
    reserveSelection: {
      type: Boolean,
      default: false
    },
    immediate: {
      type: Boolean,
      default: true
    },
    showPagination: {
      type: Boolean,
      default: true
    },
    pageSize: {
      type: Number,
      default: 10
    },
    pageSizes: {
      type: Array,
      default: () => [10, 20, 50]
    },
    pageField: {
      type: String,
      default: 'page'
    },
    pageSizeField: {
      type: String,
      default: 'pageSize'
    },
    listField: {
      type: String,
      default: 'list'
    },
    totalField: {
      type: String,
      default: 'total'
    },
    sortMapper: {
      type: Function,
      default: null
    },
    border: {
      type: Boolean,
      default: true
    },
    stripe: {
      type: Boolean,
      default: false
    },
    tableProps: {
      type: Object,
      default: () => ({})
    },
    enableColumnSetting: {
      type: Boolean,
      default: true
    },
    persistKey: {
      type: String,
      default: ''
    },
    showTotalTag: {
      type: Boolean,
      default: true
    }
  })

  const emit = defineEmits(['selection-change', 'loaded', 'error'])

  const innerLoading = ref(false)
  const rows = ref([])
  const total = ref(0)
  const selectedRows = ref([])
  const errorInfo = ref(null)
  const columnSettingVisible = ref(false)

  const pager = ref({ page: 1, pageSize: props.pageSize })
  const sortState = ref({ prop: '', order: '' })
  const configurableColumnKeys = ref([])
  const visibleColumnKeys = ref([])

  const tableLoading = computed(() => {
    if (typeof props.loading === 'boolean') {
      return props.loading
    }
    return innerLoading.value
  })

  const normalizedColumns = computed(() => {
    return (props.columns || []).map((column, index) => {
      const key = column.key || column.prop || `col_${index}`
      return {
        ...column,
        key,
        prop: column.prop || key
      }
    })
  })

  const configurableColumns = computed(() => {
    return normalizedColumns.value.filter(
      (column) =>
        column.configurable !== false &&
        column.type !== 'selection' &&
        column.type !== 'index'
    )
  })

  const visibleColumns = computed(() => {
    const map = new Map(normalizedColumns.value.map((col) => [col.key, col]))
    const pinned = normalizedColumns.value.filter(
      (col) => col.configurable === false || col.type === 'selection' || col.type === 'index'
    )

    const orderedConfigurable = configurableColumnKeys.value
      .filter((key) => visibleColumnKeys.value.includes(key))
      .map((key) => map.get(key))
      .filter(Boolean)

    return [...pinned, ...orderedConfigurable]
  })

  const storageKey = computed(() => {
    if (props.persistKey) {
      return `lab-table-pro:${props.persistKey}`
    }
    return ''
  })

  const normalizeColumnState = () => {
    const keys = configurableColumns.value.map((item) => item.key)
    const cache = readColumnState()
    const order = Array.isArray(cache.order)
      ? cache.order.filter((key) => keys.includes(key))
      : []
    const fullOrder = [...order, ...keys.filter((key) => !order.includes(key))]
    const visible = Array.isArray(cache.visible)
      ? cache.visible.filter((key) => keys.includes(key))
      : []
    const fullVisible = visible.length ? visible : [...keys]

    configurableColumnKeys.value = fullOrder
    visibleColumnKeys.value = fullVisible
  }

  const readColumnState = () => {
    if (!storageKey.value) {
      return {}
    }
    try {
      const raw = localStorage.getItem(storageKey.value)
      return raw ? JSON.parse(raw) : {}
    } catch (error) {
      return {}
    }
  }

  const persistColumnState = () => {
    if (!storageKey.value) {
      return
    }
    localStorage.setItem(
      storageKey.value,
      JSON.stringify({
        order: configurableColumnKeys.value,
        visible: visibleColumnKeys.value
      })
    )
  }

  const getByPath = (source, path) => {
    if (!path || !source) {
      return undefined
    }
    const keys = String(path).split('.').filter(Boolean)
    return keys.reduce((acc, key) => (acc == null ? undefined : acc[key]), source)
  }

  const resolveSortMeta = () => {
    const fallback = {
      sortBy: sortState.value.prop || '',
      sortOrder:
        sortState.value.order === 'ascending'
          ? 'asc'
          : sortState.value.order === 'descending'
            ? 'desc'
            : ''
    }
    if (typeof props.sortMapper !== 'function') {
      return fallback
    }
    const mapped = props.sortMapper({
      prop: sortState.value.prop || '',
      order: sortState.value.order || ''
    })
    if (!mapped || typeof mapped !== 'object') {
      return fallback
    }
    return {
      sortBy: mapped.sortBy ?? fallback.sortBy,
      sortOrder: mapped.sortOrder ?? fallback.sortOrder
    }
  }

  const buildRequestPayload = () => {
    const sortMeta = resolveSortMeta()
    const dynamicPager = {
      [props.pageField]: pager.value.page,
      [props.pageSizeField]: pager.value.pageSize
    }
    return {
      ...dynamicPager,
      page: pager.value.page,
      pageSize: pager.value.pageSize,
      sortBy: sortMeta.sortBy,
      sortOrder: sortMeta.sortOrder,
      filters: { ...(props.queryParams || {}) },
      query: { ...(props.queryParams || {}) }
    }
  }

  const fetchPage = async () => {
    innerLoading.value = true
    errorInfo.value = null
    try {
      const payload = buildRequestPayload()
      const result = await props.fetcher(payload)
      const list =
        getByPath(result, props.listField) ??
        result?.list ??
        result?.rows ??
        []
      const totalCount =
        getByPath(result, props.totalField) ??
        result?.total ??
        0

      rows.value = Array.isArray(list) ? list : []
      total.value = Number(totalCount) || 0
      emit('loaded', { rows: rows.value, total: total.value, payload })
    } catch (error) {
      errorInfo.value = error
      emit('error', error)
      ElMessage.error('表格数据加载失败')
    } finally {
      innerLoading.value = false
    }
  }

  const reload = async ({ resetPage = false } = {}) => {
    if (resetPage) {
      pager.value.page = 1
    }
    await fetchPage()
  }

  const handlePageSizeChange = async () => {
    pager.value.page = 1
    await fetchPage()
  }

  const handleSortChange = async ({ prop, order }) => {
    sortState.value = { prop: prop || '', order: order || '' }
    pager.value.page = 1
    await fetchPage()
  }

  const handleSelectionChange = (val) => {
    selectedRows.value = val || []
    emit('selection-change', selectedRows.value)
  }

  const toggleColumn = (key, checked) => {
    if (checked) {
      if (!visibleColumnKeys.value.includes(key)) {
        visibleColumnKeys.value = [...visibleColumnKeys.value, key]
      }
      return
    }
    if (visibleColumnKeys.value.length === 1 && visibleColumnKeys.value[0] === key) {
      ElMessage.warning('至少保留一列可见')
      return
    }
    visibleColumnKeys.value = visibleColumnKeys.value.filter((item) => item !== key)
  }

  const moveColumnUp = (index) => {
    if (index <= 0) {
      return
    }
    const next = [...configurableColumnKeys.value]
    const temp = next[index - 1]
    next[index - 1] = next[index]
    next[index] = temp
    configurableColumnKeys.value = next
  }

  const moveColumnDown = (index) => {
    if (index >= configurableColumnKeys.value.length - 1) {
      return
    }
    const next = [...configurableColumnKeys.value]
    const temp = next[index + 1]
    next[index + 1] = next[index]
    next[index] = temp
    configurableColumnKeys.value = next
  }

  const formatCellValue = (row, column, value) => {
    if (typeof column.formatter === 'function') {
      return column.formatter(row, column, value)
    }
    return value
  }

  watch(
    normalizedColumns,
    () => {
      normalizeColumnState()
    },
    { immediate: true }
  )

  watch(
    [visibleColumnKeys, configurableColumnKeys],
    () => {
      persistColumnState()
    },
    { deep: true }
  )

  defineExpose({
    fetchPage,
    reload
  })

  onMounted(() => {
    if (props.immediate) {
      fetchPage()
    }
  })
</script>

<style scoped>
  .lab-table-pro {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
  }

  .panel-title-wrap {
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: 600;
  }

  .panel-actions {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .table-error-wrap {
    margin-bottom: 12px;
  }

  .pager {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }

  .column-tip {
    margin-bottom: 12px;
    color: #64748b;
    font-size: 13px;
  }

  .column-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .column-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 12px;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    background: #f8fafc;
  }

  .column-order {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  @media (max-width: 768px) {
    .panel-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .pager {
      justify-content: flex-start;
    }
  }
</style>

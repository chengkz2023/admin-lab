<template>
  <!-- 多级表头 -->
  <el-table-column
    v-if="column.children && column.children.length > 0"
    v-bind="getElColumnProps(column)"
  >
    <ColumnRenderer
      v-for="childColumn in column.children"
      :key="childColumn.prop || childColumn.label"
      :column="childColumn"
      :table-slots="tableSlots"
    />
  </el-table-column>

  <!-- El-Table 内置类型列 (selection, index, expand) -->
  <el-table-column
    v-else-if="column.type === 'selection' || column.type === 'index'"
    v-bind="getElColumnProps(column)"
  />

  <!-- expand 类型需要特殊处理 -->
  <el-table-column
    v-else-if="column.type === 'expand'"
    v-bind="getElColumnProps(column)"
  >
    <template #default="{ row }">
      <!-- expand 类型使用插槽渲染展开内容 -->
      <component
        :is="tableSlots.expand"
        v-if="tableSlots.expand"
        :row="row"
      />
      <div
        v-else
        class="text-gray-400 text-sm"
      >
        未定义展开内容
      </div>
    </template>
  </el-table-column>

  <!-- 自定义类型列 -->
  <el-table-column
    v-else
    v-bind="getElColumnProps(column)"
  >
    <template #header="{ column: col }">
      <span>{{ col.label }}</span>
      <el-icon
        v-if="col.sortable"
        class="ml-1 text-gray-400 opacity-60"
      >
        <Sort />
      </el-icon>
    </template>

    <template #default="{ row }">
      <div
        class="flex items-center h-full"
        :class="getAlignClass(column.align)"
      >
        <!-- 优先使用插槽 -->
        <component
          :is="tableSlots[column.prop]"
          v-if="tableSlots[column.prop]"
          :row="row"
          :value="row[column.prop]"
          :column="column"
        />
        <!-- 使用自定义渲染器 -->
        <component
          :is="resolveCellComponent(column.type)"
          v-else
          :value="row[column.prop]"
          :row="row"
          :column="column"
        />
      </div>
    </template>
  </el-table-column>
</template>

<script setup>
import {Sort} from '@element-plus/icons-vue';
// 递归渲染子列需要导入自身
import ColumnRenderer from './ColumnRenderer.vue';
import {
  resolveCellComponent,
  getAlignClass,
  getElColumnProps,
} from '../renderer.js'; // 导入中央逻辑模块

/**
 * @file ColumnRenderer.vue
 * @description 递归渲染 ElTablePro 的列组件。它能够处理多级表头、El-Table 内置类型列
 *              (如 selection, index, expand) 以及自定义渲染器类型。
 *              这是 ElTablePro 灵活列配置的核心。
 *
 * @usage
 * 该组件通常由 ElTablePro 内部使用，不直接暴露给外部。
 * 它的作用是根据 `column` 配置对象，动态渲染不同类型的 `el-table-column`。
 *
 * 示例列配置 (在 ElTablePro 的 `columns` prop 中使用):
 *
 * // 1. 基础文本列 (默认渲染器)
 * { prop: 'name', label: '名称' },
 *
 * // 2. 自定义类型列 (使用内置渲染器，如 'status')
 * { prop: 'status', label: '状态', type: 'status' },
 *
 * // 3. 多级表头
 * {
 *   label: '网络信息',
 *   children: [
 *     { prop: 'ipAddress', label: 'IP地址', type: 'ip' },
 *     { prop: 'port', label: '端口', type: 'port' },
 *   ]
 * },
 *
 * // 4. El-Table 内置类型列
 * { type: 'selection', width: 50 },
 * { type: 'index', label: '序号', width: 60 },
 * { type: 'expand', label: '详情', width: 80 },
 */


defineProps({
  column: {type: Object, default: null},
  tableSlots: {type: Object, default: () => ({})},
});
</script>

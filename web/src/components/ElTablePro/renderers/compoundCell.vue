<template>
  <div
    class="py-1 grid w-full gap-x-2 gap-y-1 items-center"
    :style="gridStyle"
  >
    <template
      v-for="item in processedItems"
      :key="item.column.prop"
    >
      <!-- 标签部分 -->
      <span class="text-xs text-slate-500 text-right">
        {{ item.column.label }}:
      </span>

      <!--
        值部分：
        - :class 直接绑定预先计算好的 `valueContainerClass`。
        - :is 直接绑定预先计算好的 `component`。
        - 模板现在是 100% 的声明式。
      -->
      <div
        class="w-min-0"
        :class="item.valueContainerClass"
      >
        <component
          :is="item.component"
          :value="row[item.column.prop]"
          :row="row"
          :column="item.column"
        />
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { resolveCellComponent } from '../renderer.js';

/**
 * @file CompoundCell.vue
 * @description 一个“复合单元格”渲染器，可以在单个单元格内垂直展示多个子指标。
 *              它使用 CSS Grid 实现完美的动态对齐，并且内部的每个子指标都可以
 *              使用 ProTable 的任何其他渲染器类型。
 *
 * @usage
 * {
 *   prop: 'metrics', // 主 prop，通常不直接使用
 *   label: '聚合指标',
 *   type: 'compound',
 *   width: 250,
 *   align: 'center',
 *   sortBy: 'cpuUsage' // 根据哪个子指标排序
 *   // [核心] 定义要在单元格内展示的子指标列表
 *   items: [
 *     // 每个子项都是一个迷你的列配置对象
 *     { prop: 'status', label: '状态', type: 'status', align: 'right' },
 *     { prop: 'cpuUsage', label: 'CPU使用率', type: 'percentage' },
 *     { prop: 'avgLatency', label: '平均响应延迟', type: 'duration', config: { unit: 'ms' } },
 *     { prop: 'owner', label: '负责人', format: (val) => `🧑‍💻 ${val}` }
 *   ]
 * }
 */

const props = defineProps({
  value: { type: [String, Number, Boolean, Object, Array], default: null },
  row: { type: Object, default: () => ({}) },
  column: { type: Object, default: () => ({}) },
});

// --- 静态的对齐映射表 ---
const alignClassMap = {
  left: 'justify-start',
  center: 'justify-center',
  right: 'justify-end',
};


const gridStyle = computed(() => {
  const labelWidth = props.column?.config?.labelWidth || 'max-content';
  return {
    gridTemplateColumns: `${labelWidth} 1fr`,
  };
});

const processedItems = computed(() => {
  if (!props.column.items) return [];

  return props.column?.items?.map(itemColumn => {
    // 1. 获取对齐配置，提供默认值
    const align = itemColumn.align || 'center';
    // 2. 从映射表中查找对应的 class
    const alignClass = alignClassMap[align] || alignClassMap.left;

    return {
      column: itemColumn,
      // 3. 预先解析组件
      component: resolveCellComponent(itemColumn.type),
      // 4. 预先计算好容器的 class
      valueContainerClass: `flex ${alignClass}`,
    };
  });
});
</script>

<style scoped>

</style>

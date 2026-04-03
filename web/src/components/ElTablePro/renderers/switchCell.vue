<template>
  <el-switch
    v-model="switchValue"
    :active-value="activeValue"
    :inactive-value="inactiveValue"
    :active-text="activeText"
    :inactive-text="inactiveText"
    :inline-prompt="inlinePrompt"
    :size="switchSize"
    :loading="switchLoading"
    :disabled="switchDisabled"
    :before-change="handleBeforeChange"
    @change="handleChange"
  />
</template>

<script setup>
import { computed } from "vue";

/**
 * @file switchCell.vue
 * @description ElTablePro 的轻量开关型单元格。
 *              只负责渲染 el-switch，并通过 column.config 暴露必要钩子：
 *              1. active/inactive 值和文案配置
 *              2. disabled/loading 的静态值或同步函数计算
 *              3. beforeChange 异步确认
 *              4. onChange 业务回调
 *
 * @usage
 * {
 *   prop: 'enable',
 *   label: '启用',
 *   type: 'switch',
 *
 *   config: {
 *     activeValue: 1,         // 开启值，默认为 1
 *     inactiveValue: 0,       // 关闭值，默认为 0
 *     activeText: '开',       // 开启文案
 *     inactiveText: '关',     // 关闭文案
 *     inlinePrompt: true,     // 是否将文案显示在开关内部
 *     size: 'small',          // ElSwitch 尺寸
 *     loading: ({ row }) => Boolean(row.loading), // 静态 loading 或同步函数
 *     disabled: ({ row }) => row.locked,          // 静态 disabled 或同步函数
 *     beforeChange: async ({ row, nextValue, previousValue }) => {
 *       // 返回 false 可取消切换
 *       return true;
 *     },
 *     onChange: async ({ row, nextValue, previousValue }) => {
 *       // 业务提交、提示、回滚由页面自己处理
 *       await saveSwitchState(row.id, nextValue);
 *     },
 *   }
 * }
 *
 * // 示例数据:
 * // { enable: 1, locked: false } -> 正常显示为开启，可切换
 * // { enable: 0, locked: true }  -> 显示为关闭，但不可操作
 */

const props = defineProps({
  value: { type: [String, Number, Boolean], default: null },
  row: { type: Object, default: () => ({}) },
  column: { type: Object, default: () => ({}) },
});

const config = computed(() => props.column?.config || {});
const activeValue = computed(() => config.value.activeValue ?? 1);
const inactiveValue = computed(() => config.value.inactiveValue ?? 0);
const activeText = computed(() => config.value.activeText ?? "");
const inactiveText = computed(() => config.value.inactiveText ?? "");
const inlinePrompt = computed(() => Boolean(config.value.inlinePrompt));
const switchSize = computed(() => config.value.size);

const getCellValue = () => {
  const prop = props.column?.prop;

  if (prop && props.row && typeof props.row === "object") {
    return props.row[prop];
  }

  return props.value;
};

const setCellValue = (nextValue) => {
  const prop = props.column?.prop;

  if (prop && props.row && typeof props.row === "object") {
    props.row[prop] = nextValue;
  }
};

const buildContext = (nextValue = getCellValue()) => ({
  row: props.row,
  column: props.column,
  value: nextValue,
  previousValue: getCellValue(),
  nextValue,
  activeValue: activeValue.value,
  inactiveValue: inactiveValue.value,
});

const isPromiseLike = (value) =>
  value && typeof value === "object" && typeof value.then === "function";

// disabled/loading 只接受布尔值或同步函数，避免把 Promise 直接绑定到模板上。
const resolveBooleanConfig = (target, fallback = false) => {
  if (target == null) return fallback;

  try {
    const resolved =
      typeof target === "function" ? target(buildContext()) : target;

    if (isPromiseLike(resolved)) {
      return fallback;
    }

    return Boolean(resolved);
  } catch (_error) {
    return fallback;
  }
};

const switchValue = computed({
  get: () => getCellValue(),
  set: (nextValue) => {
    setCellValue(nextValue);
  },
});

const switchLoading = computed(() =>
  resolveBooleanConfig(config.value.loading),
);

const switchDisabled = computed(() => {
  if (!props.column?.prop) {
    return true;
  }

  return switchLoading.value || resolveBooleanConfig(config.value.disabled);
});

const getNextValue = (currentValue) =>
  currentValue === activeValue.value ? inactiveValue.value : activeValue.value;

// 交给 ElSwitch 的 before-change，只负责在切换前执行确认逻辑。
const handleBeforeChange = async () => {
  const beforeChange = config.value.beforeChange;

  if (typeof beforeChange !== "function") {
    return true;
  }

  return await beforeChange(buildContext(getNextValue(getCellValue())));
};

// onChange 只做事件分发，业务页自行决定是否提示、刷新或回滚。
const handleChange = async (nextValue) => {
  const onChange = config.value.onChange;

  if (typeof onChange === "function") {
    await onChange(buildContext(nextValue));
  }
};
</script>

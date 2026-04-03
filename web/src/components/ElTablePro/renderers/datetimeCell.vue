<template>
  <el-tooltip
      v-if="formatted.absolute"
      placement="top"
  >
    <template #content>
      <div class="font-mono text-xs dark:text-slate-500">
        {{ formatted.relative }}
      </div>
    </template>
    <span
        class="datetime-display cursor-help font-mono text-sm text-slate-700 dark:text-slate-300"
    >
      {{ formatted.absolute }}
    </span>
  </el-tooltip>
  <span
      v-else
      class="text-gray-400 dark:text-gray-300"
  >—</span>
</template>

<script setup>
import { computed } from 'vue';
import { ElTooltip } from 'element-plus';
import dayjs from 'dayjs'


/**
 * @file DatetimeCell.vue
 * @description 一个智能的时间渲染器，默认显示格式化的精确时间，悬停时显示相对时间。
 *
 * @usage
 * {
 *   prop: 'lastLoginTime',
 *   label: '上次登录',
 *   type: 'datetime',
 *
 *   // [可选] 提供一个 config 对象来自定义格式化行为
 *   config: {
 *     // 自定义默认显示的日期时间格式。
 *     // 如果不提供，则使用 'YYYY-MM-DD HH:mm:ss'。
 *     format: 'YYYY/MM/DD HH:mm',
 *
 *     // [可选] 是否将值解析为 Unix 时间戳（秒或毫秒）
 *     // 如果为 ture, 则会将值视为 Unix 时间戳进行解析。
 *     useUnix: true,
 *   }
 * }
 */

const props = defineProps({
  value: { type: [ String, Number, Date ], default: null },
  column: { type: Object, default: () => ({}) },
});

// --- 从列配置中读取或使用默认值 ---
const config = computed(() => props.column?.config || {});
const displayFormat = computed(() => config.value?.format || 'YYYY-MM-DD HH:mm:ss');
const useUnix = computed(() => config.value?.useUnix || props.column?.useUnix || false);

// --- 核心计算属性 ---
const formatted = computed(() => {
  const val = props.value
  if (val === null || val === undefined) {
    return { relative: null, absolute: null };
  }

  let date
  const numVal = Number(val)
  let len
  // 处理时间戳
  if (useUnix.value && !isNaN(numVal)) {
    len = String(numVal).length;
    if (len === 10) { // 秒
      date = dayjs(numVal * 1000)
    } else if (len === 13) { // 毫秒
      date = dayjs(numVal)
    } else if (len === 19) { // 纳秒
      date = dayjs(Math.floor(numVal / 1_000_000))
    } else date = dayjs(numVal);
  }
  // 处理字符串时间
  else {
    date = dayjs(val)
  }

  if (!dayjs.isDayjs(date) || !date.isValid()) return { relative: null, absolute: null };

  // 计算相对时间 (Tooltip 内容)
  const now = dayjs();
  const diffSeconds = now.diff(date, 'second')
  let relative = '';

  if (diffSeconds < 60) {
    relative = '刚刚';
  } else if (diffSeconds < 3600) {
    relative = `${Math.floor(diffSeconds / 60)} 分钟前`;
  } else if (diffSeconds < 86400) {
    relative = `${Math.floor(diffSeconds / 3600)} 小时前`;
  } else if (diffSeconds < 2592000) {
    relative = `${Math.floor(diffSeconds / 86400)} 天前`;
  } else relative = '很久以前';

  // 格式化绝对时间 (默认显示)，使用从 config 中解析出的格式
  const fmt = useUnix.value && len !== 10 ? displayFormat.value.replace(/ss/, 'ss.SSS') : displayFormat.value
  const absolute = date?.format(fmt)
  return { relative, absolute };
});

</script>

<style scoped>
/* 因被全局样式污染，所以采用 scoped css 来重置按钮样式 */
.datetime-display {
  border-bottom: 1px dashed #94a3b8
}
</style>


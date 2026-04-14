<!-- components/pro-table/renderers/UrlCell.vue -->
<template>
  <el-tooltip
    :content="value"
    placement="top"
    :disabled="!isTruncated"
  >
    <div class="group relative flex items-center w-full space-x-1">
      <!-- 链接图标，颜色与文本联动 -->
      <el-icon
        :color="'#3b82f6'"
        class="mr-2 flex-shrink-0 dark:text-blue-400"
      >
        <Link />
      </el-icon>
      <span
        ref="textRef"
        class=" truncate text-sm font-mono text-blue-600 dark:text-blue-400  dark:border-blue-500 cursor-pointer"
      >
        {{ value }}
      </span>

      <!-- 复制按钮 -->
      <button
        title="复制"
        class="copy-button flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
        @click.stop="handleCopy"
      >
        <el-icon
          :size="15"
          class="text-slate-400 hover:text-blue-600"
        >
          <CopyDocument />
        </el-icon>
      </button>
    </div>
  </el-tooltip>
</template>

<script setup>
import {ref, onMounted, onUpdated} from 'vue';
import {ElIcon, ElTooltip, ElMessage} from 'element-plus';
import {Link, CopyDocument} from '@element-plus/icons-vue';

/**
 * @file UrlCell.vue
 * @description URL 单元格渲染器，以链接的形式展示 URL，并提供复制功能。
 *              支持溢出截断和 Tooltip 提示。
 *
 * @usage
 * {
 *   prop: 'websiteUrl',
 *   label: '网站链接',
 *   type: 'url',
 * }
 *
 * // 示例数据:
 * // { websiteUrl: 'https://www.example.com/path/to/page' } -> 显示可点击的 URL，并可复制
 */

const props = defineProps({
  value: {type: String, default: null},
  column: { type: Object, default: () => ({}) },
});

const textRef = ref(null);
const isTruncated = ref(false);

const checkTruncation = () => {
  const el = textRef.value;
  if (el) {
    isTruncated.value = el.scrollWidth > el.offsetWidth;
  }
};

onMounted(checkTruncation);
onUpdated(checkTruncation);

const handleCopy = () => {
  navigator.clipboard.writeText(props.value).then(() => {
    ElMessage.success({
      message: '已复制到剪贴板',
      duration: 2000,
      showClose: true
    });
  }).catch(() => {
    ElMessage.error({
      message: '复制失败，请重试',
      duration: 2000
    });
  });
};
</script>

<style scoped>
/* 因被全局样式污染，所以采用 scoped css 来重置按钮样式 */
.url-display {
  border-bottom: 1px dashed #94a3b8
}

/* 因被全局样式污染，所以采用 scoped css 来重置按钮样式 */
.copy-button {
  /* 核心：移除背景和边框 */
  background: none;
  border: none;

  /* 移除内边距，让图表自己决定大小 */
  padding: 0;

  /* 确保鼠标悬停时显示为可点击的手型指针 */
  cursor: pointer;

  /* 垂直居中图表 */
  display: inline-flex;
  align-items: center;
}

.url-link-text {
  border-bottom: 1px dashed #60a5fa;
}
</style>

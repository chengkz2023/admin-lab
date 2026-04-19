<template>
  <el-popover
    v-model:visible="popoverVisible"
    placement="right"
    :width="800"
    trigger="click"
    popper-class="json-cell-popover"
    :show-arrow="true"
  >
    <!-- Popover 内容 -->
    <template #default>
      <div class="json-popover-content">
        <div class="json-header">
          <div class="flex items-center gap-2">
            <el-icon class="text-blue-500">
              <Document />
            </el-icon>
            <span class="font-medium text-gray-800 dark:text-gray-200">JSON 数据预览</span>
          </div>
          <div class="flex items-center gap-2">
            <el-button
              link
              type="primary"
              size="small"
              class="copy-btn"
              @click="copyJson"
            >
              <el-icon class="mr-1">
                <CopyDocument />
              </el-icon>
              复制
            </el-button>
            <el-button
              link
              size="small"
              class="close-btn"
              @click="popoverVisible = false"
            >
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
        </div>

        <div class="json-viewer-container">
          <pre class="json-pre">{{ formattedJsonString }}</pre>
        </div>
      </div>
    </template>

    <!-- 触发 Popover 的元素 -->
    <template #reference>
      <div class="json-cell-preview">
        <div class="preview-content">
          <el-button type="success" link icon="Document"></el-button>
<!--          <el-icon class="preview-icon">-->
<!--            <Document />-->
<!--          </el-icon>-->
          <span v-if="previewText" class="preview-text">{{ previewText }}</span>
        </div>
      </div>
    </template>
  </el-popover>
</template>

<script setup>
import { ref, computed } from 'vue';
import { ElPopover, ElButton, ElIcon, ElMessage } from 'element-plus';
import { Document, CopyDocument, Close, ArrowRight } from '@element-plus/icons-vue';

/**
 * @file JsonCell.vue
 * @description JSON 数据单元格渲染器，以可点击的预览形式展示 JSON 字符串或对象。
 *              点击后在 Popover 中显示格式化后的 JSON，并提供复制功能。
 *
 * @usage
 * {
 *   prop: 'jsonData',
 *   label: 'JSON 数据',
 *   type: 'json',
 *
 *   // [可选]
 *   previewText: '查看详情', // 预览显示内容，默认为 '[JSON 数据]'
 * }
 */

const props = defineProps({
  value: { type: [String, Object, Array], default: null },
  column: { type: Object, default: () => ({}) },
});

const popoverVisible = ref(false);

// --- 计算属性 ---
const previewText = computed(() => {
  return props.column?.previewText || '';
});

const formattedJsonString = computed(() => {
  let dataToFormat = props.value;
  if (typeof props.value === 'string') {
    try {
      dataToFormat = JSON.parse(props.value);
    } catch (e) {
      return props.value;
    }
  }
  return JSON.stringify(dataToFormat, null, 2);
});

// --- 方法 ---
const copyJson = () => {
  navigator.clipboard.writeText(formattedJsonString.value).then(() => {
    ElMessage.success({
      message: '✨ JSON 数据已复制到剪贴板',
      duration: 2000,
      showClose: true
    });
    popoverVisible.value = false;
  }).catch(() => {
    ElMessage.error({
      message: '复制失败，请重试',
      duration: 2000
    });
  });
};
</script>

<style>
/* Popover 全局样式 */
.json-cell-popover {
  padding: 0 !important;
  border-radius: 12px !important;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04) !important;
  border: 1px solid #e5e7eb !important;
  overflow: hidden !important;
}

.dark .json-cell-popover {
  border: 1px solid #374151 !important;
  background-color: #1f2937 !important;
}

/* Popover 内容容器 */
.json-popover-content {
  background: linear-gradient(to bottom, #ffffff, #fafbfc);
  min-height: 200px;
  max-height: 800px;
}

.dark .json-popover-content {
  background: linear-gradient(to bottom, #1f2937, #111827);
}

/* 头部样式 */
.json-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(59, 130, 246, 0.05);
  border-bottom: 1px solid #e5e7eb;
  backdrop-filter: blur(10px);
  @apply py-2 px-4;
}

.dark .json-header {
  background: rgba(59, 130, 246, 0.1);
  border-bottom: 1px solid #374151;
}

/* JSON 内容容器 */
.json-viewer-container {
  max-height: 600px;
  overflow: auto;
  background-color: #f8fafc;
  position: relative;
  @apply p-3;
}

.dark .json-viewer-container {
  background-color: #0f172a;
}

/* 自定义滚动条 */
.json-viewer-container::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.json-viewer-container::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 4px;
}

.json-viewer-container::-webkit-scrollbar-thumb {
  background: linear-gradient(45deg, #cbd5e1, #94a3b8);
  border-radius: 4px;
}

.json-viewer-container::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(45deg, #94a3b8, #64748b);
}

.dark .json-viewer-container::-webkit-scrollbar-track {
  background: #1e293b;
}

.dark .json-viewer-container::-webkit-scrollbar-thumb {
  background: linear-gradient(45deg, #475569, #64748b);
}

/* JSON 代码样式 */
.json-pre {
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  color: #1e293b;
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  position: relative;
  overflow: hidden;
}

/* 按钮样式优化 */
.copy-btn {
  border-radius: 6px !important;
  transition: all 0.3s ease !important;
  font-weight: 500 !important;
}

.close-btn {
  color: #6b7280 !important;
  transition: all 0.2s ease !important;
  border-radius: 4px !important;
  padding: 4px !important;
}

.close-btn:hover {
  color: #ef4444 !important;
  background-color: #fef2f2 !important;
}

.dark .close-btn:hover {
  background-color: #450a0a !important;
}
</style>

<style scoped>
/* 触发元素样式 - 简洁版本，适配表格样式 */
.json-cell-preview {
  display: inline-flex;
  cursor: pointer;
  transition: all 0.2s ease;
}



.dark .preview-content {
  background: #334155;
  border-color: #475569;
  color: #cbd5e1;
}

.dark .json-cell-preview:hover .preview-content {
  background: #1e3a8a;
  border-color: #3b82f6;
}

.preview-icon {
  color: #64748b;
  font-size: 14px;
  transition: color 0.2s ease;
  flex-shrink: 0;
}

.json-cell-preview:hover .preview-icon {
  color: #0ea5e9;
}

.dark .preview-icon {
  color: #94a3b8;
}

.dark .json-cell-preview:hover .preview-icon {
  color: #60a5fa;
}

.preview-text {
  font-size: 12px;
  font-weight: 400;
  color: #475569;
  transition: color 0.2s ease;
  user-select: none;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 80px;
}

.json-cell-preview:hover .preview-text {
  color: #0c4a6e;
}

.dark .preview-text {
  color: #cbd5e1;
}

.dark .json-cell-preview:hover .preview-text {
  color: #dbeafe;
}

.expand-icon {
  color: #94a3b8;
  font-size: 10px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.json-cell-preview:hover .expand-icon {
  color: #0ea5e9;
  transform: translateX(1px);
}

.dark .expand-icon {
  color: #64748b;
}

.dark .json-cell-preview:hover .expand-icon {
  color: #60a5fa;
}
</style>

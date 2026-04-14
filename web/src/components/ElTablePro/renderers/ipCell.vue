<template>
  <el-tooltip
    :content="value"
    placement="top"
    :disabled="!isTruncated"
    :show-after="800"
  >
    <div
      class="group items-center space-x-2"
      :class="[isInTable ? 'flex w-full' : 'inline-flex']"
    >
      <!-- IP 地址文本 -->
      <span
        ref="textRef"
        class="flex-1 font-mono text-sm font-medium text-slate-700 dark:text-slate-300 truncate transition-colors duration-200 hover:text-sky-600 dark:hover:text-sky-600"
      >
        <!-- 版本标签 -->
        <span
          class="px-1 py-0.5 text-xs font-medium font-mono rounded border leading-none"
          :class="versionInfo.badgeClass"
        >
          {{ versionInfo.text }}
        </span>
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
import { ref, computed, onMounted, onUpdated } from 'vue';
import { ElIcon, ElTooltip, ElMessage } from 'element-plus';
import { CopyDocument } from '@element-plus/icons-vue';


/**
 * @file IpCell.vue
 * @description IP 地址单元格渲染器，智能识别 IPv4/IPv6，并提供复制功能。
 *              悬停时显示 Tooltip，溢出时自动截断。
 *
 * @usage
 * {
 *   prop: 'ipAddress',
 *   label: 'IP地址',
 *   type: 'ip',
 * }
 *
 * // 示例数据:
 * // { ipAddress: '192.168.1.1' } -> 显示 IPv4 标签和 IP 地址
 * // { ipAddress: '2001:0db8::1' } -> 显示 IPv6 标签和 IP 地址
 */

const props = defineProps({
  value: { type: String, default: null },
  column: { type: Object, default: () => ({}) },
});

const textRef = ref(null);
const isTruncated = ref(false);

const isInTable = computed(()=>{
  return Object.keys(props.column).length > 0
})

const checkTruncation = () => {
  requestAnimationFrame(() => {
    const el = textRef.value;
    if (el) {
      isTruncated.value = el.scrollWidth > el.offsetWidth;
    }
  });
};

onMounted(checkTruncation);
onUpdated(checkTruncation);

// 容器样式
const containerClass = computed(() => props.inline ? 'flex w-full' : 'flex w-full')

function getIpVersion(ip) {
  if (!ip || typeof ip !== 'string') return 'unknown';
  return ip.includes('.') ? 'v4' : 'v6';
}

const versionInfo = computed(() => {
  const version = getIpVersion(props.value);
  const styles = {
    v4: {
      text: 'v4',
      badgeClass: 'bg-blue-100 text-blue-700 border border-blue-200 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-700',
    },
    v6: {
      text: 'v6',
      badgeClass: 'bg-emerald-100 text-emerald-700 border border-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-300 dark:border-emerald-700',
    },
    unknown: {
      text: 'N/A',
      badgeClass: 'bg-gray-100 text-gray-600 border border-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600',
    },
  };
  return styles[version];
});

const handleCopy = () => {
  navigator.clipboard.writeText(props.value).then(() => {
    ElMessage.success({
      message: '🌐 IP 地址已复制到剪贴板',
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
</style>

<template>
  <div
    v-if="protocolInfo"
    class="group inline-flex items-center space-x-2 pl-2 pr-3 py-1 rounded-full
           font-medium transition-colors duration-200"
    :class="[protocolInfo.bgClass, protocolInfo.hoverBgClass]"
  >
    <span
      class="h-1.5 w-1.5 rounded-full"
      :class="protocolInfo.dotClass"
    />


    <span
      class="text-xs"
      :class="protocolInfo.textClass"
    >
      {{ protocolInfo.label || value }}
    </span>
  </div>

  <!-- 未知协议的默认显示 -->
  <span
    v-else
    class="text-sm text-slate-600 dark:text-slate-400"
  >{{ value }}</span>
</template>

<script setup>
import {computed} from 'vue';

/**
 * @file ProtocolCell.vue
 * @description 协议类型单元格渲染器，根据协议值显示带有不同颜色和标签的徽章。
 *              支持数字协议号和常见协议名称（如 TCP, UDP, HTTPS）。
 *
 * @usage
 * {
 *   prop: 'protocol',
 *   label: '协议',
 *   type: 'protocol',
 *
 *   // [可选] 提供一个 config 对象来自定义协议映射和样式
 *   config: {
 *     // 示例: 覆盖默认的 TCP 样式或添加新的协议
 *     'TCP': {
 *       label: '传输控制协议',
 *       bgClass: 'bg-purple-100 dark:bg-purple-900',
 *       dotClass: 'bg-purple-500 dark:bg-purple-400',
 *       textClass: 'text-purple-700 dark:text-purple-300'
 *     },
 *     '255': { // 自定义一个未知协议的显示
 *       label: '自定义协议',
 *       bgClass: 'bg-gray-200 dark:bg-gray-700',
 *       dotClass: 'bg-gray-500 dark:bg-gray-400',
 *       textClass: 'text-gray-700 dark:text-gray-300'
 *     }
 *   }
 * }
 *
 * // 示例数据:
 * // { protocol: 6 }    -> 显示 TCP 徽章
 * // { protocol: 'UDP' } -> 显示 UDP 徽章
 * // { protocol: 'HTTP' } -> 显示 HTTP 徽章
 */

const props = defineProps({
  value: {type: [String, Number], default: null},
  column: { type: Object, default: () => ({}) },
});

// --- 协议样式配置 (面向政企用户的专业色彩体系) ---
// 结构: { bgClass, hoverBgClass, dotClass, textClass }
// 理念: 淡色背景 + 深色文字 + 醒目圆点 = 清晰、专业的视觉层次
const defaultProtocols = {
  // --- 核心网络协议 (蓝色系，代表稳定、技术) ---
  1: {
    label: 'TCP',
    bgClass: 'bg-blue-50 dark:bg-blue-900',
    hoverBgClass: 'hover:bg-blue-100 dark:hover:bg-blue-800',
    dotClass: 'bg-blue-500 dark:bg-blue-400',
    textClass: 'text-blue-700 dark:text-blue-300'
  },
  2: {
    label: 'UDP',
    bgClass: 'bg-sky-50 dark:bg-sky-900',
    hoverBgClass: 'hover:bg-sky-100 dark:hover:bg-sky-800',
    dotClass: 'bg-sky-500 dark:bg-sky-400',
    textClass: 'text-sky-700 dark:text-sky-300'
  },
  3: {
    label: 'SCTP',
    bgClass: 'bg-indigo-50 dark:bg-indigo-900',
    hoverBgClass: 'hover:bg-indigo-100 dark:hover:bg-indigo-800',
    dotClass: 'bg-indigo-500 dark:bg-indigo-400',
    textClass: 'text-indigo-700 dark:text-indigo-300'
  },
  4: {
    label: 'ICMP',
    bgClass: 'bg-teal-50 dark:bg-teal-900',
    hoverBgClass: 'hover:bg-teal-100 dark:hover:bg-teal-800',
    dotClass: 'bg-teal-500 dark:bg-teal-400',
    textClass: 'text-teal-700 dark:text-teal-300'
  },
  5: {
    label: 'ARP',
    bgClass: 'bg-orange-50 dark:bg-orange-900',
    hoverBgClass: 'hover:bg-orange-100 dark:hover:bg-orange-800',
    dotClass: 'bg-orange-500 dark:bg-orange-400',
    textClass: 'text-orange-700 dark:text-orange-300'
  },


  'TCP': {
    label: 'TCP',
    bgClass: 'bg-blue-50 dark:bg-blue-900',
    hoverBgClass: 'hover:bg-blue-100 dark:hover:bg-blue-800',
    dotClass: 'bg-blue-500 dark:bg-blue-400',
    textClass: 'text-blue-700 dark:text-blue-300'
  },
  'UDP': {
    label: 'UDP',
    bgClass: 'bg-sky-50 dark:bg-sky-900',
    hoverBgClass: 'hover:bg-sky-100 dark:hover:bg-sky-800',
    dotClass: 'bg-sky-500 dark:bg-sky-400',
    textClass: 'text-sky-700 dark:text-sky-300'
  },

  // --- 安全/加密协议 (绿色系，代表安全、可信) ---
  'HTTPS': {
    label: 'HTTPS',
    bgClass: 'bg-green-50 dark:bg-green-900',
    hoverBgClass: 'hover:bg-green-100 dark:hover:bg-green-800',
    dotClass: 'bg-green-500 dark:bg-green-400',
    textClass: 'text-green-700 dark:text-green-300'
  },
  'TLS': {
    label: 'TLS',
    bgClass: 'bg-teal-50 dark:bg-teal-900',
    hoverBgClass: 'hover:bg-teal-100 dark:hover:bg-teal-800',
    dotClass: 'bg-teal-500 dark:bg-teal-400',
    textClass: 'text-teal-700 dark:text-teal-300'
  },
  'SSH': {
    label: 'SSH',
    bgClass: 'bg-emerald-50 dark:bg-emerald-900',
    hoverBgClass: 'hover:bg-emerald-100 dark:hover:bg-emerald-800',
    dotClass: 'bg-emerald-500 dark:bg-emerald-400',
    textClass: 'text-emerald-700 dark:text-emerald-300'
  },

  // --- 应用层协议 (靛蓝色/紫色系，代表应用、服务) ---
  'HTTP': {
    label: 'HTTP',
    bgClass: 'bg-violet-50 dark:bg-violet-900',
    hoverBgClass: 'hover:bg-violet-100 dark:hover:bg-violet-800',
    dotClass: 'bg-violet-500 dark:bg-violet-400',
    textClass: 'text-violet-700 dark:text-violet-300'
  },
  'DNS': {
    label: 'DNS',
    bgClass: 'bg-indigo-50 dark:bg-indigo-900',
    hoverBgClass: 'hover:bg-indigo-100 dark:hover:bg-indigo-800',
    dotClass: 'bg-indigo-500 dark:bg-indigo-400',
    textClass: 'text-indigo-700 dark:text-indigo-300'
  },
  'FTP': {
    label: 'FTP',
    bgClass: 'bg-orange-50 dark:bg-orange-900',
    hoverBgClass: 'hover:bg-orange-100 dark:hover:bg-orange-800',
    dotClass: 'bg-orange-500 dark:bg-orange-400',
    textClass: 'text-orange-700 dark:text-orange-300'
  }, // FTP 用橙色提示其非加密特性

  // --- 工具/其他协议 (灰色系，代表中性、基础) ---
  255: {
    label: '非传输层协议',
    bgClass: 'bg-slate-100 dark:bg-slate-700',
    hoverBgClass: 'hover:bg-slate-200 dark:hover:bg-slate-600',
    dotClass: 'bg-slate-500 dark:bg-slate-400',
    textClass: 'text-slate-700 dark:text-slate-300'
  },

};

// --- 计算属性，用于查找并应用样式 ---
const protocolInfo = computed(() => {
  // 合并默认配置和外部传入的自定义配置
  const effectiveConfig = {
    ...defaultProtocols,
    ...(props.column?.config || {}),
  };

  const key = String(props.value);
  const upperKey = typeof props.value === 'string' ? props.value.toUpperCase() : key;

  // 优先按原始值查找，再按大写查找
  const foundConfig = effectiveConfig[key] || effectiveConfig[upperKey];

  // 如果找到配置，则返回
  if (foundConfig) {
    return foundConfig;
  }

  // 如果找不到任何已知配置，则提供一个统一的、清晰的未知协议样式
  // 确保即使 value 是一个未知的长字符串，也能正常显示
  return {
    label: key, // 直接显示原始值
    bgClass: 'bg-slate-100 dark:bg-slate-800',
    hoverBgClass: 'hover:bg-slate-200 dark:hover:bg-slate-700',
    dotClass: 'bg-slate-400 dark:bg-slate-500',
    textClass: 'text-slate-700 dark:text-slate-300',
  };
});
</script>

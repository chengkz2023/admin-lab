import { defineAsyncComponent } from "vue";

// =============================================================================
//  1. 单元格渲染器注册中心
//  ######################### 可点击具体组件查看相关用法及特性 #######################
// =============================================================================
const cellRenderers = {
  default: defineAsyncComponent(() => import("./renderers/defaultCell.vue")), // 默认，无样式，展示原始数据，支持自定义 formatter
  status: defineAsyncComponent(() => import("./renderers/statusCell.vue")), // 状态
  progress: defineAsyncComponent(() => import("./renderers/progressCell.vue")), // 进度
  signal: defineAsyncComponent(() => import("./renderers/signalCell.vue")), // 信号量（可用于百分比展示）
  unit: defineAsyncComponent(() => import("./renderers/unitCell.vue")), // 单位 ['B', 'KB', 'MB', 'GB', 'TB']
  bps: defineAsyncComponent(() => import("./renderers/bpsCell.vue")), // 比特率 ['bps', 'Kbps', 'Mbps', 'Gbps', 'Tbps']
  Bps: defineAsyncComponent(() => import("./renderers/flowRateCell.vue")), // 流量速率 ['B/s', 'KB/s', 'MB/s', 'GB/s']
  pps: defineAsyncComponent(() => import("./renderers/ppsCell.vue")), // 包速率 ['pps', 'kpps', 'Mpps', 'Gpps']
  size: defineAsyncComponent(() => import("./renderers/sizeCell.vue")), // 字节 ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB']
  byte: defineAsyncComponent(() => import("./renderers/sizeCell.vue")), // 字节 ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB']
  tag: defineAsyncComponent(() => import("./renderers/tagCell.vue")), // 标签
  icon: defineAsyncComponent(() => import("./renderers/iconCell.vue")), // 图标
  switch: defineAsyncComponent(() => import("./renderers/switchCell.vue")), // 开关
  datetime: defineAsyncComponent(() => import("./renderers/datetimeCell.vue")), // 时间日期
  duration: defineAsyncComponent(() => import("./renderers/durationCell.vue")), // 持续时间
  delay: defineAsyncComponent(() => import("./renderers/durationCell.vue")), // 延迟
  compound: defineAsyncComponent(() => import("./renderers/compoundCell.vue")), // 复合单元格，在单个单元格内垂直展示多个子指标
  url: defineAsyncComponent(() => import("./renderers/urlCell.vue")), // 链接/地址
  ip: defineAsyncComponent(() => import("./renderers/ipCell.vue")), // IP
  json: defineAsyncComponent(() => import("./renderers/jsonCell.vue")), // JSON
  protocol: defineAsyncComponent(() => import("./renderers/protocolCell.vue")), // 协议
  port: defineAsyncComponent(() => import("./renderers/portCell.vue")), // 端口号
  id: defineAsyncComponent(() => import("./renderers/idCell.vue")), // ID
  msgType: defineAsyncComponent(() => import("./renderers/msgTypeCell.vue")), // msgType
  policyType: defineAsyncComponent(() => import("./renderers/msgTypeCell.vue")), // policyType
  options: defineAsyncComponent(() => import("./renderers/optionsCell.vue")), // options
  enumTag: defineAsyncComponent(() => import("./renderers/optionsCell.vue")), // backward compatibility
  base64: defineAsyncComponent(() => import("./renderers/base64Cell.vue")), // base64
  sync_status: defineAsyncComponent(() =>
    import("./renderers/syncStatusCell.vue"),
  ), // 同步状态
};

/**
 * 根据类型字符串解析并返回对应的渲染器组件。
 * @param {string} type - 列配置中的 type 字符串。
 * @returns {Component} 解析出的 Vue 组件。
 */
export function resolveCellComponent(type) {
  return cellRenderers[type] || cellRenderers.default;
}

// =============================================================================
//  2. 辅助函数
// =============================================================================

// 内置类型
const BUILT_IN_TYPES = ["selection", "index", "expand"];

/**
 * 判断是否为 El-Table 内置类型
 * @param {string} type - 列类型
 * @returns {boolean}
 */
export function isBuiltInType(type) {
  return BUILT_IN_TYPES.includes(type);
}

/**
 * 清洗 column 对象，并处理 compound 列的排序代理。
 * @param {object} column - 原始的列配置对象。
 * @returns {object} 只包含标准 props 的对象。
 */
export function getElColumnProps(column) {
  const {
    // 自定义属性黑名单
    config,
    formatter,
    type,
    children,
    items, // compound 列的 items 也需要被过滤掉
    previewText, // json 列的 previewText 也需要过滤
    sortBy, // compound 列的 sortBy 也需要被过滤掉
    ...standardProps
  } = column;

  standardProps.align = column.align || "left";

  // --- 预先处理 el-table 的内置类型 ---
  if (isBuiltInType(type)) {
    return {
      type,
      ...standardProps,
    };
  }

  // --- 这是新增的排序代理逻辑 ---
  if (type === "compound" && sortBy) {
    // 如果是 compound 列并且配置了 sortBy
    return {
      ...standardProps,
      // 直接将 sortable 和 sortBy 传递给 el-table-column
      sortable: true,
      sortBy: sortBy,
    };
  }

  return standardProps;
}

/**
 * 获取对齐方式的 CSS 类
 * @param {string} align - 对齐方式
 * @returns {string}
 */
export function getAlignClass(align) {
  switch (align) {
    case "center":
      return "justify-center";
    case "right":
      return "justify-end";
    default:
      return "justify-start";
  }
}

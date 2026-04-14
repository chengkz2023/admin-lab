<!--
TrendChart - 趋势图表组件

完整配置示例：
<TrendChart
  ref="chartRef"    // 图表引用，用于获取图表 defineExpose 的对象

  :data="[                                  // 数据数组（必填）
            { agg_time: '2024-01-01 00:00:00', cpu: 80, memory: 1.5 },
            { agg_time: '2024-01-01 01:00:00', cpu: 70, memory: 1.8 },
            { agg_time: '2024-01-01 02:00:00', cpu: 60, memory: 1.2 },
        ]"

  :series="{
    cpu: {                                  // 数据字段项配置
      label: 'CPU使用率',                   // 图例显示名称（推荐使用label，向下兼容name）
      type: 'line',                        // 图表类型：'line'|'bar', 默认 line
      yAxis: 'left',                       // Y轴位置：'left'|'right'，默认 left
      formatter: (val) => val + '%'        // 数值格式化函数
    },
    memory: {
      label: '内存占用',
      type: 'bar',
      yAxis: 'right',
      formatter: (val) => val + 'GB'
    },
  }"

  // 或者使用数组格式配置
  :series="[
    {
      name: 'cpu',                         // 数据字段名（必填）
      label: 'CPU使用率',                   // 图例显示名称
      type: 'line',                        // 图表类型：'line'|'bar', 默认 line
      yAxis: 'left',                       // Y轴位置：'left'|'right'，默认 left
      formatter: (val) => val + '%'        // 数值格式化函数
    },
    {
      name: 'memory',
      label: '内存占用',
      type: 'bar',
      yAxis: 'right',
      formatter: (val) => val + 'GB'
    },
  ]"

  x-axis-key="agg_time"                     // X轴索引字段名（必填）

  :x-axis-formatter="(val) => val + '时'"   // X轴格式化函数
  legend="top"                              // 图例位置：'top'|'bottom'|'none'，默认 top
  y-axis-name="数值"                        // 左Y轴名称
  y-axis-right-name="百分比"                // 右Y轴名称

  start-time="2024-01-01 00:00:00"         // 开始时间，用于时间对齐
  end-time="2024-01-01 23:59:59"           // 结束时间，用于时间对齐
  granularity="15m"                        // 时间粒度，用于时间对齐
  :fill="false"                            // 时间对齐时填充方式：true=开启填充，false=零值
  fill-type="mode"                         // 填充类型：'avg'|'median'|'mode'|'adjacent_avg'

  :grid="{ top: '2%', right: '2%', bottom: '2%', left: '2%' }"  // 网格配置
  :show-y-axis="false"                                          // 是否显示Y轴
  :x-axis-style="{ axisLine: { show: true } }"                                 // X轴样式配置
  :y-axis-style="{ splitLine: { show: true, lineStyle: { color: '#eee' } } }"  // Y轴样式配置

  :theme="'techBlue'"                      // 主题配置：'default'|'classic'|'vibrant'|'corporate'|'techBlue'|'forest'|'monochrome'

  :width="800"                             // 宽度
  :height="400"                            // 高度

  @registry-event="handleEventRegistry"  // 事件注册回调
/>

// 获取处理后的数据（开启了数据填充特性）
const processedChartData = computed(() => {
    return chartRef.value ? chartRef.value.processedData : []
});

// 处理事件注册
const handleEventRegistry = (chartInstance) => {
    console.log('Chart instance ready:', chartInstance);

    // 监听整个图表区域的点击事件
    chartInstance.getZr().on('click', function (params) {
        const pointInPixel = [params.offsetX, params.offsetY];

        // 判断是否在绘图区域内
        if (chartInstance.containPixel('grid', pointInPixel)) {
            // 像素坐标转数据坐标
            const pointInGrid = chartInstance.convertFromPixel({ seriesIndex: 0 }, pointInPixel);

            // 获取最接近的x轴索引
            const xIndex = Math.round(pointInGrid[0]);

            // 判断是否在x轴索引范围内，即是否在绘图区域内，如果在绘图区域内，则获取点击时间
            if (xIndex >= 0 && xIndex < processedChartData.value.length) {
                const clickedTime = processedChartData.value[xIndex].time;
                console.log('点击时间：', clickedTime);
                alert(`点击了时间点: ${clickedTime}`);
            }
        }
    });
};
-->

<script setup>
import {
  ref,
  computed,
  onMounted,
  onUnmounted,
  watch,
  markRaw,
  nextTick,
} from "vue";
import * as echarts from "echarts";
// 确保此文件路径正确，用于生成美观的配色方案
import { generateColorStyles } from "./generateColorStyles.js";
import dayjs from "dayjs";

// ==================== Props 定义 ====================
const props = defineProps({
  // 原始数据数组
  data: {
    type: Array,
    required: true,
    default: () => [],
  },
  // 系列配置：支持 Object { key: config } 或 Array [{ name: key, ... }]
  series: {
    type: [Object, Array],
    required: true,
    default: () => ({}),
  },
  // X轴对应的数据字段名 (通常是时间字段)
  xAxisKey: {
    type: String,
    required: true,
  },
  // X轴标签格式化函数
  xAxisFormatter: {
    type: Function,
    default: null,
  },
  // 图例位置
  legend: {
    type: String,
    default: "bottom",
    validator: (v) => ["top", "bottom", "none"].includes(v),
  },
  // Y轴名称 (左轴)
  yAxisName: { type: String, default: null },
  // Y轴名称 (右轴)
  yAxisRightName: { type: String, default: null },

  // --- 时间对齐相关配置 ---
  startTime: { type: [String, Date], default: null },
  endTime: { type: [String, Date], default: null },
  granularity: { type: String, default: null }, // 如 "10m", "1h"

  // --- 数据填充配置 ---
  fill: { type: Boolean, default: false }, // 是否开启填充
  fillType: {
    type: String,
    default: "avg",
    validator: (v) => ["avg", "median", "mode", "adjacent_avg"].includes(v),
  },

  // --- 样式与布局配置 ---
  grid: { type: Object, default: () => ({}) }, // ECharts Grid 覆盖
  showYAxis: { type: Boolean, default: true },
  xAxisStyle: { type: Object, default: () => ({}) },
  yAxisStyle: { type: Object, default: () => ({}) },

  // 主题预设
  theme: {
    type: String,
    default: "default",
    validator: (v) =>
      [
        "default",
        "classic",
        "vibrant",
        "corporate",
        "techBlue",
        "forest",
        "monochrome",
      ].includes(v),
  },

  // 容器尺寸
  width: { type: [String, Number], default: "100%" },
  height: { type: [String, Number], default: 500 },

  // 高级：直接覆盖 ECharts Option
  options: { type: Object, default: () => ({}) },
});

// ==================== 状态管理 ====================

const chartContainer = ref(null);
// 使用普通变量存储 ECharts 实例，或者使用 shallowRef。
// 注意：不要使用 ref()，因为 ECharts 实例非常庞大，Vue 的深度代理会造成严重的性能问题。
let myChart = null;
let resizeObserver = null;

const emit = defineEmits(["registry-event"]);

// ==================== 数据预处理 ====================

/**
 * 规范化 Series 配置
 * 将用户传入的 Array 或 Object 格式统一转换为 { key: config } 格式
 */
const processedSeries = computed(() => {
  if (Array.isArray(props.series)) {
    return props.series?.reduce((acc, item) => {
      const { name, ...rest } = item;
      acc[name] = {
        label: item.label || name,
        type: item.type || "line",
        yAxis: item.yAxis || "left",
        formatter: item.formatter,
        ...rest,
      };
      return acc;
    }, {});
  }

  if (typeof props.series === "object" && props.series !== null) {
    const result = {};
    for (const [key, config] of Object.entries(props.series)) {
      const { name, ...rest } = config;
      result[key] = {
        label: config.label || name || key, // 优先使用 label，其次 name，最后 key
        type: config.type || "line",
        yAxis: config.yAxis || "left",
        formatter: config.formatter,
        ...rest,
      };
    }
    return result;
  }
  return props.series;
});

// 对外暴露处理后的数据和 resize 方法
defineExpose({
  processedData: computed(() => processedChartData.value),
  resize: () => myChart?.resize(),
});

// ==================== 工具函数 (纯逻辑) ====================

/**
 * 数值归一化：将 null/undefined/NaN 转为 0
 */
const normalizeValue = (val) => {
  const num = Number(val);
  return val === null || val === undefined || Number.isNaN(num) ? 0 : num;
};

/**
 * 生成对齐的时间序列
 * 解析 granularity (如 "15m") 并生成 start 到 end 之间的所有时间点
 */
const generateAlignedTimeSeries = (start, end, granularity) => {
  const match = granularity.match(/^(\d+)([smhd])$/);
  if (!match) return [];
  const align = parseInt(match[1]);
  // 映射 dayjs 单位
  const unit = { s: "second", m: "minute", h: "hour", d: "day" }[match[2]];

  const alignTime = (time) => {
    const dt = dayjs(time);
    if (align === 1) return dt.startOf(unit);
    // 对齐到整倍数逻辑 (例如 10:23 对齐到 10:20，如果 align=10m)
    const current = dt.get(unit);
    const aligned = Math.floor(current / align) * align;
    let baseTime;
    switch (unit) {
      case "second":
        baseTime = dt.startOf("minute");
        break;
      case "minute":
        baseTime = dt.startOf("hour");
        break;
      case "hour":
        baseTime = dt.startOf("day");
        break;
      case "day":
        baseTime = dt.startOf("month");
        break;
      default:
        baseTime = dt.startOf(unit);
    }
    return baseTime.add(aligned, unit);
  };

  const startTime = alignTime(start);
  const endTime = alignTime(end);
  const result = [];
  let current = startTime;

  while (current.isBefore(endTime) || current.isSame(endTime)) {
    result.push(current.format("YYYY-MM-DD HH:mm:ss"));
    current = current.add(align, unit);
  }
  return result;
};

// --- 样式生成器辅助函数 ---

const createDefaultStyleConfig = () => ({
  gradientTopAlpha: 0.2,
  gradientBottomAlpha: 0.05,
  shadowBlur: 8,
  optimizeForMixed: true,
  palette: props.theme || "default",
  enableModernTooltip: true,
  enableModernAxes: true,
});

/**
 * 为不同图表类型生成优化后的样式
 * 柱状图：增强立体感和渐变
 * 线图：增强高亮交互
 */
const createOptimizedStyle = (type, autoStyle) => {
  const [r, g, b] = autoStyle.lineColor.match(/\d+/g).map(Number);
  if (type === "bar") {
    return {
      ...autoStyle,
      itemStyle: {
        // 垂直渐变增强立体感
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: `rgba(${r}, ${g}, ${b}, 0.9)` },
          { offset: 1, color: `rgba(${r}, ${g}, ${b}, 0.7)` },
        ]),
        borderColor: `rgba(${r}, ${g}, ${b}, 0.8)`,
        borderWidth: 1,
        borderRadius: [4, 4, 0, 0], // 顶部圆角
        shadowColor: `rgba(${r}, ${g}, ${b}, 0.3)`,
        shadowBlur: 8,
        shadowOffsetY: 2,
      },
      emphasis: {
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: `rgba(${r}, ${g}, ${b}, 1)` },
            { offset: 1, color: `rgba(${r}, ${g}, ${b}, 0.8)` },
          ]),
          borderColor: `rgba(${r}, ${g}, ${b}, 1)`,
          borderWidth: 2,
          shadowBlur: 12,
          shadowOffsetY: 4,
        },
      },
    };
  }
  return {
    ...autoStyle,
    emphasis: {
      itemStyle: {
        borderColor: "rgba(0,0,0,0.2)",
        borderWidth: 8,
      },
    },
  };
};

/**
 * 构建单个 Series 配置项
 * 将 props 配置、自动生成的颜色、优化样式合并
 */
const createSeriesItem = (key, index, seriesConfig, autoColors) => {
  const cfg = seriesConfig[key];
  const {
    label,
    type,
    yAxis,
    formatter,
    emphasis: userEmphasis,
    ...rest
  } = cfg;

  const seriesValues = processedChartData.value.map((item) =>
    normalizeValue(item[key])
  );

  const finalType = type || "line";
  const autoColorStyle = autoColors[index];
  const optimizedStyle = createOptimizedStyle(finalType, autoColorStyle);

  return {
    name: label || key,
    data: seriesValues,
    type: finalType,
    yAxisIndex: yAxis === "right" ? 1 : 0, // 双轴逻辑
    smooth: finalType === "line", // 仅线图开启平滑
    symbol: finalType === "line" ? "circle" : "none",
    symbolSize: 8,
    showSymbol: false, // 默认不显示点，hover时显示
    formatter: formatter,
    emphasis: {
      focus: "series", // hover时淡出其他系列
      ...optimizedStyle.emphasis,
      ...userEmphasis,
    },
    lineStyle:
      finalType === "line"
        ? {
            ...optimizedStyle.lineStyle,
            width: 2,
            ...(rest.lineStyle || {}),
          }
        : undefined,
    itemStyle: {
      ...optimizedStyle.itemStyle,
      ...(rest.itemStyle || {}),
    },
    areaStyle:
      rest.areaStyle !== undefined
        ? rest.areaStyle
        : finalType === "line"
        ? optimizedStyle.areaStyle // 自动应用渐变区域填充
        : undefined,
    ...(finalType === "bar"
      ? { barWidth: "60%", barMaxWidth: 50, barGap: "20%" }
      : {}),
    ...rest, // 允许用户覆盖任意属性
  };
};

/**
 * 生成 Y 轴配置
 * 智能判断单轴/双轴，并应用格式化
 */
const createYAxisConfig = (
  series,
  yAxisName,
  yAxisRightName,
  showYAxis = true
) => {
  const yAxisConfig = [];
  // 如果无数据或隐藏，生成隐藏轴配置以占位
  if (!showYAxis || !hasData.value) {
    yAxisConfig.push({ type: "value", show: false });
    const usesSecondYAxis = series.some((s) => s.yAxisIndex === 1);
    if (usesSecondYAxis) {
      yAxisConfig.push({ type: "value", show: false });
    }
    return yAxisConfig;
  }

  const isSingleSeries = series.length === 1;
  // 场景1：单系列，直接使用该系列的 formatter
  if (isSingleSeries) {
    const singleSeries = series[0];
    let singleYAxisFormatter = "{value}";
    if (typeof singleSeries.formatter === "function") {
      singleYAxisFormatter = (value) => singleSeries.formatter(value);
    }
    yAxisConfig.push({
      type: "value",
      name: yAxisName || undefined,
      axisLabel: {
        formatter: singleYAxisFormatter,
        color: "#666",
      },
    });
  } else {
    // 场景2：多系列，左右轴分离
    const usesFirstYAxisItem = series.find((s) => s.yAxisIndex === 0);
    if (usesFirstYAxisItem) {
      yAxisConfig.push({
        type: "value",
        name: yAxisName || undefined,
        position: "left",
        axisLabel: {
          formatter: usesFirstYAxisItem.formatter || "{value}",
          color: "#666",
        },
      });
    }
    const usesSecondYAxisItem = series.find((s) => s.yAxisIndex === 1);
    if (usesSecondYAxisItem) {
      yAxisConfig.push({
        type: "value",
        name: yAxisRightName || undefined,
        position: "right",
        axisLabel: {
          formatter: usesSecondYAxisItem.formatter || "{value} %", // 右轴默认百分比
          color: "#666",
        },
      });
    }
  }
  // 统一应用刻度对齐
  yAxisConfig.forEach((axis) => {
    axis.alignTicks = true;
    axis.axisLine = { show: true };
  });
  return yAxisConfig;
};

// --- 统计计算函数 ---

const calculateMode = (values) => {
  if (values.length === 0) return 0;
  const frequency = {};
  values.forEach((value) => {
    frequency[value] = (frequency[value] || 0) + 1;
  });
  let maxFrequency = 0;
  let mode = values[0];
  for (const [value, freq] of Object.entries(frequency)) {
    if (freq > maxFrequency) {
      maxFrequency = freq;
      mode = Number(value);
    }
  }
  return mode;
};

const calculateMedian = (values) => {
  if (values.length === 0) return 0;
  const sorted = [...values].sort((a, b) => a - b);
  const middle = Math.floor(sorted.length / 2);
  if (sorted.length % 2 === 0) {
    return (sorted[middle - 1] + sorted[middle]) / 2;
  } else {
    return sorted[middle];
  }
};

const calculateAverage = (values) => {
  if (values.length === 0) return 0;
  return values.reduce((sum, val) => sum + val, 0) / values.length;
};

const calculateFillValue = (values, fillType) => {
  if (values.length === 0) return 0;
  switch (fillType) {
    case "median":
      return calculateMedian(values);
    case "mode":
      return calculateMode(values);
    case "adjacent_avg":
      // 邻点均值模式下，作为无法找到邻点时的回退值
      return calculateAverage(values);
    case "avg":
    default:
      return calculateAverage(values);
  }
};

const isFillDebugEnabled = () =>
  Boolean(import.meta.env.DEV) ||
  (typeof window !== "undefined" &&
    window.__TREND_CHART_DEBUG_FILL__ === true);

const createFillPreviewRows = (data, xAxisKey, seriesKeys, limit = 5) => {
  if (!Array.isArray(data)) return [];
  return data.slice(0, limit).map((item) => {
    const preview = { [xAxisKey]: item?.[xAxisKey] };
    seriesKeys.forEach((key) => {
      preview[key] = item?.[key];
    });
    return preview;
  });
};

const toValidNumberOrNull = (value) => {
  if (value === null || value === undefined) return null;
  const num = Number(value);
  return Number.isNaN(num) ? null : num;
};

const buildAdjacentFillLookup = (seriesKeys, timeSeries, dataMap) => {
  const lookup = {};
  for (let j = 0; j < seriesKeys.length; j++) {
    const key = seriesKeys[j];
    const prevValues = new Array(timeSeries.length).fill(null);
    const nextValues = new Array(timeSeries.length).fill(null);

    let prev = null;
    for (let i = 0; i < timeSeries.length; i++) {
      const item = dataMap.get(timeSeries[i]);
      const current = toValidNumberOrNull(item?.[key]);
      if (current !== null) prev = current;
      prevValues[i] = prev;
    }

    let next = null;
    for (let i = timeSeries.length - 1; i >= 0; i--) {
      const item = dataMap.get(timeSeries[i]);
      const current = toValidNumberOrNull(item?.[key]);
      if (current !== null) next = current;
      nextValues[i] = next;
    }

    lookup[key] = { prevValues, nextValues };
  }
  return lookup;
};

const resolveAdjacentFillValue = (adjacentLookup, key, index, fallbackValue) => {
  const cache = adjacentLookup?.[key];
  if (!cache) return fallbackValue;
  const prev = cache.prevValues[index];
  const next = cache.nextValues[index];
  if (prev !== null && next !== null) return (prev + next) / 2;
  if (prev !== null) return prev;
  if (next !== null) return next;
  return fallbackValue;
};

/**
 * 核心数据处理：时间对齐与填充
 * 根据 granularity 生成完整时间轴，并填充缺失数据
 */
const processTimeAlignedData = (originalData) => {
  const debugEnabled = isFillDebugEnabled();
  if (!props.startTime || !props.endTime || !props.granularity) {
    if (debugEnabled) {
      console.warn("[TrendChart][fill-debug] 跳过填充：时间对齐参数不完整", {
        startTime: props.startTime,
        endTime: props.endTime,
        granularity: props.granularity,
      });
    }
    return originalData;
  }
  try {
    const timeSeries = generateAlignedTimeSeries(
      props.startTime,
      props.endTime,
      props.granularity
    );
    if (timeSeries.length === 0) {
      if (debugEnabled) {
        console.warn(
          "[TrendChart][fill-debug] 跳过填充：生成的时间序列为空",
          {
            startTime: props.startTime,
            endTime: props.endTime,
            granularity: props.granularity,
          }
        );
      }
      return originalData;
    }

    // 使用标准化后的 series 键，兼容 Object/Array 两种传参格式
    const seriesKeys = processedSeries.value
      ? Object.keys(processedSeries.value)
      : [];
    if (seriesKeys.length === 0) {
      if (debugEnabled) {
        console.warn("[TrendChart][fill-debug] 跳过填充：seriesKeys 为空");
      }
      return originalData;
    }

    const originalRows = Array.isArray(originalData) ? originalData : [];
    if (debugEnabled) {
      console.groupCollapsed("[TrendChart][fill-debug] 输入参数");
      console.log("config", {
        fill: props.fill,
        fillType: props.fillType,
        xAxisKey: props.xAxisKey,
        startTime: props.startTime,
        endTime: props.endTime,
        granularity: props.granularity,
        inputLength: originalRows.length,
        alignedTimeCount: timeSeries.length,
        seriesKeys,
      });
      console.log(
        "inputPreview",
        createFillPreviewRows(originalRows, props.xAxisKey, seriesKeys)
      );
      console.groupEnd();
    }

    // 建立时间索引 Map
    const dataMap = new Map();
    originalRows.forEach((item) => {
      dataMap.set(
        dayjs(item[props.xAxisKey]).format("YYYY-MM-DD HH:mm:ss"),
        item
      );
    });

    // 计算填充基准值
    let seriesFillValues = {};
    const seriesValues = {};
    seriesKeys.forEach((key) => (seriesValues[key] = []));
    for (let i = 0; i < originalRows.length; i++) {
      const item = originalRows[i];
      for (let j = 0; j < seriesKeys.length; j++) {
        const key = seriesKeys[j];
        const value = toValidNumberOrNull(item[key]);
        if (value !== null) seriesValues[key].push(value);
      }
    }
    for (let j = 0; j < seriesKeys.length; j++) {
      const key = seriesKeys[j];
      seriesFillValues[key] = calculateFillValue(
        seriesValues[key],
        props.fillType
      );
    }
    const adjacentFillLookup =
      props.fillType === "adjacent_avg"
        ? buildAdjacentFillLookup(seriesKeys, timeSeries, dataMap)
        : null;
    const getFillValue = (key, index) => {
      if (!props.fill) return 0;
      const fallbackValue = seriesFillValues[key] || 0;
      if (props.fillType !== "adjacent_avg") return fallbackValue;
      return resolveAdjacentFillValue(
        adjacentFillLookup,
        key,
        index,
        fallbackValue
      );
    };

    if (debugEnabled) {
      const fillStats = seriesKeys.reduce((acc, key) => {
        acc[key] = {
          validCount: seriesValues[key].length,
          fillType: props.fillType,
          fillValue:
            props.fillType === "adjacent_avg"
              ? "dynamic(adjacent)"
              : seriesFillValues[key],
          fallbackValue:
            props.fillType === "adjacent_avg"
              ? seriesFillValues[key]
              : undefined,
          rawSample: seriesValues[key].slice(0, 5),
        };
        return acc;
      }, {});
      console.log("[TrendChart][fill-debug] 填充值统计", fillStats);
    }

    // 生成对齐数据
    const alignedData = new Array(timeSeries.length);
    let filledByTimeGapCount = 0;
    let patchedFieldCount = 0;
    const filledByTimeGapExamples = [];
    const patchedFieldExamples = [];
    for (let i = 0; i < timeSeries.length; i++) {
      const timestamp = timeSeries[i];
      const existingData = dataMap.get(timestamp);
      if (existingData) {
        // 对命中时间点做字段级补齐，避免 null/undefined 被后续渲染强制归零
        const mergedItem = { ...existingData, [props.xAxisKey]: timestamp };
        for (let j = 0; j < seriesKeys.length; j++) {
          const key = seriesKeys[j];
          const value = mergedItem[key];
          const isMissingOrInvalid =
            value === null || value === undefined || Number.isNaN(Number(value));
          if (isMissingOrInvalid) {
            const fillResult = getFillValue(key, i);
            mergedItem[key] = fillResult;
            patchedFieldCount += 1;
            if (debugEnabled && patchedFieldExamples.length < 5) {
              patchedFieldExamples.push({
                time: timestamp,
                key,
                before: value,
                after: fillResult,
                fillType: props.fillType,
              });
            }
          }
        }
        alignedData[i] = mergedItem;
      } else {
        const newItem = { [props.xAxisKey]: timestamp };
        for (let j = 0; j < seriesKeys.length; j++) {
          const key = seriesKeys[j];
          // fill=true 时按 fillType 填充，fill=false 时补零
          newItem[key] = getFillValue(key, i);
        }
        filledByTimeGapCount += 1;
        if (debugEnabled && filledByTimeGapExamples.length < 5) {
          filledByTimeGapExamples.push(
            createFillPreviewRows([newItem], props.xAxisKey, seriesKeys, 1)[0]
          );
        }
        alignedData[i] = newItem;
      }
    }
    if (debugEnabled) {
      console.groupCollapsed("[TrendChart][fill-debug] 输出结果");
      console.log("summary", {
        inputLength: originalRows.length,
        outputLength: alignedData.length,
        filledByTimeGapCount,
        patchedFieldCount,
        fillType: props.fillType,
      });
      console.log("filledByTimeGapExamples", filledByTimeGapExamples);
      console.log("patchedFieldExamples", patchedFieldExamples);
      console.log(
        "outputPreview",
        createFillPreviewRows(alignedData, props.xAxisKey, seriesKeys)
      );
      console.groupEnd();
    }
    return alignedData;
  } catch (error) {
    console.error("时间对齐处理失败:", error);
    return originalData;
  }
};

// ==================== 主要计算属性 (ECharts 配置生成) ====================

const chartContainerStyle = computed(() => {
  const formatDimension = (value) => {
    if (typeof value === "number") {
      return `${value}px`;
    }
    return String(value);
  };
  return {
    width: formatDimension(props.width),
    height: formatDimension(props.height),
  };
});

const createLegendConfig = (legendPosition, series) => {
  const baseConfig = {
    itemWidth: 18,
    itemHeight: 12,
    itemGap: 14,
    textStyle: { color: "#666", fontSize: 12, padding: [2, 0, 2, 4] },
    data: series.map((item) => item.name),
  };
  switch (legendPosition) {
    case "top":
      return {
        ...baseConfig,
        type: "scroll",
        orient: "horizontal",
        top: "0%",
        left: "center",
      };
    case "bottom":
      return {
        ...baseConfig,
        type: "scroll",
        orient: "horizontal",
        bottom: "0%",
        left: "center",
      };
    case "none":
      return { show: false };
    default:
      return {
        ...baseConfig,
        type: "scroll",
        orient: "horizontal",
        bottom: "2%",
        left: "center",
      };
  }
};

const formatXAxisDisplayValue = (value) => {
  if (props.xAxisFormatter && typeof props.xAxisFormatter === "function") {
    return props.xAxisFormatter(value);
  }

  const parsed = dayjs(value);
  if (parsed.isValid()) {
    return parsed.format("MM-DD HH:mm:ss");
  }

  return value;
};

// 触发数据对齐处理
const processedChartData = computed(() => {
  return processTimeAlignedData(props.data);
});

const hasData = computed(() => {
  const data = processedChartData.value;
  if (!data || data.length === 0) return false;
  return data.some((item) =>
    Object.keys(item).some(
      (key) => key !== props.xAxisKey && Number(item[key]) > 0
    )
  );
});

/**
 * 组装最终的 ECharts Option
 * 这是整个组件的核心，将所有逻辑汇聚于此
 */
const processedOptions = computed(() => {
  if (
    !processedChartData.value?.length ||
    !props.xAxisKey ||
    !processedSeries.value
  ) {
    return {};
  }

  const {
    xAxisKey,
    xAxisFormatter,
    legend,
    showYAxis,
    yAxisName,
    yAxisRightName,
  } = props;
  const seriesConfig = processedSeries.value;
  const xAxisData = processedChartData.value.map((item) => item[xAxisKey]);
  const seriesKeys = Object.keys(seriesConfig);

  // 1. 生成基础样式配置
  const defaultStyleConfig = createDefaultStyleConfig();
  // 2. 根据系列数量生成调色盘
  const autoColors = generateColorStyles(seriesKeys.length, {
    gradientTopAlpha: defaultStyleConfig.gradientTopAlpha,
    gradientBottomAlpha: defaultStyleConfig.gradientBottomAlpha,
    shadowAlpha: 0.3,
    shadowBlur: defaultStyleConfig.shadowBlur,
    lineAlpha: 0.8,
    useECharts: true,
    palette: defaultStyleConfig.palette,
    optimizeForMixed: defaultStyleConfig.optimizeForMixed,
  });

  // 3. 构建 Series
  const series = seriesKeys.map((key, index) =>
    createSeriesItem(key, index, seriesConfig, autoColors)
  );

  // 4. 构建 Axes
  const yAxisConfig = createYAxisConfig(
    series,
    yAxisName,
    yAxisRightName,
    props.showYAxis
  );

  // 5. 构建 Legend
  const legendConfig = createLegendConfig(legend || "bottom", series);

  // 6. 构建现代化 Tooltip (使用 HTML 自定义)
  const modernTooltip = {
    trigger: "axis",
    confine: true,
    axisPointer: { type: "line", lineStyle: { color: "#888", type: "dashed" } },
    backgroundColor: "transparent",
    borderColor: "transparent",
    borderWidth: 0,
    padding: 0,
    textStyle: { fontFamily: "Arial, sans-serif", fontSize: 14 },
    extraCssText:
      "box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1); border-radius: 8px;",
    formatter: (params) => {
      if (!params || params.length === 0) return "";
      const xAxisLabel = formatXAxisDisplayValue(params[0].axisValueLabel);
      let listItems = "";
      params.forEach((param) => {
        const seriesConfig = series[param.seriesIndex];
        const value = param.value;
        const safeVal = normalizeValue(value);
        let formattedValue = safeVal;
        if (typeof seriesConfig.formatter === "function") {
          formattedValue = seriesConfig.formatter(safeVal);
        }
        listItems += `<div style="display: flex; gap: 6px; justify-content: space-between; align-items: center; margin-top: 4px; color: #333;">
                                <span>${param.marker} ${seriesConfig.name}</span>
                                <b style="font-weight: 600;">${formattedValue}</b>
                              </div>`;
      });
      return `<div style="overflow: hidden; border-radius: 8px;">
                      <div style="background-color: #e5e7eb; padding: 4px 8px; font-size: 12px; color: #4a5568; text-align: left;">${xAxisLabel}</div>
                      <div style="background-color: #fff; padding: 4px 8px 6px;">${listItems}</div>
                    </div>`;
    },
  };

  const modernXAxisStyle = {
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: { color: "#666" },
    splitLine: { show: true, lineStyle: { color: "#eee", type: "dashed" } },
    ...props.xAxisStyle,
  };

  const modernYAxisStyle = {
    axisLine: { show: true, lineStyle: { color: "#ccc" } },
    splitLine: { show: true, lineStyle: { color: "#eee", type: "dashed" } },
    ...props.yAxisStyle,
  };

  // 7. 返回完整 Option
  return {
    tooltip: modernTooltip,
    legend: legendConfig,
    grid: {
      left: showYAxis ? "2%" : "0%",
      right: "3%",
      top: legend === "top" ? "10%" : "4%",
      bottom: legend === "bottom" || !legend ? "10%" : "4%",
      containLabel: true,
      ...props.grid,
    },
    xAxis: [
      {
        type: "category",
        boundaryGap: series.some((s) => s.type === "bar"), // 柱状图时开启间隙
        data: xAxisData,
        ...modernXAxisStyle,
        axisLabel: {
          ...modernXAxisStyle.axisLabel,
          formatter: (value, index) => {
            if (index === 0 || index === xAxisData.length - 1) return "";
            return formatXAxisDisplayValue(value);
          },
        },
      },
    ],
    yAxis: yAxisConfig.map((axis) => ({
      ...axis,
      ...modernYAxisStyle,
    })),
    series: series,
    ...(Object.keys(props.options).length > 0 && props.options),
  };
});

// ====================
// 生命周期与渲染逻辑 (解决 v-show/display:none 问题)
// ====================

/**
 * 尝试初始化图表
 *
 * 原理：
 * 当组件使用 v-show 隐藏时，DOM 存在但宽高为 0。
 * 此时如果调用 echarts.init 会导致 "Can't get DOM width or height" 报错。
 * 因此，我们必须做一个守卫，只有当宽高 > 0 时才进行初始化。
 */
const initChart = () => {
  // 1. 如果已经初始化，直接跳过，避免重复创建实例
  if (myChart) return;

  // 2. 防御性检查 DOM 引用
  if (!chartContainer.value) return;

  // 3. 关键：检查尺寸。如果是隐藏状态 (display: none)，clientWidth 为 0
  const width = chartContainer.value.clientWidth;
  const height = chartContainer.value.clientHeight;

  if (width === 0 || height === 0) {
    // console.warn("Chart container hidden, skipping init...");
    return;
  }

  // 4. 初始化实例
  // markRaw 标记该对象，告诉 Vue "不要把这个对象变成响应式的"
  // ECharts 实例内部有极其复杂的循环引用和大量属性，Vue 代理它会导致极大的性能开销甚至栈溢出
  const instance = echarts.init(chartContainer.value);
  myChart = markRaw(instance);

  // 5. 初始化后立即应用配置
  if (processedOptions.value) {
    myChart.setOption(processedOptions.value, true);
    emit("registry-event", myChart);
  }
};

/**
 * 统一的 ResizeObserver 回调
 *
 * 功能：
 * 1. 自动处理浏览器窗口缩放。
 * 2. 自动处理 v-show 从 false 变为 true 时的初始化（此时 width 从 0 变为有值）。
 */
const handleResizeEntry = (entries) => {
  for (const entry of entries) {
    const { width, height } = entry.contentRect;

    // 只有当尺寸有效（变为可见）时才处理
    if (width > 0 && height > 0) {
      if (!myChart) {
        // 场景：从 v-show=false 变为 true，首次初始化
        initChart();
      } else {
        // 场景：浏览器缩放或容器大小改变，调用 ECharts 自身的 resize
        myChart.resize();
      }
    }
  }
};

// 监听配置变化
watch(
  processedOptions,
  (newOptions) => {
    // 只有当图表已初始化（即当前是可见状态）时才更新
    if (myChart) {
      myChart.setOption(newOptions, true);
      emit("registry-event", myChart);
    }
    // 注意：如果 myChart 为空（当前隐藏），无需操作。
    // 等 ResizeObserver 检测到变为可见时，initChart 会自动读取最新的 options 进行渲染。
  },
  { deep: true }
);

onMounted(() => {
  if (chartContainer.value) {
    // 1. 启动 ResizeObserver
    // 这将自动处理初始化（当变为可见时）和重绘（当尺寸变化时）
    if (window.ResizeObserver) {
      resizeObserver = new ResizeObserver(handleResizeEntry);
      resizeObserver.observe(chartContainer.value);
    }

    // 2. 尝试立即初始化
    // 如果组件一开始就是可见的，ResizeObserver 可能有微小的延迟，立即调用可确保无闪烁
    initChart();
  }
});

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect();
    resizeObserver = null;
  }
  if (myChart) {
    myChart.dispose();
    myChart = null;
  }
});
</script>

<template>
  <div ref="chartContainer" :style="chartContainerStyle" />
</template>

<style scoped>
/* 可以在这里添加组件特定的样式，但目前主要依赖 ECharts 内部样式 */
</style>

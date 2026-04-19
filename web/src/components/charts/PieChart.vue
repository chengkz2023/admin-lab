<!--
PieChart - 饼图/柱状图组件

完整配置示例：
```vue
后端数据示例：
基础：
[
{'移动端': 50},
{'桌面': 80}
]
当需要做国际化：
[
{mobile: 50},
{desktop: 80}
...
]


<PieChart
  :data="salesData"              // 数据数组（必需）
  chart-type="bar"               // 图表类型：'pie'|'bar'
  name-key="category"            // 名称字段，默认 'name'
  value-key="amount"             // 数值字段，默认 'value'
  :show-pie-percent="true"       // 饼图是否显示百分比
  :show-bar-value="false"        // 柱状图是否在顶部显示数值
  legend="bottom"                // 图例位置：'right'|'bottom'|'top'|'left'|'none'
  :name-map="{                   // 字段名称映射（可选）
    mobile: '移动设备',
    desktop: '桌面设备',
    tablet: '平板设备'
  }"
  :style-config="{
    theme: 'vibrant'             // 主题选择
  }"
  :width="800"                   // 宽度
  :height="400"                  // 高度
  :show-y-axis="true"            // 显示Y轴（仅柱状图有效）
  :formatter="formatTraffic"     // 自定义数值格式化函数（可选）
  @chart-click="handleClick"     // 点击事件处理
/>

// 自定义格式化函数示例
const formatTraffic = (value) => {
  if (value >= 1024 * 1024 * 1024) {
    return (value / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
  } else if (value >= 1024 * 1024) {
    return (value / (1024 * 1024)).toFixed(2) + ' MB';
  } else if (value >= 1024) {
    return (value / 1024).toFixed(2) + ' KB';
  }
  return value + ' B';
};
```
-->

<script setup lang="ts">
import {
  ref,
  computed,
  onMounted,
  onUnmounted,
  watch,
  nextTick,
  markRaw,
} from "vue";
import * as echarts from "echarts";
import type { EChartsOption } from "echarts";
import { generateColorStyles } from "./generateColorStyles.js";

// 1. 声明组件可以发出的事件
const emit = defineEmits(["chart-click"]);

// 定义组件的属性
const props = defineProps({
  data: {
    type: Array,
    required: true,
    default: () => [],
  },
  // 数据字段配置
  nameKey: {
    type: String,
    default: "name",
  },
  valueKey: {
    type: String,
    default: "value",
  },
  // 显示配置
  showPiePercent: {
    type: Boolean,
    default: true,
  },
  legend: {
    type: String,
    default: "right",
    validator: (value) =>
      ["right", "bottom", "top", "left", "none"].includes(value),
  },
  // 字段名称映射配置 - 可选
  nameMap: {
    type: Object,
    default: () => ({}),
  },
  // 样式配置
  styleConfig: {
    type: Object,
    default: () => ({
      theme: "default",
    }),
  },
  // 图表尺寸配置 - 可选
  width: {
    type: [String, Number],
    default: "100%", // 默认宽度为100%
  },
  height: {
    type: [String, Number],
    default: 500, // 默认高度为500px
  },
  // 图表类型配置 - 可选
  chartType: {
    type: String,
    default: "pie", // 默认为饼图
    validator: (value) => ["pie", "bar"].includes(value),
  },
  nameMaxLength: {
    type: Number,
    default: 12,
  },
  // 是否显示Y轴 - 仅对柱状图有效
  showYAxis: {
    type: Boolean,
    default: false, // 默认不显示Y轴，保持现有行为
  },
  // 是否在柱状图顶部显示数值 - 仅对柱状图有效
  showBarValue: {
    type: Boolean,
    default: true, // 默认显示柱状图顶部数值
  },
  barDirection: {
    type: String,
    default: "vertical",
    validator: (value) => ["vertical", "horizontal"].includes(value),
  },
  // 数值格式化函数 - 用于自定义数值显示格式
  formatter: {
    type: Function,
    default: null, // 默认为 null,使用内置格式化
  },
});

const chartContainer = ref<HTMLElement | null>(null);
// 使用普通变量存储，避免 Vue 深度代理导致的性能问题
let myChart: echarts.ECharts | null = null;
let resizeObserver: ResizeObserver | null = null;

/**
 * 图表容器样式计算属性
 */
const chartContainerStyle = computed(() => {
  const formatDimension = (value: string | number): string => {
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

// ====================
// 工具函数区域 (保持原样)
// ====================

const getDisplayName = (item: any, index: number, nameKey: string): string => {
  const fieldValue = item[nameKey];
  if (props.nameMap && fieldValue && props.nameMap[fieldValue]) {
    return props.nameMap[fieldValue];
  }
  if (fieldValue) {
    return fieldValue;
  }
  return `项目${index + 1}`;
};

const getTruncatedDisplayName = (
  item: any,
  index: number,
  nameKey: string,
  maxLength: number = 12,
): string => {
  const fullName = getDisplayName(item, index, nameKey);
  if (fullName.length > maxLength) {
    return fullName.substring(0, maxLength) + "...";
  }
  return fullName;
};

// ====================
// 样式生成函数 (保持原样，保留你的渐变和阴影逻辑)
// ====================

const createPieStyleConfig = () => ({
  gradientTopAlpha: 0.8,
  gradientBottomAlpha: 0.6,
  shadowBlur: 10,
  optimizeForMixed: true,
  palette: props.styleConfig.theme || "default",
  enableModernTooltip: true,
  ...props.styleConfig,
});

const createPieDataItem = (
  item: any,
  index: number,
  nameKey: string,
  valueKey: string,
  showPercent: boolean,
  nameMaxLength: number,
  styleConfig: any,
  autoColors: any[],
) => {
  const autoColorStyle = autoColors[index];
  const [r, g, b] = autoColorStyle.lineColor.match(/\d+/g).map(Number);

  return {
    name: getTruncatedDisplayName(item, index, nameKey, nameMaxLength),
    value: item[valueKey] || 0,
    fullName: getDisplayName(item, index, nameKey),

    itemStyle: {
      color: new echarts.graphic.RadialGradient(0, 0, 1, [
        {
          offset: 0,
          color: `rgba(${r}, ${g}, ${b}, ${styleConfig.gradientTopAlpha})`,
        },
        {
          offset: 1,
          color: `rgba(${r}, ${g}, ${b}, ${styleConfig.gradientBottomAlpha})`,
        },
      ]),
      borderColor: `rgba(${r}, ${g}, ${b}, 0.8)`,
      borderWidth: 2,
      shadowColor: `rgba(${r}, ${g}, ${b}, 0.4)`,
      shadowBlur: styleConfig.shadowBlur,
      shadowOffsetY: 2,
    },
    emphasis: {
      scale: true,
      scaleSize: 10,
      itemStyle: {
        color: new echarts.graphic.RadialGradient(0, 0, 1, [
          { offset: 0, color: `rgba(${r}, ${g}, ${b}, 1)` },
          { offset: 1, color: `rgba(${r}, ${g}, ${b}, 0.8)` },
        ]),
        shadowBlur: styleConfig.shadowBlur + 5,
        shadowOffsetY: 4,
      },
    },
    label: {
      show: true,
      fontSize: 12,
      color: "#333",
      formatter: showPercent ? "{b}: {d}%" : "{b}: {c}",
      fontWeight: "bold",
    },
    labelLine: {
      show: true,
      length: 15,
      length2: 10,
      lineStyle: {
        color: `rgba(${r}, ${g}, ${b}, 0.6)`,
        width: 2,
      },
    },
  };
};

const createBarDataItem = (
  item: any,
  index: number,
  nameKey: string,
  valueKey: string,
  nameMaxLength: number,
  styleConfig: any,
  autoColors: any[],
) => {
  const autoColorStyle = autoColors[index];
  const [r, g, b] = autoColorStyle.lineColor.match(/\d+/g).map(Number);

  return {
    name: getTruncatedDisplayName(item, index, nameKey, nameMaxLength),
    value: item[valueKey] || 0,
    fullName: getDisplayName(item, index, nameKey),
    itemStyle: {
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        {
          offset: 0,
          color: `rgba(${r}, ${g}, ${b}, ${styleConfig.gradientTopAlpha})`,
        },
        {
          offset: 1,
          color: `rgba(${r}, ${g}, ${b}, ${styleConfig.gradientBottomAlpha})`,
        },
      ]),
      borderColor: `rgba(${r}, ${g}, ${b}, 0.8)`,
      borderWidth: 1,
      borderRadius: [4, 4, 0, 0],
      shadowColor: `rgba(${r}, ${g}, ${b}, 0.3)`,
      shadowBlur: styleConfig.shadowBlur,
      shadowOffsetY: 2,
    },
    emphasis: {
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: `rgba(${r}, ${g}, ${b}, 1)` },
          { offset: 1, color: `rgba(${r}, ${g}, ${b}, 0.8)` },
        ]),
        shadowBlur: styleConfig.shadowBlur + 5,
        shadowOffsetY: 4,
      },
    },
  };
};

const formatValue = (value: number): string => {
  if (props.formatter && typeof props.formatter === "function") {
    const formatted = props.formatter(value);
    return typeof formatted === "string" ? formatted : value.toLocaleString();
  }
  return value.toLocaleString();
};

// 保持原有的 Tooltip HTML 结构和样式
const createModernTooltip = (valueKey: string) => ({
  trigger: "item",
  backgroundColor: "transparent",
  borderColor: "transparent",
  borderWidth: 0,
  padding: 0,
  textStyle: {
    fontFamily: "Arial, sans-serif",
    fontSize: 14,
  },
  extraCssText:
    "box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15); border-radius: 8px;",

  formatter: (params: any) => {
    const { name, value, percent, data } = params;
    const displayName = data?.fullName || name;

    const total = props.data.reduce((sum: number, item: any) => {
      let itemValue;
      if (typeof item === "object") {
        itemValue =
          item[valueKey] ||
          item.count ||
          item.amount ||
          item.num ||
          Object.values(item).find((v) => typeof v === "number") ||
          0;
      } else {
        itemValue = Array.isArray(item) ? item[1] : item;
      }
      return sum + (Number(itemValue) || 0);
    }, 0);

    const actualPercent =
      percent || (total > 0 ? ((value / total) * 100).toFixed(1) : "0");

    return `<div style="overflow: hidden; border-radius: 8px; box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1); min-width: 140px;">
                  <div style="background-color: #e5e7eb; padding: 4px 12px; font-size: 12px; color: #4a5568; font-weight: 600; text-align: center;">${displayName}</div>
                  <div style="background-color: #fff; padding: 10px 14px;">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px;">
                      <span style="color: #6b7280; font-size: 12px;">数值</span>
                      <span style="font-weight: 700; color: #1f2937; font-size: 14px;">${formatValue(
                        value,
                      )}</span>
                    </div>
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px;">
                      <span style="color: #6b7280; font-size: 12px;">占比</span>
                      <span style="font-weight: 700; color: #059669; font-size: 14px;">${actualPercent}%</span>
                    </div>
                    <div style="border-top: 1px solid #e5e7eb; padding-top: 4px; margin-top: 2px;">
                      <div style="display: flex; justify-content: space-between; align-items: center;">
                        <span style="color: #9ca3af; font-size: 11px;">总计</span>
                        <span style="font-weight: 600; color: #6b7280; font-size: 12px;">${formatValue(
                          total,
                        )}</span>
                      </div>
                    </div>
                  </div>
                </div>`;
  },
});

const createLegendConfig = (legendPosition: string, pieData: any[]) => {
  const baseConfig = {
    itemWidth: 18,
    itemHeight: 14,
    itemGap: 14,
    textStyle: { color: "#666", fontSize: 12 },
    data: pieData.map((item) => item.fullName || item.name),
  };

  switch (legendPosition) {
    case "right":
      return {
        ...baseConfig,
        type: "scroll",
        orient: "vertical",
        right: "10%",
        top: "center",
      };
    case "bottom":
      return {
        ...baseConfig,
        type: "scroll",
        orient: "horizontal",
        bottom: "0%",
        left: "center",
      };
    case "top":
      return {
        ...baseConfig,
        type: "scroll",
        orient: "horizontal",
        top: "6%",
        left: "center",
      };
    case "left":
      return {
        ...baseConfig,
        type: "scroll",
        orient: "vertical",
        left: "6%",
        top: "center",
      };
    case "none":
      return { show: false };
    default:
      return {
        ...baseConfig,
        type: "scroll",
        orient: "vertical",
        right: "10%",
        top: "center",
      };
  }
};

// ====================
// 主要计算属性
// ====================

const processedOptions = computed((): EChartsOption => {
  if (!props.data?.length) return {};

  const defaultStyleConfig = createPieStyleConfig();
  const autoColors = generateColorStyles(props.data.length, defaultStyleConfig);

  const chartData =
    props.chartType === "bar"
      ? props.data.map((item: any, index: number) =>
          createBarDataItem(
            item,
            index,
            props.nameKey,
            props.valueKey,
            props.nameMaxLength,
            defaultStyleConfig,
            autoColors,
          ),
        )
      : props.data.map((item: any, index: number) =>
          createPieDataItem(
            item,
            index,
            props.nameKey,
            props.valueKey,
            props.showPiePercent,
            props.nameMaxLength,
            defaultStyleConfig,
            autoColors,
          ),
        );

  const modernTooltip = createModernTooltip(props.valueKey);
  const legendConfig =
    props.chartType === "pie"
      ? createLegendConfig(props.legend, chartData)
      : null;
  const isHorizontalBar =
    props.chartType === "bar" && props.barDirection === "horizontal";

  return {
    tooltip: modernTooltip,
    legend: legendConfig,
    series:
      props.chartType === "bar"
        ? [
            {
              name: "数据",
              type: "bar",
              data: chartData,
              barWidth: isHorizontalBar ? 18 : "65%",
              barMaxWidth: isHorizontalBar ? 28 : 50,
              itemStyle: {
                borderRadius: isHorizontalBar ? [0, 4, 4, 0] : [4, 4, 0, 0],
              },
              label: {
                show: props.showBarValue,
                position: isHorizontalBar ? "right" : "top",
                fontSize: 12,
                color: "#666",
                formatter: (params: any) => formatValue(params.value),
              },
              emphasis: { focus: "self" },
              animationDelay: (idx: number) => idx * 100,
            },
          ]
        : [
            {
              name: "数据",
              type: "pie",
              radius: ["40%", "70%"],
              center: ["50%", "55%"],
              data: chartData,
              roseType: false,
              animationType: "scale",
              animationEasing: "elasticOut",
              animationDelay: (idx: number) => Math.random() * 200,
              label: {
                show: true,
                position: "outside",
                fontSize: 11,
                color: "#666",
                formatter: props.showPiePercent ? "{b}\n{d}%" : "{b}\n{c}",
                fontWeight: "normal",
              },
              labelLine: {
                show: true,
                length: 20,
                length2: 15,
                smooth: true,
              },
              emphasis: {
                focus: "self",
                blurScope: "coordinateSystem",
              },
            },
          ],
    ...(props.chartType === "bar"
      ? {
          xAxis: {
            type: isHorizontalBar ? "value" : "category",
            data: isHorizontalBar ? undefined : chartData.map((item) => item.name),
            axisLine: { show: false },
            axisTick: { show: false },
            axisLabel: {
              color: "#666",
              fontSize: 12,
              interval: isHorizontalBar ? undefined : 0,
              rotate: isHorizontalBar ? 0 : chartData.length > 6 ? 45 : 0,
              formatter: isHorizontalBar
                ? (value: number) => formatValue(value)
                : undefined,
            },
            splitLine: {
              show: true,
              lineStyle: { color: "#eee", type: "dashed" },
            },
          },
          yAxis: {
            type: isHorizontalBar ? "category" : "value",
            data: isHorizontalBar ? chartData.map((item) => item.name) : undefined,
            show: props.showYAxis,
            axisLine: { show: false },
            axisTick: { show: false },
            axisLabel: {
              color: "#666",
              fontSize: 12,
              interval: 0,
              formatter: isHorizontalBar
                ? undefined
                : (value: number) => formatValue(value),
            },
          },
          grid: {
            left: isHorizontalBar ? "8%" : "2%",
            right: "2%",
            bottom: "4%",
            top: "8%",
            containLabel: true,
          },
        }
      : {}),
    backgroundColor: "transparent",
    textStyle: {
      fontFamily: 'Arial, "Microsoft YaHei", sans-serif',
    },
  };
});

// ====================
// 🚀 核心修复：生命周期与渲染逻辑 (彻底解决 v-show/display:none 问题)
// ====================

/**
 * 尝试初始化图表
 * 只有当 DOM 真正存在且 width/height > 0 时才执行
 */
const initChart = () => {
  if (myChart) return;
  if (!chartContainer.value) return;

  const width = chartContainer.value.clientWidth;
  const height = chartContainer.value.clientHeight;

  // 关键：如果宽高为0（通常因为 v-show=false），则不初始化
  if (width === 0 || height === 0) return;

  // 使用 markRaw 优化性能
  const instance = echarts.init(chartContainer.value);
  myChart = markRaw(instance);

  // 绑定事件
  myChart.on("click", (params) => {
    emit("chart-click", params);
  });

  // 应用配置
  if (processedOptions.value) {
    myChart.setOption(processedOptions.value, true);
  }
};

/**
 * ResizeObserver 回调
 * 处理显示/隐藏切换和窗口缩放
 */
const handleResizeEntry = (entries: ResizeObserverEntry[]) => {
  for (const entry of entries) {
    const { width, height } = entry.contentRect;
    if (width > 0 && height > 0) {
      if (!myChart) {
        // 从隐藏变为显示 -> 初始化
        initChart();
      } else {
        // 尺寸变化 -> Resize
        myChart.resize();
      }
    }
  }
};

// 监听配置变化
watch(
  processedOptions,
  (newOptions) => {
    if (myChart) {
      myChart.setOption(newOptions, true);
    }
  },
  { deep: true },
);

onMounted(() => {
  if (chartContainer.value) {
    // 1. 启动 ResizeObserver
    if (window.ResizeObserver) {
      resizeObserver = new ResizeObserver(handleResizeEntry);
      resizeObserver.observe(chartContainer.value);
    }
    // 2. 尝试立即初始化
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

defineExpose({
  resize: () => myChart?.resize(),
});
</script>

<template>
  <div ref="chartContainer" :style="chartContainerStyle"></div>
</template>

<style scoped>
/* 样式已内联在 JS 中或由 ECharts 管理 */
</style>

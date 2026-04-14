# 通用图表组件

## 模块归类

- 分区：复用组件
- 类型：前端通用图表组件
- 收录位置：`web/src/components/charts`

## 背景与目标

这次收录的是两类高频图表能力：

- `TrendChart.vue`：趋势分析图，适合折线图、柱状图、双轴混合趋势展示
- `PieChart.vue`：占比分布图，适合饼图、柱状图两种形态切换

它们解决的不是某个单页图表能否画出来，而是后台项目里反复出现的这些问题：

- 趋势图与分布图在多个页面中重复开发
- 图表主题、尺寸、图例、格式化规则难以统一
- 同一批数据常常需要在饼图 / 柱状图之间切换展示
- 趋势图经常需要支持多序列、双轴、时间对齐和数值格式化

在 `admin-lab` 中收录它们，是为了先沉淀一套可迁移、可扩展、可配置的图表底座。

## 当前收录范围

已收录文件：

- `web/src/components/charts/TrendChart.vue`
- `web/src/components/charts/PieChart.vue`
- `web/src/components/charts/generateColorStyles.js`

配套展示入口：

- 前端页面：`实验室 / 复用组件 / 通用图表组件`

## 组件能力

### TrendChart

当前已覆盖：

- 多序列折线 / 柱状趋势图
- 左右双 Y 轴
- 时间维度 X 轴格式化
- 图例位置切换
- 主题切换
- 时间区间对齐与数据填充
- 容器宽高配置
- 实例暴露与事件注册

适合场景：

- 指标趋势分析
- 监控类时序图
- 攻防态势趋势看板
- 同时展示数量与比例的双轴图

### PieChart

当前已覆盖：

- 饼图 / 柱状图双形态展示
- 图例方位切换
- 百分比显示
- 名称映射与名称截断
- 柱状图方向切换（横向 / 纵向）
- 柱顶数值显示
- 自定义数值格式化
- 点击事件透出

适合场景：

- 渠道分布
- 风险等级占比
- 类型统计
- 枚举分组结果可视化

## 最小接入方式

### TrendChart

最小接入需要：

1. 一组数组数据
2. 一个 X 轴字段
3. 一份 `series` 配置

示例：

```vue
<TrendChart
  :data="trendData"
  :series="{
    orderCount: { label: '订单量', type: 'bar', yAxis: 'left' },
    successRate: { label: '成功率', type: 'line', yAxis: 'right', formatter: (v) => v + '%' }
  }"
  x-axis-key="date"
  y-axis-name="数量"
  y-axis-right-name="比例"
  height="320"
/>
```

### PieChart

最小接入需要：

1. 一组名称 + 数值数据
2. 指定 `chart-type`
3. 按需补充显示配置

示例：

```vue
<PieChart
  :data="distributionData"
  chart-type="pie"
  name-key="name"
  value-key="value"
  legend="bottom"
  :show-pie-percent="true"
  height="320"
/>
```

## 迁移到内网时建议同步的内容

建议至少同步这些文件：

- `web/src/components/charts/TrendChart.vue`
- `web/src/components/charts/PieChart.vue`
- `web/src/components/charts/generateColorStyles.js`
- 使用它们的业务示例页或封装页

如果内网项目已存在统一图表主题，可保留组件逻辑，仅替换：

- 主题默认值
- 色板生成策略
- Tooltip / Label 文案
- 事件透出协议

## 迁移时要特别注意

- 内网项目是否已统一使用 `echarts` 及对应版本
- 主题色和暗黑模式是否需要与现有设计体系对齐
- 趋势图的时间字段格式是否统一
- 数值格式化函数是否应抽到公共 `utils`
- 点击事件透出的 payload 是否要和现有页面协议统一

## 当前边界

这次在 `admin-lab` 里补齐的是：

- 组件目录收录
- 实验室展示页
- 菜单路由入口
- 迁移说明文档

当前还没有做的是：

- 与后端图表协议的统一适配层
- 图表导出图片 / 导出数据能力
- 单测与更细粒度的 props 文档页
- 和 `security-echarts` 这类业务图表容器的统一协议整合

如果后续继续沉淀，建议下一步补：

- 图表统一 adapter 协议
- 主题与配色中心化配置
- 图表空态 / 加载态 / 错态容器
- 典型后台场景的更多组合示例

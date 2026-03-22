package response

// SecurityDashboardPayload 是可视化面板输出给前端的统一契约。
// 前端可先做 adapter 归一化，再稳定绑定到复用图表组件。
type SecurityDashboardPayload struct {
	Summary []SecuritySummaryCard      `json:"summary"`
	Charts  SecurityChartBindingGroup  `json:"charts"`
	Meta    SecurityDashboardExtraMeta `json:"meta"`
}

// SecuritySummaryCard 表示顶部 KPI 卡片。
type SecuritySummaryCard struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Value string `json:"value"`
	Sub   string `json:"sub"`
}

// SecurityChartBindingGroup 聚合 bar/line/pie/mix 四类图表绑定。
type SecurityChartBindingGroup struct {
	Bar  SecurityChartBinding `json:"bar"`
	Line SecurityChartBinding `json:"line"`
	Pie  SecurityChartBinding `json:"pie"`
	Mix  SecurityChartBinding `json:"mix"`
}

// SecurityChartBinding 是单个图表的统一绑定结构。
// Data 为图表主体数据；Dimension* 用于维度切换。
type SecurityChartBinding struct {
	Data             any                       `json:"data"`
	DimensionOptions []SecurityDimensionOption `json:"dimensionOptions,omitempty"`
	DimensionDataMap map[string]any            `json:"dimensionDataMap,omitempty"`
	DimensionValue   string                    `json:"dimensionValue,omitempty"`
}

// SecurityDimensionOption 定义维度候选项。
type SecurityDimensionOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// SecurityDashboardExtraMeta 承载非图表直绘数据（如表格视图）。
type SecurityDashboardExtraMeta struct {
	MixTableData []SecurityMixTableRow `json:"mixTableData"`
}

// SecuritySeries 仅表达业务语义和数值，不携带样式信息。
// 前端组件会根据 metricKey/name 做统一样式映射。
type SecuritySeries struct {
	Name      string    `json:"name"`
	MetricKey string    `json:"metricKey,omitempty"`
	Values    []float64 `json:"values"`
}

// SecurityBarLineData 是柱状图与折线图的通用结构。
type SecurityBarLineData struct {
	Categories []string         `json:"categories"`
	Series     []SecuritySeries `json:"series"`
}

// SecurityPieSlice 表示饼图单个扇区。
type SecurityPieSlice struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// SecurityPieData 表示饼图结构。
type SecurityPieData struct {
	SeriesName string             `json:"seriesName"`
	Data       []SecurityPieSlice `json:"data"`
}

// SecurityMixData 表示双轴混合图结构。
type SecurityMixData struct {
	Categories []string         `json:"categories"`
	LeftName   string           `json:"leftName"`
	RightName  string           `json:"rightName"`
	BarSeries  []SecuritySeries `json:"barSeries"`
	LineSeries []SecuritySeries `json:"lineSeries"`
}

// SecurityMixTableRow 是混合图对应的表格行。
type SecurityMixTableRow struct {
	Time         string  `json:"time"`
	EventCount   int     `json:"eventCount"`
	BlockedCount int     `json:"blockedCount"`
	HandleRate   float64 `json:"handleRate"`
}

// SecurityDrilldownRow 定义下钻抽屉的明细行结构。
type SecurityDrilldownRow struct {
	EventID    string `json:"eventId"`
	AttackType string `json:"attackType"`
	SourceIP   string `json:"sourceIp"`
	Status     string `json:"status"`
	Level      string `json:"level"`
	OccurAt    string `json:"occurAt"`
}

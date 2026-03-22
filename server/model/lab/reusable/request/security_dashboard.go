package request

// SecurityDashboardQuery 定义可视化面板查询参数。
// 这里的字段保持“组件视角”与“业务筛选视角”兼容，便于前端复用。
type SecurityDashboardQuery struct {
	Keyword        string `form:"keyword" json:"keyword"`
	ThreatLevels   string `form:"threatLevels" json:"threatLevels"`
	DataSource     string `form:"dataSource" json:"dataSource"`
	TimeFrom       string `form:"timeFrom" json:"timeFrom"`
	TimeTo         string `form:"timeTo" json:"timeTo"`
	TrendDimension string `form:"trendDimension" json:"trendDimension"`
}

// SecurityDrilldownQuery 定义图表点击下钻时的查询参数。
// 除图元上下文外，额外保留当前筛选条件，保证下钻结果与主视图一致。
type SecurityDrilldownQuery struct {
	Category     string `form:"category" json:"category"`
	Series       string `form:"series" json:"series"`
	Dimension    string `form:"dimension" json:"dimension"`
	Keyword      string `form:"keyword" json:"keyword"`
	ThreatLevels string `form:"threatLevels" json:"threatLevels"`
	DataSource   string `form:"dataSource" json:"dataSource"`
	TimeFrom     string `form:"timeFrom" json:"timeFrom"`
	TimeTo       string `form:"timeTo" json:"timeTo"`
}

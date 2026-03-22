package reusable

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strings"

	reusableReq "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/request"
	reusableRes "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/response"
)

type SecurityDashboardService struct{}

// GetPanel 生成可视化面板数据契约。
// 当前示例使用可重复的伪随机数据来模拟后端，后续可替换为真实仓储查询。
func (s *SecurityDashboardService) GetPanel(query reusableReq.SecurityDashboardQuery) reusableRes.SecurityDashboardPayload {
	r := newDeterministicRand(buildSeedKey(query.Keyword, query.ThreatLevels, query.DataSource, query.TimeFrom, query.TimeTo))

	barDimensionOptions := []reusableRes.SecurityDimensionOption{
		{Label: "近24小时", Value: "24h"},
		{Label: "近7天", Value: "7d"},
	}
	lineDimensionOptions := []reusableRes.SecurityDimensionOption{
		{Label: "小时", Value: "hour"},
		{Label: "天", Value: "day"},
		{Label: "周", Value: "week"},
	}

	mixData, mixTable := buildMixData(r)
	trendDimension := query.TrendDimension
	if trendDimension == "" {
		trendDimension = "hour"
	}

	return reusableRes.SecurityDashboardPayload{
		Summary: buildSummaryCards(r),
		Charts: reusableRes.SecurityChartBindingGroup{
			Bar: reusableRes.SecurityChartBinding{
				Data:             buildBarData(r),
				DimensionOptions: barDimensionOptions,
				DimensionDataMap: buildBarDimensionDataMap(r),
			},
			Line: reusableRes.SecurityChartBinding{
				Data:             buildLineData(r),
				DimensionOptions: lineDimensionOptions,
				DimensionDataMap: buildLineDimensionDataMap(r),
				DimensionValue:   trendDimension,
			},
			Pie: reusableRes.SecurityChartBinding{
				Data:             buildPieData(r),
				DimensionOptions: barDimensionOptions,
				DimensionDataMap: buildPieDimensionDataMap(r),
			},
			Mix: reusableRes.SecurityChartBinding{
				Data: mixData,
			},
		},
		Meta: reusableRes.SecurityDashboardExtraMeta{
			MixTableData: mixTable,
		},
	}
}

// GetDrilldown 根据图表点击上下文返回下钻明细。
// 通过筛选参数参与种子计算，确保同条件下结果稳定，便于联调演示。
func (s *SecurityDashboardService) GetDrilldown(query reusableReq.SecurityDrilldownQuery) []reusableRes.SecurityDrilldownRow {
	attackType := strings.TrimSpace(query.Category)
	if attackType == "" {
		attackType = "未知类别"
	}

	r := newDeterministicRand(buildSeedKey(query.Category, query.Series, query.Dimension, query.Keyword, query.ThreatLevels, query.DataSource, query.TimeFrom, query.TimeTo))
	levels := []string{"高危", "中危", "低危"}
	rows := make([]reusableRes.SecurityDrilldownRow, 0, 8)
	for index := 0; index < 8; index++ {
		status := "已处置"
		if r.Float64() > 0.56 {
			status = "待处置"
		}
		rows = append(rows, reusableRes.SecurityDrilldownRow{
			EventID:    fmt.Sprintf("EVT-%d-%02d", randomRange(r, 20260100, 20269999), index+1),
			AttackType: attackType,
			SourceIP: fmt.Sprintf(
				"%d.%d.%d.%d",
				randomRange(r, 10, 223),
				randomRange(r, 1, 254),
				randomRange(r, 1, 254),
				randomRange(r, 1, 254),
			),
			Status:  status,
			Level:   levels[randomRange(r, 0, len(levels)-1)],
			OccurAt: fmt.Sprintf("2026-03-%02d %02d:%02d:00", randomRange(r, 1, 22), randomRange(r, 0, 23), randomRange(r, 0, 59)),
		})
	}
	return rows
}

// buildSummaryCards 生成顶部 KPI 卡片。
func buildSummaryCards(r *rand.Rand) []reusableRes.SecuritySummaryCard {
	return []reusableRes.SecuritySummaryCard{
		{Key: "alerts", Title: "总告警量", Value: fmt.Sprintf("%d", randomRange(r, 12000, 18000)), Sub: "较昨日 +6.4%"},
		{Key: "high", Title: "高危事件", Value: fmt.Sprintf("%d", randomRange(r, 420, 980)), Sub: "需优先处置"},
		{Key: "blocked", Title: "已拦截攻击", Value: fmt.Sprintf("%d", randomRange(r, 8600, 11000)), Sub: "拦截率 98%+"},
		{Key: "assets", Title: "受影响资产", Value: fmt.Sprintf("%d", randomRange(r, 70, 180)), Sub: "含服务器/终端/业务系统"},
	}
}

// buildBarData 生成柱状图基础数据。
func buildBarData(r *rand.Rand) reusableRes.SecurityBarLineData {
	// 后端仅返回语义字段 metricKey，颜色/面积等样式由前端统一映射。
	return reusableRes.SecurityBarLineData{
		Categories: []string{"恶意扫描", "暴力破解", "漏洞利用", "木马回连", "横向移动", "数据外传"},
		Series: []reusableRes.SecuritySeries{
			{Name: "已处置", MetricKey: "handled_count", Values: randomSeriesValues(r, 6, 180, 980)},
			{Name: "待处置", MetricKey: "pending_count", Values: randomSeriesValues(r, 6, 80, 640)},
		},
	}
}

// buildLineData 生成折线图基础数据。
func buildLineData(r *rand.Rand) reusableRes.SecurityBarLineData {
	return reusableRes.SecurityBarLineData{
		Categories: []string{"00:00", "04:00", "08:00", "12:00", "16:00", "20:00", "24:00"},
		Series: []reusableRes.SecuritySeries{
			{Name: "外网攻击", MetricKey: "external_attack_count", Values: randomSeriesValues(r, 7, 420, 2300)},
			{Name: "内网异常", MetricKey: "internal_anomaly_count", Values: randomSeriesValues(r, 7, 180, 980)},
		},
	}
}

// buildPieData 生成饼图基础数据。
func buildPieData(r *rand.Rand) reusableRes.SecurityPieData {
	return reusableRes.SecurityPieData{
		SeriesName: "风险等级",
		Data: []reusableRes.SecurityPieSlice{
			{Name: "高危", Value: float64(randomRange(r, 180, 420))},
			{Name: "中危", Value: float64(randomRange(r, 420, 960))},
			{Name: "低危", Value: float64(randomRange(r, 900, 2200))},
		},
	}
}

// buildLineDimensionDataMap 生成折线图多维度映射数据。
func buildLineDimensionDataMap(r *rand.Rand) map[string]any {
	return map[string]any{
		"hour": reusableRes.SecurityBarLineData{
			Categories: []string{"00:00", "04:00", "08:00", "12:00", "16:00", "20:00", "24:00"},
			Series: []reusableRes.SecuritySeries{
				{Name: "外网攻击", MetricKey: "external_attack_count", Values: randomSeriesValues(r, 7, 420, 2300)},
				{Name: "内网异常", MetricKey: "internal_anomaly_count", Values: randomSeriesValues(r, 7, 180, 980)},
			},
		},
		"day": reusableRes.SecurityBarLineData{
			Categories: []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"},
			Series: []reusableRes.SecuritySeries{
				{Name: "外网攻击", MetricKey: "external_attack_count", Values: randomSeriesValues(r, 7, 1800, 6900)},
				{Name: "内网异常", MetricKey: "internal_anomaly_count", Values: randomSeriesValues(r, 7, 800, 3900)},
			},
		},
		"week": reusableRes.SecurityBarLineData{
			Categories: []string{"第1周", "第2周", "第3周", "第4周", "第5周", "第6周"},
			Series: []reusableRes.SecuritySeries{
				{Name: "外网攻击", MetricKey: "external_attack_count", Values: randomSeriesValues(r, 6, 12000, 33000)},
				{Name: "内网异常", MetricKey: "internal_anomaly_count", Values: randomSeriesValues(r, 6, 7000, 19000)},
			},
		},
	}
}

// buildBarDimensionDataMap 生成柱状图多维度映射数据。
func buildBarDimensionDataMap(r *rand.Rand) map[string]any {
	categories := []string{"DDoS", "漏洞探测", "口令攻击", "恶意IP", "木马通讯"}
	return map[string]any{
		"24h": reusableRes.SecurityBarLineData{
			Categories: categories,
			Series: []reusableRes.SecuritySeries{
				{Name: "已处置", MetricKey: "handled_count", Values: randomSeriesValues(r, 5, 80, 620)},
				{Name: "待处置", MetricKey: "pending_count", Values: randomSeriesValues(r, 5, 60, 540)},
			},
		},
		"7d": reusableRes.SecurityBarLineData{
			Categories: categories,
			Series: []reusableRes.SecuritySeries{
				{Name: "已处置", MetricKey: "handled_count", Values: randomSeriesValues(r, 5, 1200, 5400)},
				{Name: "待处置", MetricKey: "pending_count", Values: randomSeriesValues(r, 5, 900, 4200)},
			},
		},
	}
}

// buildPieDimensionDataMap 生成饼图多维度映射数据。
func buildPieDimensionDataMap(r *rand.Rand) map[string]any {
	return map[string]any{
		"24h": reusableRes.SecurityPieData{
			SeriesName: "资产类型占比（24h）",
			Data: []reusableRes.SecurityPieSlice{
				{Name: "服务器", Value: float64(randomRange(r, 40, 90))},
				{Name: "终端", Value: float64(randomRange(r, 50, 120))},
				{Name: "网络设备", Value: float64(randomRange(r, 20, 80))},
			},
		},
		"7d": reusableRes.SecurityPieData{
			SeriesName: "资产类型占比（7d）",
			Data: []reusableRes.SecurityPieSlice{
				{Name: "服务器", Value: float64(randomRange(r, 180, 420))},
				{Name: "终端", Value: float64(randomRange(r, 200, 540))},
				{Name: "网络设备", Value: float64(randomRange(r, 110, 300))},
			},
		},
	}
}

// buildMixData 生成双轴混合图与配套表格数据。
func buildMixData(r *rand.Rand) (reusableRes.SecurityMixData, []reusableRes.SecurityMixTableRow) {
	categories := []string{"00:00", "04:00", "08:00", "12:00", "16:00", "20:00", "24:00"}
	events := make([]float64, len(categories))
	blocked := make([]float64, len(categories))
	rates := make([]float64, len(categories))
	rows := make([]reusableRes.SecurityMixTableRow, 0, len(categories))

	for idx := range categories {
		eventCount := float64(randomRange(r, 420, 2200))
		blockedCount := float64(int(eventCount * float64(randomRange(r, 72, 98)) / 100))
		rate := roundFloat(blockedCount/eventCount*100, 2)
		events[idx] = eventCount
		blocked[idx] = blockedCount
		rates[idx] = rate

		rows = append(rows, reusableRes.SecurityMixTableRow{
			Time:         categories[idx],
			EventCount:   int(eventCount),
			BlockedCount: int(blockedCount),
			HandleRate:   rate,
		})
	}

	return reusableRes.SecurityMixData{
		Categories: categories,
		LeftName:   "事件量",
		RightName:  "处置率",
		BarSeries: []reusableRes.SecuritySeries{
			{Name: "事件量", MetricKey: "event_count", Values: events},
			{Name: "拦截量", MetricKey: "blocked_count", Values: blocked},
		},
		LineSeries: []reusableRes.SecuritySeries{
			{Name: "处置率", MetricKey: "handle_rate", Values: rates},
		},
	}, rows
}

// randomSeriesValues 用于批量生成序列值。
func randomSeriesValues(r *rand.Rand, count int, min int, max int) []float64 {
	values := make([]float64, count)
	for idx := 0; idx < count; idx++ {
		values[idx] = float64(randomRange(r, min, max))
	}
	return values
}

// randomRange 返回闭区间随机数。
func randomRange(r *rand.Rand, min int, max int) int {
	if max <= min {
		return min
	}
	return r.Intn(max-min+1) + min
}

// roundFloat 按指定精度四舍五入。
func roundFloat(v float64, scale int) float64 {
	if scale <= 0 {
		return float64(int(v))
	}
	pow := 1.0
	for idx := 0; idx < scale; idx++ {
		pow *= 10
	}
	return float64(int(v*pow+0.5)) / pow
}

// buildSeedKey 将筛选条件序列化为稳定种子字符串。
func buildSeedKey(parts ...string) string {
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		values = append(values, strings.TrimSpace(part))
	}
	return strings.Join(values, "|")
}

// newDeterministicRand 生成稳定随机源，保证同条件下输出可复现。
func newDeterministicRand(seedKey string) *rand.Rand {
	h := fnv.New64a()
	_, _ = h.Write([]byte(seedKey))
	seed := int64(h.Sum64())
	if seed == 0 {
		seed = 20260322
	}
	return rand.New(rand.NewSource(seed))
}

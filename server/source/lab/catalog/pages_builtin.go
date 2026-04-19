package catalog

func init() {
	registerPage(PageSpec{
		Key:  "simulation-overview",
		Area: AreaSimulation,
		Menu: MenuItem{
			Path:      "overview",
			Name:      "labSimulationOverview",
			Component: "view/lab/simulation/overview.vue",
			Sort:      1,
			Title:     "概览",
			Icon:      "tickets",
		},
	})

	registerPage(PageSpec{
		Key:  "simulation-base-data-io",
		Area: AreaSimulation,
		Menu: MenuItem{
			Path:      "base-data-io",
			Name:      "labSimulationBaseDataIO",
			Component: "view/lab/simulation/base-data-io.vue",
			Sort:      2,
			Title:     "基础数据导入导出仿真",
			Icon:      "document-copy",
		},
		APIs: []APIItem{
			{APIGroup: "BaseDataSimulation", Method: "GET", Path: "/baseDataSimulation/templates", Description: "获取仿真模板列表"},
			{APIGroup: "BaseDataSimulation", Method: "GET", Path: "/baseDataSimulation/template", Description: "下载仿真模板"},
			{APIGroup: "BaseDataSimulation", Method: "GET", Path: "/baseDataSimulation/export", Description: "导出仿真数据"},
			{APIGroup: "BaseDataSimulation", Method: "POST", Path: "/baseDataSimulation/import", Description: "导入仿真数据"},
		},
	})

	registerPage(PageSpec{
		Key:  "simulation-cucode",
		Area: AreaSimulation,
		Menu: MenuItem{
			Path:      "cucode",
			Name:      "labSimulationCucode",
			Component: "view/lab/simulation/cucode.vue",
			Sort:      3,
			Title:     "代码表管理",
			Icon:      "memo",
		},
	})

	registerPage(PageSpec{
		Key:  "component-demo-overview",
		Area: AreaComponentDemo,
		Menu: MenuItem{
			Path:      "overview",
			Name:      "labComponentDemoOverview",
			Component: "view/lab/component-demo/overview.vue",
			Sort:      1,
			Title:     "概览",
			Icon:      "tickets",
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-overview",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "overview",
			Name:      "labReusableOverview",
			Component: "view/lab/reusable/overview.vue",
			Sort:      1,
			Title:     "概览",
			Icon:      "tickets",
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-excel-io",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "excel-io",
			Name:      "labReusableExcelIO",
			Component: "view/lab/reusable/excel-io.vue",
			Sort:      2,
			Title:     "Excel 实验面板",
			Icon:      "document-copy",
		},
		APIs: []APIItem{
			{APIGroup: "ExcelIO", Method: "GET", Path: "/excelIO/templates", Description: "获取 Excel 模板列表"},
			{APIGroup: "ExcelIO", Method: "GET", Path: "/excelIO/template", Description: "下载 Excel 导入模板"},
			{APIGroup: "ExcelIO", Method: "GET", Path: "/excelIO/export", Description: "导出 Excel 示例数据"},
			{APIGroup: "ExcelIO", Method: "POST", Path: "/excelIO/import", Description: "导入并解析 Excel"},
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-crud-form-dialog",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "crud-form-dialog",
			Name:      "labReusableCrudFormDialog",
			Component: "view/lab/reusable/crud-form-dialog.vue",
			Sort:      3,
			Title:     "新增编辑弹窗",
			Icon:      "edit-pen",
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-security-echarts",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "security-echarts",
			Name:      "labReusableSecurityEcharts",
			Component: "view/lab/reusable/security-echarts.vue",
			Sort:      4,
			Title:     "网安可视化面板",
			Icon:      "trend-charts",
		},
		APIs: []APIItem{
			{APIGroup: "SecurityDashboard", Method: "GET", Path: "/securityDashboard/panel", Description: "获取网安可视化面板数据"},
			{APIGroup: "SecurityDashboard", Method: "GET", Path: "/securityDashboard/drilldown", Description: "获取网安可视化下钻明细"},
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-list-query-bar",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "list-query-bar",
			Name:      "labReusableListQueryBar",
			Component: "view/lab/reusable/list-query-bar.vue",
			Sort:      5,
			Title:     "列表查询栏",
			Icon:      "search",
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-reliable-upload",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "reliable-upload",
			Name:      "labReusableReliableUpload",
			Component: "view/lab/reusable/reliable-upload.vue",
			Sort:      6,
			Title:     "可靠上报框架",
			Icon:      "upload-filled",
		},
		APIs: []APIItem{
			{APIGroup: "ReliableUpload", Method: "GET", Path: "/reliableUpload/profile", Description: "获取可靠上报框架资料"},
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-table-pro",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "table-pro",
			Name:      "labReusableTablePro",
			Component: "view/lab/reusable/table-pro.vue",
			Sort:      7,
			Title:     "Table Pro",
			Icon:      "grid",
		},
		APIs: []APIItem{
			{APIGroup: "TablePro", Method: "POST", Path: "/tablePro/page", Description: "Get table pro page"},
			{APIGroup: "TablePro", Method: "POST", Path: "/tablePro/export", Description: "Export table pro data"},
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-dict-usage",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "dict-usage",
			Name:      "labReusableDictUsage",
			Component: "view/lab/reusable/dict-usage.vue",
			Sort:      8,
			Title:     "字典消费组件",
			Icon:      "collection-tag",
		},
	})

	registerPage(PageSpec{
		Key:  "reusable-biz-log",
		Area: AreaReusable,
		Menu: MenuItem{
			Path:      "biz-log",
			Name:      "labReusableBizLog",
			Component: "view/lab/reusable/biz-log.vue",
			Sort:      9,
			Title:     "业务操作日志",
			Icon:      "document",
		},
		APIs: []APIItem{
			{APIGroup: "BizLog", Method: "GET", Path: "/bizLog/list", Description: "查询业务操作日志"},
			{APIGroup: "BizLog", Method: "POST", Path: "/bizLog/writeDemo", Description: "写入测试业务日志"},
		},
	})
}

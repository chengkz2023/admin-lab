package simulation

import simulationRes "github.com/flipped-aurora/gin-vue-admin/server/model/lab/simulation/response"

type CustomerDetailSimulationService struct{}

func (c *CustomerDetailSimulationService) GetDetail() simulationRes.CustomerDetailSimulation {
	return simulationRes.CustomerDetailSimulation{
		Breadcrumb: []string{"客户管理", "客户列表", "北京云讯科技有限公司"},
		Status:     "服务中",
		Customer: simulationRes.CustomerSummary{
			Name:         "北京云讯科技有限公司",
			ShortName:    "云讯",
			CustomerNo:   "C-20240318-0042",
			UnifiedCode:  "91110108MA01XXXXXX",
			ServiceType:  "互联网应用服务",
			ServiceStart: "2024-03-18",
		},
		Metrics: []simulationRes.CustomerMetric{
			{Label: "总带宽", Value: "700 Mbps", Sub: "跨 2 个机房"},
			{Label: "占用机架", Value: "5 个", Sub: "BJ-01 × 3 · SH-02 × 2"},
			{Label: "IP 地址段", Value: "8 段", Sub: "静态 5 · 动态 2 · 专线 1"},
			{Label: "应用服务", Value: "3 项", Sub: "2 对外 / 1 内部"},
		},
		BasicInfo: []simulationRes.KeyValueItem{
			{Key: "客户属性", Value: "互联网应用服务", Type: "tag-blue"},
			{Key: "单位属性", Value: "企业法人"},
			{Key: "证件类型", Value: "营业执照"},
			{Key: "证件号码", Value: "91110108MA01XXXXXX", Type: "mono"},
			{Key: "单位地址", Value: "北京市海淀区中关村大街 1 号"},
			{Key: "邮编", Value: "100080"},
			{Key: "服务开通", Value: "2024-03-18"},
		},
		SecurityOwner: []simulationRes.KeyValueItem{
			{Key: "姓名", Value: "张伟"},
			{Key: "职务", Value: "技术总监"},
			{Key: "手机", Value: "138-0000-0001", Type: "mono"},
			{Key: "邮箱", Value: "zhangwei@yunxun.com"},
		},
		RoomNodes: []simulationRes.CustomerRoomNode{
			{
				Code:    "IDC-BJ-01",
				Name:    "北京一号机房",
				SubName: "IDC-BJ-01 · 自建机房",
				Tags:    []string{"IDC", "ISP"},
				Racks:   []string{"A-02-05", "A-02-06", "B-01-12"},
				Connections: []simulationRes.RoomConnection{
					{Label: "IDC", Content: "300 Mbps · 2024-03-18"},
					{Label: "ISP", Content: "ISP-BJ-03 · 200 Mbps · 2024-04-01"},
				},
				Details: []simulationRes.RoomDetailGroup{
					{
						Label: "机房信息",
						Items: []simulationRes.KeyValueItem{
							{Key: "地址", Value: "海淀区西北旺路 10 号"},
							{Key: "机房性质", Value: "自建机房"},
							{Key: "责任人", Value: "李明 · 138-0000-0099"},
						},
					},
					{
						Label: "IDC 接入",
						Items: []simulationRes.KeyValueItem{
							{Key: "带宽", Value: "300 Mbps"},
							{Key: "分配时间", Value: "2024-03-18"},
							{Key: "机架编号", Value: "A-02-05 / A-02-06 / B-01-12"},
						},
					},
					{
						Label: "ISP 接入（ISP-BJ-03）",
						Items: []simulationRes.KeyValueItem{
							{Key: "带宽", Value: "200 Mbps"},
							{Key: "分配时间", Value: "2024-04-01"},
							{Key: "客户IP", Value: "202.96.1.100", Type: "mono"},
						},
					},
				},
			},
			{
				Code:    "IDC-SH-02",
				Name:    "上海二号机房",
				SubName: "IDC-SH-02 · 托管机房",
				Tags:    []string{"IDC"},
				Racks:   []string{"C-01-03", "C-01-04"},
				Connections: []simulationRes.RoomConnection{
					{Label: "IDC", Content: "200 Mbps · 2024-07-15"},
				},
				Details: []simulationRes.RoomDetailGroup{
					{
						Label: "机房信息",
						Items: []simulationRes.KeyValueItem{
							{Key: "地址", Value: "浦东新区张江高科技园区"},
							{Key: "机房性质", Value: "托管机房"},
							{Key: "责任人", Value: "王芳 · 139-0000-0088"},
						},
					},
					{
						Label: "IDC 接入",
						Items: []simulationRes.KeyValueItem{
							{Key: "带宽", Value: "200 Mbps"},
							{Key: "分配时间", Value: "2024-07-15"},
							{Key: "机架编号", Value: "C-01-03 / C-01-04"},
						},
					},
					{
						Label: "上联链路",
						Items: []simulationRes.KeyValueItem{
							{Key: "链路编号", Value: "LK-SH-02-001", Type: "mono"},
							{Key: "网关IP", Value: "10.16.0.1", Type: "mono"},
							{Key: "接入单位", Value: "中国电信上海分公司"},
						},
					},
				},
			},
		},
		IPSegments: []simulationRes.CustomerIPSegment{
			{
				Source:      "IDC",
				Range:       "10.8.12.0 — .255",
				Mask:        "/24",
				AddressNum:  256,
				Usage:       "静态",
				NodeCode:    "IDC-BJ-01",
				AllocatedAt: "2024-03-18",
				NATMappings: []simulationRes.NATMapping{
					{From: "10.8.12.1", To: "202.96.1.100"},
					{From: "10.8.12.2", To: "202.96.1.101"},
					{From: "10.8.12.3", To: "202.96.1.102"},
				},
			},
			{
				Source:      "IDC",
				Range:       "10.8.13.0 — .127",
				Mask:        "/25",
				AddressNum:  128,
				Usage:       "动态",
				NodeCode:    "IDC-BJ-01",
				AllocatedAt: "2024-06-01",
			},
			{
				Source:      "ISP",
				Range:       "202.96.1.96 — .103",
				Mask:        "/29",
				AddressNum:  8,
				Usage:       "静态",
				NodeCode:    "ISP-BJ-03",
				AllocatedAt: "2024-04-01",
			},
			{
				Source:      "ISP",
				Range:       "202.96.2.0 — .63",
				Mask:        "/26",
				AddressNum:  64,
				Usage:       "专线",
				NodeCode:    "ISP-BJ-03",
				AllocatedAt: "2024-04-01",
			},
			{
				Source:      "IDC",
				Range:       "10.16.0.0 — .255",
				Mask:        "/24",
				AddressNum:  256,
				Usage:       "静态",
				NodeCode:    "IDC-SH-02",
				AllocatedAt: "2024-07-15",
			},
			{
				Source:      "IDC",
				Range:       "10.16.1.0 — .63",
				Mask:        "/26",
				AddressNum:  64,
				Usage:       "动态",
				NodeCode:    "IDC-SH-02",
				AllocatedAt: "2024-09-01",
				Hidden:      true,
			},
			{
				Source:      "IDC",
				Range:       "10.16.2.0 — .255",
				Mask:        "/24",
				AddressNum:  256,
				Usage:       "静态",
				NodeCode:    "IDC-SH-02",
				AllocatedAt: "2024-10-20",
				Hidden:      true,
				NATMappings: []simulationRes.NATMapping{
					{From: "10.16.2.1", To: "121.43.88.10"},
				},
			},
			{
				Source:      "IDC",
				Range:       "172.16.4.0 — .255",
				Mask:        "/24",
				AddressNum:  256,
				Usage:       "静态",
				NodeCode:    "IDC-BJ-01",
				AllocatedAt: "2024-08-10",
				Hidden:      true,
			},
		},
		Applications: []simulationRes.CustomerApplicationEntry{
			{
				ServiceType: "对外应用",
				Permit:      "京ICP备2024XXXXXX号",
				AccessMode:  "专线接入",
				Content:     "电商交易平台",
				Domain:      "www.yunxun.com",
			},
			{
				ServiceType: "电信业务",
				Permit:      "B1-XXXXXXXXXX",
				AccessMode:  "宽带接入",
				Content:     "CDN 加速服务",
				Domain:      "cdn.yunxun.com",
			},
			{
				ServiceType: "内部应用",
				Permit:      "—",
				AccessMode:  "专线接入",
				Content:     "内部 OA 系统",
				Domain:      "oa.yunxun.internal",
			},
		},
	}
}

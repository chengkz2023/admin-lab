package response

// TableProColumn defines configurable table columns.
type TableProColumn struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// TableProRow defines a table row.
type TableProRow struct {
	ID           int     `json:"id"`
	OrderNo      string  `json:"orderNo"`
	CustomerName string  `json:"customerName"`
	Status       string  `json:"status"`
	Priority     string  `json:"priority"`
	Owner        string  `json:"owner"`
	Source       string  `json:"source"`
	Amount       float64 `json:"amount"`
	CreatedAt    string  `json:"createdAt"`
}

// TableProPageResult defines paged response payload.
type TableProPageResult struct {
	List        []TableProRow    `json:"list"`
	Total       int64            `json:"total"`
	Page        int              `json:"page"`
	PageSize    int              `json:"pageSize"`
	Columns     []TableProColumn `json:"columns"`
	Statuses    []string         `json:"statuses"`
	Owners      []string         `json:"owners"`
	DefaultSort string           `json:"defaultSort"`
}

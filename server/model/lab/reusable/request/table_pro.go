package request

// TableProPageQuery defines server-side pagination and filtering fields.
type TableProPageQuery struct {
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
	Keyword   string `json:"keyword" form:"keyword"`
	Status    string `json:"status" form:"status"`
	Owner     string `json:"owner" form:"owner"`
	StartDate string `json:"startDate" form:"startDate"`
	EndDate   string `json:"endDate" form:"endDate"`
	SortBy    string `json:"sortBy" form:"sortBy"`
	SortOrder string `json:"sortOrder" form:"sortOrder"`
}

// TableProExportQuery defines fields used by export endpoint.
type TableProExportQuery struct {
	Keyword        string   `json:"keyword" form:"keyword"`
	Status         string   `json:"status" form:"status"`
	Owner          string   `json:"owner" form:"owner"`
	StartDate      string   `json:"startDate" form:"startDate"`
	EndDate        string   `json:"endDate" form:"endDate"`
	SortBy         string   `json:"sortBy" form:"sortBy"`
	SortOrder      string   `json:"sortOrder" form:"sortOrder"`
	VisibleColumns []string `json:"visibleColumns" form:"visibleColumns"`
}

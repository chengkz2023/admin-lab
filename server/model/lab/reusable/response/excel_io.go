package response

type ExcelImportRow struct {
	RowNumber   int               `json:"rowNumber"`
	Values      map[string]string `json:"values"`
	ErrorFields []string          `json:"errorFields"`
}

type ExcelImportResult struct {
	Columns      []string         `json:"columns"`
	Rows         []ExcelImportRow `json:"rows"`
	TotalRows    int              `json:"totalRows"`
	SuccessRows  int              `json:"successRows"`
	FailedRows   int              `json:"failedRows"`
	TemplateKey  string           `json:"templateKey"`
	TemplateName string           `json:"templateName"`
}

type ExcelTemplateOption struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	FileName    string   `json:"fileName"`
	Description string   `json:"description"`
	Scene       string   `json:"scene"`
	Columns     []string `json:"columns"`
	IsDefault   bool     `json:"isDefault"`
}

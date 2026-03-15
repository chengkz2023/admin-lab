package example

import (
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	exampleRes "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

type ExcelIOService struct{}

type excelTemplateDefinition struct {
	Key             string
	Name            string
	FileName        string
	Description     string
	Scene           string
	SheetName       string
	GuideSheet      string
	Columns         []string
	GuideRows       [][]string
	SampleRows      [][]string
	AllowedStatuses []string
}

var defaultExcelTemplateKey = "basic-user"

var excelTemplates = []excelTemplateDefinition{
	{
		Key:         "basic-user",
		Name:        "基础用户导入模板",
		FileName:    "admin-lab-user-template.xlsx",
		Description: "适合用户、客户联系人、组织成员等基础信息导入。",
		Scene:       "适合先把最常见的人员资料导入链路跑通。",
		SheetName:   "数据",
		GuideSheet:  "说明",
		Columns:     []string{"姓名", "部门", "手机号", "邮箱", "状态", "备注"},
		GuideRows: [][]string{
			{"字段", "说明"},
			{"姓名", "必填，建议和真实业务字段保持一致"},
			{"部门", "必填，用于模拟组织或业务归属"},
			{"手机号", "选填，演示常见字符串字段"},
			{"邮箱", "选填，演示联系信息校验"},
			{"状态", "必填，可填写启用或停用"},
			{"备注", "选填，用于承载额外业务描述"},
		},
		SampleRows: [][]string{
			{"张三", "运营中心", "13800000001", "zhangsan@example.com", "启用", "第一批导入示例"},
			{"李四", "销售一部", "13800000002", "lisi@example.com", "停用", "用于演示状态字段"},
			{"王五", "产品部", "13800000003", "wangwu@example.com", "启用", "可直接改成业务模板"},
		},
		AllowedStatuses: []string{"启用", "停用"},
	},
	{
		Key:         "customer-order",
		Name:        "客户订单导入模板",
		FileName:    "admin-lab-customer-order-template.xlsx",
		Description: "适合批量导入客户订单、开票申请、发货任务等业务单据。",
		Scene:       "更接近真实业务仿真，适合先验证字段映射和状态流转。",
		SheetName:   "订单数据",
		GuideSheet:  "填写说明",
		Columns:     []string{"客户名称", "订单编号", "业务类型", "金额", "负责人", "状态", "备注"},
		GuideRows: [][]string{
			{"字段", "说明"},
			{"客户名称", "必填，用于识别客户主体"},
			{"订单编号", "必填，建议在真实系统中做唯一性校验"},
			{"业务类型", "必填，可扩展为字典映射"},
			{"金额", "必填，适合验证数值和精度处理"},
			{"负责人", "选填，用于匹配内部成员"},
			{"状态", "必填，可模拟待审核、执行中、已完成"},
			{"备注", "选填，承载补充说明"},
		},
		SampleRows: [][]string{
			{"华北客户A", "SO20260315001", "续费", "12800", "陈晨", "待审核", "用于导入订单预演"},
			{"华东客户B", "SO20260315002", "新签", "35600", "刘洋", "执行中", "适合搬去内网改业务校验"},
		},
		AllowedStatuses: []string{"待审核", "执行中", "已完成"},
	},
	{
		Key:         "product-catalog",
		Name:        "商品资料导入模板",
		FileName:    "admin-lab-product-template.xlsx",
		Description: "适合商品、SKU、物料、仓储档案等基础资料导入。",
		Scene:       "适合高频基础数据维护场景。",
		SheetName:   "商品资料",
		GuideSheet:  "模板说明",
		Columns:     []string{"商品编码", "商品名称", "分类", "单位", "单价", "状态", "备注"},
		GuideRows: [][]string{
			{"字段", "说明"},
			{"商品编码", "必填，通常需要全局唯一"},
			{"商品名称", "必填，用于展示和检索"},
			{"分类", "必填，可对接字典或分类树"},
			{"单位", "选填，如件、箱、套"},
			{"单价", "选填，演示金额字段格式"},
			{"状态", "必填，可填写上架或下架"},
			{"备注", "选填，补充商品说明"},
		},
		SampleRows: [][]string{
			{"SKU1001", "实验商品A", "数码配件", "件", "59.90", "上架", "适合测试商品批量导入"},
			{"SKU1002", "实验商品B", "办公耗材", "箱", "120.00", "下架", "可扩展库存字段"},
		},
		AllowedStatuses: []string{"上架", "下架"},
	},
}

func (e *ExcelIOService) ListTemplates() []exampleRes.ExcelTemplateOption {
	options := make([]exampleRes.ExcelTemplateOption, 0, len(excelTemplates))
	for _, template := range excelTemplates {
		options = append(options, exampleRes.ExcelTemplateOption{
			Key:         template.Key,
			Name:        template.Name,
			FileName:    template.FileName,
			Description: template.Description,
			Scene:       template.Scene,
			Columns:     template.Columns,
			IsDefault:   template.Key == defaultExcelTemplateKey,
		})
	}
	return options
}

func (e *ExcelIOService) ExportTemplate(templateKey string) ([]byte, string, error) {
	template, err := findExcelTemplate(templateKey)
	if err != nil {
		return nil, "", err
	}

	file := excelize.NewFile()
	file.SetSheetName("Sheet1", template.SheetName)
	file.NewSheet(template.GuideSheet)

	for idx, header := range template.Columns {
		cell, _ := excelize.CoordinatesToCellName(idx+1, 1)
		file.SetCellValue(template.SheetName, cell, header)
	}

	for rowIndex, row := range template.GuideRows {
		for colIndex, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			file.SetCellValue(template.GuideSheet, cell, value)
		}
	}

	return writeExcelFile(file, template.FileName, len(template.Columns))
}

func (e *ExcelIOService) ExportSample() ([]byte, string, error) {
	template, err := findExcelTemplate(defaultExcelTemplateKey)
	if err != nil {
		return nil, "", err
	}

	file := excelize.NewFile()
	file.SetSheetName("Sheet1", template.SheetName)

	for idx, header := range template.Columns {
		cell, _ := excelize.CoordinatesToCellName(idx+1, 1)
		file.SetCellValue(template.SheetName, cell, header)
	}

	for rowIndex, row := range template.SampleRows {
		for colIndex, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			file.SetCellValue(template.SheetName, cell, value)
		}
	}

	return writeExcelFile(file, "admin-lab-excel-sample.xlsx", len(template.Columns))
}

func (e *ExcelIOService) ImportExcel(templateKey string, header *multipart.FileHeader) (exampleRes.ExcelImportResult, error) {
	template, err := findExcelTemplate(templateKey)
	if err != nil {
		return exampleRes.ExcelImportResult{}, err
	}

	result := exampleRes.ExcelImportResult{
		Columns:      template.Columns,
		TemplateKey:  template.Key,
		TemplateName: template.Name,
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".xlsx" {
		return result, errors.New("仅支持 .xlsx 格式的 Excel 文件")
	}

	file, err := header.Open()
	if err != nil {
		return result, errors.Wrap(err, "打开上传文件失败")
	}
	defer file.Close()

	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		return result, errors.Wrap(err, "解析 Excel 文件失败")
	}
	defer func() {
		_ = excelFile.Close()
	}()

	sheetName := template.SheetName
	if index, _ := excelFile.GetSheetIndex(sheetName); index == -1 {
		sheets := excelFile.GetSheetList()
		if len(sheets) == 0 {
			return result, errors.New("Excel 中未找到可读取的工作表")
		}
		sheetName = sheets[0]
	}

	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return result, errors.Wrap(err, "读取 Excel 行数据失败")
	}
	if len(rows) <= 1 {
		return result, errors.New("Excel 中没有可导入的数据行")
	}

	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]
		if isEmptyRow(row) {
			continue
		}

		values := make(map[string]string, len(template.Columns))
		for index, column := range template.Columns {
			values[column] = safeCell(row, index)
		}

		item := exampleRes.ExcelImportRow{
			RowNumber:   rowIndex + 1,
			Values:      values,
			ErrorFields: validateImportRow(template, values),
		}

		if len(item.ErrorFields) == 0 {
			result.SuccessRows++
		} else {
			result.FailedRows++
		}

		result.Rows = append(result.Rows, item)
	}

	result.TotalRows = len(result.Rows)
	if result.TotalRows == 0 {
		return result, errors.New("Excel 中没有有效的数据行")
	}

	return result, nil
}

func writeExcelFile(file *excelize.File, fileName string, columnCount int) ([]byte, string, error) {
	sheetName := file.GetSheetList()[0]
	endCell, _ := excelize.CoordinatesToCellName(columnCount, 1)

	styleID, err := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E8F3FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err == nil {
		_ = file.SetCellStyle(sheetName, "A1", endCell, styleID)
	}
	_ = file.SetColWidth(sheetName, "A", "H", 18)

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, "", errors.Wrap(err, "生成 Excel 文件失败")
	}
	return buffer.Bytes(), fileName, nil
}

func findExcelTemplate(templateKey string) (excelTemplateDefinition, error) {
	if templateKey == "" {
		templateKey = defaultExcelTemplateKey
	}

	for _, template := range excelTemplates {
		if template.Key == templateKey {
			return template, nil
		}
	}

	return excelTemplateDefinition{}, errors.New("未找到对应的 Excel 模板")
}

func safeCell(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func validateImportRow(template excelTemplateDefinition, values map[string]string) []string {
	switch template.Key {
	case "customer-order":
		return validateCustomerOrderRow(template, values)
	case "product-catalog":
		return validateProductCatalogRow(template, values)
	default:
		return validateBasicUserRow(template, values)
	}
}

func validateBasicUserRow(template excelTemplateDefinition, values map[string]string) []string {
	var errorFields []string
	if values["姓名"] == "" {
		errorFields = append(errorFields, "姓名不能为空")
	}
	if values["部门"] == "" {
		errorFields = append(errorFields, "部门不能为空")
	}
	if statusError := validateStatus(values["状态"], template.AllowedStatuses, "状态"); statusError != "" {
		errorFields = append(errorFields, statusError)
	}
	if email := values["邮箱"]; email != "" && !strings.Contains(email, "@") {
		errorFields = append(errorFields, "邮箱格式不正确")
	}
	if phone := values["手机号"]; phone != "" && len([]rune(phone)) < 6 {
		errorFields = append(errorFields, "手机号长度过短")
	}
	return errorFields
}

func validateCustomerOrderRow(template excelTemplateDefinition, values map[string]string) []string {
	var errorFields []string
	if values["客户名称"] == "" {
		errorFields = append(errorFields, "客户名称不能为空")
	}
	if values["订单编号"] == "" {
		errorFields = append(errorFields, "订单编号不能为空")
	}
	if values["业务类型"] == "" {
		errorFields = append(errorFields, "业务类型不能为空")
	}
	if amount := values["金额"]; amount == "" {
		errorFields = append(errorFields, "金额不能为空")
	} else if !isPositiveNumber(amount) {
		errorFields = append(errorFields, "金额必须是大于 0 的数字")
	}
	if statusError := validateStatus(values["状态"], template.AllowedStatuses, "状态"); statusError != "" {
		errorFields = append(errorFields, statusError)
	}
	return errorFields
}

func validateProductCatalogRow(template excelTemplateDefinition, values map[string]string) []string {
	var errorFields []string
	if values["商品编码"] == "" {
		errorFields = append(errorFields, "商品编码不能为空")
	}
	if values["商品名称"] == "" {
		errorFields = append(errorFields, "商品名称不能为空")
	}
	if values["分类"] == "" {
		errorFields = append(errorFields, "分类不能为空")
	}
	if price := values["单价"]; price != "" && !isPositiveNumber(price) {
		errorFields = append(errorFields, "单价必须是合法数字")
	}
	if statusError := validateStatus(values["状态"], template.AllowedStatuses, "状态"); statusError != "" {
		errorFields = append(errorFields, statusError)
	}
	return errorFields
}

func validateStatus(value string, allowed []string, field string) string {
	if value == "" {
		return field + "不能为空"
	}
	for _, item := range allowed {
		if value == item {
			return ""
		}
	}
	return field + "不在允许范围内"
}

func isPositiveNumber(value string) bool {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	return number > 0
}

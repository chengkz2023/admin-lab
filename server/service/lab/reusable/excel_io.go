package reusable

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	reusableRes "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/response"
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
	Dropdowns       []excelDropdownValidation
	Cascade         *excelCascadeValidation
}

type excelDropdownValidation struct {
	Column  string
	Options []string
	Prompt  string
}

type excelCascadeValidation struct {
	ProvinceColumn string
	CityColumn     string
	CountyColumn   string
	Regions        map[string]map[string][]string
}

const (
	templateValidationMaxRow = 200
	templateOptionSheetName  = "选项数据"
)

var defaultExcelTemplateKey = "basic-user"

var excelTemplates = []excelTemplateDefinition{
	{
		Key:         "basic-user",
		Name:        "基础用户导入模板",
		FileName:    "admin-lab-user-template.xlsx",
		Description: "适合用户、客户联系人、组织成员等基础信息导入。",
		Scene:       "示范了普通枚举下拉和省/市/区县级联下拉两类高频场景。",
		SheetName:   "数据",
		GuideSheet:  "说明",
		Columns:     []string{"姓名", "部门", "手机号", "邮箱", "角色类型", "状态", "省", "市", "区县", "备注"},
		GuideRows: [][]string{
			{"字段", "说明"},
			{"姓名", "必填，建议和真实业务字段保持一致"},
			{"部门", "必填，模板内已配置下拉列表"},
			{"手机号", "选填，演示常见字符串字段"},
			{"邮箱", "选填，演示联系信息校验"},
			{"角色类型", "必填，模板内已配置下拉列表"},
			{"状态", "必填，模板内已配置下拉列表"},
			{"省 / 市 / 区县", "演示级联下拉，先选省，再选市，最后选区县"},
			{"备注", "选填，用于承载额外业务描述"},
		},
		SampleRows: [][]string{
			{"张三", "运营中心", "13800000001", "zhangsan@example.com", "运营", "启用", "浙江省", "杭州市", "西湖区", "第一批导入示例"},
			{"李四", "销售一部", "13800000002", "lisi@example.com", "销售", "停用", "江苏省", "苏州市", "工业园区", "用于演示状态字段"},
			{"王五", "产品部", "13800000003", "wangwu@example.com", "管理员", "启用", "广东省", "深圳市", "南山区", "可直接改成业务模板"},
		},
		AllowedStatuses: []string{"启用", "停用"},
		Dropdowns: []excelDropdownValidation{
			{Column: "部门", Options: []string{"运营中心", "销售一部", "产品部", "客服部"}, Prompt: "请选择部门"},
			{Column: "角色类型", Options: []string{"管理员", "运营", "销售", "访客"}, Prompt: "请选择角色类型"},
			{Column: "状态", Options: []string{"启用", "停用"}, Prompt: "请选择状态"},
		},
		Cascade: &excelCascadeValidation{
			ProvinceColumn: "省",
			CityColumn:     "市",
			CountyColumn:   "区县",
			Regions: map[string]map[string][]string{
				"浙江省": {
					"杭州市": {"西湖区", "余杭区", "滨江区"},
					"宁波市": {"海曙区", "鄞州区", "江北区"},
				},
				"江苏省": {
					"南京市": {"玄武区", "建邺区", "鼓楼区"},
					"苏州市": {"姑苏区", "工业园区", "吴中区"},
				},
				"广东省": {
					"广州市": {"天河区", "越秀区", "海珠区"},
					"深圳市": {"南山区", "福田区", "龙岗区"},
				},
			},
		},
	},
	{
		Key:         "customer-order",
		Name:        "客户订单导入模板",
		FileName:    "admin-lab-customer-order-template.xlsx",
		Description: "适合批量导入客户订单、开票申请、发货任务等业务单据。",
		Scene:       "适合先验证业务类型、状态、负责人等枚举字段的下拉配置。",
		SheetName:   "订单数据",
		GuideSheet:  "填写说明",
		Columns:     []string{"客户名称", "订单编号", "业务类型", "金额", "负责人", "状态", "备注"},
		GuideRows: [][]string{
			{"字段", "说明"},
			{"客户名称", "必填，用于识别客户主体"},
			{"订单编号", "必填，建议在真实系统中做唯一性校验"},
			{"业务类型", "必填，模板内已配置下拉列表"},
			{"金额", "必填，适合验证数值和精度处理"},
			{"负责人", "选填，模板内已配置下拉列表"},
			{"状态", "必填，模板内已配置下拉列表"},
			{"备注", "选填，承载补充说明"},
		},
		SampleRows: [][]string{
			{"华北客户A", "SO20260315001", "续费", "12800", "陈晨", "待审核", "用于导入订单预演"},
			{"华东客户B", "SO20260315002", "新签", "35600", "刘洋", "执行中", "适合搬去内网改业务校验"},
		},
		AllowedStatuses: []string{"待审核", "执行中", "已完成"},
		Dropdowns: []excelDropdownValidation{
			{Column: "业务类型", Options: []string{"新签", "续费", "升级", "补单"}, Prompt: "请选择业务类型"},
			{Column: "负责人", Options: []string{"陈晨", "刘洋", "王欣", "赵雷"}, Prompt: "请选择负责人"},
			{Column: "状态", Options: []string{"待审核", "执行中", "已完成"}, Prompt: "请选择状态"},
		},
	},
	{
		Key:         "product-catalog",
		Name:        "商品资料导入模板",
		FileName:    "admin-lab-product-template.xlsx",
		Description: "适合商品、SKU、物料、仓储档案等基础资料导入。",
		Scene:       "适合高频基础数据维护场景，常见于分类、单位、状态等下拉录入。",
		SheetName:   "商品资料",
		GuideSheet:  "模板说明",
		Columns:     []string{"商品编码", "商品名称", "分类", "单位", "单价", "状态", "备注"},
		GuideRows: [][]string{
			{"字段", "说明"},
			{"商品编码", "必填，通常需要全局唯一"},
			{"商品名称", "必填，用于展示和检索"},
			{"分类", "必填，模板内已配置下拉列表"},
			{"单位", "选填，模板内已配置下拉列表"},
			{"单价", "选填，演示金额字段格式"},
			{"状态", "必填，模板内已配置下拉列表"},
			{"备注", "选填，补充商品说明"},
		},
		SampleRows: [][]string{
			{"SKU1001", "实验商品A", "数码配件", "件", "59.90", "上架", "适合测试商品批量导入"},
			{"SKU1002", "实验商品B", "办公耗材", "箱", "120.00", "下架", "可扩展库存字段"},
		},
		AllowedStatuses: []string{"上架", "下架"},
		Dropdowns: []excelDropdownValidation{
			{Column: "分类", Options: []string{"数码配件", "办公耗材", "家居用品", "测试物料"}, Prompt: "请选择商品分类"},
			{Column: "单位", Options: []string{"件", "箱", "套", "个"}, Prompt: "请选择单位"},
			{Column: "状态", Options: []string{"上架", "下架"}, Prompt: "请选择状态"},
		},
	},
}

func (e *ExcelIOService) ListTemplates() []reusableRes.ExcelTemplateOption {
	options := make([]reusableRes.ExcelTemplateOption, 0, len(excelTemplates))
	for _, template := range excelTemplates {
		options = append(options, reusableRes.ExcelTemplateOption{
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

	if err = applyTemplateValidations(file, template); err != nil {
		return nil, "", err
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

	if err = applyTemplateValidations(file, template); err != nil {
		return nil, "", err
	}

	return writeExcelFile(file, "admin-lab-excel-sample.xlsx", len(template.Columns))
}

func (e *ExcelIOService) ImportExcel(templateKey string, header *multipart.FileHeader) (reusableRes.ExcelImportResult, error) {
	template, err := findExcelTemplate(templateKey)
	if err != nil {
		return reusableRes.ExcelImportResult{}, err
	}

	result := reusableRes.ExcelImportResult{
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

		item := reusableRes.ExcelImportRow{
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

func applyTemplateValidations(file *excelize.File, template excelTemplateDefinition) error {
	if len(template.Dropdowns) == 0 && template.Cascade == nil {
		return nil
	}

	optionSheet := templateOptionSheetName
	if index, _ := file.GetSheetIndex(optionSheet); index == -1 {
		file.NewSheet(optionSheet)
	}

	optionColumnIndex := 1
	for _, dropdown := range template.Dropdowns {
		if len(dropdown.Options) == 0 {
			continue
		}
		if err := applyDropdownValidation(file, template, optionSheet, dropdown, &optionColumnIndex); err != nil {
			return err
		}
	}

	if template.Cascade != nil {
		if err := applyCascadeValidation(file, template, optionSheet, *template.Cascade, &optionColumnIndex); err != nil {
			return err
		}
	}

	return file.SetSheetVisible(optionSheet, false, true)
}

func applyDropdownValidation(file *excelize.File, template excelTemplateDefinition, optionSheet string, dropdown excelDropdownValidation, optionColumnIndex *int) error {
	targetColumnIndex := findColumnIndex(template.Columns, dropdown.Column)
	if targetColumnIndex == 0 {
		return nil
	}

	optionColumnName, _ := excelize.ColumnNumberToName(*optionColumnIndex)
	file.SetCellValue(optionSheet, optionColumnName+"1", dropdown.Column)
	for index, item := range dropdown.Options {
		file.SetCellValue(optionSheet, fmt.Sprintf("%s%d", optionColumnName, index+2), item)
	}

	targetColumnName, _ := excelize.ColumnNumberToName(targetColumnIndex)
	dv := excelize.NewDataValidation(true)
	dv.Sqref = fmt.Sprintf("%s2:%s%d", targetColumnName, targetColumnName, templateValidationMaxRow)
	dv.SetSqrefDropList(fmt.Sprintf("'%s'!$%s$2:$%s$%d", optionSheet, optionColumnName, optionColumnName, len(dropdown.Options)+1))
	dv.SetInput("请选择", dropdown.Prompt)
	dv.SetError(excelize.DataValidationErrorStyleStop, "输入有误", "请从下拉列表中选择有效值")
	if err := file.AddDataValidation(template.SheetName, dv); err != nil {
		return err
	}

	*optionColumnIndex = *optionColumnIndex + 1
	return nil
}

func applyCascadeValidation(file *excelize.File, template excelTemplateDefinition, optionSheet string, cascade excelCascadeValidation, optionColumnIndex *int) error {
	provinceColumnIndex := findColumnIndex(template.Columns, cascade.ProvinceColumn)
	cityColumnIndex := findColumnIndex(template.Columns, cascade.CityColumn)
	countyColumnIndex := findColumnIndex(template.Columns, cascade.CountyColumn)
	if provinceColumnIndex == 0 || cityColumnIndex == 0 || countyColumnIndex == 0 {
		return nil
	}

	provinceOptionColumn, _ := excelize.ColumnNumberToName(*optionColumnIndex)
	file.SetCellValue(optionSheet, provinceOptionColumn+"1", "省份")

	provinces := orderedKeys(cascade.Regions)
	for index, province := range provinces {
		file.SetCellValue(optionSheet, fmt.Sprintf("%s%d", provinceOptionColumn, index+2), province)
	}

	*optionColumnIndex = *optionColumnIndex + 1
	for _, province := range provinces {
		cityOptionColumn, _ := excelize.ColumnNumberToName(*optionColumnIndex)
		defineName := buildDefinedName("R", province)
		file.SetCellValue(optionSheet, cityOptionColumn+"1", defineName)

		cities := orderedKeys(cascade.Regions[province])
		for index, city := range cities {
			file.SetCellValue(optionSheet, fmt.Sprintf("%s%d", cityOptionColumn, index+2), city)
		}

		if len(cities) > 0 {
			if err := file.SetDefinedName(&excelize.DefinedName{
				Name:     defineName,
				RefersTo: fmt.Sprintf("'%s'!$%s$2:$%s$%d", optionSheet, cityOptionColumn, cityOptionColumn, len(cities)+1),
			}); err != nil {
				return err
			}
		}

		*optionColumnIndex = *optionColumnIndex + 1
		for _, city := range cities {
			countyOptionColumn, _ := excelize.ColumnNumberToName(*optionColumnIndex)
			countyDefineName := buildDefinedName("C", province, city)
			file.SetCellValue(optionSheet, countyOptionColumn+"1", countyDefineName)

			counties := cascade.Regions[province][city]
			for index, county := range counties {
				file.SetCellValue(optionSheet, fmt.Sprintf("%s%d", countyOptionColumn, index+2), county)
			}

			if len(counties) > 0 {
				if err := file.SetDefinedName(&excelize.DefinedName{
					Name:     countyDefineName,
					RefersTo: fmt.Sprintf("'%s'!$%s$2:$%s$%d", optionSheet, countyOptionColumn, countyOptionColumn, len(counties)+1),
				}); err != nil {
					return err
				}
			}

			*optionColumnIndex = *optionColumnIndex + 1
		}
	}

	provinceColumnName, _ := excelize.ColumnNumberToName(provinceColumnIndex)
	cityColumnName, _ := excelize.ColumnNumberToName(cityColumnIndex)
	countyColumnName, _ := excelize.ColumnNumberToName(countyColumnIndex)

	provinceDV := excelize.NewDataValidation(true)
	provinceDV.Sqref = fmt.Sprintf("%s2:%s%d", provinceColumnName, provinceColumnName, templateValidationMaxRow)
	provinceDV.SetSqrefDropList(fmt.Sprintf("'%s'!$%s$2:$%s$%d", optionSheet, provinceOptionColumn, provinceOptionColumn, len(provinces)+1))
	provinceDV.SetInput("请选择省份", "先选择省份，再继续选择市和区县")
	provinceDV.SetError(excelize.DataValidationErrorStyleStop, "输入有误", "请选择有效的省份")
	if err := file.AddDataValidation(template.SheetName, provinceDV); err != nil {
		return err
	}

	cityDV := excelize.NewDataValidation(true)
	cityDV.Sqref = fmt.Sprintf("%s2:%s%d", cityColumnName, cityColumnName, templateValidationMaxRow)
	if err := cityDV.SetDropList([]string{fmt.Sprintf(`=INDIRECT("R_"&SUBSTITUTE($%s2," ",""))`, provinceColumnName)}); err != nil {
		return err
	}
	cityDV.SetInput("请选择城市", "城市列表会跟随省份变化")
	cityDV.SetError(excelize.DataValidationErrorStyleStop, "输入有误", "请先选择省份，再从下拉列表中选择城市")
	if err := file.AddDataValidation(template.SheetName, cityDV); err != nil {
		return err
	}

	countyDV := excelize.NewDataValidation(true)
	countyDV.Sqref = fmt.Sprintf("%s2:%s%d", countyColumnName, countyColumnName, templateValidationMaxRow)
	if err := countyDV.SetDropList([]string{fmt.Sprintf(`=INDIRECT("C_"&SUBSTITUTE($%s2," ","")&"_"&SUBSTITUTE($%s2," ",""))`, provinceColumnName, cityColumnName)}); err != nil {
		return err
	}
	countyDV.SetInput("请选择区县", "区县列表会跟随省份和城市变化")
	countyDV.SetError(excelize.DataValidationErrorStyleStop, "输入有误", "请先选择省份和城市，再从下拉列表中选择区县")
	return file.AddDataValidation(template.SheetName, countyDV)
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
	_ = file.SetColWidth(sheetName, "A", "J", 18)

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

func findColumnIndex(columns []string, target string) int {
	for index, column := range columns {
		if column == target {
			return index + 1
		}
	}
	return 0
}

func buildDefinedName(prefix string, values ...string) string {
	parts := []string{prefix}
	for _, value := range values {
		parts = append(parts, sanitizeDefinedNamePart(value))
	}
	return strings.Join(parts, "_")
}

func sanitizeDefinedNamePart(value string) string {
	replacer := strings.NewReplacer(
		" ", "",
		"-", "_",
		"/", "_",
		"\\", "_",
		"（", "_",
		"）", "_",
		"(", "_",
		")", "_",
		".", "_",
		",", "_",
		"，", "_",
		"、", "_",
	)
	return replacer.Replace(value)
}

func orderedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	// Keep a stable order for template output.
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
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
	if values["角色类型"] == "" {
		errorFields = append(errorFields, "角色类型不能为空")
	}
	if statusError := validateStatus(values["状态"], template.AllowedStatuses, "状态"); statusError != "" {
		errorFields = append(errorFields, statusError)
	}
	if values["省"] == "" {
		errorFields = append(errorFields, "省不能为空")
	}
	if values["市"] == "" {
		errorFields = append(errorFields, "市不能为空")
	}
	if values["区县"] == "" {
		errorFields = append(errorFields, "区县不能为空")
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

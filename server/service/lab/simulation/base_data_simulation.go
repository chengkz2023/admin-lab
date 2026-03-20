package simulation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/response"
	"mime/multipart"
	"net"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type BaseDataSimulationService struct{}

type baseDataTemplate struct {
	Key         string
	Name        string
	FileName    string
	Description string
	Scene       string
	SheetName   string
	Columns     []string
	Required    map[string]bool
	MockRows    [][]string
}

type addressNodeRaw struct {
	Code     any              `json:"code"`
	Name     string           `json:"name"`
	Children []addressNodeRaw `json:"children"`
}

type addressWrapper struct {
	Address []addressNodeRaw `json:"address"`
}

type addressNode struct {
	Code     string
	Name     string
	Children []addressNode
}

const (
	defaultBaseDataTemplate = "machine-room"
	addressJSONPath         = `C:\Users\Administrator\Downloads\address.json`
	templateDataStartRow    = 2
	templateDataEndRow      = 5000
	addressSheetName        = "AddressDict"
)

var baseDataTemplates = []baseDataTemplate{
	{
		Key:         "machine-room",
		Name:        "机房模板",
		FileName:    "base-data-machine-room-template.xlsx",
		Description: "基础数据机房导入模板",
		Scene:       "基础数据 > 机房",
		SheetName:   "MachineRoom",
		Columns:     []string{"机房名称", "机房编码", "所在城市", "地址", "联系人", "联系电话", "状态", "备注"},
		Required: map[string]bool{
			"机房名称": true,
			"机房编码": true,
			"所在城市": true,
			"状态":   true,
		},
		MockRows: [][]string{
			{"浦东机房A", "IDC-PD-001", "上海", "浦东新区1号", "张三", "13800001111", "启用", "仿真数据"},
			{"海淀机房B", "IDC-HD-002", "北京", "海淀区2号", "李四", "13800002222", "停用", "仿真数据"},
		},
	},
	{
		Key:         "service-user",
		Name:        "用户模板（服务用户）",
		FileName:    "base-data-service-user-template.xlsx",
		Description: "服务用户导入模板（含服务、IP等）",
		Scene:       "基础数据 > 用户（服务用户）",
		SheetName:   "ServiceUser",
		Columns:     []string{"用户名称", "用户账号", "所属部门", "省", "市", "县", "服务名称", "IP地址", "链路名称", "状态", "备注"},
		Required: map[string]bool{
			"用户名称": true,
			"用户账号": true,
			"服务名称": true,
			"IP地址": true,
			"状态":   true,
		},
		MockRows: [][]string{
			{"政企客户A", "corp_a", "政企一部", "110000-北京市", "110100-市辖区", "110101-东城区", "云专线", "10.10.1.11", "华东主链路", "启用", "仿真数据"},
			{"政企客户B", "corp_b", "政企二部", "310000-上海市", "310100-市辖区", "310101-黄浦区", "云防护", "10.10.1.12", "华北主链路", "停用", "仿真数据"},
		},
	},
	{
		Key:         "other-user",
		Name:        "用户模板（其他用户）",
		FileName:    "base-data-other-user-template.xlsx",
		Description: "其他用户导入模板",
		Scene:       "基础数据 > 用户（其他用户）",
		SheetName:   "OtherUser",
		Columns:     []string{"用户名称", "用户账号", "所属部门", "省", "市", "县", "联系电话", "邮箱", "状态", "备注"},
		Required: map[string]bool{
			"用户名称": true,
			"用户账号": true,
			"所属部门": true,
			"状态":   true,
		},
		MockRows: [][]string{
			{"内部用户A", "inner_a", "综合部", "110000-北京市", "110100-市辖区", "110101-东城区", "13800002222", "inner_a@example.com", "启用", "仿真数据"},
			{"内部用户B", "inner_b", "运营部", "310000-上海市", "310100-市辖区", "310101-黄浦区", "13800003333", "inner_b@example.com", "停用", "仿真数据"},
		},
	},
}

func (s *BaseDataSimulationService) ListTemplates() []response.ExcelTemplateOption {
	options := make([]response.ExcelTemplateOption, 0, len(baseDataTemplates))
	for _, tpl := range baseDataTemplates {
		options = append(options, response.ExcelTemplateOption{
			Key:         tpl.Key,
			Name:        tpl.Name,
			FileName:    tpl.FileName,
			Description: tpl.Description,
			Scene:       tpl.Scene,
			Columns:     tpl.Columns,
			IsDefault:   tpl.Key == defaultBaseDataTemplate,
		})
	}
	return options
}

func (s *BaseDataSimulationService) DownloadTemplate(templateKey string) ([]byte, string, error) {
	tpl, err := findBaseDataTemplate(templateKey)
	if err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	if err = f.SetSheetName("Sheet1", tpl.SheetName); err != nil {
		return nil, "", errors.Wrap(err, "设置模板页失败")
	}
	guide := "Guide"
	if _, err = f.NewSheet(guide); err != nil {
		return nil, "", errors.Wrap(err, "创建说明页失败")
	}

	for i, col := range tpl.Columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		header := col
		if tpl.Required[col] {
			header += " *"
		}
		if err = f.SetCellValue(tpl.SheetName, cell, header); err != nil {
			return nil, "", errors.Wrap(err, "写入模板表头失败")
		}
	}

	if err = f.SetCellValue(guide, "A1", "字段"); err != nil {
		return nil, "", errors.Wrap(err, "写入说明页失败")
	}
	if err = f.SetCellValue(guide, "B1", "是否必填"); err != nil {
		return nil, "", errors.Wrap(err, "写入说明页失败")
	}
	if err = f.SetCellValue(guide, "C1", "说明"); err != nil {
		return nil, "", errors.Wrap(err, "写入说明页失败")
	}
	for i, col := range tpl.Columns {
		row := i + 2
		c1, _ := excelize.CoordinatesToCellName(1, row)
		c2, _ := excelize.CoordinatesToCellName(2, row)
		c3, _ := excelize.CoordinatesToCellName(3, row)
		if err = f.SetCellValue(guide, c1, col); err != nil {
			return nil, "", errors.Wrap(err, "写入说明页字段失败")
		}
		if tpl.Required[col] {
			if err = f.SetCellValue(guide, c2, "是"); err != nil {
				return nil, "", errors.Wrap(err, "写入说明页必填标识失败")
			}
		} else {
			if err = f.SetCellValue(guide, c2, "否"); err != nil {
				return nil, "", errors.Wrap(err, "写入说明页必填标识失败")
			}
		}
		desc := "仿真模板字段，正式校验规则待后续实现"
		if col == "省" || col == "市" || col == "县" {
			desc = "下拉选择，格式为 编号-名称；市/县会根据上级联动过滤"
		}
		if err = f.SetCellValue(guide, c3, desc); err != nil {
			return nil, "", errors.Wrap(err, "写入说明页描述失败")
		}
	}

	if isUserTemplate(tpl.Key) {
		if err = addAddressCascadeValidation(f, tpl); err != nil {
			return nil, "", err
		}
	}

	return writeBaseDataExcel(f, tpl.FileName, len(tpl.Columns))
}

func (s *BaseDataSimulationService) ExportData(templateKey string) ([]byte, string, error) {
	tpl, err := findBaseDataTemplate(templateKey)
	if err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	if err = f.SetSheetName("Sheet1", tpl.SheetName); err != nil {
		return nil, "", errors.Wrap(err, "设置导出页失败")
	}
	for i, col := range tpl.Columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err = f.SetCellValue(tpl.SheetName, cell, col); err != nil {
			return nil, "", errors.Wrap(err, "写入导出表头失败")
		}
	}
	for r, row := range tpl.MockRows {
		for c, val := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			if err = f.SetCellValue(tpl.SheetName, cell, val); err != nil {
				return nil, "", errors.Wrap(err, "写入导出示例数据失败")
			}
		}
	}

	fileName := strings.Replace(tpl.FileName, "-template.xlsx", "-export-demo.xlsx", 1)
	return writeBaseDataExcel(f, fileName, len(tpl.Columns))
}

func (s *BaseDataSimulationService) ImportData(templateKey string, header *multipart.FileHeader) (response.ExcelImportResult, error) {
	tpl, err := findBaseDataTemplate(templateKey)
	if err != nil {
		return response.ExcelImportResult{}, err
	}

	result := response.ExcelImportResult{
		Columns:      tpl.Columns,
		Rows:         []response.ExcelImportRow{},
		TemplateKey:  tpl.Key,
		TemplateName: tpl.Name,
	}

	if strings.ToLower(filepath.Ext(header.Filename)) != ".xlsx" {
		return result, errors.New("仅支持 .xlsx 文件")
	}

	file, err := header.Open()
	if err != nil {
		return result, errors.Wrap(err, "打开上传文件失败")
	}
	defer file.Close()

	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		return result, errors.Wrap(err, "解析 Excel 失败")
	}
	defer func() {
		_ = excelFile.Close()
	}()

	sheetName := tpl.SheetName
	if index, _ := excelFile.GetSheetIndex(sheetName); index == -1 {
		sheets := excelFile.GetSheetList()
		if len(sheets) == 0 {
			return result, errors.New("Excel 中没有可读取的工作表")
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

		values := make(map[string]string, len(tpl.Columns))
		for colIndex, colName := range tpl.Columns {
			values[colName] = safeCell(row, colIndex)
		}

		item := response.ExcelImportRow{
			RowNumber:   rowIndex + 1,
			Values:      values,
			ErrorFields: validateSimulationRow(tpl, values),
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

func findBaseDataTemplate(templateKey string) (baseDataTemplate, error) {
	key := templateKey
	if key == "" {
		key = defaultBaseDataTemplate
	}
	for _, tpl := range baseDataTemplates {
		if tpl.Key == key {
			return tpl, nil
		}
	}
	return baseDataTemplate{}, errors.New("未找到对应模板")
}

func writeBaseDataExcel(f *excelize.File, fileName string, columnCount int) ([]byte, string, error) {
	sheet := f.GetSheetList()[0]
	end, _ := excelize.CoordinatesToCellName(columnCount, 1)
	style, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E8F3FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err == nil {
		_ = f.SetCellStyle(sheet, "A1", end, style)
	}
	lastCol, _ := excelize.ColumnNumberToName(columnCount)
	_ = f.SetColWidth(sheet, "A", lastCol, 18)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", errors.Wrap(err, "生成 Excel 失败")
	}
	return buf.Bytes(), fileName, nil
}

func validateSimulationRow(tpl baseDataTemplate, values map[string]string) []string {
	errs := make([]string, 0)

	for field, required := range tpl.Required {
		if required && strings.TrimSpace(values[field]) == "" {
			errs = append(errs, field+"不能为空")
		}
	}

	for _, field := range tpl.Columns {
		val := values[field]
		maxLen := 100
		if strings.Contains(field, "备注") {
			maxLen = 200
		}
		if utf8.RuneCountInString(val) > maxLen {
			errs = append(errs, field+"长度超出限制")
		}
	}

	switch tpl.Key {
	case "machine-room":
		errs = append(errs, validateMachineRoomRow(tpl, values)...)
	case "service-user":
		errs = append(errs, validateServiceUserRow(tpl, values)...)
	case "other-user":
		errs = append(errs, validateOtherUserRow(tpl, values)...)
	}
	return errs
}

func validateMachineRoomRow(tpl baseDataTemplate, values map[string]string) []string {
	errs := make([]string, 0)
	if len(tpl.Columns) < 8 {
		return errs
	}
	code := values[tpl.Columns[1]]
	if code != "" && !regexp.MustCompile(`^[A-Za-z0-9-]{3,32}$`).MatchString(code) {
		errs = append(errs, tpl.Columns[1]+"格式不合法（3-32位字母/数字/-）")
	}
	phone := values[tpl.Columns[5]]
	if phone != "" && !regexp.MustCompile(`^1\d{10}$`).MatchString(phone) {
		errs = append(errs, tpl.Columns[5]+"格式不合法（需11位手机号）")
	}
	status := values[tpl.Columns[6]]
	if status != "" && status != "启用" && status != "停用" {
		errs = append(errs, tpl.Columns[6]+"仅支持：启用/停用")
	}
	return errs
}

func validateServiceUserRow(tpl baseDataTemplate, values map[string]string) []string {
	errs := make([]string, 0)
	if len(tpl.Columns) < 11 {
		return errs
	}
	account := values[tpl.Columns[1]]
	if account != "" && !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{2,31}$`).MatchString(account) {
		errs = append(errs, tpl.Columns[1]+"格式不合法（字母开头，3-32位）")
	}
	ip := values[tpl.Columns[7]]
	if ip != "" && net.ParseIP(ip) == nil {
		errs = append(errs, tpl.Columns[7]+"格式不合法（IP地址）")
	}
	status := values[tpl.Columns[9]]
	if status != "" && status != "启用" && status != "停用" {
		errs = append(errs, tpl.Columns[9]+"仅支持：启用/停用")
	}
	return errs
}

func validateOtherUserRow(tpl baseDataTemplate, values map[string]string) []string {
	errs := make([]string, 0)
	if len(tpl.Columns) < 10 {
		return errs
	}
	account := values[tpl.Columns[1]]
	if account != "" && !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{2,31}$`).MatchString(account) {
		errs = append(errs, tpl.Columns[1]+"格式不合法（字母开头，3-32位）")
	}
	phone := values[tpl.Columns[6]]
	if phone != "" && !regexp.MustCompile(`^1\d{10}$`).MatchString(phone) {
		errs = append(errs, tpl.Columns[6]+"格式不合法（需11位手机号）")
	}
	email := values[tpl.Columns[7]]
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			errs = append(errs, tpl.Columns[7]+"格式不合法（邮箱）")
		}
	}
	status := values[tpl.Columns[8]]
	if status != "" && status != "启用" && status != "停用" {
		errs = append(errs, tpl.Columns[8]+"仅支持：启用/停用")
	}
	return errs
}

func isUserTemplate(templateKey string) bool {
	return templateKey == "service-user" || templateKey == "other-user"
}

func addAddressCascadeValidation(f *excelize.File, tpl baseDataTemplate) error {
	provCol := findColumnIndex(tpl.Columns, "省")
	cityCol := findColumnIndex(tpl.Columns, "市")
	countyCol := findColumnIndex(tpl.Columns, "县")
	if provCol == -1 || cityCol == -1 || countyCol == -1 {
		return nil
	}

	tree, err := loadAddressTree()
	if err != nil {
		return err
	}

	if _, err = f.NewSheet(addressSheetName); err != nil {
		return errors.Wrap(err, "创建地址字典页失败")
	}
	if err = fillAddressDictSheet(f, tree); err != nil {
		return err
	}
	_ = f.SetSheetVisible(addressSheetName, false, true)

	provLetter, _ := excelize.ColumnNumberToName(provCol)
	cityLetter, _ := excelize.ColumnNumberToName(cityCol)

	provRange := buildSqref(provCol, templateDataStartRow, templateDataEndRow)
	cityRange := buildSqref(cityCol, templateDataStartRow, templateDataEndRow)
	countyRange := buildSqref(countyCol, templateDataStartRow, templateDataEndRow)

	if err = addDropListValidation(f, tpl.SheetName, provRange, "=ADDR_PROVINCES", "请选择省", "省份格式：编号-名称"); err != nil {
		return err
	}
	cityFormula := fmt.Sprintf(`=INDIRECT("ADDR_P_"&LEFT($%s2,6))`, provLetter)
	if err = addDropListValidation(f, tpl.SheetName, cityRange, cityFormula, "请选择市", "先选择省，再选择市"); err != nil {
		return err
	}
	countyFormula := fmt.Sprintf(`=INDIRECT("ADDR_C_"&LEFT($%s2,6))`, cityLetter)
	if err = addDropListValidation(f, tpl.SheetName, countyRange, countyFormula, "请选择县", "先选择市，再选择县"); err != nil {
		return err
	}

	return nil
}

func buildSqref(col, startRow, endRow int) string {
	start, _ := excelize.CoordinatesToCellName(col, startRow)
	end, _ := excelize.CoordinatesToCellName(col, endRow)
	return start + ":" + end
}

func addDropListValidation(f *excelize.File, sheet, sqref, formula, title, msg string) error {
	dv := excelize.NewDataValidation(true)
	dv.Sqref = sqref
	if err := dv.SetDropList([]string{formula}); err != nil {
		return errors.Wrap(err, "设置下拉数据源失败")
	}
	dv.SetError(excelize.DataValidationErrorStyleStop, title, msg)
	dv.SetInput(title, msg)
	if err := f.AddDataValidation(sheet, dv); err != nil {
		return errors.Wrap(err, "添加下拉校验失败")
	}
	return nil
}

func findColumnIndex(columns []string, target string) int {
	for idx, col := range columns {
		if col == target {
			return idx + 1
		}
	}
	return -1
}

func fillAddressDictSheet(f *excelize.File, tree []addressNode) error {
	if err := f.SetCellValue(addressSheetName, "A1", "省"); err != nil {
		return errors.Wrap(err, "写入地址字典失败")
	}
	if err := f.SetCellValue(addressSheetName, "B1", "市"); err != nil {
		return errors.Wrap(err, "写入地址字典失败")
	}
	if err := f.SetCellValue(addressSheetName, "C1", "县"); err != nil {
		return errors.Wrap(err, "写入地址字典失败")
	}

	provinceRow := 2
	cityRow := 2
	countyRow := 2

	for _, province := range tree {
		if err := f.SetCellValue(addressSheetName, fmt.Sprintf("A%d", provinceRow), formatAreaOption(province)); err != nil {
			return errors.Wrap(err, "写入省份字典失败")
		}
		provinceRow++

		cityStart := cityRow
		for _, city := range province.Children {
			if err := f.SetCellValue(addressSheetName, fmt.Sprintf("B%d", cityRow), formatAreaOption(city)); err != nil {
				return errors.Wrap(err, "写入城市字典失败")
			}
			cityRow++

			countyStart := countyRow
			for _, county := range city.Children {
				if err := f.SetCellValue(addressSheetName, fmt.Sprintf("C%d", countyRow), formatAreaOption(county)); err != nil {
					return errors.Wrap(err, "写入区县字典失败")
				}
				countyRow++
			}
			if countyRow > countyStart {
				if err := setDefinedName(f, "ADDR_C_"+firstSix(city.Code), "C", countyStart, countyRow-1); err != nil {
					return err
				}
			}
		}
		if cityRow > cityStart {
			if err := setDefinedName(f, "ADDR_P_"+firstSix(province.Code), "B", cityStart, cityRow-1); err != nil {
				return err
			}
		}
	}

	if provinceRow <= 2 {
		return errors.New("地址字典为空，无法设置省市县下拉")
	}
	if err := setDefinedName(f, "ADDR_PROVINCES", "A", 2, provinceRow-1); err != nil {
		return err
	}
	_ = f.SetColWidth(addressSheetName, "A", "C", 28)
	return nil
}

func setDefinedName(f *excelize.File, name, col string, startRow, endRow int) error {
	if startRow > endRow {
		return nil
	}
	ref := fmt.Sprintf("'%s'!$%s$%d:$%s$%d", addressSheetName, col, startRow, col, endRow)
	if err := f.SetDefinedName(&excelize.DefinedName{Name: name, RefersTo: ref}); err != nil {
		return errors.Wrapf(err, "设置命名范围失败: %s", name)
	}
	return nil
}

func loadAddressTree() ([]addressNode, error) {
	raw, err := os.ReadFile(addressJSONPath)
	if err != nil {
		return nil, errors.Wrap(err, "读取地址字典文件失败")
	}

	raw = bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF})
	if !utf8.Valid(raw) {
		decoded, _, decodeErr := transform.Bytes(simplifiedchinese.GBK.NewDecoder(), raw)
		if decodeErr == nil && utf8.Valid(decoded) {
			raw = decoded
		}
	}

	normalized := normalizeAddressJSON(raw)
	var wrap addressWrapper
	if err = json.Unmarshal(normalized, &wrap); err != nil {
		return nil, errors.Wrap(err, "解析地址字典失败")
	}

	tree := normalizeAddressNodes(wrap.Address)
	if len(tree) == 0 {
		return nil, errors.New("地址字典没有可用省市县数据")
	}
	return tree, nil
}

func normalizeAddressJSON(raw []byte) []byte {
	s := string(raw)
	re := regexp.MustCompile(`([\{\[,]\s*)([A-Za-z_][A-Za-z0-9_]*)\s*:`)
	s = re.ReplaceAllString(s, `${1}"${2}":`)
	return []byte(s)
}

func normalizeAddressNodes(rawNodes []addressNodeRaw) []addressNode {
	nodes := make([]addressNode, 0, len(rawNodes))
	for _, item := range rawNodes {
		code := normalizeAreaCode(item.Code)
		name := strings.TrimSpace(item.Name)
		if code == "" || name == "" {
			continue
		}
		node := addressNode{
			Code: code,
			Name: name,
		}
		node.Children = normalizeAddressNodes(item.Children)
		nodes = append(nodes, node)
	}
	return nodes
}

func normalizeAreaCode(value any) string {
	switch v := value.(type) {
	case string:
		return cleanupCodeString(v)
	case float64:
		return cleanupCodeString(strconv.FormatFloat(v, 'f', 0, 64))
	case json.Number:
		return cleanupCodeString(v.String())
	default:
		return ""
	}
}

func cleanupCodeString(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	if f, err := strconv.ParseFloat(code, 64); err == nil {
		if strings.ContainsAny(strings.ToLower(code), "e.") {
			code = strconv.FormatFloat(f, 'f', 0, 64)
		}
	}
	code = regexp.MustCompile(`\D`).ReplaceAllString(code, "")
	if len(code) == 0 {
		return ""
	}
	if len(code) < 6 {
		return strings.Repeat("0", 6-len(code)) + code
	}
	return code
}

func firstSix(code string) string {
	if len(code) >= 6 {
		return code[:6]
	}
	if len(code) == 0 {
		return "000000"
	}
	return code + strings.Repeat("0", 6-len(code))
}

func formatAreaOption(node addressNode) string {
	return fmt.Sprintf("%s-%s", node.Code, node.Name)
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

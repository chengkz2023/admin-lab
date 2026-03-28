package reusable

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"sort"
	"strconv"
	"strings"

	reusableReq "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/request"
	reusableRes "github.com/flipped-aurora/gin-vue-admin/server/model/lab/reusable/response"
	"github.com/pkg/errors"
)

type TableProService struct{}

var tableProColumns = []reusableRes.TableProColumn{
	{Key: "orderNo", Label: "Order No"},
	{Key: "customerName", Label: "Customer"},
	{Key: "status", Label: "Status"},
	{Key: "priority", Label: "Priority"},
	{Key: "owner", Label: "Owner"},
	{Key: "source", Label: "Source"},
	{Key: "amount", Label: "Amount"},
	{Key: "createdAt", Label: "Created At"},
}

func (t *TableProService) Page(query reusableReq.TableProPageQuery) reusableRes.TableProPageResult {
	normalized := normalizePageQuery(query)
	rows := tableProMockRows()
	filtered := applyTableProFilter(rows, normalized.Keyword, normalized.Status, normalized.Owner, normalized.StartDate, normalized.EndDate)
	sorted := applyTableProSort(filtered, normalized.SortBy, normalized.SortOrder)
	paged := paginateTableProRows(sorted, normalized.Page, normalized.PageSize)

	return reusableRes.TableProPageResult{
		List:        paged,
		Total:       int64(len(filtered)),
		Page:        normalized.Page,
		PageSize:    normalized.PageSize,
		Columns:     tableProColumns,
		Statuses:    []string{"pending", "processing", "done"},
		Owners:      []string{"Alice", "Bob", "Carol", "David"},
		DefaultSort: "createdAt:desc",
	}
}

func (t *TableProService) Export(query reusableReq.TableProExportQuery) ([]byte, string, error) {
	rows := tableProMockRows()
	filtered := applyTableProFilter(rows, strings.TrimSpace(query.Keyword), strings.TrimSpace(query.Status), strings.TrimSpace(query.Owner), strings.TrimSpace(query.StartDate), strings.TrimSpace(query.EndDate))
	sorted := applyTableProSort(filtered, strings.TrimSpace(query.SortBy), strings.TrimSpace(query.SortOrder))

	columnMap := make(map[string]string, len(tableProColumns))
	defaultColumns := make([]string, 0, len(tableProColumns))
	for _, column := range tableProColumns {
		columnMap[column.Key] = column.Label
		defaultColumns = append(defaultColumns, column.Key)
	}

	visibleColumns := normalizeVisibleColumns(query.VisibleColumns, columnMap, defaultColumns)

	buffer := bytes.NewBuffer(nil)
	writer := csv.NewWriter(buffer)

	headers := make([]string, 0, len(visibleColumns))
	for _, key := range visibleColumns {
		headers = append(headers, columnMap[key])
	}
	if err := writer.Write(headers); err != nil {
		return nil, "", errors.Wrap(err, "write csv header failed")
	}

	for _, row := range sorted {
		record := make([]string, 0, len(visibleColumns))
		for _, key := range visibleColumns {
			record = append(record, tableProFieldValue(row, key))
		}
		if err := writer.Write(record); err != nil {
			return nil, "", errors.Wrap(err, "write csv row failed")
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", errors.Wrap(err, "flush csv failed")
	}

	return buffer.Bytes(), "table-pro-export.csv", nil
}

func normalizePageQuery(query reusableReq.TableProPageQuery) reusableReq.TableProPageQuery {
	query.Keyword = strings.TrimSpace(query.Keyword)
	query.Status = strings.TrimSpace(query.Status)
	query.Owner = strings.TrimSpace(query.Owner)
	query.StartDate = strings.TrimSpace(query.StartDate)
	query.EndDate = strings.TrimSpace(query.EndDate)
	query.SortBy = strings.TrimSpace(query.SortBy)
	query.SortOrder = strings.TrimSpace(query.SortOrder)

	if query.Page <= 0 {
		query.Page = 1
	}
	switch {
	case query.PageSize <= 0:
		query.PageSize = 10
	case query.PageSize > 100:
		query.PageSize = 100
	}
	if query.SortBy == "" {
		query.SortBy = "createdAt"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	return query
}

func tableProMockRows() []reusableRes.TableProRow {
	statuses := []string{"pending", "processing", "done"}
	priorities := []string{"P0", "P1", "P2"}
	owners := []string{"Alice", "Bob", "Carol", "David"}
	sources := []string{"ops-platform", "finance-center", "order-hub", "crm-sync"}

	rows := make([]reusableRes.TableProRow, 0, 96)
	for index := 0; index < 96; index++ {
		day := (index % 28) + 1
		amount := 1200 + float64((index%17)*230) + float64(index%5)*18.8
		rows = append(rows, reusableRes.TableProRow{
			ID:           index + 1,
			OrderNo:      fmt.Sprintf("SO-2026-%05d", index+1),
			CustomerName: fmt.Sprintf("Customer-%03d", (index%57)+1),
			Status:       statuses[index%len(statuses)],
			Priority:     priorities[index%len(priorities)],
			Owner:        owners[index%len(owners)],
			Source:       sources[index%len(sources)],
			Amount:       amount,
			CreatedAt:    fmt.Sprintf("2026-03-%02d", day),
		})
	}
	return rows
}

func applyTableProFilter(rows []reusableRes.TableProRow, keyword string, status string, owner string, startDate string, endDate string) []reusableRes.TableProRow {
	if keyword == "" && status == "" && owner == "" && startDate == "" && endDate == "" {
		return rows
	}

	filtered := make([]reusableRes.TableProRow, 0, len(rows))
	for _, row := range rows {
		if status != "" && row.Status != status {
			continue
		}
		if owner != "" && row.Owner != owner {
			continue
		}
		if startDate != "" && row.CreatedAt < startDate {
			continue
		}
		if endDate != "" && row.CreatedAt > endDate {
			continue
		}
		if keyword != "" {
			hit := strings.Contains(strings.ToLower(row.OrderNo), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(row.CustomerName), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(row.Source), strings.ToLower(keyword))
			if !hit {
				continue
			}
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func applyTableProSort(rows []reusableRes.TableProRow, sortBy string, sortOrder string) []reusableRes.TableProRow {
	sorted := append([]reusableRes.TableProRow(nil), rows...)
	desc := strings.EqualFold(sortOrder, "desc")

	sort.SliceStable(sorted, func(i, j int) bool {
		left := sorted[i]
		right := sorted[j]

		var compare int
		switch sortBy {
		case "orderNo":
			compare = strings.Compare(left.OrderNo, right.OrderNo)
		case "customerName":
			compare = strings.Compare(left.CustomerName, right.CustomerName)
		case "status":
			compare = strings.Compare(left.Status, right.Status)
		case "priority":
			compare = strings.Compare(left.Priority, right.Priority)
		case "owner":
			compare = strings.Compare(left.Owner, right.Owner)
		case "source":
			compare = strings.Compare(left.Source, right.Source)
		case "amount":
			if left.Amount < right.Amount {
				compare = -1
			} else if left.Amount > right.Amount {
				compare = 1
			}
		default:
			compare = strings.Compare(left.CreatedAt, right.CreatedAt)
		}

		if desc {
			return compare > 0
		}
		return compare < 0
	})

	return sorted
}

func paginateTableProRows(rows []reusableRes.TableProRow, page int, pageSize int) []reusableRes.TableProRow {
	start := (page - 1) * pageSize
	if start >= len(rows) {
		return []reusableRes.TableProRow{}
	}
	end := start + pageSize
	if end > len(rows) {
		end = len(rows)
	}
	return rows[start:end]
}

func normalizeVisibleColumns(columns []string, columnMap map[string]string, fallback []string) []string {
	if len(columns) == 0 {
		return fallback
	}
	result := make([]string, 0, len(columns))
	for _, key := range columns {
		if _, ok := columnMap[key]; ok {
			result = append(result, key)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}

func tableProFieldValue(row reusableRes.TableProRow, key string) string {
	switch key {
	case "orderNo":
		return row.OrderNo
	case "customerName":
		return row.CustomerName
	case "status":
		return row.Status
	case "priority":
		return row.Priority
	case "owner":
		return row.Owner
	case "source":
		return row.Source
	case "amount":
		return strconv.FormatFloat(row.Amount, 'f', 2, 64)
	case "createdAt":
		return row.CreatedAt
	default:
		return ""
	}
}

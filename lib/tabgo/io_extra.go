package tabgo

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ToJSON serializes a DataFrame to a JSON array of objects.
func ToJSON(df *DataFrame) (string, error) {
	records := ToRecords(df)
	b, err := json.Marshal(records)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadJSON parses a JSON array of objects and returns a DataFrame.
func ReadJSON(data string) (*DataFrame, error) {
	var records []map[string]any
	if err := json.Unmarshal([]byte(data), &records); err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return NewDataFrameFromRows(nil, nil), nil
	}

	// Collect all column names, sorted for determinism.
	colSet := make(map[string]bool)
	for _, rec := range records {
		for k := range rec {
			colSet[k] = true
		}
	}
	columns := make([]string, 0, len(colSet))
	for k := range colSet {
		columns = append(columns, k)
	}
	sort.Strings(columns)

	rows := make([][]any, len(records))
	for i, rec := range records {
		row := make([]any, len(columns))
		for ci, col := range columns {
			v, ok := rec[col]
			if !ok {
				row[ci] = nil
			} else {
				// JSON numbers come as float64; keep them as-is.
				row[ci] = v
			}
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(columns, rows), nil
}

// ToHTML renders a DataFrame as an HTML table string.
func ToHTML(df *DataFrame) string {
	names := df.Columns()
	nRows := df.Len()

	allVals := make([][]any, len(names))
	for i, n := range names {
		allVals[i] = df.Column(n).Values()
	}

	var sb strings.Builder
	sb.WriteString("<table>\n<thead>\n<tr>")
	for _, n := range names {
		sb.WriteString("<th>")
		sb.WriteString(htmlEscape(n))
		sb.WriteString("</th>")
	}
	sb.WriteString("</tr>\n</thead>\n<tbody>\n")

	for r := 0; r < nRows; r++ {
		sb.WriteString("<tr>")
		for ci := range names {
			sb.WriteString("<td>")
			sb.WriteString(htmlEscape(fmt.Sprintf("%v", allVals[ci][r])))
			sb.WriteString("</td>")
		}
		sb.WriteString("</tr>\n")
	}
	sb.WriteString("</tbody>\n</table>")
	return sb.String()
}

// htmlEscape escapes basic HTML special characters.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// ToDict converts a DataFrame to a map of column name to slice of values.
func ToDict(df *DataFrame) map[string][]any {
	names := df.Columns()
	result := make(map[string][]any, len(names))
	for _, n := range names {
		result[n] = df.Column(n).Values()
	}
	return result
}

// ToRecords converts a DataFrame to a slice of row maps.
func ToRecords(df *DataFrame) []map[string]any {
	names := df.Columns()
	nRows := df.Len()
	allVals := make([][]any, len(names))
	for i, n := range names {
		allVals[i] = df.Column(n).Values()
	}

	records := make([]map[string]any, nRows)
	for r := 0; r < nRows; r++ {
		rec := make(map[string]any, len(names))
		for ci, n := range names {
			rec[n] = allVals[ci][r]
		}
		records[r] = rec
	}
	return records
}

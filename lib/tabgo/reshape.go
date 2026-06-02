package tabgo

import (
	"fmt"
	"sort"
)

// Melt unpivots a DataFrame from wide to long format.
// idVars are the identifier columns to keep. valueVars are the columns to unpivot.
// The result has idVars + "variable" + "value" columns.
func Melt(df *DataFrame, idVars, valueVars []string) (*DataFrame, error) {
	// Validate columns exist.
	for _, c := range idVars {
		if _, ok := df.index[c]; !ok {
			return nil, fmt.Errorf("tabgo: column %q not found", c)
		}
	}
	for _, c := range valueVars {
		if _, ok := df.index[c]; !ok {
			return nil, fmt.Errorf("tabgo: column %q not found", c)
		}
	}

	nRows := df.Len()
	resultNames := make([]string, 0, len(idVars)+2)
	resultNames = append(resultNames, idVars...)
	resultNames = append(resultNames, "variable", "value")

	// Pre-fetch id column values.
	idVals := make([][]any, len(idVars))
	for i, c := range idVars {
		idVals[i] = df.Column(c).Values()
	}

	var rows [][]any
	for _, varName := range valueVars {
		vals := df.Column(varName).Values()
		for r := 0; r < nRows; r++ {
			row := make([]any, len(resultNames))
			for ci := range idVars {
				row[ci] = idVals[ci][r]
			}
			row[len(idVars)] = varName
			row[len(idVars)+1] = vals[r]
			rows = append(rows, row)
		}
	}
	return NewDataFrameFromRows(resultNames, rows), nil
}

// PivotTable creates a pivot table with aggregation.
// index: column for row labels, columns: column for new column headers,
// values: column to aggregate, aggFunc: one of "sum", "mean", "count", "min", "max".
func PivotTable(df *DataFrame, index, columns, values string, aggFunc string) (*DataFrame, error) {
	// Validate columns exist.
	for _, c := range []string{index, columns, values} {
		if _, ok := df.index[c]; !ok {
			return nil, fmt.Errorf("tabgo: column %q not found", c)
		}
	}

	indexVals := df.Column(index).Values()
	colVals := df.Column(columns).Values()
	valVals := df.Column(values).Values()

	// Group values by (indexVal, colVal).
	type groupKey struct {
		idx string
		col string
	}
	grouped := make(map[groupKey][]float64)
	indexOrder := make(map[string]bool)
	colOrder := make(map[string]bool)

	for r := 0; r < df.Len(); r++ {
		idxStr := fmt.Sprintf("%v", indexVals[r])
		colStr := fmt.Sprintf("%v", colVals[r])
		k := groupKey{idx: idxStr, col: colStr}
		grouped[k] = append(grouped[k], toFloat64(valVals[r]))
		indexOrder[idxStr] = true
		colOrder[colStr] = true
	}

	// Sort unique index and column values.
	uniqueIdx := make([]string, 0, len(indexOrder))
	for k := range indexOrder {
		uniqueIdx = append(uniqueIdx, k)
	}
	sort.Strings(uniqueIdx)

	uniqueCols := make([]string, 0, len(colOrder))
	for k := range colOrder {
		uniqueCols = append(uniqueCols, k)
	}
	sort.Strings(uniqueCols)

	// Build result columns: index column + one column per unique column value.
	resultNames := make([]string, 0, 1+len(uniqueCols))
	resultNames = append(resultNames, index)
	resultNames = append(resultNames, uniqueCols...)

	rows := make([][]any, len(uniqueIdx))
	for i, idx := range uniqueIdx {
		row := make([]any, len(resultNames))
		row[0] = idx
		for ci, col := range uniqueCols {
			k := groupKey{idx: idx, col: col}
			vals := grouped[k]
			row[ci+1] = aggregate(vals, aggFunc)
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows), nil
}

// aggregate applies the named aggregation function to a slice of float64 values.
func aggregate(vals []float64, aggFunc string) float64 {
	if len(vals) == 0 {
		return 0
	}
	switch aggFunc {
	case "sum":
		var s float64
		for _, v := range vals {
			s += v
		}
		return s
	case "mean":
		var s float64
		for _, v := range vals {
			s += v
		}
		return s / float64(len(vals))
	case "count":
		return float64(len(vals))
	case "min":
		m := vals[0]
		for _, v := range vals[1:] {
			if v < m {
				m = v
			}
		}
		return m
	case "max":
		m := vals[0]
		for _, v := range vals[1:] {
			if v > m {
				m = v
			}
		}
		return m
	default:
		return 0
	}
}

// Crosstab computes a cross-tabulation (frequency count) between two columns.
func Crosstab(df *DataFrame, row, col string) (*DataFrame, error) {
	if _, ok := df.index[row]; !ok {
		return nil, fmt.Errorf("tabgo: column %q not found", row)
	}
	if _, ok := df.index[col]; !ok {
		return nil, fmt.Errorf("tabgo: column %q not found", col)
	}

	rowVals := df.Column(row).Values()
	colVals := df.Column(col).Values()

	// Count frequencies.
	type pair struct{ r, c string }
	counts := make(map[pair]int)
	rowSet := make(map[string]bool)
	colSet := make(map[string]bool)

	for i := 0; i < df.Len(); i++ {
		rs := fmt.Sprintf("%v", rowVals[i])
		cs := fmt.Sprintf("%v", colVals[i])
		counts[pair{rs, cs}]++
		rowSet[rs] = true
		colSet[cs] = true
	}

	uniqueRows := make([]string, 0, len(rowSet))
	for k := range rowSet {
		uniqueRows = append(uniqueRows, k)
	}
	sort.Strings(uniqueRows)

	uniqueCols := make([]string, 0, len(colSet))
	for k := range colSet {
		uniqueCols = append(uniqueCols, k)
	}
	sort.Strings(uniqueCols)

	resultNames := make([]string, 0, 1+len(uniqueCols))
	resultNames = append(resultNames, row)
	resultNames = append(resultNames, uniqueCols...)

	rows := make([][]any, len(uniqueRows))
	for i, rv := range uniqueRows {
		r := make([]any, len(resultNames))
		r[0] = rv
		for ci, cv := range uniqueCols {
			r[ci+1] = counts[pair{rv, cv}]
		}
		rows[i] = r
	}
	return NewDataFrameFromRows(resultNames, rows), nil
}

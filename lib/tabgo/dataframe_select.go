package tabgo

import (
	"fmt"
	"math/rand"
	"sort"
)

// Iloc returns a new DataFrame by selecting rows and columns by integer position.
// If rows is nil, all rows are selected. If cols is nil, all columns are selected.
func (df *DataFrame) Iloc(rows []int, cols []int) *DataFrame {
	names := df.Columns()
	nRows := df.Len()

	// Determine selected columns.
	var selectedCols []int
	if cols == nil {
		selectedCols = make([]int, len(names))
		for i := range names {
			selectedCols[i] = i
		}
	} else {
		selectedCols = cols
	}

	// Determine selected rows.
	var selectedRows []int
	if rows == nil {
		selectedRows = make([]int, nRows)
		for i := 0; i < nRows; i++ {
			selectedRows[i] = i
		}
	} else {
		selectedRows = rows
	}

	// Pre-fetch values for selected columns.
	colNames := make([]string, len(selectedCols))
	colVals := make([][]any, len(selectedCols))
	for i, ci := range selectedCols {
		if ci < 0 || ci >= len(names) {
			panic(fmt.Sprintf("tabgo: column index %d out of range [0, %d)", ci, len(names)))
		}
		colNames[i] = names[ci]
		colVals[i] = df.columns[ci].Values()
	}

	// Build new column data.
	newColData := make([][]any, len(selectedCols))
	for i := range selectedCols {
		data := make([]any, len(selectedRows))
		for j, ri := range selectedRows {
			if ri < 0 || ri >= nRows {
				panic(fmt.Sprintf("tabgo: row index %d out of range [0, %d)", ri, nRows))
			}
			data[j] = colVals[i][ri]
		}
		newColData[i] = data
	}

	newCols := make([]*Series, len(selectedCols))
	newIdx := make(map[string]int, len(selectedCols))
	for i, name := range colNames {
		newCols[i] = NewSeries(name, newColData[i])
		newIdx[name] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Sample returns a random sample of n rows from the DataFrame.
func (df *DataFrame) Sample(n int, seed int64) *DataFrame {
	nRows := df.Len()
	if n > nRows {
		n = nRows
	}
	if n <= 0 {
		return NewDataFrameFromRows(df.Columns(), nil)
	}

	rng := rand.New(rand.NewSource(seed))
	// Fisher-Yates partial shuffle to pick n indices.
	indices := make([]int, nRows)
	for i := range indices {
		indices[i] = i
	}
	for i := 0; i < n; i++ {
		j := i + rng.Intn(nRows-i)
		indices[i], indices[j] = indices[j], indices[i]
	}
	selected := indices[:n]
	// Sort for stable output order.
	sort.Ints(selected)

	return df.Iloc(selected, nil)
}

// Nlargest returns the top n rows sorted by the given column in descending order.
// The column must contain numeric values.
func (df *DataFrame) Nlargest(n int, column string) *DataFrame {
	vals := df.Column(column).Values()
	nRows := df.Len()
	if n > nRows {
		n = nRows
	}

	type indexedVal struct {
		idx int
		val float64
	}
	iv := make([]indexedVal, nRows)
	for i, v := range vals {
		iv[i] = indexedVal{idx: i, val: toFloat64(v)}
	}
	sort.Slice(iv, func(i, j int) bool {
		return iv[i].val > iv[j].val
	})

	rows := make([]int, n)
	for i := 0; i < n; i++ {
		rows[i] = iv[i].idx
	}
	return df.Iloc(rows, nil)
}

// Nsmallest returns the bottom n rows sorted by the given column in ascending order.
// The column must contain numeric values.
func (df *DataFrame) Nsmallest(n int, column string) *DataFrame {
	vals := df.Column(column).Values()
	nRows := df.Len()
	if n > nRows {
		n = nRows
	}

	type indexedVal struct {
		idx int
		val float64
	}
	iv := make([]indexedVal, nRows)
	for i, v := range vals {
		iv[i] = indexedVal{idx: i, val: toFloat64(v)}
	}
	sort.Slice(iv, func(i, j int) bool {
		return iv[i].val < iv[j].val
	})

	rows := make([]int, n)
	for i := 0; i < n; i++ {
		rows[i] = iv[i].idx
	}
	return df.Iloc(rows, nil)
}

// Where returns a new DataFrame where values in rows satisfying the condition
// are kept, and values in rows not satisfying the condition are replaced with other.
func (df *DataFrame) Where(condition func(row map[string]any) bool, other any) *DataFrame {
	nRows := df.Len()
	names := df.Columns()
	allVals := make([][]any, len(names))
	for i, n := range names {
		allVals[i] = df.Column(n).Values()
	}

	row := make(map[string]any, len(names))
	newColData := make([][]any, len(names))
	for i := range names {
		newColData[i] = make([]any, nRows)
	}

	for r := 0; r < nRows; r++ {
		for i, n := range names {
			row[n] = allVals[i][r]
		}
		if condition(row) {
			for i := range names {
				newColData[i][r] = allVals[i][r]
			}
		} else {
			for i := range names {
				newColData[i][r] = other
			}
		}
	}

	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		newCols[i] = NewSeries(n, newColData[i])
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

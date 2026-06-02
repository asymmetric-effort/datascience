package tabgo

import (
	"fmt"
	"sort"
)

// DataFrame is an ordered collection of named Series (columns) of equal length.
type DataFrame struct {
	columns []*Series      // preserves insertion order
	index   map[string]int // column name → index into columns slice
}

// NewDataFrame creates a DataFrame from a map of column name to *Series.
// Column iteration order is sorted alphabetically for determinism.
func NewDataFrame(columns map[string]*Series) *DataFrame {
	names := make([]string, 0, len(columns))
	for n := range columns {
		names = append(names, n)
	}
	sort.Strings(names)

	df := &DataFrame{
		columns: make([]*Series, 0, len(names)),
		index:   make(map[string]int, len(names)),
	}
	for i, n := range names {
		s := columns[n]
		// ensure the series carries the map key as its name
		df.columns = append(df.columns, NewSeries(n, s.Values()))
		df.index[n] = i
	}
	return df
}

// NewDataFrameFromRows builds a DataFrame from row-oriented data.
func NewDataFrameFromRows(columnNames []string, rows [][]any) *DataFrame {
	nCols := len(columnNames)
	colData := make([][]any, nCols)
	for i := range colData {
		colData[i] = make([]any, 0, len(rows))
	}
	for _, row := range rows {
		for c := 0; c < nCols; c++ {
			if c < len(row) {
				colData[c] = append(colData[c], row[c])
			} else {
				colData[c] = append(colData[c], nil)
			}
		}
	}
	df := &DataFrame{
		columns: make([]*Series, nCols),
		index:   make(map[string]int, nCols),
	}
	for i, name := range columnNames {
		df.columns[i] = NewSeries(name, colData[i])
		df.index[name] = i
	}
	return df
}

// Columns returns the column names in order.
func (df *DataFrame) Columns() []string {
	names := make([]string, len(df.columns))
	for i, s := range df.columns {
		names[i] = s.Name()
	}
	return names
}

// Len returns the number of rows. Returns 0 if the DataFrame has no columns.
func (df *DataFrame) Len() int {
	if len(df.columns) == 0 {
		return 0
	}
	return df.columns[0].Len()
}

// Column returns the named Series or panics if not found.
func (df *DataFrame) Column(name string) *Series {
	idx, ok := df.index[name]
	if !ok {
		panic(fmt.Sprintf("tabgo: column %q not found", name))
	}
	return df.columns[idx]
}

// Select returns a new DataFrame containing only the specified columns.
func (df *DataFrame) Select(columns ...string) *DataFrame {
	cols := make([]*Series, 0, len(columns))
	idx := make(map[string]int, len(columns))
	for i, name := range columns {
		s := df.Column(name) // panics if missing
		cols = append(cols, NewSeries(name, s.Values()))
		idx[name] = i
	}
	return &DataFrame{columns: cols, index: idx}
}

// Filter returns a new DataFrame containing only rows where fn returns true.
func (df *DataFrame) Filter(fn func(row map[string]any) bool) *DataFrame {
	nRows := df.Len()
	names := df.Columns()
	// pre-fetch all column values
	allVals := make([][]any, len(names))
	for i, n := range names {
		allVals[i] = df.Column(n).Values()
	}

	// collect matching row indices
	var kept []int
	row := make(map[string]any, len(names))
	for r := 0; r < nRows; r++ {
		for i, n := range names {
			row[n] = allVals[i][r]
		}
		if fn(row) {
			kept = append(kept, r)
		}
	}

	// build new column data
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		data := make([]any, len(kept))
		for j, r := range kept {
			data[j] = allVals[i][r]
		}
		newCols[i] = NewSeries(n, data)
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Head returns the first n rows.
func (df *DataFrame) Head(n int) *DataFrame {
	total := df.Len()
	if n > total {
		n = total
	}
	return df.sliceRows(0, n)
}

// Tail returns the last n rows.
func (df *DataFrame) Tail(n int) *DataFrame {
	total := df.Len()
	if n > total {
		n = total
	}
	return df.sliceRows(total-n, total)
}

// Copy returns a deep copy of the DataFrame.
func (df *DataFrame) Copy() *DataFrame {
	return df.sliceRows(0, df.Len())
}

// sliceRows returns a new DataFrame with rows [start, end).
func (df *DataFrame) sliceRows(start, end int) *DataFrame {
	names := df.Columns()
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		vals := df.Column(n).Values()
		newCols[i] = NewSeries(n, vals[start:end])
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

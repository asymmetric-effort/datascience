package tabgo

import (
	"fmt"
	"sort"
)

// Assign returns a new DataFrame with the specified column added or replaced.
// If the column already exists, its values are replaced. Otherwise a new column is appended.
func (df *DataFrame) Assign(column string, values []any) *DataFrame {
	names := df.Columns()
	nRows := df.Len()

	// Validate length.
	if nRows > 0 && len(values) != nRows {
		panic(fmt.Sprintf("tabgo: Assign: values length %d does not match DataFrame length %d", len(values), nRows))
	}

	// Check if column already exists.
	if _, exists := df.index[column]; exists {
		newCols := make([]*Series, len(names))
		newIdx := make(map[string]int, len(names))
		for i, n := range names {
			if n == column {
				newCols[i] = NewSeries(n, values)
			} else {
				newCols[i] = NewSeries(n, df.Column(n).Values())
			}
			newIdx[n] = i
		}
		return &DataFrame{columns: newCols, index: newIdx}
	}

	// Append new column.
	newNames := make([]string, len(names)+1)
	copy(newNames, names)
	newNames[len(names)] = column

	newCols := make([]*Series, len(newNames))
	newIdx := make(map[string]int, len(newNames))
	for i, n := range names {
		newCols[i] = NewSeries(n, df.Column(n).Values())
		newIdx[n] = i
	}
	newCols[len(names)] = NewSeries(column, values)
	newIdx[column] = len(names)
	return &DataFrame{columns: newCols, index: newIdx}
}

// Insert inserts a column at the given position (0-based). Returns error if loc is out of range
// or the column name already exists.
func (df *DataFrame) Insert(loc int, column string, values []any) (*DataFrame, error) {
	names := df.Columns()
	nCols := len(names)

	if loc < 0 || loc > nCols {
		return nil, fmt.Errorf("tabgo: Insert: loc %d out of range [0, %d]", loc, nCols)
	}
	if _, exists := df.index[column]; exists {
		return nil, fmt.Errorf("tabgo: Insert: column %q already exists", column)
	}
	nRows := df.Len()
	if nRows > 0 && len(values) != nRows {
		return nil, fmt.Errorf("tabgo: Insert: values length %d does not match DataFrame length %d", len(values), nRows)
	}

	newNames := make([]string, nCols+1)
	copy(newNames[:loc], names[:loc])
	newNames[loc] = column
	copy(newNames[loc+1:], names[loc:])

	newCols := make([]*Series, nCols+1)
	newIdx := make(map[string]int, nCols+1)
	for i, n := range newNames {
		if n == column {
			newCols[i] = NewSeries(column, values)
		} else {
			newCols[i] = NewSeries(n, df.Column(n).Values())
		}
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}, nil
}

// Drop returns a new DataFrame without the specified columns.
func (df *DataFrame) Drop(columns ...string) *DataFrame {
	dropSet := make(map[string]bool, len(columns))
	for _, c := range columns {
		dropSet[c] = true
	}

	names := df.Columns()
	var keepNames []string
	for _, n := range names {
		if !dropSet[n] {
			keepNames = append(keepNames, n)
		}
	}

	if len(keepNames) == 0 {
		return NewDataFrameFromRows(nil, nil)
	}
	return df.Select(keepNames...)
}

// Rename returns a new DataFrame with columns renamed according to mapping.
// Keys are old names, values are new names.
func (df *DataFrame) Rename(mapping map[string]string) *DataFrame {
	names := df.Columns()
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		newName := n
		if mapped, ok := mapping[n]; ok {
			newName = mapped
		}
		newCols[i] = NewSeries(newName, df.Column(n).Values())
		newIdx[newName] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// SortValues returns a new DataFrame sorted by the given column.
func (df *DataFrame) SortValues(by string, ascending bool) *DataFrame {
	vals := df.Column(by).Values()
	nRows := df.Len()

	indices := make([]int, nRows)
	for i := range indices {
		indices[i] = i
	}

	sort.SliceStable(indices, func(i, j int) bool {
		a := toFloat64(vals[indices[i]])
		b := toFloat64(vals[indices[j]])
		if ascending {
			return a < b
		}
		return a > b
	})

	return df.Iloc(indices, nil)
}

// Append appends the rows of other to this DataFrame and returns the result.
// Both DataFrames must have the same columns.
func (df *DataFrame) Append(other *DataFrame) (*DataFrame, error) {
	return Concat([]*DataFrame{df, other})
}

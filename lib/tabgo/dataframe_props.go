package tabgo

import (
	"fmt"
	"strings"
)

// Shape returns [rows, cols].
func (df *DataFrame) Shape() [2]int {
	return [2]int{df.Len(), len(df.columns)}
}

// Values returns a 2D slice of all values (row-major).
func (df *DataFrame) Values() [][]any {
	nRows := df.Len()
	names := df.Columns()
	allVals := make([][]any, len(names))
	for i, n := range names {
		allVals[i] = df.Column(n).Values()
	}
	rows := make([][]any, nRows)
	for r := 0; r < nRows; r++ {
		row := make([]any, len(names))
		for c := range names {
			row[c] = allVals[c][r]
		}
		rows[r] = row
	}
	return rows
}

// Size returns rows * cols.
func (df *DataFrame) Size() int {
	return df.Len() * len(df.columns)
}

// Ndim returns the number of dimensions (always 2).
func (df *DataFrame) Ndim() int {
	return 2
}

// Empty returns true if the DataFrame has no rows or no columns.
func (df *DataFrame) Empty() bool {
	return len(df.columns) == 0 || df.Len() == 0
}

// T returns the transpose of the DataFrame.
// Original column names become the first row-index; row indices become column names.
// Columns in the transposed frame are named "0", "1", "2", etc.
func (df *DataFrame) T() *DataFrame {
	nRows := df.Len()
	names := df.Columns()
	nCols := len(names)

	// Each original column becomes a row in the transposed frame.
	// Each original row becomes a column named by its index.
	newColNames := make([]string, nRows)
	for i := 0; i < nRows; i++ {
		newColNames[i] = fmt.Sprintf("%d", i)
	}

	allVals := make([][]any, nCols)
	for i, n := range names {
		allVals[i] = df.Column(n).Values()
	}

	// Build rows for transposed frame: one row per original column.
	rows := make([][]any, nCols)
	for c := 0; c < nCols; c++ {
		row := make([]any, nRows)
		for r := 0; r < nRows; r++ {
			row[r] = allVals[c][r]
		}
		rows[c] = row
	}

	return NewDataFrameFromRows(newColNames, rows)
}

// ToString returns a formatted string representation of the DataFrame.
func (df *DataFrame) ToString() string {
	names := df.Columns()
	if len(names) == 0 {
		return "Empty DataFrame"
	}

	nRows := df.Len()
	allVals := make([][]any, len(names))
	for i, n := range names {
		allVals[i] = df.Column(n).Values()
	}

	// Compute column widths.
	widths := make([]int, len(names))
	for i, n := range names {
		widths[i] = len(n)
	}
	for r := 0; r < nRows; r++ {
		for c := range names {
			s := fmt.Sprintf("%v", allVals[c][r])
			if len(s) > widths[c] {
				widths[c] = len(s)
			}
		}
	}

	var b strings.Builder

	// Header.
	for i, n := range names {
		if i > 0 {
			b.WriteString("  ")
		}
		b.WriteString(fmt.Sprintf("%-*s", widths[i], n))
	}
	b.WriteString("\n")

	// Rows.
	for r := 0; r < nRows; r++ {
		for c := range names {
			if c > 0 {
				b.WriteString("  ")
			}
			s := fmt.Sprintf("%v", allVals[c][r])
			b.WriteString(fmt.Sprintf("%-*s", widths[c], s))
		}
		b.WriteString("\n")
	}

	return b.String()
}

package tabgo

import (
	"errors"
	"fmt"
	"strings"
)

// Merge performs a join between left and right DataFrames on the specified columns.
// Supported values for how: "inner", "left", "right", "outer".
func Merge(left, right *DataFrame, on []string, how string) (*DataFrame, error) {
	switch how {
	case "inner", "left", "right", "outer":
		// supported
	default:
		return nil, fmt.Errorf("tabgo: unsupported merge type %q (supported: inner, left, right, outer)", how)
	}
	if len(on) == 0 {
		return nil, errors.New("tabgo: merge requires at least one join column")
	}

	// Validate that join columns exist in both DataFrames.
	for _, col := range on {
		if _, ok := left.index[col]; !ok {
			return nil, fmt.Errorf("tabgo: column %q not found in left DataFrame", col)
		}
		if _, ok := right.index[col]; !ok {
			return nil, fmt.Errorf("tabgo: column %q not found in right DataFrame", col)
		}
	}

	// Determine result column names: all left columns + right columns not in 'on'.
	onSet := make(map[string]bool, len(on))
	for _, c := range on {
		onSet[c] = true
	}

	leftNames := left.Columns()
	rightNames := right.Columns()
	var rightExtra []string
	for _, n := range rightNames {
		if !onSet[n] {
			rightExtra = append(rightExtra, n)
		}
	}
	resultNames := make([]string, 0, len(leftNames)+len(rightExtra))
	resultNames = append(resultNames, leftNames...)
	resultNames = append(resultNames, rightExtra...)

	// Pre-fetch right join column values and build an index: key -> []rowIndex.
	rightKeyVals := make([][]any, len(on))
	for i, c := range on {
		rightKeyVals[i] = right.Column(c).Values()
	}
	rightIndex := make(map[string][]int)
	for r := 0; r < right.Len(); r++ {
		parts := make([]string, len(on))
		for i, vals := range rightKeyVals {
			parts[i] = fmt.Sprintf("%v", vals[r])
		}
		k := strings.Join(parts, "|")
		rightIndex[k] = append(rightIndex[k], r)
	}

	// Pre-fetch all column values.
	leftAllVals := make([][]any, len(leftNames))
	for i, n := range leftNames {
		leftAllVals[i] = left.Column(n).Values()
	}
	rightExtraVals := make([][]any, len(rightExtra))
	for i, n := range rightExtra {
		rightExtraVals[i] = right.Column(n).Values()
	}

	leftKeyVals := make([][]any, len(on))
	for i, c := range on {
		leftKeyVals[i] = left.Column(c).Values()
	}

	// Build left key index for right/outer joins.
	leftIndex := make(map[string][]int)
	for lr := 0; lr < left.Len(); lr++ {
		parts := make([]string, len(on))
		for i, vals := range leftKeyVals {
			parts[i] = fmt.Sprintf("%v", vals[lr])
		}
		k := strings.Join(parts, "|")
		leftIndex[k] = append(leftIndex[k], lr)
	}

	// Perform the join.
	var rows [][]any
	rightMatched := make(map[int]bool) // track matched right rows for right/outer

	for lr := 0; lr < left.Len(); lr++ {
		parts := make([]string, len(on))
		for i, vals := range leftKeyVals {
			parts[i] = fmt.Sprintf("%v", vals[lr])
		}
		k := strings.Join(parts, "|")
		rightRows, ok := rightIndex[k]
		if !ok || len(rightRows) == 0 {
			// No match on right side.
			if how == "left" || how == "outer" {
				row := make([]any, len(resultNames))
				for ci := range leftNames {
					row[ci] = leftAllVals[ci][lr]
				}
				for ci := range rightExtra {
					row[len(leftNames)+ci] = nil
				}
				rows = append(rows, row)
			}
			continue
		}
		for _, rr := range rightRows {
			rightMatched[rr] = true
			row := make([]any, len(resultNames))
			for ci := range leftNames {
				row[ci] = leftAllVals[ci][lr]
			}
			for ci := range rightExtra {
				row[len(leftNames)+ci] = rightExtraVals[ci][rr]
			}
			rows = append(rows, row)
		}
	}

	// For right/outer: add unmatched right rows.
	if how == "right" || how == "outer" {
		for rr := 0; rr < right.Len(); rr++ {
			if rightMatched[rr] {
				continue
			}
			row := make([]any, len(resultNames))
			// Fill left columns with nil, except for the join columns which come from right.
			for ci, n := range leftNames {
				if onSet[n] {
					row[ci] = rightKeyVals[onIndex(on, n)][rr]
				} else {
					row[ci] = nil
				}
			}
			for ci := range rightExtra {
				row[len(leftNames)+ci] = rightExtraVals[ci][rr]
			}
			rows = append(rows, row)
		}
	}

	return NewDataFrameFromRows(resultNames, rows), nil
}

// onIndex returns the index of name in the on slice.
func onIndex(on []string, name string) int {
	for i, n := range on {
		if n == name {
			return i
		}
	}
	return -1
}

// Concat vertically concatenates multiple DataFrames.
// All DataFrames must have the same columns (in any order).
func Concat(frames []*DataFrame) (*DataFrame, error) {
	if len(frames) == 0 {
		return NewDataFrameFromRows(nil, nil), nil
	}

	refNames := frames[0].Columns()
	refSet := make(map[string]bool, len(refNames))
	for _, n := range refNames {
		refSet[n] = true
	}

	// Validate all frames have the same columns.
	for i, f := range frames[1:] {
		fNames := f.Columns()
		if len(fNames) != len(refNames) {
			return nil, fmt.Errorf("tabgo: frame %d has %d columns, expected %d", i+1, len(fNames), len(refNames))
		}
		for _, n := range fNames {
			if !refSet[n] {
				return nil, fmt.Errorf("tabgo: frame %d has unexpected column %q", i+1, n)
			}
		}
	}

	// Concatenate data using the column order from the first frame.
	totalRows := 0
	for _, f := range frames {
		totalRows += f.Len()
	}

	colData := make([][]any, len(refNames))
	for i := range colData {
		colData[i] = make([]any, 0, totalRows)
	}

	for _, f := range frames {
		for ci, n := range refNames {
			colData[ci] = append(colData[ci], f.Column(n).Values()...)
		}
	}

	newCols := make([]*Series, len(refNames))
	newIdx := make(map[string]int, len(refNames))
	for i, n := range refNames {
		newCols[i] = NewSeries(n, colData[i])
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}, nil
}

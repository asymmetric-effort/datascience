package tabgo

import (
	"fmt"
	"sort"
	"strings"
)

// GroupBy holds the source DataFrame and the columns to group by.
type GroupBy struct {
	df      *DataFrame
	columns []string
}

// GroupBy returns a GroupBy object for the given column names.
func (df *DataFrame) GroupBy(columns ...string) *GroupBy {
	return &GroupBy{df: df, columns: columns}
}

// groupKey builds the "|"-joined key for a single row.
func (g *GroupBy) groupKey(row int, colVals [][]any) string {
	parts := make([]string, len(g.columns))
	for i, vals := range colVals {
		v := vals[row]
		if v == nil {
			parts[i] = "<nil>"
		} else {
			parts[i] = fmt.Sprintf("%v", v)
		}
	}
	return strings.Join(parts, "|")
}

// groupColVals returns the pre-fetched values for the group columns.
func (g *GroupBy) groupColVals() [][]any {
	colVals := make([][]any, len(g.columns))
	for i, c := range g.columns {
		colVals[i] = g.df.Column(c).Values()
	}
	return colVals
}

// buildGroups partitions row indices by group key.
func (g *GroupBy) buildGroups() ([]string, map[string][]int) {
	colVals := g.groupColVals()
	nRows := g.df.Len()
	groups := make(map[string][]int)
	var keys []string
	for r := 0; r < nRows; r++ {
		k := g.groupKey(r, colVals)
		if _, exists := groups[k]; !exists {
			keys = append(keys, k)
		}
		groups[k] = append(groups[k], r)
	}
	sort.Strings(keys)
	return keys, groups
}

// Groups returns a map from group key to sub-DataFrame.
func (g *GroupBy) Groups() map[string]*DataFrame {
	keys, groups := g.buildGroups()
	allNames := g.df.Columns()
	allVals := make([][]any, len(allNames))
	for i, n := range allNames {
		allVals[i] = g.df.Column(n).Values()
	}

	result := make(map[string]*DataFrame, len(keys))
	for _, k := range keys {
		indices := groups[k]
		rows := make([][]any, len(indices))
		for ri, idx := range indices {
			row := make([]any, len(allNames))
			for ci := range allNames {
				row[ci] = allVals[ci][idx]
			}
			rows[ri] = row
		}
		result[k] = NewDataFrameFromRows(allNames, rows)
	}
	return result
}

// Count returns a DataFrame with group columns and a "count" column.
func (g *GroupBy) Count() *DataFrame {
	keys, groups := g.buildGroups()
	colVals := g.groupColVals()

	// Build result rows: group columns + count
	resultNames := make([]string, 0, len(g.columns)+1)
	resultNames = append(resultNames, g.columns...)
	resultNames = append(resultNames, "count")

	rows := make([][]any, len(keys))
	for i, k := range keys {
		firstRow := groups[k][0]
		row := make([]any, len(resultNames))
		for ci, vals := range colVals {
			row[ci] = vals[firstRow]
		}
		row[len(g.columns)] = len(groups[k])
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows)
}

// Sum returns a DataFrame with group columns and the sum of each specified numeric column.
func (g *GroupBy) Sum(columns ...string) *DataFrame {
	keys, groups := g.buildGroups()
	colVals := g.groupColVals()

	// Pre-fetch numeric column values
	numVals := make([][]any, len(columns))
	for i, c := range columns {
		numVals[i] = g.df.Column(c).Values()
	}

	resultNames := make([]string, 0, len(g.columns)+len(columns))
	resultNames = append(resultNames, g.columns...)
	resultNames = append(resultNames, columns...)

	rows := make([][]any, len(keys))
	for i, k := range keys {
		indices := groups[k]
		firstRow := indices[0]
		row := make([]any, len(resultNames))
		for ci, vals := range colVals {
			row[ci] = vals[firstRow]
		}
		for ci, vals := range numVals {
			var sum float64
			for _, idx := range indices {
				sum += toFloat64(vals[idx])
			}
			row[len(g.columns)+ci] = sum
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows)
}

// Mean returns a DataFrame with group columns and the mean of each specified numeric column.
func (g *GroupBy) Mean(columns ...string) *DataFrame {
	keys, groups := g.buildGroups()
	colVals := g.groupColVals()

	numVals := make([][]any, len(columns))
	for i, c := range columns {
		numVals[i] = g.df.Column(c).Values()
	}

	resultNames := make([]string, 0, len(g.columns)+len(columns))
	resultNames = append(resultNames, g.columns...)
	resultNames = append(resultNames, columns...)

	rows := make([][]any, len(keys))
	for i, k := range keys {
		indices := groups[k]
		firstRow := indices[0]
		row := make([]any, len(resultNames))
		for ci, vals := range colVals {
			row[ci] = vals[firstRow]
		}
		for ci, vals := range numVals {
			var sum float64
			for _, idx := range indices {
				sum += toFloat64(vals[idx])
			}
			row[len(g.columns)+ci] = sum / float64(len(indices))
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows)
}

// Apply calls fn on each group sub-DataFrame and concatenates the results.
func (g *GroupBy) Apply(fn func(*DataFrame) *DataFrame) *DataFrame {
	keys, _ := g.buildGroups()
	groupMap := g.Groups()

	var frames []*DataFrame
	for _, k := range keys {
		sub := groupMap[k]
		result := fn(sub)
		if result != nil && result.Len() > 0 {
			frames = append(frames, result)
		}
	}
	if len(frames) == 0 {
		return NewDataFrameFromRows(nil, nil)
	}
	out, _ := Concat(frames)
	return out
}

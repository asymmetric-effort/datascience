package tabgo

import (
	"math"
	"sort"
)

// Std returns a DataFrame with group columns and the sample standard deviation
// of each specified numeric column.
func (g *GroupBy) Std(columns ...string) *DataFrame {
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
			n := len(indices)
			if n < 2 {
				row[len(g.columns)+ci] = 0.0
				continue
			}
			var sum float64
			for _, idx := range indices {
				sum += toFloat64(vals[idx])
			}
			mean := sum / float64(n)
			var sumSq float64
			for _, idx := range indices {
				d := toFloat64(vals[idx]) - mean
				sumSq += d * d
			}
			row[len(g.columns)+ci] = math.Sqrt(sumSq / float64(n-1))
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows)
}

// Var returns a DataFrame with group columns and the sample variance
// of each specified numeric column.
func (g *GroupBy) Var(columns ...string) *DataFrame {
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
			n := len(indices)
			if n < 2 {
				row[len(g.columns)+ci] = 0.0
				continue
			}
			var sum float64
			for _, idx := range indices {
				sum += toFloat64(vals[idx])
			}
			mean := sum / float64(n)
			var sumSq float64
			for _, idx := range indices {
				d := toFloat64(vals[idx]) - mean
				sumSq += d * d
			}
			row[len(g.columns)+ci] = sumSq / float64(n-1)
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows)
}

// Min returns a DataFrame with group columns and the minimum
// of each specified numeric column.
func (g *GroupBy) Min(columns ...string) *DataFrame {
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
			minVal := math.Inf(1)
			for _, idx := range indices {
				f := toFloat64(vals[idx])
				if f < minVal {
					minVal = f
				}
			}
			row[len(g.columns)+ci] = minVal
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows)
}

// Max returns a DataFrame with group columns and the maximum
// of each specified numeric column.
func (g *GroupBy) Max(columns ...string) *DataFrame {
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
			maxVal := math.Inf(-1)
			for _, idx := range indices {
				f := toFloat64(vals[idx])
				if f > maxVal {
					maxVal = f
				}
			}
			row[len(g.columns)+ci] = maxVal
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows)
}

// Median returns a DataFrame with group columns and the median
// of each specified numeric column.
func (g *GroupBy) Median(columns ...string) *DataFrame {
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
			sorted := make([]float64, len(indices))
			for j, idx := range indices {
				sorted[j] = toFloat64(vals[idx])
			}
			sort.Float64s(sorted)
			n := len(sorted)
			var med float64
			if n%2 == 0 {
				med = (sorted[n/2-1] + sorted[n/2]) / 2
			} else {
				med = sorted[n/2]
			}
			row[len(g.columns)+ci] = med
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows)
}

// First returns a DataFrame with the first row of each group.
func (g *GroupBy) First() *DataFrame {
	keys, groups := g.buildGroups()
	allNames := g.df.Columns()
	allVals := make([][]any, len(allNames))
	for i, n := range allNames {
		allVals[i] = g.df.Column(n).Values()
	}

	rows := make([][]any, len(keys))
	for i, k := range keys {
		idx := groups[k][0]
		row := make([]any, len(allNames))
		for ci := range allNames {
			row[ci] = allVals[ci][idx]
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(allNames, rows)
}

// Last returns a DataFrame with the last row of each group.
func (g *GroupBy) Last() *DataFrame {
	keys, groups := g.buildGroups()
	allNames := g.df.Columns()
	allVals := make([][]any, len(allNames))
	for i, n := range allNames {
		allVals[i] = g.df.Column(n).Values()
	}

	rows := make([][]any, len(keys))
	for i, k := range keys {
		indices := groups[k]
		idx := indices[len(indices)-1]
		row := make([]any, len(allNames))
		for ci := range allNames {
			row[ci] = allVals[ci][idx]
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(allNames, rows)
}

// Size returns a DataFrame with group columns and a "size" column.
func (g *GroupBy) Size() *DataFrame {
	keys, groups := g.buildGroups()
	colVals := g.groupColVals()

	resultNames := make([]string, 0, len(g.columns)+1)
	resultNames = append(resultNames, g.columns...)
	resultNames = append(resultNames, "size")

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

// Ngroups returns the number of groups.
func (g *GroupBy) Ngroups() int {
	keys, _ := g.buildGroups()
	return len(keys)
}

// Filter returns a DataFrame containing only groups where fn returns true.
func (g *GroupBy) Filter(fn func(*DataFrame) bool) *DataFrame {
	groupMap := g.Groups()
	keys, _ := g.buildGroups()

	var frames []*DataFrame
	for _, k := range keys {
		sub := groupMap[k]
		if fn(sub) {
			frames = append(frames, sub)
		}
	}
	if len(frames) == 0 {
		return NewDataFrameFromRows(g.df.Columns(), nil)
	}
	out, _ := Concat(frames)
	return out
}

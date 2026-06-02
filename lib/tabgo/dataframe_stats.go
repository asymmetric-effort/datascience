package tabgo

import (
	"fmt"
	"math"
	"sort"
)

// Corr returns a pairwise Pearson correlation matrix for all numeric columns.
func Corr(df *DataFrame) (*DataFrame, error) {
	cols := dfNumericColumns(df)
	if len(cols) == 0 {
		return nil, fmt.Errorf("tabgo: no numeric columns found")
	}

	nRows := df.Len()
	// Pre-fetch float values.
	data := make([][]float64, len(cols))
	for i, c := range cols {
		vals := df.Column(c).Values()
		floats := make([]float64, nRows)
		for j, v := range vals {
			if v == nil {
				floats[j] = math.NaN()
			} else {
				floats[j] = toFloat64(v)
			}
		}
		data[i] = floats
	}

	// Build result: first column is the row label, then one column per numeric column.
	resultNames := make([]string, 0, 1+len(cols))
	resultNames = append(resultNames, "")
	resultNames = append(resultNames, cols...)

	rows := make([][]any, len(cols))
	for i := range cols {
		row := make([]any, 1+len(cols))
		row[0] = cols[i]
		for j := range cols {
			row[j+1] = pearson(data[i], data[j])
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows), nil
}

// Cov returns a pairwise covariance matrix for all numeric columns.
func Cov(df *DataFrame) (*DataFrame, error) {
	cols := dfNumericColumns(df)
	if len(cols) == 0 {
		return nil, fmt.Errorf("tabgo: no numeric columns found")
	}

	nRows := df.Len()
	data := make([][]float64, len(cols))
	for i, c := range cols {
		vals := df.Column(c).Values()
		floats := make([]float64, nRows)
		for j, v := range vals {
			if v == nil {
				floats[j] = math.NaN()
			} else {
				floats[j] = toFloat64(v)
			}
		}
		data[i] = floats
	}

	resultNames := make([]string, 0, 1+len(cols))
	resultNames = append(resultNames, "")
	resultNames = append(resultNames, cols...)

	rows := make([][]any, len(cols))
	for i := range cols {
		row := make([]any, 1+len(cols))
		row[0] = cols[i]
		for j := range cols {
			row[j+1] = covariance(data[i], data[j])
		}
		rows[i] = row
	}
	return NewDataFrameFromRows(resultNames, rows), nil
}

// pearson computes the Pearson correlation between two float64 slices,
// skipping pairs where either value is NaN.
func pearson(x, y []float64) float64 {
	var xs, ys []float64
	for i := range x {
		if !math.IsNaN(x[i]) && !math.IsNaN(y[i]) {
			xs = append(xs, x[i])
			ys = append(ys, y[i])
		}
	}
	if len(xs) < 2 {
		return 0
	}
	mx, my := mean(xs), mean(ys)
	var sxy, sx2, sy2 float64
	for i := range xs {
		dx := xs[i] - mx
		dy := ys[i] - my
		sxy += dx * dy
		sx2 += dx * dx
		sy2 += dy * dy
	}
	d := math.Sqrt(sx2 * sy2)
	if d == 0 {
		return 0
	}
	return sxy / d
}

// covariance computes the sample covariance between two float64 slices.
func covariance(x, y []float64) float64 {
	var xs, ys []float64
	for i := range x {
		if !math.IsNaN(x[i]) && !math.IsNaN(y[i]) {
			xs = append(xs, x[i])
			ys = append(ys, y[i])
		}
	}
	if len(xs) < 2 {
		return 0
	}
	mx, my := mean(xs), mean(ys)
	var s float64
	for i := range xs {
		s += (xs[i] - mx) * (ys[i] - my)
	}
	return s / float64(len(xs)-1)
}

// Cumsum returns a DataFrame with cumulative sums of numeric columns.
// Non-numeric columns are dropped.
func Cumsum(df *DataFrame) *DataFrame {
	cols := dfNumericColumns(df)
	nRows := df.Len()

	resultNames := make([]string, len(cols))
	copy(resultNames, cols)

	colData := make([][]any, len(cols))
	for i, c := range cols {
		vals := df.Column(c).Values()
		data := make([]any, nRows)
		var sum float64
		for r := 0; r < nRows; r++ {
			if vals[r] != nil {
				sum += toFloat64(vals[r])
			}
			data[r] = sum
		}
		colData[i] = data
	}

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for i, n := range cols {
		newCols[i] = NewSeries(n, colData[i])
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Cumprod returns a DataFrame with cumulative products of numeric columns.
func Cumprod(df *DataFrame) *DataFrame {
	cols := dfNumericColumns(df)
	nRows := df.Len()

	colData := make([][]any, len(cols))
	for i, c := range cols {
		vals := df.Column(c).Values()
		data := make([]any, nRows)
		prod := 1.0
		for r := 0; r < nRows; r++ {
			if vals[r] != nil {
				prod *= toFloat64(vals[r])
			}
			data[r] = prod
		}
		colData[i] = data
	}

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for i, n := range cols {
		newCols[i] = NewSeries(n, colData[i])
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Diff returns a DataFrame with the discrete difference of numeric columns.
// The first `periods` rows will be nil.
func Diff(df *DataFrame, periods int) *DataFrame {
	cols := dfNumericColumns(df)
	nRows := df.Len()

	colData := make([][]any, len(cols))
	for i, c := range cols {
		vals := df.Column(c).Values()
		data := make([]any, nRows)
		for r := 0; r < nRows; r++ {
			if r < periods {
				data[r] = nil
			} else {
				curr := toFloat64(vals[r])
				prev := toFloat64(vals[r-periods])
				data[r] = curr - prev
			}
		}
		colData[i] = data
	}

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for i, n := range cols {
		newCols[i] = NewSeries(n, colData[i])
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// PctChange returns a DataFrame with percentage change of numeric columns.
// The first `periods` rows will be nil.
func PctChange(df *DataFrame, periods int) *DataFrame {
	cols := dfNumericColumns(df)
	nRows := df.Len()

	colData := make([][]any, len(cols))
	for i, c := range cols {
		vals := df.Column(c).Values()
		data := make([]any, nRows)
		for r := 0; r < nRows; r++ {
			if r < periods {
				data[r] = nil
			} else {
				prev := toFloat64(vals[r-periods])
				if prev == 0 {
					data[r] = nil
				} else {
					curr := toFloat64(vals[r])
					data[r] = (curr - prev) / prev
				}
			}
		}
		colData[i] = data
	}

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for i, n := range cols {
		newCols[i] = NewSeries(n, colData[i])
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Rank returns a DataFrame with rank values per numeric column.
// Ties receive the average rank.
func Rank(df *DataFrame) *DataFrame {
	cols := dfNumericColumns(df)
	nRows := df.Len()

	colData := make([][]any, len(cols))
	for ci, c := range cols {
		vals := df.Column(c).Values()

		type indexedVal struct {
			idx int
			val float64
		}
		items := make([]indexedVal, 0, nRows)
		for r := 0; r < nRows; r++ {
			if vals[r] != nil {
				items = append(items, indexedVal{idx: r, val: toFloat64(vals[r])})
			}
		}
		sort.Slice(items, func(a, b int) bool {
			return items[a].val < items[b].val
		})

		data := make([]any, nRows)
		for i := 0; i < len(items); {
			j := i
			for j < len(items) && items[j].val == items[i].val {
				j++
			}
			avgRank := float64(i+j+1) / 2.0
			for k := i; k < j; k++ {
				data[items[k].idx] = avgRank
			}
			i = j
		}
		for r := 0; r < nRows; r++ {
			if vals[r] == nil {
				data[r] = nil
			}
		}
		colData[ci] = data
	}

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for i, n := range cols {
		newCols[i] = NewSeries(n, colData[i])
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

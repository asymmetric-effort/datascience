package tabgo

import (
	"math"
	"sort"
)

// dfNumericColumns returns the names of columns where all non-nil values are numeric.
func dfNumericColumns(df *DataFrame) []string {
	names := df.Columns()
	var numeric []string
	for _, n := range names {
		vals := df.Column(n).Values()
		allNumeric := true
		hasValue := false
		for _, v := range vals {
			if v == nil {
				continue
			}
			hasValue = true
			if !isNumeric(v) {
				allNumeric = false
				break
			}
		}
		if allNumeric && hasValue {
			numeric = append(numeric, n)
		}
	}
	return numeric
}

// numericValues returns the non-nil float64 values for a column.
func numericValues(s *Series) []float64 {
	vals := s.Values()
	out := make([]float64, 0, len(vals))
	for _, v := range vals {
		if v != nil && isNumeric(v) {
			out = append(out, toFloat64(v))
		}
	}
	return out
}

// Sum returns the sum of each numeric column.
func (df *DataFrame) Sum() map[string]float64 {
	cols := dfNumericColumns(df)
	result := make(map[string]float64, len(cols))
	for _, n := range cols {
		vals := numericValues(df.Column(n))
		var s float64
		for _, v := range vals {
			s += v
		}
		result[n] = s
	}
	return result
}

// MeanAll returns the mean of each numeric column.
func (df *DataFrame) MeanAll() map[string]float64 {
	cols := dfNumericColumns(df)
	result := make(map[string]float64, len(cols))
	for _, n := range cols {
		vals := numericValues(df.Column(n))
		if len(vals) == 0 {
			result[n] = 0
			continue
		}
		var s float64
		for _, v := range vals {
			s += v
		}
		result[n] = s / float64(len(vals))
	}
	return result
}

// MedianAll returns the median of each numeric column.
func (df *DataFrame) MedianAll() map[string]float64 {
	cols := dfNumericColumns(df)
	result := make(map[string]float64, len(cols))
	for _, n := range cols {
		vals := numericValues(df.Column(n))
		result[n] = median(vals)
	}
	return result
}

// StdAll returns the sample standard deviation of each numeric column.
func (df *DataFrame) StdAll() map[string]float64 {
	cols := dfNumericColumns(df)
	result := make(map[string]float64, len(cols))
	for _, n := range cols {
		vals := numericValues(df.Column(n))
		result[n] = stddev(vals)
	}
	return result
}

// VarAll returns the sample variance of each numeric column.
func (df *DataFrame) VarAll() map[string]float64 {
	cols := dfNumericColumns(df)
	result := make(map[string]float64, len(cols))
	for _, n := range cols {
		vals := numericValues(df.Column(n))
		result[n] = variance(vals)
	}
	return result
}

// MinAll returns the minimum of each numeric column.
func (df *DataFrame) MinAll() map[string]float64 {
	cols := dfNumericColumns(df)
	result := make(map[string]float64, len(cols))
	for _, n := range cols {
		vals := numericValues(df.Column(n))
		if len(vals) == 0 {
			result[n] = 0
			continue
		}
		m := vals[0]
		for _, v := range vals[1:] {
			if v < m {
				m = v
			}
		}
		result[n] = m
	}
	return result
}

// MaxAll returns the maximum of each numeric column.
func (df *DataFrame) MaxAll() map[string]float64 {
	cols := dfNumericColumns(df)
	result := make(map[string]float64, len(cols))
	for _, n := range cols {
		vals := numericValues(df.Column(n))
		if len(vals) == 0 {
			result[n] = 0
			continue
		}
		m := vals[0]
		for _, v := range vals[1:] {
			if v > m {
				m = v
			}
		}
		result[n] = m
	}
	return result
}

// Count returns the non-nil count per column (all columns, not just numeric).
func (df *DataFrame) Count() map[string]int {
	names := df.Columns()
	result := make(map[string]int, len(names))
	for _, n := range names {
		vals := df.Column(n).Values()
		c := 0
		for _, v := range vals {
			if v != nil {
				c++
			}
		}
		result[n] = c
	}
	return result
}

// Describe returns summary statistics for numeric columns:
// count, mean, std, min, 25%, 50%, 75%, max.
func (df *DataFrame) Describe() *DataFrame {
	cols := dfNumericColumns(df)
	statNames := []string{"count", "mean", "std", "min", "25%", "50%", "75%", "max"}

	// Result columns: "stat" + one column per numeric column.
	resultCols := make([]string, 1+len(cols))
	resultCols[0] = "stat"
	copy(resultCols[1:], cols)

	rows := make([][]any, len(statNames))
	for si, stat := range statNames {
		row := make([]any, 1+len(cols))
		row[0] = stat
		for ci, n := range cols {
			vals := numericValues(df.Column(n))
			sorted := sortedCopy(vals)
			switch stat {
			case "count":
				row[ci+1] = float64(len(vals))
			case "mean":
				row[ci+1] = mean(vals)
			case "std":
				row[ci+1] = stddev(vals)
			case "min":
				if len(sorted) == 0 {
					row[ci+1] = 0.0
				} else {
					row[ci+1] = sorted[0]
				}
			case "25%":
				row[ci+1] = percentile(sorted, 25)
			case "50%":
				row[ci+1] = percentile(sorted, 50)
			case "75%":
				row[ci+1] = percentile(sorted, 75)
			case "max":
				if len(sorted) == 0 {
					row[ci+1] = 0.0
				} else {
					row[ci+1] = sorted[len(sorted)-1]
				}
			}
		}
		rows[si] = row
	}

	return NewDataFrameFromRows(resultCols, rows)
}

// Apply applies fn to each column Series and returns a new DataFrame.
func (df *DataFrame) Apply(fn func(*Series) *Series) *DataFrame {
	names := df.Columns()
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		s := df.Column(n)
		result := fn(NewSeries(n, s.Values()))
		newCols[i] = NewSeries(n, result.Values())
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Applymap applies fn element-wise to every value in the DataFrame.
func (df *DataFrame) Applymap(fn func(any) any) *DataFrame {
	names := df.Columns()
	nRows := df.Len()
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		vals := df.Column(n).Values()
		newVals := make([]any, nRows)
		for j, v := range vals {
			newVals[j] = fn(v)
		}
		newCols[i] = NewSeries(n, newVals)
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// --- helper math functions ---

// isNumeric returns true if the value can be converted to float64 meaningfully.
func isNumeric(v any) bool {
	switch v.(type) {
	case float64, float32, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var s float64
	for _, v := range vals {
		s += v
	}
	return s / float64(len(vals))
}

func variance(vals []float64) float64 {
	if len(vals) < 2 {
		return 0
	}
	m := mean(vals)
	var ss float64
	for _, v := range vals {
		d := v - m
		ss += d * d
	}
	return ss / float64(len(vals)-1)
}

func stddev(vals []float64) float64 {
	return math.Sqrt(variance(vals))
}

func sortedCopy(vals []float64) []float64 {
	cp := make([]float64, len(vals))
	copy(cp, vals)
	sort.Float64s(cp)
	return cp
}

func median(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sorted := sortedCopy(vals)
	n := len(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// percentile computes the p-th percentile using linear interpolation.
func percentile(vals []float64, p float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sorted := sortedCopy(vals)
	n := float64(len(sorted))
	rank := (p / 100) * (n - 1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	if lo == hi || hi >= len(sorted) {
		return sorted[lo]
	}
	frac := rank - float64(lo)
	return sorted[lo]*(1-frac) + sorted[hi]*frac
}

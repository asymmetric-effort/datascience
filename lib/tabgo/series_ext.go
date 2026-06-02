package tabgo

import (
	"math"
	"sort"
)

// Sum returns the sum of all numeric values in the series.
func (s *Series) Sum() float64 {
	var sum float64
	for _, v := range s.values {
		if v != nil {
			sum += toFloat64(v)
		}
	}
	return sum
}

// Mean returns the mean of all numeric values in the series.
func (s *Series) Mean() float64 {
	count := s.Count()
	if count == 0 {
		return 0
	}
	return s.Sum() / float64(count)
}

// Var returns the population variance of the series.
func (s *Series) Var() float64 {
	count := s.Count()
	if count < 2 {
		return 0
	}
	mean := s.Mean()
	var sumSq float64
	for _, v := range s.values {
		if v != nil {
			d := toFloat64(v) - mean
			sumSq += d * d
		}
	}
	// sample variance (ddof=1)
	return sumSq / float64(count-1)
}

// Std returns the sample standard deviation of the series.
func (s *Series) Std() float64 {
	return math.Sqrt(s.Var())
}

// Min returns the minimum numeric value in the series.
func (s *Series) Min() float64 {
	min := math.Inf(1)
	found := false
	for _, v := range s.values {
		if v != nil {
			f := toFloat64(v)
			if !found || f < min {
				min = f
				found = true
			}
		}
	}
	if !found {
		return 0
	}
	return min
}

// Max returns the maximum numeric value in the series.
func (s *Series) Max() float64 {
	max := math.Inf(-1)
	found := false
	for _, v := range s.values {
		if v != nil {
			f := toFloat64(v)
			if !found || f > max {
				max = f
				found = true
			}
		}
	}
	if !found {
		return 0
	}
	return max
}

// Median returns the median numeric value.
func (s *Series) Median() float64 {
	var vals []float64
	for _, v := range s.values {
		if v != nil {
			vals = append(vals, toFloat64(v))
		}
	}
	if len(vals) == 0 {
		return 0
	}
	sort.Float64s(vals)
	n := len(vals)
	if n%2 == 0 {
		return (vals[n/2-1] + vals[n/2]) / 2
	}
	return vals[n/2]
}

// Count returns the number of non-nil values.
func (s *Series) Count() int {
	count := 0
	for _, v := range s.values {
		if v != nil {
			count++
		}
	}
	return count
}

// Describe returns a summary statistics map with keys:
// "count", "mean", "std", "min", "25%", "50%", "75%", "max".
func (s *Series) Describe() map[string]float64 {
	var vals []float64
	for _, v := range s.values {
		if v != nil {
			vals = append(vals, toFloat64(v))
		}
	}
	count := float64(len(vals))

	result := map[string]float64{
		"count": count,
		"mean":  s.Mean(),
		"std":   s.Std(),
		"min":   s.Min(),
		"25%":   percentile(vals, 25),
		"50%":   percentile(vals, 50),
		"75%":   percentile(vals, 75),
		"max":   s.Max(),
	}
	return result
}

// Apply returns a new Series with fn applied element-wise.
func (s *Series) Apply(fn func(any) any) *Series {
	out := make([]any, len(s.values))
	for i, v := range s.values {
		out[i] = fn(v)
	}
	return &Series{name: s.name, values: out}
}

// Map is an alias for Apply.
func (s *Series) Map(fn func(any) any) *Series {
	return s.Apply(fn)
}

// Replace returns a new Series with all occurrences of old replaced with new_.
func (s *Series) Replace(old, new_ any) *Series {
	out := make([]any, len(s.values))
	for i, v := range s.values {
		if v == old {
			out[i] = new_
		} else {
			out[i] = v
		}
	}
	return &Series{name: s.name, values: out}
}

// Clip returns a new Series with numeric values clipped to [lower, upper].
func (s *Series) Clip(lower, upper float64) *Series {
	out := make([]any, len(s.values))
	for i, v := range s.values {
		if v == nil {
			out[i] = nil
			continue
		}
		f := toFloat64(v)
		if f < lower {
			f = lower
		}
		if f > upper {
			f = upper
		}
		out[i] = f
	}
	return &Series{name: s.name, values: out}
}

// Abs returns a new Series with absolute values of numeric elements.
func (s *Series) Abs() *Series {
	out := make([]any, len(s.values))
	for i, v := range s.values {
		if v == nil {
			out[i] = nil
			continue
		}
		out[i] = math.Abs(toFloat64(v))
	}
	return &Series{name: s.name, values: out}
}

// Round returns a new Series with numeric values rounded to the given decimals.
func (s *Series) Round(decimals int) *Series {
	pow := math.Pow(10, float64(decimals))
	out := make([]any, len(s.values))
	for i, v := range s.values {
		if v == nil {
			out[i] = nil
			continue
		}
		out[i] = math.Round(toFloat64(v)*pow) / pow
	}
	return &Series{name: s.name, values: out}
}

// Sort returns a new Series with values sorted.
func (s *Series) Sort(ascending bool) *Series {
	var vals []float64
	for _, v := range s.values {
		if v == nil {
			vals = append(vals, math.NaN())
		} else {
			vals = append(vals, toFloat64(v))
		}
	}
	if ascending {
		sort.Float64s(vals)
	} else {
		sort.Sort(sort.Reverse(sort.Float64Slice(vals)))
	}
	out := make([]any, len(vals))
	for i, f := range vals {
		if math.IsNaN(f) {
			out[i] = nil
		} else {
			out[i] = f
		}
	}
	return &Series{name: s.name, values: out}
}

// Rank returns a new Series where each value is replaced by its rank (1-based).
// Ties receive the average rank.
func (s *Series) Rank() *Series {
	type indexedVal struct {
		idx int
		val float64
	}
	n := len(s.values)
	items := make([]indexedVal, 0, n)
	for i, v := range s.values {
		if v != nil {
			items = append(items, indexedVal{idx: i, val: toFloat64(v)})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].val < items[j].val
	})

	ranks := make([]any, n)
	// assign average rank for ties
	for i := 0; i < len(items); {
		j := i
		for j < len(items) && items[j].val == items[i].val {
			j++
		}
		avgRank := float64(i+j+1) / 2.0 // average of 1-based positions i+1..j
		for k := i; k < j; k++ {
			ranks[items[k].idx] = avgRank
		}
		i = j
	}
	// nil values get nil rank
	for i, v := range s.values {
		if v == nil {
			ranks[i] = nil
		}
	}
	return &Series{name: s.name, values: ranks}
}

// Isin returns a boolean slice indicating whether each value is in the given set.
func (s *Series) Isin(values []any) []bool {
	set := make(map[any]bool, len(values))
	for _, v := range values {
		set[v] = true
	}
	out := make([]bool, len(s.values))
	for i, v := range s.values {
		out[i] = set[v]
	}
	return out
}

// Between returns a boolean slice where true indicates the value is in [lower, upper].
func (s *Series) Between(lower, upper float64) []bool {
	out := make([]bool, len(s.values))
	for i, v := range s.values {
		if v == nil {
			out[i] = false
			continue
		}
		f := toFloat64(v)
		out[i] = f >= lower && f <= upper
	}
	return out
}

// Corr returns the Pearson correlation coefficient between two series.
func (s *Series) Corr(other *Series) float64 {
	n := s.Len()
	if n != other.Len() || n == 0 {
		return 0
	}
	// Use only positions where both are non-nil.
	var xs, ys []float64
	for i := 0; i < n; i++ {
		if s.values[i] != nil && other.values[i] != nil {
			xs = append(xs, toFloat64(s.values[i]))
			ys = append(ys, toFloat64(other.values[i]))
		}
	}
	if len(xs) < 2 {
		return 0
	}
	meanX, meanY := 0.0, 0.0
	for i := range xs {
		meanX += xs[i]
		meanY += ys[i]
	}
	meanX /= float64(len(xs))
	meanY /= float64(len(ys))

	var sumXY, sumX2, sumY2 float64
	for i := range xs {
		dx := xs[i] - meanX
		dy := ys[i] - meanY
		sumXY += dx * dy
		sumX2 += dx * dx
		sumY2 += dy * dy
	}
	denom := math.Sqrt(sumX2 * sumY2)
	if denom == 0 {
		return 0
	}
	return sumXY / denom
}

// DropNA returns a new Series with nil values removed.
func (s *Series) DropNA() *Series {
	var out []any
	for _, v := range s.values {
		if v != nil {
			out = append(out, v)
		}
	}
	return &Series{name: s.name, values: out}
}

// FillNA returns a new Series with nil values replaced by value.
func (s *Series) FillNA(value any) *Series {
	out := make([]any, len(s.values))
	for i, v := range s.values {
		if v == nil {
			out[i] = value
		} else {
			out[i] = v
		}
	}
	return &Series{name: s.name, values: out}
}

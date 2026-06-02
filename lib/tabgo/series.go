package tabgo

import "fmt"

// Series represents a named column of data.
type Series struct {
	name   string
	values []any
}

// NewSeries creates a new Series with the given name and data.
func NewSeries(name string, data []any) *Series {
	cp := make([]any, len(data))
	copy(cp, data)
	return &Series{name: name, values: cp}
}

// Name returns the series name.
func (s *Series) Name() string { return s.name }

// Len returns the number of elements.
func (s *Series) Len() int { return len(s.values) }

// Values returns a copy of the underlying data.
func (s *Series) Values() []any {
	cp := make([]any, len(s.values))
	copy(cp, s.values)
	return cp
}

// Float64 converts each element to float64.
// Supported source types: float64, float32, int, int8–int64, uint–uint64, string (via fmt).
func (s *Series) Float64() []float64 {
	out := make([]float64, len(s.values))
	for i, v := range s.values {
		out[i] = toFloat64(v)
	}
	return out
}

// Int converts each element to int.
func (s *Series) Int() []int {
	out := make([]int, len(s.values))
	for i, v := range s.values {
		out[i] = toInt(v)
	}
	return out
}

// ValueCounts returns a map from value to occurrence count.
func (s *Series) ValueCounts() map[any]int {
	m := make(map[any]int, len(s.values))
	for _, v := range s.values {
		m[v]++
	}
	return m
}

// Unique returns the distinct values in order of first appearance.
func (s *Series) Unique() []any {
	seen := make(map[any]bool, len(s.values))
	var out []any
	for _, v := range s.values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

// NUnique returns the number of distinct values.
func (s *Series) NUnique() int {
	return len(s.Unique())
}

// toFloat64 converts a single value to float64.
func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case uint:
		return float64(n)
	case uint8:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case string:
		var f float64
		_, _ = fmt.Sscan(n, &f)
		return f
	default:
		return 0
	}
}

// toInt converts a single value to int.
func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int8:
		return int(n)
	case int16:
		return int(n)
	case int32:
		return int(n)
	case int64:
		return int(n)
	case uint:
		return int(n)
	case uint8:
		return int(n)
	case uint16:
		return int(n)
	case uint32:
		return int(n)
	case uint64:
		return int(n)
	case float64:
		return int(n)
	case float32:
		return int(n)
	case string:
		var i int
		_, _ = fmt.Sscan(n, &i)
		return i
	default:
		return 0
	}
}

package tabgo

import "fmt"

// Dtype returns a string describing the predominant type of the Series values.
// Returns "float64", "int", "string", "bool", "mixed", or "empty".
func (s *Series) Dtype() string {
	if len(s.values) == 0 {
		return "empty"
	}
	counts := make(map[string]int)
	for _, v := range s.values {
		if v == nil {
			continue
		}
		switch v.(type) {
		case float64, float32:
			counts["float64"]++
		case int, int8, int16, int32, int64:
			counts["int"]++
		case uint, uint8, uint16, uint32, uint64:
			counts["int"]++
		case string:
			counts["string"]++
		case bool:
			counts["bool"]++
		default:
			counts["object"]++
		}
	}
	if len(counts) == 0 {
		return "empty"
	}
	if len(counts) == 1 {
		for k := range counts {
			return k
		}
	}
	// If mix of int and float, report float64
	if len(counts) == 2 {
		_, hasInt := counts["int"]
		_, hasFloat := counts["float64"]
		if hasInt && hasFloat {
			return "float64"
		}
	}
	return "mixed"
}

// Shape returns a single-element array with the length of the Series.
func (s *Series) Shape() [1]int {
	return [1]int{len(s.values)}
}

// Empty returns true if the Series has no elements.
func (s *Series) Empty() bool {
	return len(s.values) == 0
}

// Index returns a slice of integer indices [0, 1, 2, ...] for the Series.
func (s *Series) Index() []int {
	idx := make([]int, len(s.values))
	for i := range idx {
		idx[i] = i
	}
	return idx
}

// Head returns a new Series with the first n elements.
func (s *Series) Head(n int) *Series {
	if n > len(s.values) {
		n = len(s.values)
	}
	if n < 0 {
		n = 0
	}
	return &Series{name: s.name, values: copySlice(s.values[:n])}
}

// Tail returns a new Series with the last n elements.
func (s *Series) Tail(n int) *Series {
	total := len(s.values)
	if n > total {
		n = total
	}
	if n < 0 {
		n = 0
	}
	return &Series{name: s.name, values: copySlice(s.values[total-n:])}
}

// Loc returns a new Series by selecting elements at the given integer positions.
func (s *Series) Loc(indices []int) *Series {
	out := make([]any, len(indices))
	for i, idx := range indices {
		if idx < 0 || idx >= len(s.values) {
			panic(fmt.Sprintf("tabgo: Series.Loc: index %d out of range [0, %d)", idx, len(s.values)))
		}
		out[i] = s.values[idx]
	}
	return &Series{name: s.name, values: out}
}

// Iloc returns a new Series by selecting elements at the given integer positions.
// This is an alias for Loc on Series since Series uses positional indexing.
func (s *Series) Iloc(indices []int) *Series {
	return s.Loc(indices)
}

// Where returns a new Series where elements at positions where cond is true
// are kept, and elements where cond is false are replaced with other.
// cond must have the same length as the Series.
func (s *Series) Where(cond []bool, other any) *Series {
	if len(cond) != len(s.values) {
		panic(fmt.Sprintf("tabgo: Series.Where: cond length %d does not match Series length %d", len(cond), len(s.values)))
	}
	out := make([]any, len(s.values))
	for i, v := range s.values {
		if cond[i] {
			out[i] = v
		} else {
			out[i] = other
		}
	}
	return &Series{name: s.name, values: out}
}

// Mask returns a new Series where elements at positions where cond is true
// are replaced with other, and elements where cond is false are kept.
// This is the inverse of Where.
func (s *Series) Mask(cond []bool, other any) *Series {
	if len(cond) != len(s.values) {
		panic(fmt.Sprintf("tabgo: Series.Mask: cond length %d does not match Series length %d", len(cond), len(s.values)))
	}
	out := make([]any, len(s.values))
	for i, v := range s.values {
		if cond[i] {
			out[i] = other
		} else {
			out[i] = v
		}
	}
	return &Series{name: s.name, values: out}
}

// Astype converts all elements in the Series to the specified type.
// Supported dtypes: "float64", "int", "string", "bool".
func (s *Series) Astype(dtype string) *Series {
	out := make([]any, len(s.values))
	for i, v := range s.values {
		if v == nil {
			out[i] = nil
			continue
		}
		switch dtype {
		case "float64":
			out[i] = toFloat64(v)
		case "int":
			out[i] = toInt(v)
		case "string":
			out[i] = fmt.Sprintf("%v", v)
		case "bool":
			out[i] = toBool(v)
		default:
			out[i] = v
		}
	}
	return &Series{name: s.name, values: out}
}

// Cumsum returns a new Series with cumulative sums.
func (s *Series) Cumsum() *Series {
	out := make([]any, len(s.values))
	var sum float64
	for i, v := range s.values {
		if v != nil {
			sum += toFloat64(v)
		}
		out[i] = sum
	}
	return &Series{name: s.name, values: out}
}

// Cumprod returns a new Series with cumulative products.
func (s *Series) Cumprod() *Series {
	out := make([]any, len(s.values))
	prod := 1.0
	for i, v := range s.values {
		if v != nil {
			prod *= toFloat64(v)
		}
		out[i] = prod
	}
	return &Series{name: s.name, values: out}
}

// Cummax returns a new Series with cumulative maximums.
func (s *Series) Cummax() *Series {
	out := make([]any, len(s.values))
	first := true
	var max float64
	for i, v := range s.values {
		if v != nil {
			f := toFloat64(v)
			if first || f > max {
				max = f
				first = false
			}
		}
		if first {
			out[i] = nil
		} else {
			out[i] = max
		}
	}
	return &Series{name: s.name, values: out}
}

// Cummin returns a new Series with cumulative minimums.
func (s *Series) Cummin() *Series {
	out := make([]any, len(s.values))
	first := true
	var min float64
	for i, v := range s.values {
		if v != nil {
			f := toFloat64(v)
			if first || f < min {
				min = f
				first = false
			}
		}
		if first {
			out[i] = nil
		} else {
			out[i] = min
		}
	}
	return &Series{name: s.name, values: out}
}

// Diff returns a new Series with the discrete difference.
// The first `periods` elements will be nil.
func (s *Series) Diff(periods int) *Series {
	out := make([]any, len(s.values))
	for i := range s.values {
		if i < periods {
			out[i] = nil
		} else {
			curr := toFloat64(s.values[i])
			prev := toFloat64(s.values[i-periods])
			out[i] = curr - prev
		}
	}
	return &Series{name: s.name, values: out}
}

// PctChange returns a new Series with the percentage change.
// The first `periods` elements will be nil.
func (s *Series) PctChange(periods int) *Series {
	out := make([]any, len(s.values))
	for i := range s.values {
		if i < periods {
			out[i] = nil
		} else {
			prev := toFloat64(s.values[i-periods])
			if prev == 0 {
				out[i] = nil
			} else {
				curr := toFloat64(s.values[i])
				out[i] = (curr - prev) / prev
			}
		}
	}
	return &Series{name: s.name, values: out}
}

// Nlargest returns a new Series with the n largest values, sorted descending.
func (s *Series) Nlargest(n int) *Series {
	sorted := s.Sort(false)
	if n > len(sorted.values) {
		n = len(sorted.values)
	}
	return &Series{name: s.name, values: copySlice(sorted.values[:n])}
}

// Nsmallest returns a new Series with the n smallest values, sorted ascending.
func (s *Series) Nsmallest(n int) *Series {
	sorted := s.Sort(true)
	if n > len(sorted.values) {
		n = len(sorted.values)
	}
	return &Series{name: s.name, values: copySlice(sorted.values[:n])}
}

// copySlice returns a copy of a []any slice.
func copySlice(src []any) []any {
	cp := make([]any, len(src))
	copy(cp, src)
	return cp
}

// toBool converts a value to bool.
func toBool(v any) bool {
	switch b := v.(type) {
	case bool:
		return b
	case int:
		return b != 0
	case float64:
		return b != 0
	case string:
		return b != "" && b != "0" && b != "false" && b != "False"
	default:
		return v != nil
	}
}

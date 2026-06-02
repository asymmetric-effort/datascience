package tabgo

import (
	"fmt"
	"strings"
)

// StrAccessor provides vectorized string operations on a Series.
type StrAccessor struct {
	series *Series
}

// Str returns a StrAccessor for vectorized string operations.
// Non-string values are converted to their string representation.
func (s *Series) Str() *StrAccessor {
	return &StrAccessor{series: s}
}

// toString converts a Series value to string.
func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// Lower returns a new Series with all string values lowercased.
func (sa *StrAccessor) Lower() *Series {
	out := make([]any, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		out[i] = strings.ToLower(toString(v))
	}
	return &Series{name: sa.series.name, values: out}
}

// Upper returns a new Series with all string values uppercased.
func (sa *StrAccessor) Upper() *Series {
	out := make([]any, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		out[i] = strings.ToUpper(toString(v))
	}
	return &Series{name: sa.series.name, values: out}
}

// Contains returns a boolean slice indicating whether each element contains substr.
func (sa *StrAccessor) Contains(substr string) []bool {
	out := make([]bool, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = false
			continue
		}
		out[i] = strings.Contains(toString(v), substr)
	}
	return out
}

// Replace returns a new Series with all occurrences of old replaced by new_ in each element.
func (sa *StrAccessor) Replace(old, new_ string) *Series {
	out := make([]any, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		out[i] = strings.ReplaceAll(toString(v), old, new_)
	}
	return &Series{name: sa.series.name, values: out}
}

// Split returns a new Series where each element is a []string from splitting by sep.
func (sa *StrAccessor) Split(sep string) *Series {
	out := make([]any, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		out[i] = strings.Split(toString(v), sep)
	}
	return &Series{name: sa.series.name, values: out}
}

// Strip returns a new Series with leading and trailing whitespace removed.
func (sa *StrAccessor) Strip() *Series {
	out := make([]any, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		out[i] = strings.TrimSpace(toString(v))
	}
	return &Series{name: sa.series.name, values: out}
}

// Len returns a new Series with the length of each string element.
func (sa *StrAccessor) Len() *Series {
	out := make([]any, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		out[i] = len(toString(v))
	}
	return &Series{name: sa.series.name, values: out}
}

// StartsWith returns a boolean slice indicating whether each element starts with prefix.
func (sa *StrAccessor) StartsWith(prefix string) []bool {
	out := make([]bool, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = false
			continue
		}
		out[i] = strings.HasPrefix(toString(v), prefix)
	}
	return out
}

// EndsWith returns a boolean slice indicating whether each element ends with suffix.
func (sa *StrAccessor) EndsWith(suffix string) []bool {
	out := make([]bool, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = false
			continue
		}
		out[i] = strings.HasSuffix(toString(v), suffix)
	}
	return out
}

// Slice returns a new Series with each string element sliced from start to end.
// Negative indices are not supported. If end exceeds the string length, it is clamped.
func (sa *StrAccessor) Slice(start, end int) *Series {
	out := make([]any, len(sa.series.values))
	for i, v := range sa.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		s := toString(v)
		if start > len(s) {
			start = len(s)
		}
		e := end
		if e > len(s) {
			e = len(s)
		}
		if start > e {
			out[i] = ""
		} else {
			out[i] = s[start:e]
		}
	}
	return &Series{name: sa.series.name, values: out}
}

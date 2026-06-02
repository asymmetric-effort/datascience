package tabgo

import (
	"fmt"
	"sort"
)

// CatAccessor provides categorical operations on a Series.
// This is an initial implementation that treats the Series values as categorical data.
type CatAccessor struct {
	series *Series
}

// Cat returns a CatAccessor for categorical operations on the Series.
func (s *Series) Cat() *CatAccessor {
	return &CatAccessor{series: s}
}

// Categories returns the unique category values in sorted order.
func (ca *CatAccessor) Categories() []any {
	seen := make(map[string]any)
	var keys []string
	for _, v := range ca.series.values {
		if v == nil {
			continue
		}
		k := fmt.Sprintf("%v", v)
		if _, exists := seen[k]; !exists {
			seen[k] = v
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	out := make([]any, len(keys))
	for i, k := range keys {
		out[i] = seen[k]
	}
	return out
}

// Codes returns a new Series where each value is replaced by its integer category code.
// Categories are assigned codes in sorted order of their string representation.
func (ca *CatAccessor) Codes() *Series {
	cats := ca.Categories()
	catMap := make(map[string]int, len(cats))
	for i, c := range cats {
		catMap[fmt.Sprintf("%v", c)] = i
	}
	out := make([]any, len(ca.series.values))
	for i, v := range ca.series.values {
		if v == nil {
			out[i] = -1
		} else {
			out[i] = catMap[fmt.Sprintf("%v", v)]
		}
	}
	return &Series{name: ca.series.name, values: out}
}

// RenameCategories returns a new Series with category values renamed according to the mapping.
// Keys in the mapping are the string representations of old values.
func (ca *CatAccessor) RenameCategories(mapping map[string]any) *Series {
	out := make([]any, len(ca.series.values))
	for i, v := range ca.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		k := fmt.Sprintf("%v", v)
		if newVal, ok := mapping[k]; ok {
			out[i] = newVal
		} else {
			out[i] = v
		}
	}
	return &Series{name: ca.series.name, values: out}
}

// AddCategories is a no-op in this implementation since categories are derived from data.
// It is provided for API compatibility. Returns the accessor unchanged.
func (ca *CatAccessor) AddCategories(_ []any) *CatAccessor {
	return ca
}

// RemoveCategories returns a new Series with values matching the given categories set to nil.
func (ca *CatAccessor) RemoveCategories(removals []any) *Series {
	removeSet := make(map[string]bool, len(removals))
	for _, r := range removals {
		removeSet[fmt.Sprintf("%v", r)] = true
	}
	out := make([]any, len(ca.series.values))
	for i, v := range ca.series.values {
		if v == nil {
			out[i] = nil
			continue
		}
		if removeSet[fmt.Sprintf("%v", v)] {
			out[i] = nil
		} else {
			out[i] = v
		}
	}
	return &Series{name: ca.series.name, values: out}
}

// Ordered returns false. This initial implementation does not support ordered categories.
func (ca *CatAccessor) Ordered() bool {
	return false
}

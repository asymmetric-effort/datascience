package tabgo

import (
	"fmt"
	"sort"
	"strings"
)

// MultiIndex represents a hierarchical index with multiple levels.
type MultiIndex struct {
	levels [][]any  // unique values at each level
	codes  [][]int  // codes indexing into levels for each element
	names  []string // name for each level
}

// NewMultiIndex creates a MultiIndex from explicit levels, codes, and names.
// levels[i] contains the unique values for level i.
// codes[i] contains indices into levels[i] for each row.
// names contains the name for each level.
// All codes slices must have the same length. Each code value must be a valid
// index into the corresponding levels slice.
func NewMultiIndex(levels [][]any, codes [][]int, names []string) (*MultiIndex, error) {
	if len(levels) == 0 {
		return nil, fmt.Errorf("tabgo: MultiIndex requires at least one level")
	}
	if len(levels) != len(codes) {
		return nil, fmt.Errorf("tabgo: MultiIndex levels count (%d) != codes count (%d)", len(levels), len(codes))
	}
	if len(names) != len(levels) {
		return nil, fmt.Errorf("tabgo: MultiIndex names count (%d) != levels count (%d)", len(names), len(levels))
	}

	// Validate all code slices have equal length.
	n := len(codes[0])
	for i, c := range codes {
		if len(c) != n {
			return nil, fmt.Errorf("tabgo: MultiIndex codes[%d] has length %d, expected %d", i, len(c), n)
		}
	}

	// Validate code values are within range.
	for i := range codes {
		for j, c := range codes[i] {
			if c < 0 || c >= len(levels[i]) {
				return nil, fmt.Errorf("tabgo: MultiIndex codes[%d][%d]=%d out of range [0,%d)", i, j, c, len(levels[i]))
			}
		}
	}

	// Deep copy inputs.
	cpLevels := make([][]any, len(levels))
	for i, lv := range levels {
		cpLevels[i] = make([]any, len(lv))
		copy(cpLevels[i], lv)
	}
	cpCodes := make([][]int, len(codes))
	for i, cd := range codes {
		cpCodes[i] = make([]int, len(cd))
		copy(cpCodes[i], cd)
	}
	cpNames := make([]string, len(names))
	copy(cpNames, names)

	return &MultiIndex{
		levels: cpLevels,
		codes:  cpCodes,
		names:  cpNames,
	}, nil
}

// NewMultiIndexFromArrays creates a MultiIndex from arrays of values for each level.
// Each element in arrays is a []any representing the values at one level.
// All arrays must have the same length. Unique values and codes are computed automatically.
func NewMultiIndexFromArrays(arrays [][]any, names []string) (*MultiIndex, error) {
	if len(arrays) == 0 {
		return nil, fmt.Errorf("tabgo: MultiIndex requires at least one level")
	}
	n := len(arrays[0])
	for i, arr := range arrays {
		if len(arr) != n {
			return nil, fmt.Errorf("tabgo: MultiIndex level %d has length %d, expected %d", i, len(arr), n)
		}
	}

	nLevels := len(arrays)
	levels := make([][]any, nLevels)
	codes := make([][]int, nLevels)

	for lvl := 0; lvl < nLevels; lvl++ {
		uniqueMap := make(map[string]int)
		var uniqueVals []any
		levelCodes := make([]int, n)

		for i, v := range arrays[lvl] {
			k := fmt.Sprintf("%v", v)
			if idx, exists := uniqueMap[k]; exists {
				levelCodes[i] = idx
			} else {
				idx := len(uniqueVals)
				uniqueMap[k] = idx
				uniqueVals = append(uniqueVals, v)
				levelCodes[i] = idx
			}
		}
		levels[lvl] = uniqueVals
		codes[lvl] = levelCodes
	}

	if names == nil {
		names = make([]string, nLevels)
		for i := range names {
			names[i] = fmt.Sprintf("level_%d", i)
		}
	}

	return &MultiIndex{
		levels: levels,
		codes:  codes,
		names:  names,
	}, nil
}

// Nlevels returns the number of levels in the MultiIndex.
func (mi *MultiIndex) Nlevels() int {
	return len(mi.levels)
}

// Len returns the number of elements (rows) in the MultiIndex.
func (mi *MultiIndex) Len() int {
	if len(mi.codes) == 0 || len(mi.codes[0]) == 0 {
		return 0
	}
	return len(mi.codes[0])
}

// Names returns the names of each level.
func (mi *MultiIndex) Names() []string {
	cp := make([]string, len(mi.names))
	copy(cp, mi.names)
	return cp
}

// Levels returns all levels (unique values at each level).
func (mi *MultiIndex) Levels() [][]any {
	cp := make([][]any, len(mi.levels))
	for i, lv := range mi.levels {
		cp[i] = make([]any, len(lv))
		copy(cp[i], lv)
	}
	return cp
}

// GetLevel returns the unique values for the level with the given name.
func (mi *MultiIndex) GetLevel(name string) ([]any, error) {
	for i, n := range mi.names {
		if n == name {
			cp := make([]any, len(mi.levels[i]))
			copy(cp, mi.levels[i])
			return cp, nil
		}
	}
	return nil, fmt.Errorf("tabgo: MultiIndex.GetLevel: level %q not found", name)
}

// GetLevelByIndex returns the unique values at the specified level index.
func (mi *MultiIndex) GetLevelByIndex(level int) []any {
	if level < 0 || level >= len(mi.levels) {
		return nil
	}
	cp := make([]any, len(mi.levels[level]))
	copy(cp, mi.levels[level])
	return cp
}

// Get returns all index values for a given row.
func (mi *MultiIndex) Get(row int) []any {
	if row < 0 || row >= mi.Len() {
		return nil
	}
	vals := make([]any, len(mi.levels))
	for lvl := range mi.levels {
		vals[lvl] = mi.levels[lvl][mi.codes[lvl][row]]
	}
	return vals
}

// GetValues returns the values at position i as a tuple (alias for Get).
func (mi *MultiIndex) GetValues(i int) []any {
	return mi.Get(i)
}

// Set updates the index values for a given row.
// The values slice must have one entry per level. Each value must already
// exist in the corresponding level's unique values.
func (mi *MultiIndex) Set(row int, values []any) error {
	if row < 0 || row >= mi.Len() {
		return fmt.Errorf("tabgo: MultiIndex.Set: row %d out of range [0,%d)", row, mi.Len())
	}
	if len(values) != len(mi.levels) {
		return fmt.Errorf("tabgo: MultiIndex.Set: expected %d values, got %d", len(mi.levels), len(values))
	}
	for lvl, v := range values {
		k := fmt.Sprintf("%v", v)
		found := false
		for idx, lv := range mi.levels[lvl] {
			if fmt.Sprintf("%v", lv) == k {
				mi.codes[lvl][row] = idx
				found = true
				break
			}
		}
		if !found {
			// Add the new value to the level.
			newIdx := len(mi.levels[lvl])
			mi.levels[lvl] = append(mi.levels[lvl], v)
			mi.codes[lvl][row] = newIdx
		}
	}
	return nil
}

// Equals returns true if two MultiIndex values are structurally identical.
func (mi *MultiIndex) Equals(other *MultiIndex) bool {
	if other == nil {
		return false
	}
	if mi.Len() != other.Len() || len(mi.levels) != len(other.levels) {
		return false
	}
	// Compare resolved values row by row.
	for row := 0; row < mi.Len(); row++ {
		a := mi.Get(row)
		b := other.Get(row)
		for i := range a {
			if fmt.Sprintf("%v", a[i]) != fmt.Sprintf("%v", b[i]) {
				return false
			}
		}
	}
	// Compare names.
	for i := range mi.names {
		if mi.names[i] != other.names[i] {
			return false
		}
	}
	return true
}

// Copy returns a deep copy of the MultiIndex.
func (mi *MultiIndex) Copy() *MultiIndex {
	cpLevels := make([][]any, len(mi.levels))
	for i, lv := range mi.levels {
		cpLevels[i] = make([]any, len(lv))
		copy(cpLevels[i], lv)
	}
	cpCodes := make([][]int, len(mi.codes))
	for i, cd := range mi.codes {
		cpCodes[i] = make([]int, len(cd))
		copy(cpCodes[i], cd)
	}
	cpNames := make([]string, len(mi.names))
	copy(cpNames, mi.names)
	return &MultiIndex{
		levels: cpLevels,
		codes:  cpCodes,
		names:  cpNames,
	}
}

// GetLevelSeries returns a Series with values from the specified level index.
func (mi *MultiIndex) GetLevelSeries(level int) *Series {
	if level < 0 || level >= len(mi.levels) {
		return NewSeries("", nil)
	}
	n := mi.Len()
	vals := make([]any, n)
	for i := 0; i < n; i++ {
		vals[i] = mi.levels[level][mi.codes[level][i]]
	}
	name := mi.names[level]
	return NewSeries(name, vals)
}

// String returns a string representation of the MultiIndex.
func (mi *MultiIndex) String() string {
	n := mi.Len()
	if n == 0 {
		return "MultiIndex([])"
	}
	var sb strings.Builder
	sb.WriteString("MultiIndex([\n")
	limit := n
	if limit > 10 {
		limit = 10
	}
	for i := 0; i < limit; i++ {
		vals := mi.Get(i)
		sb.WriteString(fmt.Sprintf("  %v", vals))
		if i < limit-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	if n > 10 {
		sb.WriteString(fmt.Sprintf("  ... (%d more)\n", n-10))
	}
	sb.WriteString(fmt.Sprintf("], names=%v)", mi.names))
	return sb.String()
}

// Droplevel returns a new MultiIndex with the specified level removed.
func (mi *MultiIndex) Droplevel(level int) (*MultiIndex, error) {
	if level < 0 || level >= len(mi.levels) {
		return nil, fmt.Errorf("tabgo: MultiIndex.Droplevel: level %d out of range [0, %d)", level, len(mi.levels))
	}
	if len(mi.levels) <= 1 {
		return nil, fmt.Errorf("tabgo: MultiIndex.Droplevel: cannot drop the only level")
	}

	newLevels := make([][]any, 0, len(mi.levels)-1)
	newCodes := make([][]int, 0, len(mi.codes)-1)
	newNames := make([]string, 0, len(mi.names)-1)

	for i := range mi.levels {
		if i == level {
			continue
		}
		newLevels = append(newLevels, mi.levels[i])
		newCodes = append(newCodes, mi.codes[i])
		newNames = append(newNames, mi.names[i])
	}

	return &MultiIndex{
		levels: newLevels,
		codes:  newCodes,
		names:  newNames,
	}, nil
}

// Sortlevel returns a new MultiIndex sorted by the specified level.
func (mi *MultiIndex) Sortlevel(level int) (*MultiIndex, error) {
	if level < 0 || level >= len(mi.levels) {
		return nil, fmt.Errorf("tabgo: MultiIndex.Sortlevel: level %d out of range", level)
	}

	n := mi.Len()
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	sort.SliceStable(indices, func(a, b int) bool {
		va := fmt.Sprintf("%v", mi.levels[level][mi.codes[level][indices[a]]])
		vb := fmt.Sprintf("%v", mi.levels[level][mi.codes[level][indices[b]]])
		return va < vb
	})

	newCodes := make([][]int, len(mi.codes))
	for lvl := range mi.codes {
		nc := make([]int, n)
		for i, idx := range indices {
			nc[i] = mi.codes[lvl][idx]
		}
		newCodes[lvl] = nc
	}

	return &MultiIndex{
		levels: mi.levels,
		codes:  newCodes,
		names:  mi.names,
	}, nil
}

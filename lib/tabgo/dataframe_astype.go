package tabgo

import "fmt"

// Astype returns a new DataFrame with the specified column cast to the given dtype.
// Supported dtypes: "float64", "int", "string", "bool".
func (df *DataFrame) Astype(column, dtype string) *DataFrame {
	names := df.Columns()
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		if n == column {
			newCols[i] = df.Column(n).Astype(dtype)
		} else {
			newCols[i] = NewSeries(n, df.Column(n).Values())
		}
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// ConvertDtypes returns a new DataFrame where each column's values are converted
// to the most appropriate type. Numeric strings become float64, boolean strings
// become bool, and others remain strings.
func (df *DataFrame) ConvertDtypes() *DataFrame {
	names := df.Columns()
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		vals := df.Column(n).Values()
		converted := convertColumnDtypes(vals)
		newCols[i] = NewSeries(n, converted)
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// convertColumnDtypes attempts to convert a column of values to a uniform type.
func convertColumnDtypes(vals []any) []any {
	// Check if all non-nil values are strings.
	allStrings := true
	for _, v := range vals {
		if v == nil {
			continue
		}
		if _, ok := v.(string); !ok {
			allStrings = false
			break
		}
	}
	if !allStrings {
		// Already typed; return as-is.
		cp := make([]any, len(vals))
		copy(cp, vals)
		return cp
	}

	// Try converting all to float64.
	floats := make([]any, len(vals))
	canFloat := true
	for i, v := range vals {
		if v == nil {
			floats[i] = nil
			continue
		}
		s := v.(string)
		if s == "" {
			floats[i] = nil
			continue
		}
		var f float64
		n, _ := fmt.Sscan(s, &f)
		if n == 0 {
			canFloat = false
			break
		}
		floats[i] = f
	}
	if canFloat {
		return floats
	}

	// Try converting all to bool.
	bools := make([]any, len(vals))
	canBool := true
	for i, v := range vals {
		if v == nil {
			bools[i] = nil
			continue
		}
		s := v.(string)
		switch s {
		case "true", "True", "TRUE", "1":
			bools[i] = true
		case "false", "False", "FALSE", "0":
			bools[i] = false
		default:
			canBool = false
		}
		if !canBool {
			break
		}
	}
	if canBool {
		return bools
	}

	// Keep as strings.
	cp := make([]any, len(vals))
	copy(cp, vals)
	return cp
}

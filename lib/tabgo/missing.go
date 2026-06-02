package tabgo

// IsNA returns a boolean slice where true indicates the value is nil.
func (s *Series) IsNA() []bool {
	out := make([]bool, s.Len())
	for i, v := range s.values {
		out[i] = v == nil
	}
	return out
}

// DropNA returns a new DataFrame with rows removed where any column value is nil.
func (df *DataFrame) DropNA() *DataFrame {
	return df.Filter(func(row map[string]any) bool {
		for _, v := range row {
			if v == nil {
				return false
			}
		}
		return true
	})
}

// FillNA returns a new DataFrame where all nil values are replaced with value.
func (df *DataFrame) FillNA(value any) *DataFrame {
	names := df.Columns()
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		vals := df.Column(n).Values()
		for j, v := range vals {
			if v == nil {
				vals[j] = value
			}
		}
		newCols[i] = NewSeries(n, vals)
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// FillNAColumn returns a new DataFrame where nil values in the specified column
// are replaced with value. Other columns are unchanged.
func (df *DataFrame) FillNAColumn(column string, value any) *DataFrame {
	names := df.Columns()
	newCols := make([]*Series, len(names))
	newIdx := make(map[string]int, len(names))
	for i, n := range names {
		vals := df.Column(n).Values()
		if n == column {
			for j, v := range vals {
				if v == nil {
					vals[j] = value
				}
			}
		}
		newCols[i] = NewSeries(n, vals)
		newIdx[n] = i
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

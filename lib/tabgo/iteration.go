package tabgo

// IterRows returns a slice of maps, one per row, mapping column name to value.
func IterRows(df *DataFrame) []map[string]any {
	names := df.Columns()
	nRows := df.Len()
	allVals := make([][]any, len(names))
	for i, n := range names {
		allVals[i] = df.Column(n).Values()
	}

	rows := make([]map[string]any, nRows)
	for r := 0; r < nRows; r++ {
		row := make(map[string]any, len(names))
		for ci, n := range names {
			row[n] = allVals[ci][r]
		}
		rows[r] = row
	}
	return rows
}

// IterCols returns a map from column name to *Series.
func IterCols(df *DataFrame) map[string]*Series {
	names := df.Columns()
	result := make(map[string]*Series, len(names))
	for _, n := range names {
		s := df.Column(n)
		result[n] = NewSeries(n, s.Values())
	}
	return result
}

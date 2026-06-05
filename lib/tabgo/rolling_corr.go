package tabgo

import "math"

// RollingCorr computes an NxN rolling correlation matrix for each time step.
// Returns a slice of DataFrames, one per row. Rows before the window is full are nil.
func RollingCorr(df *DataFrame, window int) []*DataFrame {
	cols := dfNumericColumns(df)
	nRows := df.Len()
	nCols := len(cols)

	// Pre-extract all column values as float64 slices.
	colVals := make([][]any, nCols)
	for ci, c := range cols {
		colVals[ci] = df.Column(c).Values()
	}

	result := make([]*DataFrame, nRows)
	for row := 0; row < nRows; row++ {
		if row < window-1 {
			result[row] = nil
			continue
		}
		// Build the NxN correlation matrix as a DataFrame.
		seriesMap := make(map[string]*Series, nCols)
		for ci, cName := range cols {
			data := make([]any, nCols)
			wi := extractWindow(colVals[ci], row, window)
			for cj := range cols {
				if ci == cj {
					data[cj] = 1.0
					continue
				}
				wj := extractWindow(colVals[cj], row, window)
				data[cj] = pearsonCorr(wi, wj)
			}
			seriesMap[cName] = NewSeries(cName, data)
		}
		result[row] = NewDataFrame(seriesMap)
	}
	return result
}

// RollingCorrPair computes pairwise rolling Pearson correlation between two series.
func RollingCorrPair(s1, s2 *Series, window int) *Series {
	v1 := s1.Values()
	v2 := s2.Values()
	n := len(v1)
	if len(v2) < n {
		n = len(v2)
	}
	data := make([]any, n)
	for i := 0; i < n; i++ {
		if i < window-1 {
			data[i] = nil
			continue
		}
		w1 := extractWindow(v1, i, window)
		w2 := extractWindow(v2, i, window)
		if len(w1) < 2 || len(w2) < 2 || len(w1) != len(w2) {
			data[i] = nil
			continue
		}
		data[i] = pearsonCorr(w1, w2)
	}
	return NewSeries(s1.Name()+"_"+s2.Name()+"_corr", data)
}

// EWMCorr computes exponentially weighted correlation matrices.
// Returns a slice of DataFrames, one per row.
func EWMCorr(df *DataFrame, span int) []*DataFrame {
	cols := dfNumericColumns(df)
	nRows := df.Len()
	nCols := len(cols)
	alpha := 2.0 / (float64(span) + 1.0)

	// Pre-extract column values.
	colVals := make([][]float64, nCols)
	for ci, c := range cols {
		colVals[ci] = df.Column(c).Float64()
	}

	// Track EWM means, variances and covariances.
	ewmMean := make([]float64, nCols)
	// Store covariance matrix as flat nCols*nCols.
	ewmCov := make([]float64, nCols*nCols)

	result := make([]*DataFrame, nRows)

	for row := 0; row < nRows; row++ {
		if row == 0 {
			// Initialize means; variances and covariances start at 0.
			for ci := range cols {
				ewmMean[ci] = colVals[ci][0]
			}
			// First row: correlation matrix is identity-like (variance=0).
			seriesMap := make(map[string]*Series, nCols)
			for ci, cName := range cols {
				data := make([]any, nCols)
				for cj := range cols {
					if ci == cj {
						data[cj] = 1.0
					} else {
						data[cj] = 0.0
					}
				}
				seriesMap[cName] = NewSeries(cName, data)
			}
			result[row] = NewDataFrame(seriesMap)
			continue
		}

		// Update EWM covariance matrix.
		diffs := make([]float64, nCols)
		for ci := range cols {
			diffs[ci] = colVals[ci][row] - ewmMean[ci]
		}
		for ci := range cols {
			ewmMean[ci] = alpha*colVals[ci][row] + (1-alpha)*ewmMean[ci]
		}
		for ci := range cols {
			for cj := range cols {
				idx := ci*nCols + cj
				ewmCov[idx] = (1 - alpha) * (ewmCov[idx] + alpha*diffs[ci]*diffs[cj])
			}
		}

		// Build correlation matrix from covariance.
		seriesMap := make(map[string]*Series, nCols)
		for ci, cName := range cols {
			data := make([]any, nCols)
			for cj := range cols {
				if ci == cj {
					data[cj] = 1.0
				} else {
					vi := ewmCov[ci*nCols+ci]
					vj := ewmCov[cj*nCols+cj]
					denom := math.Sqrt(vi * vj)
					if denom == 0 {
						data[cj] = 0.0
					} else {
						data[cj] = ewmCov[ci*nCols+cj] / denom
					}
				}
			}
			seriesMap[cName] = NewSeries(cName, data)
		}
		result[row] = NewDataFrame(seriesMap)
	}
	return result
}

// pearsonCorr computes Pearson correlation between two equal-length float64 slices.
func pearsonCorr(x, y []float64) float64 {
	n := len(x)
	if n < 2 {
		return 0
	}
	mx := mean(x)
	my := mean(y)
	var sxy, sxx, syy float64
	for i := range x {
		dx := x[i] - mx
		dy := y[i] - my
		sxy += dx * dy
		sxx += dx * dx
		syy += dy * dy
	}
	denom := math.Sqrt(sxx * syy)
	if denom == 0 {
		return 0
	}
	return sxy / denom
}

package tabgo

import "math"

// Rolling represents a rolling window operation on a DataFrame.
type Rolling struct {
	df     *DataFrame
	window int
}

// Rolling returns a Rolling object for rolling window calculations.
func (df *DataFrame) Rolling(window int) *Rolling {
	return &Rolling{df: df, window: window}
}

// Mean returns a DataFrame with the rolling mean of numeric columns.
func (r *Rolling) Mean() *DataFrame {
	return r.apply(func(vals []float64) float64 {
		if len(vals) == 0 {
			return 0
		}
		var s float64
		for _, v := range vals {
			s += v
		}
		return s / float64(len(vals))
	})
}

// Sum returns a DataFrame with the rolling sum of numeric columns.
func (r *Rolling) Sum() *DataFrame {
	return r.apply(func(vals []float64) float64 {
		var s float64
		for _, v := range vals {
			s += v
		}
		return s
	})
}

// Std returns a DataFrame with the rolling standard deviation of numeric columns.
func (r *Rolling) Std() *DataFrame {
	return r.apply(func(vals []float64) float64 {
		if len(vals) < 2 {
			return 0
		}
		return stddev(vals)
	})
}

// Min returns a DataFrame with the rolling minimum of numeric columns.
func (r *Rolling) Min() *DataFrame {
	return r.apply(func(vals []float64) float64 {
		if len(vals) == 0 {
			return 0
		}
		m := vals[0]
		for _, v := range vals[1:] {
			if v < m {
				m = v
			}
		}
		return m
	})
}

// Max returns a DataFrame with the rolling maximum of numeric columns.
func (r *Rolling) Max() *DataFrame {
	return r.apply(func(vals []float64) float64 {
		if len(vals) == 0 {
			return 0
		}
		m := vals[0]
		for _, v := range vals[1:] {
			if v > m {
				m = v
			}
		}
		return m
	})
}

// Count returns a DataFrame with the rolling non-nil count of numeric columns.
func (r *Rolling) Count() *DataFrame {
	return r.apply(func(vals []float64) float64 {
		return float64(len(vals))
	})
}

// apply applies fn to each rolling window for each numeric column.
func (r *Rolling) apply(fn func([]float64) float64) *DataFrame {
	cols := dfNumericColumns(r.df)
	nRows := r.df.Len()

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for ci, c := range cols {
		vals := r.df.Column(c).Values()
		data := make([]any, nRows)
		for row := 0; row < nRows; row++ {
			if row < r.window-1 {
				data[row] = nil
				continue
			}
			window := make([]float64, 0, r.window)
			for w := row - r.window + 1; w <= row; w++ {
				if vals[w] != nil {
					window = append(window, toFloat64(vals[w]))
				}
			}
			data[row] = fn(window)
		}
		newCols[ci] = NewSeries(c, data)
		newIdx[c] = ci
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Expanding represents an expanding window operation on a DataFrame.
type Expanding struct {
	df *DataFrame
}

// Expanding returns an Expanding object for expanding window calculations.
func (df *DataFrame) Expanding() *Expanding {
	return &Expanding{df: df}
}

// Mean returns a DataFrame with the expanding mean of numeric columns.
func (e *Expanding) Mean() *DataFrame {
	return e.apply(func(vals []float64) float64 {
		if len(vals) == 0 {
			return 0
		}
		var s float64
		for _, v := range vals {
			s += v
		}
		return s / float64(len(vals))
	})
}

// Sum returns a DataFrame with the expanding sum of numeric columns.
func (e *Expanding) Sum() *DataFrame {
	return e.apply(func(vals []float64) float64 {
		var s float64
		for _, v := range vals {
			s += v
		}
		return s
	})
}

// Std returns a DataFrame with the expanding standard deviation of numeric columns.
func (e *Expanding) Std() *DataFrame {
	return e.apply(func(vals []float64) float64 {
		if len(vals) < 2 {
			return 0
		}
		return stddev(vals)
	})
}

// Min returns a DataFrame with the expanding minimum of numeric columns.
func (e *Expanding) Min() *DataFrame {
	return e.apply(func(vals []float64) float64 {
		if len(vals) == 0 {
			return 0
		}
		m := vals[0]
		for _, v := range vals[1:] {
			if v < m {
				m = v
			}
		}
		return m
	})
}

// Max returns a DataFrame with the expanding maximum of numeric columns.
func (e *Expanding) Max() *DataFrame {
	return e.apply(func(vals []float64) float64 {
		if len(vals) == 0 {
			return 0
		}
		m := vals[0]
		for _, v := range vals[1:] {
			if v > m {
				m = v
			}
		}
		return m
	})
}

// apply applies fn to the expanding window for each numeric column.
func (e *Expanding) apply(fn func([]float64) float64) *DataFrame {
	cols := dfNumericColumns(e.df)
	nRows := e.df.Len()

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for ci, c := range cols {
		vals := e.df.Column(c).Values()
		data := make([]any, nRows)
		window := make([]float64, 0, nRows)
		for row := 0; row < nRows; row++ {
			if vals[row] != nil {
				window = append(window, toFloat64(vals[row]))
			}
			data[row] = fn(window)
		}
		newCols[ci] = NewSeries(c, data)
		newIdx[c] = ci
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// EWM represents exponentially weighted moving calculations on a DataFrame.
type EWM struct {
	df    *DataFrame
	alpha float64
}

// EWM returns an EWM object for exponentially weighted calculations.
// span is the decay in terms of span: alpha = 2 / (span + 1).
func (df *DataFrame) EWM(span float64) *EWM {
	alpha := 2.0 / (span + 1.0)
	return &EWM{df: df, alpha: alpha}
}

// Mean returns a DataFrame with the exponentially weighted mean of numeric columns.
func (ew *EWM) Mean() *DataFrame {
	cols := dfNumericColumns(ew.df)
	nRows := ew.df.Len()

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for ci, c := range cols {
		vals := ew.df.Column(c).Values()
		data := make([]any, nRows)
		first := true
		var ewmVal float64
		for row := 0; row < nRows; row++ {
			if vals[row] == nil {
				data[row] = nil
				continue
			}
			f := toFloat64(vals[row])
			if first {
				ewmVal = f
				first = false
			} else {
				ewmVal = ew.alpha*f + (1-ew.alpha)*ewmVal
			}
			data[row] = ewmVal
		}
		newCols[ci] = NewSeries(c, data)
		newIdx[c] = ci
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Std returns a DataFrame with the exponentially weighted standard deviation.
func (ew *EWM) Std() *DataFrame {
	cols := dfNumericColumns(ew.df)
	nRows := ew.df.Len()

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for ci, c := range cols {
		vals := ew.df.Column(c).Values()
		data := make([]any, nRows)
		first := true
		var ewmMean, ewmVar float64
		for row := 0; row < nRows; row++ {
			if vals[row] == nil {
				data[row] = nil
				continue
			}
			f := toFloat64(vals[row])
			if first {
				ewmMean = f
				ewmVar = 0
				first = false
				data[row] = 0.0
			} else {
				diff := f - ewmMean
				ewmMean = ew.alpha*f + (1-ew.alpha)*ewmMean
				ewmVar = (1 - ew.alpha) * (ewmVar + ew.alpha*diff*diff)
				data[row] = math.Sqrt(ewmVar)
			}
		}
		newCols[ci] = NewSeries(c, data)
		newIdx[c] = ci
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

// Var returns a DataFrame with the exponentially weighted variance.
func (ew *EWM) Var() *DataFrame {
	cols := dfNumericColumns(ew.df)
	nRows := ew.df.Len()

	newCols := make([]*Series, len(cols))
	newIdx := make(map[string]int, len(cols))
	for ci, c := range cols {
		vals := ew.df.Column(c).Values()
		data := make([]any, nRows)
		first := true
		var ewmMean, ewmVar float64
		for row := 0; row < nRows; row++ {
			if vals[row] == nil {
				data[row] = nil
				continue
			}
			f := toFloat64(vals[row])
			if first {
				ewmMean = f
				ewmVar = 0
				first = false
				data[row] = 0.0
			} else {
				diff := f - ewmMean
				ewmMean = ew.alpha*f + (1-ew.alpha)*ewmMean
				ewmVar = (1 - ew.alpha) * (ewmVar + ew.alpha*diff*diff)
				data[row] = ewmVar
			}
		}
		newCols[ci] = NewSeries(c, data)
		newIdx[c] = ci
	}
	return &DataFrame{columns: newCols, index: newIdx}
}

package tabgo

// RollingBeta computes the rolling regression beta: Cov(asset, benchmark) / Var(benchmark).
func RollingBeta(asset, benchmark *Series, window int) *Series {
	va := asset.Values()
	vb := benchmark.Values()
	n := len(va)
	if len(vb) < n {
		n = len(vb)
	}
	data := make([]any, n)
	for i := 0; i < n; i++ {
		if i < window-1 {
			data[i] = nil
			continue
		}
		wa := extractWindow(va, i, window)
		wb := extractWindow(vb, i, window)
		if len(wa) < 2 || len(wb) < 2 || len(wa) != len(wb) {
			data[i] = nil
			continue
		}
		varB := variance(wb)
		if varB == 0 {
			data[i] = 0.0
			continue
		}
		covAB := covariance(wa, wb)
		data[i] = covAB / varB
	}
	return NewSeries(asset.Name()+"_beta", data)
}

// RollingAlpha computes the rolling regression alpha: mean(asset) - beta * mean(benchmark).
func RollingAlpha(asset, benchmark *Series, window int) *Series {
	va := asset.Values()
	vb := benchmark.Values()
	n := len(va)
	if len(vb) < n {
		n = len(vb)
	}
	data := make([]any, n)
	for i := 0; i < n; i++ {
		if i < window-1 {
			data[i] = nil
			continue
		}
		wa := extractWindow(va, i, window)
		wb := extractWindow(vb, i, window)
		if len(wa) < 2 || len(wb) < 2 || len(wa) != len(wb) {
			data[i] = nil
			continue
		}
		varB := variance(wb)
		mA := mean(wa)
		mB := mean(wb)
		if varB == 0 {
			data[i] = mA
			continue
		}
		beta := covariance(wa, wb) / varB
		data[i] = mA - beta*mB
	}
	return NewSeries(asset.Name()+"_alpha", data)
}

// RollingR2 computes the rolling coefficient of determination (R-squared).
func RollingR2(asset, benchmark *Series, window int) *Series {
	va := asset.Values()
	vb := benchmark.Values()
	n := len(va)
	if len(vb) < n {
		n = len(vb)
	}
	data := make([]any, n)
	for i := 0; i < n; i++ {
		if i < window-1 {
			data[i] = nil
			continue
		}
		wa := extractWindow(va, i, window)
		wb := extractWindow(vb, i, window)
		if len(wa) < 2 || len(wb) < 2 || len(wa) != len(wb) {
			data[i] = nil
			continue
		}
		r := pearsonCorr(wa, wb)
		data[i] = r * r
	}
	return NewSeries(asset.Name()+"_r2", data)
}

// RollingResidual computes rolling regression residuals: asset - (alpha + beta * benchmark).
func RollingResidual(asset, benchmark *Series, window int) *Series {
	va := asset.Values()
	vb := benchmark.Values()
	n := len(va)
	if len(vb) < n {
		n = len(vb)
	}
	data := make([]any, n)
	for i := 0; i < n; i++ {
		if i < window-1 {
			data[i] = nil
			continue
		}
		wa := extractWindow(va, i, window)
		wb := extractWindow(vb, i, window)
		if len(wa) < 2 || len(wb) < 2 || len(wa) != len(wb) {
			data[i] = nil
			continue
		}
		varB := variance(wb)
		mA := mean(wa)
		mB := mean(wb)
		var beta, alpha float64
		if varB == 0 {
			alpha = mA
		} else {
			beta = covariance(wa, wb) / varB
			alpha = mA - beta*mB
		}
		a := toFloat64(va[i])
		b := toFloat64(vb[i])
		data[i] = a - (alpha + beta*b)
	}
	return NewSeries(asset.Name()+"_residual", data)
}

// RollingMultiBeta computes rolling OLS betas for multiple factors per window.
// Returns a DataFrame where each column is the rolling beta for that factor.
func RollingMultiBeta(asset *Series, factors *DataFrame, window int) *DataFrame {
	factorNames := dfNumericColumns(factors)
	nFactors := len(factorNames)
	va := asset.Values()
	n := len(va)

	// Pre-extract factor values.
	factorVals := make([][]any, nFactors)
	for fi, fn := range factorNames {
		factorVals[fi] = factors.Column(fn).Values()
	}

	// Output columns for each factor beta.
	betaData := make([][]any, nFactors)
	for fi := range betaData {
		betaData[fi] = make([]any, n)
	}

	for i := 0; i < n; i++ {
		if i < window-1 {
			for fi := range betaData {
				betaData[fi][i] = nil
			}
			continue
		}
		wa := extractWindow(va, i, window)
		wf := make([][]float64, nFactors)
		valid := true
		for fi := range factorNames {
			wf[fi] = extractWindow(factorVals[fi], i, window)
			if len(wf[fi]) < 2 || len(wf[fi]) != len(wa) {
				valid = false
				break
			}
		}
		if !valid {
			for fi := range betaData {
				betaData[fi][i] = nil
			}
			continue
		}

		// Solve OLS using normal equations: beta = (X'X)^-1 X'y
		betas := solveOLS(wa, wf)
		for fi := range factorNames {
			if fi < len(betas) {
				betaData[fi][i] = betas[fi]
			} else {
				betaData[fi][i] = nil
			}
		}
	}

	seriesMap := make(map[string]*Series, nFactors)
	for fi, fn := range factorNames {
		seriesMap[fn] = NewSeries(fn, betaData[fi])
	}
	return NewDataFrame(seriesMap)
}

// solveOLS solves the ordinary least squares problem y = X*beta.
// y is the dependent variable, cols are the independent variable columns.
// Returns betas for each factor.
func solveOLS(y []float64, cols [][]float64) []float64 {
	nFactors := len(cols)
	n := len(y)

	// Build X'X matrix (nFactors x nFactors) and X'y vector.
	xtx := make([]float64, nFactors*nFactors)
	xty := make([]float64, nFactors)

	// Demean for regression through mean (no intercept in betas, intercept implicit).
	yMean := mean(y)
	colMeans := make([]float64, nFactors)
	for fi := range cols {
		colMeans[fi] = mean(cols[fi])
	}

	// Build centered versions.
	yc := make([]float64, n)
	xc := make([][]float64, nFactors)
	for fi := range cols {
		xc[fi] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		yc[i] = y[i] - yMean
		for fi := range cols {
			xc[fi][i] = cols[fi][i] - colMeans[fi]
		}
	}

	for fi := 0; fi < nFactors; fi++ {
		for fj := 0; fj < nFactors; fj++ {
			var s float64
			for k := 0; k < n; k++ {
				s += xc[fi][k] * xc[fj][k]
			}
			xtx[fi*nFactors+fj] = s
		}
		var s float64
		for k := 0; k < n; k++ {
			s += xc[fi][k] * yc[k]
		}
		xty[fi] = s
	}

	// Solve via Gaussian elimination with partial pivoting.
	return gaussianSolve(xtx, xty, nFactors)
}

// gaussianSolve solves A*x = b using Gaussian elimination with partial pivoting.
func gaussianSolve(A []float64, b []float64, n int) []float64 {
	// Make copies.
	a := make([]float64, n*n)
	copy(a, A)
	x := make([]float64, n)
	rhs := make([]float64, n)
	copy(rhs, b)

	// Forward elimination with partial pivoting.
	for col := 0; col < n; col++ {
		// Find pivot.
		maxVal := 0.0
		maxRow := col
		for row := col; row < n; row++ {
			v := a[row*n+col]
			if v < 0 {
				v = -v
			}
			if v > maxVal {
				maxVal = v
				maxRow = row
			}
		}
		if maxVal < 1e-15 {
			// Singular; return zeros.
			return x
		}
		// Swap rows.
		if maxRow != col {
			for j := 0; j < n; j++ {
				a[col*n+j], a[maxRow*n+j] = a[maxRow*n+j], a[col*n+j]
			}
			rhs[col], rhs[maxRow] = rhs[maxRow], rhs[col]
		}
		// Eliminate below.
		for row := col + 1; row < n; row++ {
			factor := a[row*n+col] / a[col*n+col]
			for j := col; j < n; j++ {
				a[row*n+j] -= factor * a[col*n+j]
			}
			rhs[row] -= factor * rhs[col]
		}
	}

	// Back substitution.
	for i := n - 1; i >= 0; i-- {
		x[i] = rhs[i]
		for j := i + 1; j < n; j++ {
			x[i] -= a[i*n+j] * x[j]
		}
		if a[i*n+i] != 0 {
			x[i] /= a[i*n+i]
		}
	}
	return x
}

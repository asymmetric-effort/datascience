package tabgo

import (
	"math"
	"sort"
)

// RollingPCAResult holds the results of a rolling PCA computation.
type RollingPCAResult struct {
	// Components contains a DataFrame per time step with the principal component loadings.
	// Each DataFrame has shape (nComponents x nFeatures). Nil for rows before window is full.
	Components []*DataFrame
	// ExplainedVar contains a Series per time step with the fraction of variance explained.
	ExplainedVar []*Series
	// Eigenvalues contains a Series per time step with the eigenvalues.
	Eigenvalues []*Series
}

// RollingPCA performs rolling PCA on a DataFrame of numeric columns.
// window is the rolling window size; nComponents is the number of principal components to keep.
// Uses eigendecomposition of the windowed covariance matrix.
func RollingPCA(df *DataFrame, window, nComponents int) *RollingPCAResult {
	cols := dfNumericColumns(df)
	nRows := df.Len()
	nCols := len(cols)

	if nComponents > nCols {
		nComponents = nCols
	}
	if nComponents <= 0 {
		nComponents = nCols
	}

	// Pre-extract column data.
	colVals := make([][]any, nCols)
	for ci, c := range cols {
		colVals[ci] = df.Column(c).Values()
	}

	result := &RollingPCAResult{
		Components:   make([]*DataFrame, nRows),
		ExplainedVar: make([]*Series, nRows),
		Eigenvalues:  make([]*Series, nRows),
	}

	for row := 0; row < nRows; row++ {
		if row < window-1 {
			result.Components[row] = nil
			result.ExplainedVar[row] = nil
			result.Eigenvalues[row] = nil
			continue
		}

		// Extract window data for each column.
		windowData := make([][]float64, nCols)
		for ci := range cols {
			windowData[ci] = extractWindow(colVals[ci], row, window)
		}

		// Compute covariance matrix.
		covMat := make([][]float64, nCols)
		for i := 0; i < nCols; i++ {
			covMat[i] = make([]float64, nCols)
			for j := 0; j < nCols; j++ {
				covMat[i][j] = covariance(windowData[i], windowData[j])
			}
		}

		// Eigendecomposition using Jacobi iteration (symmetric matrix).
		eigenvalues, eigenvectors := jacobiEigen(covMat, nCols)

		// Sort eigenvalues/eigenvectors by descending eigenvalue.
		type eigPair struct {
			val float64
			vec []float64
		}
		pairs := make([]eigPair, nCols)
		for i := 0; i < nCols; i++ {
			vec := make([]float64, nCols)
			for j := 0; j < nCols; j++ {
				vec[j] = eigenvectors[j][i]
			}
			pairs[i] = eigPair{val: eigenvalues[i], vec: vec}
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].val > pairs[j].val
		})

		// Extract top nComponents.
		totalVar := 0.0
		for _, p := range pairs {
			if p.val > 0 {
				totalVar += p.val
			}
		}

		eigVals := make([]any, nComponents)
		explVar := make([]any, nComponents)
		compSeriesMap := make(map[string]*Series, nComponents)

		for pc := 0; pc < nComponents; pc++ {
			eigVals[pc] = pairs[pc].val
			if totalVar > 0 {
				ev := pairs[pc].val
				if ev < 0 {
					ev = 0
				}
				explVar[pc] = ev / totalVar
			} else {
				explVar[pc] = 0.0
			}

			// Component loadings as a Series with column names.
			compData := make([]any, nCols)
			for ci := range cols {
				compData[ci] = pairs[pc].vec[ci]
			}
			compSeriesMap[cols[pc]] = NewSeries(cols[pc], compData)
		}

		// Build components DataFrame: rows=components, cols=features
		// Actually we store component loadings where each column is a feature
		// and each row is a component.
		compColMap := make(map[string]*Series, nCols)
		for ci, cName := range cols {
			loadings := make([]any, nComponents)
			for pc := 0; pc < nComponents; pc++ {
				loadings[pc] = pairs[pc].vec[ci]
			}
			compColMap[cName] = NewSeries(cName, loadings)
		}
		result.Components[row] = NewDataFrame(compColMap)
		result.Eigenvalues[row] = NewSeries("eigenvalues", eigVals)
		result.ExplainedVar[row] = NewSeries("explained_var", explVar)
	}

	return result
}

// jacobiEigen performs Jacobi eigenvalue decomposition for a symmetric matrix.
// Returns eigenvalues and eigenvectors (column-major).
func jacobiEigen(mat [][]float64, n int) ([]float64, [][]float64) {
	// Copy matrix.
	a := make([][]float64, n)
	for i := range a {
		a[i] = make([]float64, n)
		copy(a[i], mat[i])
	}

	// Initialize eigenvector matrix to identity.
	v := make([][]float64, n)
	for i := range v {
		v[i] = make([]float64, n)
		v[i][i] = 1.0
	}

	maxIter := 100 * n * n
	if maxIter < 1000 {
		maxIter = 1000
	}

	for iter := 0; iter < maxIter; iter++ {
		// Find largest off-diagonal element.
		p, q := 0, 1
		maxOff := math.Abs(a[0][1])
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if math.Abs(a[i][j]) > maxOff {
					maxOff = math.Abs(a[i][j])
					p, q = i, j
				}
			}
		}

		if maxOff < 1e-12 {
			break
		}

		// Compute rotation.
		var theta, t, c, s float64
		if math.Abs(a[p][p]-a[q][q]) < 1e-15 {
			theta = math.Pi / 4
			c = math.Cos(theta)
			s = math.Sin(theta)
		} else {
			tau := (a[q][q] - a[p][p]) / (2 * a[p][q])
			if tau >= 0 {
				t = 1.0 / (tau + math.Sqrt(1+tau*tau))
			} else {
				t = -1.0 / (-tau + math.Sqrt(1+tau*tau))
			}
			c = 1.0 / math.Sqrt(1+t*t)
			s = t * c
		}

		// Apply rotation.
		for i := 0; i < n; i++ {
			if i == p || i == q {
				continue
			}
			aip := a[i][p]
			aiq := a[i][q]
			a[i][p] = c*aip - s*aiq
			a[p][i] = a[i][p]
			a[i][q] = s*aip + c*aiq
			a[q][i] = a[i][q]
		}

		app := a[p][p]
		aqq := a[q][q]
		apq := a[p][q]
		a[p][p] = c*c*app - 2*s*c*apq + s*s*aqq
		a[q][q] = s*s*app + 2*s*c*apq + c*c*aqq
		a[p][q] = 0
		a[q][p] = 0

		// Update eigenvectors.
		for i := 0; i < n; i++ {
			vip := v[i][p]
			viq := v[i][q]
			v[i][p] = c*vip - s*viq
			v[i][q] = s*vip + c*viq
		}
	}

	eigenvalues := make([]float64, n)
	for i := 0; i < n; i++ {
		eigenvalues[i] = a[i][i]
	}

	return eigenvalues, v
}

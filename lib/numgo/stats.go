package numgo

import (
	"fmt"
	"math"
	"sort"
)

// Min reduces the array by taking the minimum along the given axes.
// If no axes are given, it returns the global minimum as a scalar (1-D, length-1) array.
func Min(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		m := math.Inf(1)
		for _, v := range a.data {
			if v < m {
				m = v
			}
		}
		return FromSlice([]float64{m})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		m := math.Inf(1)
		for _, v := range vals {
			if v < m {
				m = v
			}
		}
		return m
	})
}

// Mean reduces the array by computing the arithmetic mean along the given axes.
// If no axes are given, it returns the global mean.
func Mean(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		total := 0.0
		for _, v := range a.data {
			total += v
		}
		return FromSlice([]float64{total / float64(a.Size())})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		s := 0.0
		for _, v := range vals {
			s += v
		}
		return s / float64(len(vals))
	})
}

// Var computes the population variance along the given axes.
// If no axes are given, it returns the global variance.
func Var(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		mean := 0.0
		for _, v := range a.data {
			mean += v
		}
		mean /= float64(a.Size())
		variance := 0.0
		for _, v := range a.data {
			d := v - mean
			variance += d * d
		}
		return FromSlice([]float64{variance / float64(a.Size())})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))
		variance := 0.0
		for _, v := range vals {
			d := v - mean
			variance += d * d
		}
		return variance / float64(len(vals))
	})
}

// Std computes the population standard deviation along the given axes.
// If no axes are given, it returns the global standard deviation.
func Std(a *NDArray, axes ...int) *NDArray {
	v := Var(a, axes...)
	data := v.Data()
	for i := range data {
		data[i] = math.Sqrt(data[i])
	}
	return NewNDArray(v.Shape(), data)
}

// Cumsum returns the cumulative sum along the given axis.
func Cumsum(a *NDArray, axis int) *NDArray {
	if axis < 0 || axis >= a.Ndim() {
		panic(fmt.Sprintf("numgo.Cumsum: axis %d out of range for %d dimensions", axis, a.Ndim()))
	}
	result := a.Copy()
	axisLen := a.shape[axis]

	iterateAlongAxis(result, axis, func(indices []int) {
		for i := 1; i < axisLen; i++ {
			indices[axis] = i - 1
			prev := result.data[result.flatIndex(indices)]
			indices[axis] = i
			fi := result.flatIndex(indices)
			result.data[fi] += prev
		}
	})
	return result
}

// Cumprod returns the cumulative product along the given axis.
func Cumprod(a *NDArray, axis int) *NDArray {
	if axis < 0 || axis >= a.Ndim() {
		panic(fmt.Sprintf("numgo.Cumprod: axis %d out of range for %d dimensions", axis, a.Ndim()))
	}
	result := a.Copy()
	axisLen := a.shape[axis]

	iterateAlongAxis(result, axis, func(indices []int) {
		for i := 1; i < axisLen; i++ {
			indices[axis] = i - 1
			prev := result.data[result.flatIndex(indices)]
			indices[axis] = i
			fi := result.flatIndex(indices)
			result.data[fi] *= prev
		}
	})
	return result
}

// Percentile returns the q-th percentile of the array along the given axes.
// q must be in [0, 100]. Uses linear interpolation.
func Percentile(a *NDArray, q float64, axes ...int) *NDArray {
	if q < 0 || q > 100 {
		panic(fmt.Sprintf("numgo.Percentile: q=%g out of range [0, 100]", q))
	}
	return quantileImpl(a, q/100.0, axes...)
}

// Quantile returns the q-th quantile of the array along the given axes.
// q must be in [0, 1]. Uses linear interpolation.
func Quantile(a *NDArray, q float64, axes ...int) *NDArray {
	if q < 0 || q > 1 {
		panic(fmt.Sprintf("numgo.Quantile: q=%g out of range [0, 1]", q))
	}
	return quantileImpl(a, q, axes...)
}

func quantileImpl(a *NDArray, q float64, axes ...int) *NDArray {
	if len(axes) == 0 {
		sorted := make([]float64, a.Size())
		copy(sorted, a.data)
		sort.Float64s(sorted)
		return FromSlice([]float64{interpolatedQuantile(sorted, q)})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		tmp := make([]float64, len(vals))
		copy(tmp, vals)
		sort.Float64s(tmp)
		return interpolatedQuantile(tmp, q)
	})
}

func interpolatedQuantile(sorted []float64, q float64) float64 {
	n := len(sorted)
	if n == 0 {
		return math.NaN()
	}
	if n == 1 {
		return sorted[0]
	}
	pos := q * float64(n-1)
	lo := int(math.Floor(pos))
	hi := lo + 1
	if hi >= n {
		return sorted[n-1]
	}
	frac := pos - float64(lo)
	return sorted[lo]*(1-frac) + sorted[hi]*frac
}

// Median returns the median of the array along the given axes.
func Median(a *NDArray, axes ...int) *NDArray {
	return quantileImpl(a, 0.5, axes...)
}

// Average computes the weighted average along the given axes.
// If weights is nil, all weights are equal (equivalent to Mean).
// When axes are specified, weights must have length equal to the axis size.
func Average(a *NDArray, weights *NDArray, axes ...int) *NDArray {
	if weights == nil {
		return Mean(a, axes...)
	}
	if len(axes) == 0 {
		if weights.Size() != a.Size() {
			panic(fmt.Sprintf("numgo.Average: weights length %d != array size %d", weights.Size(), a.Size()))
		}
		num := 0.0
		den := 0.0
		for i, v := range a.data {
			w := weights.data[i]
			num += v * w
			den += w
		}
		return FromSlice([]float64{num / den})
	}
	// With axis: weights must match axis length.
	if len(axes) != 1 {
		panic("numgo.Average: weighted average supports at most one axis")
	}
	ax := axes[0]
	if ax < 0 || ax >= a.Ndim() {
		panic(fmt.Sprintf("numgo.Average: axis %d out of range for %d dimensions", ax, a.Ndim()))
	}
	axLen := a.shape[ax]
	if weights.Size() != axLen {
		panic(fmt.Sprintf("numgo.Average: weights length %d != axis size %d", weights.Size(), axLen))
	}
	wdata := weights.Data()
	return reduceAxis(a, axes, func(vals []float64) float64 {
		num := 0.0
		den := 0.0
		for i, v := range vals {
			num += v * wdata[i]
			den += wdata[i]
		}
		return num / den
	})
}

// Nanmean computes the arithmetic mean, ignoring NaN values.
func Nanmean(a *NDArray, axes ...int) *NDArray {
	return nanReduce(a, axes, func(vals []float64) float64 {
		s := 0.0
		n := 0
		for _, v := range vals {
			if !math.IsNaN(v) {
				s += v
				n++
			}
		}
		if n == 0 {
			return math.NaN()
		}
		return s / float64(n)
	})
}

// Nanstd computes the population standard deviation, ignoring NaN values.
func Nanstd(a *NDArray, axes ...int) *NDArray {
	v := Nanvar(a, axes...)
	data := v.Data()
	for i := range data {
		data[i] = math.Sqrt(data[i])
	}
	return NewNDArray(v.Shape(), data)
}

// Nanvar computes the population variance, ignoring NaN values.
func Nanvar(a *NDArray, axes ...int) *NDArray {
	return nanReduce(a, axes, func(vals []float64) float64 {
		s := 0.0
		n := 0
		for _, v := range vals {
			if !math.IsNaN(v) {
				s += v
				n++
			}
		}
		if n == 0 {
			return math.NaN()
		}
		mean := s / float64(n)
		variance := 0.0
		for _, v := range vals {
			if !math.IsNaN(v) {
				d := v - mean
				variance += d * d
			}
		}
		return variance / float64(n)
	})
}

// Nanmin returns the minimum, ignoring NaN values.
func Nanmin(a *NDArray, axes ...int) *NDArray {
	return nanReduce(a, axes, func(vals []float64) float64 {
		m := math.Inf(1)
		found := false
		for _, v := range vals {
			if !math.IsNaN(v) && v < m {
				m = v
				found = true
			}
		}
		if !found {
			return math.NaN()
		}
		return m
	})
}

// Nanmax returns the maximum, ignoring NaN values.
func Nanmax(a *NDArray, axes ...int) *NDArray {
	return nanReduce(a, axes, func(vals []float64) float64 {
		m := math.Inf(-1)
		found := false
		for _, v := range vals {
			if !math.IsNaN(v) && v > m {
				m = v
				found = true
			}
		}
		if !found {
			return math.NaN()
		}
		return m
	})
}

// Nansum returns the sum, treating NaN as zero.
func Nansum(a *NDArray, axes ...int) *NDArray {
	return nanReduce(a, axes, func(vals []float64) float64 {
		s := 0.0
		for _, v := range vals {
			if !math.IsNaN(v) {
				s += v
			}
		}
		return s
	})
}

// Nanprod returns the product, treating NaN as one.
func Nanprod(a *NDArray, axes ...int) *NDArray {
	return nanReduce(a, axes, func(vals []float64) float64 {
		p := 1.0
		for _, v := range vals {
			if !math.IsNaN(v) {
				p *= v
			}
		}
		return p
	})
}

// nanReduce is a helper that applies a reduction function, used by Nan* functions.
func nanReduce(a *NDArray, axes []int, fn func([]float64) float64) *NDArray {
	if len(axes) == 0 {
		return FromSlice([]float64{fn(a.data)})
	}
	return reduceAxis(a, axes, fn)
}

// Histogram computes a histogram of a flat array.
// Returns counts (length bins) and edges (length bins+1).
func Histogram(a *NDArray, bins int) (counts, edges *NDArray) {
	if bins <= 0 {
		panic("numgo.Histogram: bins must be > 0")
	}
	mn := math.Inf(1)
	mx := math.Inf(-1)
	for _, v := range a.data {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	if mn == mx {
		mx = mn + 1
	}

	edgeData := make([]float64, bins+1)
	step := (mx - mn) / float64(bins)
	for i := 0; i <= bins; i++ {
		edgeData[i] = mn + float64(i)*step
	}

	countData := make([]float64, bins)
	for _, v := range a.data {
		idx := int((v - mn) / step)
		if idx >= bins {
			idx = bins - 1
		}
		if idx < 0 {
			idx = 0
		}
		countData[idx]++
	}

	return FromSlice(countData), FromSlice(edgeData)
}

// Bincount counts the number of occurrences of each non-negative integer value.
// Values are truncated to integers. The result length is max(a)+1.
func Bincount(a *NDArray) *NDArray {
	mx := 0
	for _, v := range a.data {
		iv := int(v)
		if iv < 0 {
			panic("numgo.Bincount: negative value encountered")
		}
		if iv > mx {
			mx = iv
		}
	}
	counts := make([]float64, mx+1)
	for _, v := range a.data {
		counts[int(v)]++
	}
	return FromSlice(counts)
}

// Corrcoef returns the Pearson correlation coefficient matrix for x and y.
// x and y must be 1-D arrays of the same length.
// Returns a 2x2 correlation matrix.
func Corrcoef(x, y *NDArray) (*NDArray, error) {
	if x.Ndim() != 1 || y.Ndim() != 1 {
		return nil, fmt.Errorf("numgo.Corrcoef: inputs must be 1-D")
	}
	if x.Size() != y.Size() {
		return nil, fmt.Errorf("numgo.Corrcoef: inputs must have same length")
	}

	n := float64(x.Size())
	mx := 0.0
	my := 0.0
	for i := 0; i < x.Size(); i++ {
		mx += x.data[i]
		my += y.data[i]
	}
	mx /= n
	my /= n

	sxx := 0.0
	syy := 0.0
	sxy := 0.0
	for i := 0; i < x.Size(); i++ {
		dx := x.data[i] - mx
		dy := y.data[i] - my
		sxx += dx * dx
		syy += dy * dy
		sxy += dx * dy
	}

	r := sxy / math.Sqrt(sxx*syy)
	return NewNDArray([]int{2, 2}, []float64{1, r, r, 1}), nil
}

// Cov returns the covariance matrix for a 2-D array where each row is a variable
// and each column is an observation. For a 1-D array, returns a 1x1 matrix.
func Cov(x *NDArray) (*NDArray, error) {
	if x.Ndim() == 1 {
		v := Var(x)
		// Use sample variance (N-1)
		n := float64(x.Size())
		if n <= 1 {
			return FromSlice([]float64{0}).Reshape(1, 1), nil
		}
		popVar := v.Data()[0]
		sampleVar := popVar * n / (n - 1)
		return NewNDArray([]int{1, 1}, []float64{sampleVar}), nil
	}
	if x.Ndim() != 2 {
		return nil, fmt.Errorf("numgo.Cov: input must be 1-D or 2-D")
	}

	nVars := x.shape[0]
	nObs := x.shape[1]
	if nObs <= 1 {
		return NewNDArray([]int{nVars, nVars}, make([]float64, nVars*nVars)), nil
	}

	// Compute means for each variable.
	means := make([]float64, nVars)
	for i := 0; i < nVars; i++ {
		s := 0.0
		for j := 0; j < nObs; j++ {
			s += x.Get(i, j)
		}
		means[i] = s / float64(nObs)
	}

	// Compute covariance matrix.
	covData := make([]float64, nVars*nVars)
	for i := 0; i < nVars; i++ {
		for j := i; j < nVars; j++ {
			s := 0.0
			for k := 0; k < nObs; k++ {
				s += (x.Get(i, k) - means[i]) * (x.Get(j, k) - means[j])
			}
			c := s / float64(nObs-1)
			covData[i*nVars+j] = c
			covData[j*nVars+i] = c
		}
	}

	return NewNDArray([]int{nVars, nVars}, covData), nil
}

// Correlate computes the cross-correlation of two 1-D arrays (full mode).
// The result has length len(a) + len(v) - 1.
func Correlate(a, v *NDArray) (*NDArray, error) {
	if a.Ndim() != 1 || v.Ndim() != 1 {
		return nil, fmt.Errorf("numgo.Correlate: inputs must be 1-D")
	}
	na := a.Size()
	nv := v.Size()
	outLen := na + nv - 1
	data := make([]float64, outLen)

	for k := 0; k < outLen; k++ {
		s := 0.0
		for j := 0; j < nv; j++ {
			ai := k - (nv - 1) + j
			if ai >= 0 && ai < na {
				s += a.data[ai] * v.data[nv-1-j]
			}
		}
		data[k] = s
	}
	return FromSlice(data), nil
}

// Convolve computes the discrete linear convolution of two 1-D arrays (full mode).
// The result has length len(a) + len(v) - 1.
func Convolve(a, v *NDArray) (*NDArray, error) {
	if a.Ndim() != 1 || v.Ndim() != 1 {
		return nil, fmt.Errorf("numgo.Convolve: inputs must be 1-D")
	}
	na := a.Size()
	nv := v.Size()
	outLen := na + nv - 1
	data := make([]float64, outLen)

	for k := 0; k < outLen; k++ {
		s := 0.0
		for j := 0; j < nv; j++ {
			ai := k - j
			if ai >= 0 && ai < na {
				s += a.data[ai] * v.data[j]
			}
		}
		data[k] = s
	}
	return FromSlice(data), nil
}
